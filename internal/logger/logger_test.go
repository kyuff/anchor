package logger_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/kyuff/anchor/internal/assert"
	"github.com/kyuff/anchor/internal/logger"
)

func TestSlog(t *testing.T) {
	t.Run("Info", func(t *testing.T) {
		t.Run("log a formatted message", func(t *testing.T) {
			// arrange
			var (
				handler = &HandlerMock{}
				sut     = logger.NewSlog(slog.New(handler))
			)

			handler.EnabledFunc = func(contextMoqParam context.Context, level slog.Level) bool {
				return true
			}
			handler.HandleFunc = func(ctx context.Context, record slog.Record) error {
				assert.Equal(t, "hello test", record.Message)
				return nil
			}

			// act
			sut.InfofCtx(t.Context(), "hello %s", "test")

			// assert
			assert.Equal(t, 1, len(handler.HandleCalls()))
		})

		t.Run("ignore disabled handler", func(t *testing.T) {
			// arrange
			var (
				handler = &HandlerMock{}
				sut     = logger.NewSlog(slog.New(handler))
			)

			handler.EnabledFunc = func(contextMoqParam context.Context, level slog.Level) bool {
				return false
			}
			handler.HandleFunc = func(ctx context.Context, record slog.Record) error {
				return nil
			}

			// act
			sut.InfofCtx(t.Context(), "hello %s", "test")

			// assert
			assert.Equal(t, 0, len(handler.HandleCalls()))
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("log a formatted message", func(t *testing.T) {
			// arrange
			var (
				handler = &HandlerMock{}
				sut     = logger.NewSlog(slog.New(handler))
			)

			handler.EnabledFunc = func(contextMoqParam context.Context, level slog.Level) bool {
				return true
			}
			handler.HandleFunc = func(ctx context.Context, record slog.Record) error {
				assert.Equal(t, "hello test", record.Message)
				return nil
			}

			// act
			sut.ErrorfCtx(t.Context(), "hello %s", "test")

			// assert
			assert.Equal(t, 1, len(handler.HandleCalls()))
		})

		t.Run("ignore disabled handler", func(t *testing.T) {
			// arrange
			var (
				handler = &HandlerMock{}
				sut     = logger.NewSlog(slog.New(handler))
			)

			handler.EnabledFunc = func(contextMoqParam context.Context, level slog.Level) bool {
				return false
			}
			handler.HandleFunc = func(ctx context.Context, record slog.Record) error {
				return nil
			}

			// act
			sut.ErrorfCtx(t.Context(), "hello %s", "test")

			// assert
			assert.Equal(t, 0, len(handler.HandleCalls()))
		})
	})

	t.Run("Don't break on noop", func(t *testing.T) {
		// arrange
		var (
			sut = logger.Noop{}
		)

		// act
		sut.InfofCtx(t.Context(), "hello %s", "test")
		sut.ErrorfCtx(t.Context(), "hello %s", "test")

	})
}
