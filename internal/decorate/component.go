package decorate

import (
	"context"
	"fmt"
)

type setupper interface {
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

type namer interface {
	starter
	Name() string
}

type fullComponent interface {
	setupper
	starter
	closer
	namer
}

func New(component starter) *Component {
	return &Component{
		start: component.Start,
		setup: func() func(ctx context.Context) error {
			if c, ok := component.(setupper); ok {
				return c.Setup
			}
			return func(ctx context.Context) error { return nil }
		}(),
		close: func() func() error {
			if c, ok := component.(closer); ok {
				return c.Close
			}
			return func() error { return nil }
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
	close func() error
}

func (c *Component) Start(ctx context.Context) error {
	return c.start(ctx)
}

func (c *Component) Setup(ctx context.Context) error {
	return c.setup(ctx)
}

func (c *Component) Close() error {
	return c.close()
}

func (c *Component) Name() string {
	return c.name()
}
