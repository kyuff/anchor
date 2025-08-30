package anchor

import (
	"context"
	"log/slog"
	"time"

	"github.com/kyuff/anchor/internal/logger"
)

// Option for an Anchor configuration.
//
// See config.go for defaults.
type Option func(cfg *config)

// WithLogger sets the Logger for the application.
//
// Default: No logging is done.
func WithLogger(logger Logger) Option {
	return func(opt *config) {
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

// WithAnchorContext runs the Anchor in the given Context. If it is
// canceled, the Anchor will shutdown.
//
// Default: context.Background()
func WithAnchorContext(ctx context.Context) Option {
	return func(cfg *config) {
		cfg.anchorCtx = ctx
	}
}

// WithSetupTimeout fails an Anchor if all Components have not been Setup within
// the timeout provided.
//
// Default: No timeout
func WithSetupTimeout(timeout time.Duration) Option {
	return func(cfg *config) {
		cfg.setupTimeout = timeout
	}
}

// WithStartTimeout fails an Anchor if all Components have not been Started within
// the timeout provided.
//
// Default: No timeout
func WithStartTimeout(timeout time.Duration) Option {
	return func(cfg *config) {
		cfg.startTimeout = timeout
	}
}

// WithCloseTimeout is the combined time components have to perform a graceful shutdown.
//
// Default: No timeout
func WithCloseTimeout(timeout time.Duration) Option {
	return func(cfg *config) {
		cfg.closeTimeout = timeout
	}
}

// WithReadyCallback is called when all Components have been Setup and Started and succeeded any Probe.
//
// This is useful for enabling features that require all Components to be ready before they can be used.
// If the function returns an error, the Anchor will fail to start and go into a shutdown state.
func WithReadyCallback(fn func(ctx context.Context) error) Option {
	return func(cfg *config) {
		cfg.onReady = fn
	}
}

// WithReadyCheckBackoff is called when a Component fails a Probe.
//
// The function should return the amount of time to wait before retrying.
// If the function returns an error, the Anchor will fail to start and go into a shutdown state.
//
// Default: linear backoff with 100 millisecond increment
func WithReadyCheckBackoff(fn func(ctx context.Context, attempt int) (time.Duration, error)) Option {
	return func(cfg *config) {
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
