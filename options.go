package anchor

import (
	"context"
	"log/slog"
	"time"

	"github.com/kyuff/anchor/internal/logger"
)

type Option func(cfg *Config)

// WithLogger sets the Logger for the application.
//
// Default: No logging is done.
func WithLogger(logger Logger) Option {
	return func(opt *Config) {
		opt.logger = logger
	}
}

// WithNoopLogger disables logging
func WithNoopLogger() Option {
	return WithLogger(logger.Noop{})
}

// WithDefaultSlog uses slog.Default() for logging
func WithDefaultSlog() Option {
	return WithSlog(slog.Default())
}

// WithSlog uses the give *slog.Logger
func WithSlog(log *slog.Logger) Option {
	return WithLogger(
		logger.NewSlog(log),
	)
}

// WithContext runs the Anchor in the given Context. If it is
// cancelled, the Anchor will shutdown.
//
// Default: context.Background()
func WithContext(ctx context.Context) Option {
	return func(cfg *Config) {
		cfg.rootCtx = ctx
	}
}

// WithSetupTimeout fails an Anchor if all Components have not been Setup within
// the timeout provided.
//
// Default: No timeout
func WithSetupTimeout(timeout time.Duration) Option {
	return func(cfg *Config) {
		cfg.setupTimeout = timeout
	}
}

// WithStartTimeout fails an Anchor if all Components have not been Started within
// the timeout provided.
//
// Default: No timeout
func WithStartTimeout(timeout time.Duration) Option {
	return func(cfg *Config) {
		cfg.startTimeout = timeout
	}
}

// WithCloseTimeout is the combined time components have to perform a graceful shutdown.
//
// Default: No timeout
func WithCloseTimeout(timeout time.Duration) Option {
	return func(cfg *Config) {
		cfg.closeTimeout = timeout
	}
}

// WithReady sets a function that is called when all Components have been detected as started.
//
// This is useful for enabling features that require all Components to be ready before they can be used.
// If the function returns an error, the Anchor will fail to start and go into a shutdown state.
func WithReady(fn func(ctx context.Context) error) Option {
	return func(cfg *Config) {
		cfg.onReady = fn
	}
}

func WithReadyCheckBackoff(fn func(ctx context.Context, attempt int) (time.Duration, error)) Option {
	return func(cfg *Config) {
		cfg.readyCheckBackoff = fn
	}
}

// WithFixedReadyCheckBackoff waits a fixed amount of time between retries.
func WithFixedReadyCheckBackoff(d time.Duration) Option {
	return WithReadyCheckBackoff(func(_ context.Context, _ int) (time.Duration, error) {
		return d, nil
	})
}

// WithLinearReadyCheckBackoff increases the wait time linearly with each retry.
func WithLinearReadyCheckBackoff(increment time.Duration) Option {
	return WithReadyCheckBackoff(func(_ context.Context, retries int) (time.Duration, error) {
		return increment * time.Duration(retries), nil
	})
}

// WithExponentialReadyCheckBackoff doubles the wait time with each retry.
func WithExponentialReadyCheckBackoff(base time.Duration) Option {
	return WithReadyCheckBackoff(func(_ context.Context, retries int) (time.Duration, error) {
		return base * time.Duration(1<<retries), nil
	})
}
