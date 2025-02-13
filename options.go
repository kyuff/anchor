package anchor

import (
	"context"
	"log/slog"

	"github.com/kyuff/anchor/internal/logger"
)

type Option func(cfg *Config)

func WithLogger(logger Logger) Option {
	return func(opt *Config) {
		opt.logger = logger
	}
}
func WithNoopLogger() Option {
	return WithLogger(logger.Noop{})
}

func WithDefaultSlog() Option {
	return WithSlog(slog.Default())
}

func WithSlog(log *slog.Logger) Option {
	return WithLogger(
		logger.NewSlog(log),
	)
}

// WithContext runs the Anchor in the given Context. If it is
// cancelled, the Anchor will shutdown.
func WithContext(ctx context.Context) Option {
	return func(cfg *Config) {
		cfg.ctx = ctx
	}
}

func withBackgroundContext() Option {
	return WithContext(context.Background())
}
