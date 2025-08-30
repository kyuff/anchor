package decorate

import (
	"context"
	"fmt"
	"reflect"
	"sync/atomic"
	"time"
)

func makeComponent[T starter](c *Component, name string, setup func() (T, error), probe func(ctx context.Context) error) *Component {
	var (
		startTime = &atomic.Int64{}
	)

	c.setup = makeSetup(c, setup)
	c.name = makeName(name)
	c.close = makeClose(c)
	c.start = makeStart(c, startTime)
	c.probe = makeProbe(startTime, probe)

	return c
}

func makeSetup[T starter](c *Component, setup func() (T, error)) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		if setup == nil {
			return fmt.Errorf("nil setup func")
		}

		var err error
		c.inner, err = setup()
		if err != nil {
			return err
		}
		if isNil(c.inner) {
			err = fmt.Errorf("nil make component")
			return err
		}

		fn := func(component any) error {
			if c, ok := component.(contextSetupper); ok {
				return c.Setup(ctx)
			}

			if c, ok := component.(setupper); ok {
				return c.Setup()
			}

			return nil
		}

		return fn(c.inner)
	}
}

func makeSetupCall(c *Component) func(ctx context.Context) error {
	if c, ok := c.inner.(setupper); ok {
		return func(_ context.Context) error {
			return c.Setup()
		}
	}
	if c, ok := c.inner.(contextSetupper); ok {
		return c.Setup
	}
	return func(ctx context.Context) error { return nil }
}

func makeClose(c *Component) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		if isNil(c.inner) {
			return fmt.Errorf("nil make component")
		}

		fn := func(component any) error {
			if c, ok := component.(contextCloser); ok {
				return c.Close(ctx)
			}

			if c, ok := component.(closer); ok {
				return c.Close()
			}

			return nil
		}

		return fn(c.inner)
	}
}

func makeStart(c *Component, startTime *atomic.Int64) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		if isNil(c.inner) {
			return fmt.Errorf("nil make component")
		}

		startTime.Store(time.Now().UnixMilli())
		return c.inner.Start(ctx)
	}
}

func makeProbe(startTime *atomic.Int64, probe func(ctx context.Context) error) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		if !probeIsReady(startTime.Load()) {
			return fmt.Errorf("component is not started yet")
		}

		return probe(ctx)
	}
}

func probeInner(c *Component) func(ctx context.Context) error {
	return func(ctx context.Context) error {

		if c, ok := c.inner.(contextProber); ok {
			return c.Probe(ctx)
		}

		return nil
	}
}

func makeName(name string) func() string {
	return func() string {
		return name
	}
}

func makeNameCall(c *Component) func() string {
	return func() string {
		if c, ok := c.inner.(namer); ok {
			return c.Name()
		}
		return fmt.Sprintf("%T", c.inner)
	}
}

func isNil(a any) bool {
	defer func() { recover() }()
	return a == nil || reflect.ValueOf(a).IsNil()
}
