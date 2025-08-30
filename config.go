package anchor

import (
	"context"
	"time"
)

type config struct {
	logger Logger
	// anchorCtx is used to derive the setup, start and close contexts.
	anchorCtx         context.Context
	setupTimeout      time.Duration
	startTimeout      time.Duration
	closeTimeout      time.Duration
	onReady           func(ctx context.Context) error
	readyCheckBackoff func(ctx context.Context, attempt int) (time.Duration, error)
}

func defaultOptions() *config {
	return applyOptions(&config{},
		// add default options here
		WithNoopLogger(),
		WithAnchorContext(context.Background()),
		WithSetupTimeout(0), // no timeout
		WithStartTimeout(0), // no timeout
		WithCloseTimeout(10*time.Second),
		WithReadyCallback(func(ctx context.Context) error { return nil }),
		WithLinearReadyCheckBackoff(time.Millisecond*100),
	)

}

func applyOptions(options *config, opts ...Option) *config {
	for _, opt := range opts {
		opt(options)
	}

	return options
}
