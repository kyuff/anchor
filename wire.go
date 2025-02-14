package anchor

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// Wire keeps the application running.
// When the returned Context is cancelled, all Components are Closed and the CancelFunc is called.
type Wire interface {
	Wire(ctx context.Context) (context.Context, context.CancelFunc)
}
type WireFunc func(ctx context.Context) (context.Context, context.CancelFunc)

func (fn WireFunc) Wire(ctx context.Context) (context.Context, context.CancelFunc) {
	return fn(ctx)
}

func SignalWire(sig os.Signal, sigs ...os.Signal) Wire {
	return WireFunc(func(ctx context.Context) (context.Context, context.CancelFunc) {
		signals := append([]os.Signal{sig}, sigs...)
		return signal.NotifyContext(ctx, signals...)
	})
}

func DefaultSignalWire() Wire {
	return SignalWire(syscall.SIGINT, syscall.SIGTERM)
}

func TestingWire(m TestingM) Wire {
	return WireFunc(func(ctx context.Context) (context.Context, context.CancelFunc) {
		wireCtx, cancel := context.WithCancel(ctx)
		go func() {
			_ = m.Run()
			cancel()
		}()
		return wireCtx, cancel
	})
}
