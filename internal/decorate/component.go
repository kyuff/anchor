package decorate

import "context"

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

type Component struct {
	name  func() string
	start func(ctx context.Context) error
	probe func(ctx context.Context) error
	setup func(ctx context.Context) error
	close func(ctx context.Context) error

	inner starter
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
