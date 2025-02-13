package anchor_test

import (
	"context"
	"testing"

	"github.com/kyuff/anchor"
	"github.com/kyuff/anchor/internal/assert"
)

func TestAnchor(t *testing.T) {
	t.Run("should start all components", func(t *testing.T) {
		// arrange
		var (
			components = []*ComponentMock{{}, {}, {}}
			sut        = anchor.New()
		)

		for _, component := range components {
			component.StartFunc = func(ctx context.Context) error {
				return nil
			}
			sut.Add(component)
		}

		// act
		code := sut.Run()

		// assert
		assert.Equal(t, 0, code)
		for _, component := range components {
			assert.Equal(t, 1, len(component.StartCalls()))
		}
	})
}
