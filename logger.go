package anchor

import "context"

// Logger is an interface for logging by the Anchor.
// There is default implementation for slog.Logger.
type Logger interface {
	// InfofCtx logs on an INFO level
	InfofCtx(ctx context.Context, template string, args ...any)
	// ErrorfCtx logs on an ERROR level
	ErrorfCtx(ctx context.Context, template string, args ...any)
}
