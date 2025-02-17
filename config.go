package anchor

import (
	"context"
	"time"
)

type Config struct {
	logger Logger
	// rootCtx is used to derive the setup, start and close contexts.
	rootCtx      context.Context
	setupTimeout time.Duration
	closeTimeout time.Duration
}

func defaultOptions() *Config {
	return applyOptions(&Config{},
		// add default options here
		WithNoopLogger(),
		WithContext(context.Background()),
		WithSetupTimeout(0), // no timeout
		WithCloseTimeout(10*time.Second),
	)

}

func applyOptions(options *Config, opts ...Option) *Config {
	for _, opt := range opts {
		opt(options)
	}

	return options
}
