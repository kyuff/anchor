package decorate_test

import (
	"context"
	"testing"
	"time"

	"github.com/kyuff/anchor/internal/assert"
	"github.com/kyuff/anchor/internal/decorate"
)

func TestSetup(t *testing.T) {
	t.Run("call setup constructor", func(t *testing.T) {
		// arrange
		var (
			called = false
		)

		// act
		sut := decorate.Setup("TEST NAME", func() error {
			called = true
			return nil
		})

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

func TestSetupContext(t *testing.T) {
	t.Run("call setup constructor with context", func(t *testing.T) {
		// arrange
		var (
			called  = false
			gotCtx  context.Context
		)

		// act
		sut := decorate.SetupContext("TEST NAME", func(ctx context.Context) error {
			called = true
			gotCtx = ctx
			return nil
		})

		// assert
		assert.NoError(t, sut.Setup(t.Context()))
		assert.NoError(t, sut.Start(t.Context()))
		assert.NoErrorEventually(t, time.Second, func() error {
			return sut.Probe(t.Context())
		})
		assert.NoError(t, sut.Close(t.Context()))
		assert.Equal(t, "TEST NAME", sut.Name())
		assert.Truef(t, called, "not called")
		assert.Truef(t, gotCtx == t.Context(), "expected context to be forwarded")
	})
}
