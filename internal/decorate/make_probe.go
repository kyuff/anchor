package decorate

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

func MakeProbe[T starter](name string, setup func() (T, error), probe func(ctx context.Context) error) *Component {
	var (
		startTime atomic.Int64
		inner     starter
		err       error
	)
	return &Component{
		name: func() string {
			return name
		},
		setup: func(ctx context.Context) error {
			if setup == nil {
				return fmt.Errorf("nil setup func")
			}

			if probe == nil {
				return fmt.Errorf("nil probe func")
			}

			inner, err = setup()
			if err != nil {
				return err
			}
			if isNil(inner) {
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

			return fn(inner)
		},

		close: func(ctx context.Context) error {
			if isNil(inner) {
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

			return fn(inner)
		},

		start: func(ctx context.Context) error {
			if isNil(inner) {
				return fmt.Errorf("nil make component")
			}

			startTime.Store(time.Now().UnixMilli())
			return inner.Start(ctx)
		},

		probe: func() func(ctx context.Context) error {
			return func(ctx context.Context) error {
				if !probeIsReady(startTime.Load()) {
					return fmt.Errorf("component %s is not started yet", name)
				}

				return probe(ctx)
			}
		}(),
	}
}
