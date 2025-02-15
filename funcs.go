package anchor

import "context"

func Start(name string, start func(ctx context.Context) error) Component {
	return &decoratedComponent{
		name:  func() string { return name },
		start: start,
		setup: func(ctx context.Context) error { return nil },
		close: func() error { return nil },
	}
}

func Setup(name string, setup func(ctx context.Context) error) Component {
	return &decoratedComponent{
		name:  func() string { return name },
		start: func(ctx context.Context) error { return nil },
		setup: setup,
		close: func() error { return nil },
	}
}

func Close(name string, close func() error) Component {
	return &decoratedComponent{
		name:  func() string { return name },
		start: func(ctx context.Context) error { return nil },
		setup: func(ctx context.Context) error { return nil },
		close: close,
	}
}
