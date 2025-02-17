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

// WithCloseTimeout is the combined time components have to perform a graceful shutdown.
//
// Default: No timeout
func WithCloseTimeout(timeout time.Duration) Option {
	return func(cfg *Config) {
		cfg.closeTimeout = timeout
	}
}
