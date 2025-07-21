package anchor

import (
	"context"
	"time"
)

type Config struct {
	logger Logger
	// rootCtx is used to derive the setup, start and close contexts.
	rootCtx           context.Context
	setupTimeout      time.Duration
	startTimeout      time.Duration
	closeTimeout      time.Duration
	onReady           func(ctx context.Context) error
	readyCheckBackoff func(ctx context.Context, attempt int) (time.Duration, error)
}

func defaultOptions() *Config {
	return applyOptions(&Config{},
		// add default options here
		WithNoopLogger(),
		WithContext(context.Background()),
		WithSetupTimeout(0), // no timeout
		WithStartTimeout(0), // no timeout
		WithCloseTimeout(10*time.Second),
		WithReadyCallback(func(ctx context.Context) error { return nil }),
		WithLinearReadyCheckBackoff(time.Millisecond*100),
	)

}

func applyOptions(options *Config, opts ...Option) *Config {
	for _, opt := range opts {
		opt(options)
	}

	return options
}
