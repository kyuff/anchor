package decorate

import "context"

func Setup(name string, setup func() error) *Component {
	return &Component{
		name: func() string {
			return name
		},
		start: func(ctx context.Context) error { return nil },
		setup: func(_ context.Context) error {
			return setup()
		},
		close: func(_ context.Context) error { return nil },
	}
}
