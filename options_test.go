package anchor

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/kyuff/anchor/internal/assert"
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
		{
			name: "WithReadyCallback",
			option: WithReadyCallback(func(ctx context.Context) error {
				return errors.New("WithReadyError")
			}),
			assert: func(t *testing.T, cfg *Config) {
				if assert.NotNil(t, cfg.onReady) {
					assert.Equal(t, "WithReadyError", cfg.onReady(context.Background()).Error())
				}
			},
		},
		{
			name: "WithReadyCheckBackoff",
			option: WithReadyCheckBackoff(func(ctx context.Context, attempt int) (time.Duration, error) {
				return time.Second * time.Duration(attempt), errors.New("WithReadyCheckBackoff")
			}),
			assert: func(t *testing.T, cfg *Config) {
				if assert.NotNil(t, cfg.readyCheckBackoff) {
					backoff, err := cfg.readyCheckBackoff(context.Background(), 5)
					assert.Equal(t, time.Second*5, backoff)
					assert.Equal(t, "WithReadyCheckBackoff", err.Error())
				}
			},
		},
		{
			name:   "WithFixedReadyCheckBackoff",
			option: WithFixedReadyCheckBackoff(42 * time.Second),
			assert: func(t *testing.T, cfg *Config) {
				if assert.NotNil(t, cfg.readyCheckBackoff) {
					backoff, err := cfg.readyCheckBackoff(context.Background(), 5)
					assert.Equal(t, time.Second*42, backoff)
					assert.NoError(t, err)
				}
			},
		},
		{
			name:   "WithLinearReadyCheckBackoff",
			option: WithLinearReadyCheckBackoff(2 * time.Second),
			assert: func(t *testing.T, cfg *Config) {
				if assert.NotNil(t, cfg.readyCheckBackoff) {
					backoff, err := cfg.readyCheckBackoff(context.Background(), 5)
					assert.Equal(t, time.Second*2*5, backoff)
					assert.NoError(t, err)
				}
			},
		},
		{
			name:   "WithExponentialReadyCheckBackoff",
			option: WithExponentialReadyCheckBackoff(time.Second),
			assert: func(t *testing.T, cfg *Config) {
				if assert.NotNil(t, cfg.readyCheckBackoff) {
					backoff, err := cfg.readyCheckBackoff(context.Background(), 5)
					assert.Equal(t, time.Second*32, backoff) // 2^5 = 32
					assert.NoError(t, err)
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
