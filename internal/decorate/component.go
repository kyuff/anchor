package decorate

import "context"

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
