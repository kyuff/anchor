package decorate_test

import (
	"testing"

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
		assert.NoError(t, sut.Close(t.Context()))
		assert.Equal(t, "TEST NAME", sut.Name())
		assert.Truef(t, called, "not called")
	})
}
