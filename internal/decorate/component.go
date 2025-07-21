package decorate

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

type setupper interface {
	starter
	Setup() error
}

type contextSetupper interface {
	starter
	Setup(ctx context.Context) error
}

type starter interface {
	Start(ctx context.Context) error
}

type closer interface {
	starter
	Close() error
}

type contextCloser interface {
	starter
	Close(ctx context.Context) error
}

type contextProber interface {
	starter
	Probe(ctx context.Context) error
}

type namer interface {
	starter
	Name() string
}

type fullComponent interface {
	contextSetupper
	starter
	contextProber
	contextCloser
	namer
}

var delay = (time.Millisecond * 15).Milliseconds()

func probeIsReady(startTime int64) bool {
	var (
		now       = time.Now().UnixMilli()
		readyTime = startTime + delay
	)

	return startTime > 0 && readyTime < now
}

func New(component starter) *Component {
	var startTime atomic.Int64
	var c *Component
	c = &Component{
		start: func(ctx context.Context) error {
			startTime.Store(time.Now().UnixMilli())
			return component.Start(ctx)
		},
		probe: func() func(ctx context.Context) error {
			return func(ctx context.Context) error {
				if !probeIsReady(startTime.Load()) {
					return fmt.Errorf("component %s is not started yet", c.Name())
				}

				if c, ok := component.(contextProber); ok {
					return c.Probe(ctx)
				}
				return nil
			}
		}(),
		setup: func() func(ctx context.Context) error {
			if c, ok := component.(setupper); ok {
				return func(_ context.Context) error {
					return c.Setup()
				}
			}
			if c, ok := component.(contextSetupper); ok {
				return c.Setup
			}
			return func(ctx context.Context) error { return nil }
		}(),
		close: func() func(ctx context.Context) error {
			if c, ok := component.(closer); ok {
				return func(_ context.Context) error {
					return c.Close()
				}
			}
			if c, ok := component.(contextCloser); ok {
				return c.Close
			}

			return func(_ context.Context) error { return nil }
		}(),
		name: func() func() string {
			if c, ok := component.(namer); ok {
				return c.Name
			}
			return func() string { return fmt.Sprintf("%T", component) }
		}(),
	}

	return c
}

type Component struct {
	name  func() string
	start func(ctx context.Context) error
	probe func(ctx context.Context) error
	setup func(ctx context.Context) error
	close func(ctx context.Context) error
}

func (c *Component) Start(ctx context.Context) error {
	return c.start(ctx)
}

func (c *Component) Probe(ctx context.Context) error { return c.probe(ctx) }

func (c *Component) Setup(ctx context.Context) error {
	return c.setup(ctx)
}

func (c *Component) Close(ctx context.Context) error {
	return c.close(ctx)
}

func (c *Component) Name() string {
	return c.name()
}
