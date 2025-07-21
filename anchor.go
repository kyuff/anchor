package anchor

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/kyuff/anchor/internal/decorate"
	"golang.org/x/sync/errgroup"
)

func New(wire Wire, opts ...Option) *Anchor {
	return &Anchor{
		cfg:       applyOptions(defaultOptions(), opts...),
		wire:      wire,
		closeChan: make(chan int, 1),
	}
}

type Anchor struct {
	// Wire controls the time that the application runs.
	// The application shuts down when the context returned is cancelled.
	wire       Wire
	cfg        *Config
	components []fullComponent
	running    atomic.Bool
	// index to the last component that was setup
	// used to be able to close in reverse order
	setupIndex int

	closeChan chan int
}

// Add will manage the Component list by the Anchor.
//
// When Run is called, all Components will be started in the order they were
// given to the Anchor.
func (a *Anchor) Add(components ...Component) *Anchor {
	if a.running.Load() {
		// even though panic is frowned upon, this is one of the few places
		// it makes sense to do, as it's part of application setup nad should fail fast.
		panic("cannot add components after Run is called")
	}

	for _, component := range components {
		if component == nil {
			panic("cannot add nil component")
		}

		a.components = append(a.components, decorate.New(component))
	}

	return a
}

func (a *Anchor) Run() int {
	if !a.running.CompareAndSwap(false, true) {
		panic("anchor is already running")
	}

	if len(a.components) == 0 {
		a.cfg.logger.ErrorfCtx(a.cfg.rootCtx, "No components added. Aborting ...")
		return OK
	}

	// wire the anchor context
	ctx, cancel := a.wire.Wire(a.cfg.rootCtx)
	defer cancel()

	closed := make(chan int)
	// monitor closeChan
	go func() {
		var code int
		select {
		case <-ctx.Done():
			code = OK
		case c := <-a.closeChan:
			code = c
		}

		closeCode := a.closeAll(context.Background())
		if code != OK {
			closed <- code
		} else {
			closed <- closeCode
		}
	}()

	if code := a.setupAll(ctx); code != OK {
		a.signalClose(code)
	} else {
		go a.startAll(ctx)
	}

	return <-closed
}

func (a *Anchor) signalClose(code int) {
	a.closeChan <- code
}

func (a *Anchor) startAll(ctx context.Context) {
	g, startCtx := errgroup.WithContext(ctx)

	for _, component := range a.components {
		g.Go(func() (err error) {
			return a.startComponent(startCtx, component)
		})
	}

	err := a.probeAll(ctx)
	if err != nil {
		a.cfg.logger.ErrorfCtx(ctx, "[anchor] Ready check failed: %v", err)
		a.signalClose(Internal)
		return
	}

	err = a.cfg.onReady(ctx)
	if err != nil {
		a.signalClose(Internal)
		return
	}

	err = g.Wait()
	if err != nil {
		a.signalClose(Internal)
	} else {
		a.signalClose(OK)
	}
}

func (a *Anchor) startComponent(ctx context.Context, component fullComponent) (err error) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			a.cfg.logger.ErrorfCtx(ctx, "[anchor] Start panic for %q: %v", component.Name(), panicErr)
			err = errors.Join(err, fmt.Errorf("%s", panicErr))
		}
	}()

	a.cfg.logger.InfofCtx(ctx, "[anchor] Start %q", component.Name())
	err = component.Start(ctx)
	if err != nil {
		a.cfg.logger.ErrorfCtx(ctx, "[anchor] Start failed for %q: %v", component.Name(), err)
		return err
	}

	a.cfg.logger.InfofCtx(ctx, "[anchor] Component exit: %q", component.Name())
	return nil
}

func (a *Anchor) probeAll(ctx context.Context) error {
	var cancel context.CancelFunc = func() {}
	if a.cfg.startTimeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, a.cfg.startTimeout)
	}

	g, probeCtx := errgroup.WithContext(ctx)
	defer cancel()

	for _, component := range a.components {
		g.Go(func() (err error) {
			return a.probeComponent(probeCtx, component)
		})
	}

	return g.Wait()
}

func (a *Anchor) probeComponent(ctx context.Context, component fullComponent) (err error) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			a.cfg.logger.ErrorfCtx(ctx, "[anchor] Probe panic for %q: %v", component.Name(), panicErr)
			err = errors.Join(err, fmt.Errorf("%s", panicErr))
		}
	}()

	var attempts int
	var backoff time.Duration
	for {
		attempts++
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := component.Probe(ctx)
			if err == nil {
				return nil
			}

			backoff, err = a.cfg.readyCheckBackoff(ctx, attempts)
			if err != nil {
				return err
			}

			time.Sleep(backoff)
		}
	}

}

func (a *Anchor) setupAll(ctx context.Context) int {
	var cancel context.CancelFunc = func() {}
	if a.cfg.setupTimeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, a.cfg.setupTimeout)
	}
	defer cancel()

	for index := 0; index < len(a.components); index++ {
		a.setupIndex = index
		code := a.setupComponent(ctx, a.components[index])
		if code != OK {
			return code
		}
	}

	return OK
}

func (a *Anchor) setupComponent(ctx context.Context, component fullComponent) (code int) {

	done := make(chan int, 1)
	go func() {
		defer func() {
			if panicErr := recover(); panicErr != nil {
				a.cfg.logger.ErrorfCtx(a.cfg.rootCtx, "[anchor] Setup %q panic: %v", component.Name(), panicErr)
				done <- SetupFailed
			}
		}()

		err := component.Setup(ctx)
		if err != nil {
			a.cfg.logger.ErrorfCtx(a.cfg.rootCtx, "[anchor] Setup %q failed: %v", component.Name(), err)
			done <- SetupFailed
			return
		}

		a.cfg.logger.InfofCtx(a.cfg.rootCtx, "[anchor] Setup %q", component.Name())
		done <- OK
	}()

	select {
	case code = <-done:
		return code
	case <-ctx.Done():
		return Interrupted
	}

}

func (a *Anchor) closeAll(ctx context.Context) int {
	done := make(chan int, 1)
	ctx, cancel := context.WithTimeout(ctx, a.cfg.closeTimeout)
	go func() {
		defer cancel()

		for index := a.setupIndex; index >= 0; index-- {
			a.closeComponent(ctx, a.components[index])
			a.setupIndex = index
		}

		done <- OK
	}()

	select {
	case code := <-done:
		return code
	case <-ctx.Done():
		return Interrupted
	}
}

func (a *Anchor) closeComponent(ctx context.Context, component fullComponent) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			a.cfg.logger.ErrorfCtx(a.cfg.rootCtx, "[anchor] Close %q panic: %v", component.Name(), panicErr)
		}
	}()

	err := component.Close(ctx)
	if err != nil {
		a.cfg.logger.ErrorfCtx(a.cfg.rootCtx, "[anchor] Closed component %s: %v", component.Name(), err)
	}

	a.cfg.logger.InfofCtx(a.cfg.rootCtx, "[anchor] Closed component %s", component.Name())
}
