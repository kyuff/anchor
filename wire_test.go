package anchor_test

import (
	"context"
	"os"
	"syscall"
	"testing"

	"github.com/kyuff/anchor"
	"github.com/kyuff/anchor/internal/assert"
)

func TestWire(t *testing.T) {
	t.Run("TestingWire", func(t *testing.T) {
		// arrange
		var (
			m   = &TestingMMock{}
			sut = anchor.TestingWire(m)
		)

		m.RunFunc = func() int {
			return 0
		}

		// act
		ctx, cancel := sut.Wire(t.Context())

		// assert
		t.Cleanup(cancel)
		<-ctx.Done()
		// test will block forever if the context is not cancelled
		assert.Equal(t, 1, len(m.RunCalls()))
	})

	t.Run("SignalWire", func(t *testing.T) {
		t.Run("exit when cancel func is called", func(t *testing.T) {
			// arrange
			var (
				testCtx, testCancel = context.WithCancel(t.Context())
				sut                 = anchor.DefaultSignalWire()
			)

			go func() {
				testCancel()
			}()

			// acct
			ctx, cancel := sut.Wire(testCtx)

			// assert
			defer cancel()
			<-ctx.Done()
		})

		t.Run("exit when signal is made", func(t *testing.T) {
			// arrange
			var (
				testCtx, testCancel = context.WithCancel(t.Context())
				sut                 = anchor.SignalWire(syscall.SIGALRM)
			)

			go func() {
				p, err := os.FindProcess(os.Getpid())
				assert.NoError(t, err)
				assert.NoError(t, p.Signal(syscall.SIGALRM))
			}()

			// acct
			ctx, cancel := sut.Wire(testCtx)

			// assert
			defer cancel()
			defer testCancel()
			<-ctx.Done()
		})
	})
}
