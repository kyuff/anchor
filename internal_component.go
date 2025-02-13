package anchor

import (
	"context"
	"fmt"
)

func decorateComponent(c Component) *decoratedComponent {
	return &decoratedComponent{
		start: c.Start,
		setup: func() func(ctx context.Context) error {
			if setupper, ok := c.(setupComponent); ok {
				return setupper.Setup
			}
			return func(ctx context.Context) error { return nil }
		}(),
		close: func() func() error {
			if closer, ok := c.(closeComponent); ok {
				return closer.Close
			}
			return func() error { return nil }
		}(),
		name: func() func() string {
			if namer, ok := c.(namedComponent); ok {
				return namer.Name
			}
			return func() string { return fmt.Sprintf("%T", c) }
		}(),
	}
}

var _ fullComponent = (*decoratedComponent)(nil)

type decoratedComponent struct {
	name  func() string
	start func(ctx context.Context) error
	setup func(ctx context.Context) error
	close func() error
}

func (c *decoratedComponent) Start(ctx context.Context) error {
	return c.start(ctx)
}

func (c *decoratedComponent) Setup(ctx context.Context) error {
	return c.setup(ctx)
}

func (c *decoratedComponent) Close() error {
	return c.close()
}

func (c *decoratedComponent) Name() string {
	return c.name()
}
