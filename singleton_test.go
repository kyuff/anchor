package anchor_test

import (
	"errors"
	"math/rand/v2"
	"testing"

	"github.com/kyuff/anchor"
	"github.com/kyuff/anchor/internal/assert"
)

func TestSingleton(t *testing.T) {
	var (
		newValue = rand.Int
	)
	t.Run("create only once", func(t *testing.T) {
		// arrange
		var (
			calls = 0
			sut   = anchor.Singleton(func() (int, error) {
				calls++
				return newValue(), nil
			})
		)

		// act
		for range 10 {
			_ = sut()
		}

		// act
		assert.Equal(t, 1, calls)
	})

	t.Run("return created value", func(t *testing.T) {
		// arrange
		var (
			value = newValue()
			sut   = anchor.Singleton(func() (int, error) {
				return value, nil
			})
		)

		// act
		got := sut()

		// act
		assert.Equal(t, value, got)
	})

	t.Run("panic on error", func(t *testing.T) {
		// arrange
		var (
			didPanic = false
			sut      = anchor.Singleton(func() (int, error) {
				return 0, errors.New("TEST")
			})
		)

		defer func() {
			if err := recover(); err != nil {
				didPanic = true
			}
		}()

		// act
		_ = sut()

		// act
		assert.Truef(t, didPanic, "panic done")
	})

	t.Run("panic on panic", func(t *testing.T) {
		// arrange
		var (
			didPanic = false
			sut      = anchor.Singleton(func() (int, error) {
				panic("TEST")
			})
		)

		defer func() {
			if err := recover(); err != nil {
				didPanic = true
			}
		}()

		// act
		_ = sut()

		// act
		assert.Truef(t, didPanic, "panic done")
	})
}
