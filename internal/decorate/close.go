package decorate

import (
	"context"
	"io"
)

func Close(name string, closer io.Closer) *Component {
	return &Component{
		name: func() string {
			return name
		},
		start: func(ctx context.Context) error { return nil },
		setup: func(_ context.Context) error { return nil },
		close: func(_ context.Context) error {
			return closer.Close()
		},
		probe: func(ctx context.Context) error { return nil },
	}
}
