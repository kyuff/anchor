package anchor

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// Wire keeps the Anchor running.
//
// When the returned Context is canceled, all Components are Closed and the CancelFunc is called.
type Wire interface {
	Wire(ctx context.Context) (context.Context, context.CancelFunc)
}

// WireFunc is a convenience type for creating a Wire from a function.
type WireFunc func(ctx context.Context) (context.Context, context.CancelFunc)

// Wire calls WireFunc.
func (fn WireFunc) Wire(ctx context.Context) (context.Context, context.CancelFunc) {
	return fn(ctx)
}

// SignalWire returns a Wire that listens for the given os signals.
//
// Use [DefaultSignalWire] to listen for SIGINT and SIGTERM.
func SignalWire(sig os.Signal, sigs ...os.Signal) Wire {
	return WireFunc(func(ctx context.Context) (context.Context, context.CancelFunc) {
		signals := append([]os.Signal{sig}, sigs...)
		return signal.NotifyContext(ctx, signals...)
	})
}

// DefaultSignalWire returns a Wire that listens for SIGINT and SIGTERM.
func DefaultSignalWire() Wire {
	return SignalWire(syscall.SIGINT, syscall.SIGTERM)
}

// TestingWire returns a Wire for use in testing.
// It will Run the tests and and then signal the application to shutdown.
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
