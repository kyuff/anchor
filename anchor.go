package anchor

import (
	"sync/atomic"

	"golang.org/x/sync/errgroup"
)

func New(opts ...Option) *Anchor {
	return &Anchor{
		cfg: applyOptions(defaultOptions(), opts...),
	}
}

type Anchor struct {
	cfg        *Config
	components []fullComponent
	running    atomic.Bool
	// index to the last component that was setup
	// used to be able to close in reverse order
	setupIndex int
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

		a.components = append(a.components, decorateComponent(component))
	}

	return a
}

func (a *Anchor) Run() int {
	if !a.running.CompareAndSwap(false, true) {
		panic("anchor is already running")
	}

	err := a.setupAll()
	if err != nil {
		return 1
	}

	err = a.startAll()
	if err != nil {
		return 1
	}

	err = a.closeAll()
	if err != nil {
		return 1
	}

	return 0
}

func (a *Anchor) startAll() error {
	g, ctx := errgroup.WithContext(a.cfg.ctx)

	for _, component := range a.components {
		g.Go(func() error {
			a.cfg.logger.InfofCtx(a.cfg.ctx, "[anchor] Starting component %s: %v", component.Name())
			err := component.Start(ctx)
			if err != nil {
				a.cfg.logger.ErrorfCtx(a.cfg.ctx, "[anchor] Failed to start component %s: %v", component.Name(), err)
				return err
			}

			a.cfg.logger.InfofCtx(a.cfg.ctx, "[anchor] Component exit: %s", component.Name())
			return nil
		})
	}

	return g.Wait()
}

func (a *Anchor) setupAll() error {
	for ; a.setupIndex < len(a.components); a.setupIndex++ {
		component := a.components[a.setupIndex]
		err := component.Setup(a.cfg.ctx)
		if err != nil {
			a.cfg.logger.ErrorfCtx(a.cfg.ctx, "[anchor] Setup component %s: %v", component.Name(), err)
			return err
		}

		a.cfg.logger.ErrorfCtx(a.cfg.ctx, "[anchor] Setup component %s", component.Name())
	}

	return nil
}

func (a *Anchor) closeAll() error {
	for ; a.setupIndex > 0; a.setupIndex-- {
		component := a.components[a.setupIndex-1]
		err := component.Close()
		if err != nil {
			a.cfg.logger.ErrorfCtx(a.cfg.ctx, "[anchor] Closed component %s: %v", component.Name(), err)
			continue
		}

		a.cfg.logger.InfofCtx(a.cfg.ctx, "[anchor] Closed component %s", component.Name())
	}

	return nil
}
