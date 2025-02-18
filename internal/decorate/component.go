package decorate

import (
	"context"
	"fmt"
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

type namer interface {
	starter
	Name() string
}

type fullComponent interface {
	contextSetupper
	starter
	contextCloser
	namer
}

func New(component starter) *Component {
	return &Component{
		start: component.Start,
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
}

var _ fullComponent = (*Component)(nil)

type Component struct {
	name  func() string
	start func(ctx context.Context) error
	setup func(ctx context.Context) error
	close func(ctx context.Context) error
}

func (c *Component) Start(ctx context.Context) error {
	return c.start(ctx)
}

func (c *Component) Setup(ctx context.Context) error {
	return c.setup(ctx)
}

func (c *Component) Close(ctx context.Context) error {
	return c.close(ctx)
}

func (c *Component) Name() string {
	return c.name()
}
