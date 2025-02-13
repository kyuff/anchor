package anchor_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/kyuff/anchor"
	"github.com/kyuff/anchor/internal/assert"
)

func TestAnchor(t *testing.T) {
	t.Run("should setup, start and close all components", func(t *testing.T) {
		// arrange
		var (
			components = []*fullComponentMock{{}, {}, {}}
			sut        = anchor.New()
		)

		for i, component := range components {
			component.SetupFunc = func(ctx context.Context) error {
				return nil
			}
			component.StartFunc = func(ctx context.Context) error {
				return nil
			}
			component.CloseFunc = func() error { return nil }
			component.NameFunc = func() string { return fmt.Sprintf("mock-%d", i) }

			sut.Add(component)
		}

		// act
		code := sut.Run()

		// assert
		assert.Equal(t, 0, code)
		for _, component := range components {
			assert.Equal(t, 1, len(component.SetupCalls()))
			assert.Equal(t, 1, len(component.StartCalls()))
			assert.Equal(t, 1, len(component.CloseCalls()))
		}
	})
}
