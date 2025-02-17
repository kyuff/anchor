package anchor

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/kyuff/anchor/internal/logger"
)

func TestOptions(t *testing.T) {
	testCases := []struct {
		name   string
		option Option
		assert func(t *testing.T, cfg *Config)
	}{
		{
			name:   "WithLogger",
			option: WithLogger(logger.Noop{}),
			assert: func(t *testing.T, cfg *Config) {
				if cfg.logger == nil {
					t.Error("expected logger")
				}
			},
		},
		{
			name:   "WithNoopLogger",
			option: WithNoopLogger(),
			assert: func(t *testing.T, cfg *Config) {
				if cfg.logger == nil {
					t.Error("expected logger")
				}
			},
		},
		{
			name:   "WithDefaultSlog",
			option: WithDefaultSlog(),
			assert: func(t *testing.T, cfg *Config) {
				if cfg.logger == nil {
					t.Error("expected logger")
				}
			},
		},
		{
			name:   "WithSlog",
			option: WithSlog(slog.Default()),
			assert: func(t *testing.T, cfg *Config) {
				if cfg.logger == nil {
					t.Error("expected logger")
				}
			},
		},
		{
			name:   "WithContext",
			option: WithContext(context.Background()),
			assert: func(t *testing.T, cfg *Config) {
				if cfg.rootCtx == nil {
					t.Error("expected root context")
				}
			},
		},
		{
			name:   "WithSetupTimeout",
			option: WithSetupTimeout(time.Hour),
			assert: func(t *testing.T, cfg *Config) {
				if cfg.setupTimeout != time.Hour {
					t.Error("expected setup timeout")
				}
			},
		},
		{
			name:   "WithCloseTimeout",
			option: WithCloseTimeout(time.Hour),
			assert: func(t *testing.T, cfg *Config) {
				if cfg.closeTimeout != time.Hour {
					t.Error("expected close timeout")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			var (
				cfg = &Config{}
			)

			// act
			tc.option(cfg)

			// assert
			tc.assert(t, cfg)
		})
	}
}
