package logger

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"time"
)

type Noop struct{}

func (Noop) InfofCtx(ctx context.Context, template string, args ...any)  {}
func (Noop) ErrorfCtx(ctx context.Context, template string, args ...any) {}

func NewSlog(logger *slog.Logger) *Slog {
	return &Slog{
		logger: logger,
	}
}

type Slog struct {
	logger *slog.Logger
}

func (log *Slog) InfofCtx(ctx context.Context, format string, args ...any) {
	if !log.logger.Enabled(ctx, slog.LevelInfo) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(2, pcs[:]) // skip [Callers, InfofCtx]
	r := slog.NewRecord(time.Now(), slog.LevelInfo, fmt.Sprintf(format, args...), pcs[0])
	_ = log.logger.Handler().Handle(ctx, r)
}

func (log *Slog) ErrorfCtx(ctx context.Context, format string, args ...any) {
	if !log.logger.Enabled(ctx, slog.LevelError) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(2, pcs[:]) // skip [Callers, ErrorfCtx]
	r := slog.NewRecord(time.Now(), slog.LevelError, fmt.Sprintf(format, args...), pcs[0])
	_ = log.logger.Handler().Handle(ctx, r)
}
