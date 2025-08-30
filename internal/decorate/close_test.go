package decorate_test

import (
	"testing"
	"time"

	"github.com/kyuff/anchor/internal/assert"
	"github.com/kyuff/anchor/internal/decorate"
)

type closeFunc func() error

func (fn closeFunc) Close() error {
	return fn()
}

func TestClose(t *testing.T) {
	t.Run("call setup constructor", func(t *testing.T) {
		// arrange
		var (
			called = false
		)

		// act
		sut := decorate.Close("TEST NAME", closeFunc(func() error {
			called = true
			return nil
		}))

		// assert
		assert.NoError(t, sut.Setup(t.Context()))
		assert.NoError(t, sut.Start(t.Context()))
		assert.NoErrorEventually(t, time.Second, func() error {
			return sut.Probe(t.Context())
		})
		assert.NoError(t, sut.Close(t.Context()))
		assert.Equal(t, "TEST NAME", sut.Name())
		assert.Truef(t, called, "not called")
	})
}
