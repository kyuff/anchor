package decorate

import "context"

func MakeProbe[T starter](name string, setup func() (T, error), probe func(ctx context.Context) error) *Component {
	return nil
}
