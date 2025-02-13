package anchor

import "context"

type Logger interface {
	InfofCtx(ctx context.Context, template string, args ...any)
	ErrorfCtx(ctx context.Context, template string, args ...any)
}
