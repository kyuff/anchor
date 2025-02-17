package anchor_test

import (
	"context"
	"testing"

	"github.com/kyuff/anchor"
	"github.com/kyuff/anchor/internal/assert"
)

func TestComponent(t *testing.T) {
	t.Run("Setup", func(t *testing.T) {
		// arrange
		var (
			called = false
		)

		// act
		sut := anchor.Setup("TEST NAME", func() error {
			called = true
			return nil
		})

		// assert
		assert.NoError(t, sut.Start(t.Context()))
		component, ok := sut.(interface {
			Setup(ctx context.Context) error
			Name() string
		})
		if assert.Truef(t, ok, "expected Setup() method") {
			assert.NoError(t, component.Setup(t.Context()))
			assert.Truef(t, called, "not called")
			assert.Equal(t, "TEST NAME", component.Name())
		}
	})
}
