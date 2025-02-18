package decorate

import (
	"context"
	"fmt"
	"reflect"
)

func Make[T starter](name string, setup func() (T, error)) *Component {
	var (
		inner starter
		err   error
	)
	return &Component{
		name: func() string {
			return name
		},
		setup: func(ctx context.Context) error {
			if setup == nil {
				return fmt.Errorf("nil setup func")
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

			return inner.Start(ctx)
		},
	}
}

func isNil(a any) bool {
	defer func() { recover() }()
	return a == nil || reflect.ValueOf(a).IsNil()
}
