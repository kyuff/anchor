package decorate_test

import (
	"context"
	"errors"
	"testing"

	"github.com/kyuff/anchor/internal/assert"
	"github.com/kyuff/anchor/internal/decorate"
)

func TestComponent(t *testing.T) {
	t.Run("call start error", func(t *testing.T) {
		var (
			component = &starterMock{}
			sut       = decorate.New(component)
		)

		component.StartFunc = func(ctx context.Context) error {
			return errors.New("TEST")
		}

		// act
		sut = decorate.New(component)

		// assert
		assert.NoError(t, sut.Setup(t.Context()))
		assert.Error(t, sut.Start(t.Context()))
		assert.NoError(t, sut.Close())
		assert.Equal(t, "*decorate_test.starterMock", sut.Name())
	})

	t.Run("call start no error", func(t *testing.T) {
		var (
			component = &starterMock{}
			sut       = decorate.New(component)
		)

		component.StartFunc = func(ctx context.Context) error {
			return nil
		}

		// act
		sut = decorate.New(component)

		// assert
		assert.NoError(t, sut.Setup(t.Context()))
		assert.NoError(t, sut.Start(t.Context()))
		assert.NoError(t, sut.Close())
		assert.Equal(t, "*decorate_test.starterMock", sut.Name())
	})

	t.Run("call setup error", func(t *testing.T) {
		var (
			component = &setupperMock{}
			sut       = decorate.New(component)
		)

		component.StartFunc = func(ctx context.Context) error {
			return nil
		}
		component.SetupFunc = func(ctx context.Context) error {
			return errors.New("error")
		}

		// act
		sut = decorate.New(component)

		// assert
		assert.Error(t, sut.Setup(t.Context()))
		assert.NoError(t, sut.Start(t.Context()))
		assert.NoError(t, sut.Close())
		assert.Equal(t, "*decorate_test.setupperMock", sut.Name())
	})

	t.Run("call setup no error", func(t *testing.T) {
		var (
			component = &setupperMock{}
			sut       = decorate.New(component)
		)

		component.StartFunc = func(ctx context.Context) error {
			return nil
		}
		component.SetupFunc = func(ctx context.Context) error {
			return nil
		}

		// act
		sut = decorate.New(component)

		// assert
		assert.NoError(t, sut.Setup(t.Context()))
		assert.NoError(t, sut.Start(t.Context()))
		assert.NoError(t, sut.Close())
		assert.Equal(t, "*decorate_test.setupperMock", sut.Name())
	})

	t.Run("call close error", func(t *testing.T) {
		var (
			component = &closerMock{}
			sut       = decorate.New(component)
		)

		component.StartFunc = func(ctx context.Context) error {
			return nil
		}
		component.CloseFunc = func() error {
			return errors.New("TEST")
		}

		// act
		sut = decorate.New(component)

		// assert
		assert.NoError(t, sut.Setup(t.Context()))
		assert.NoError(t, sut.Start(t.Context()))
		assert.Error(t, sut.Close())
		assert.Equal(t, "*decorate_test.closerMock", sut.Name())
	})

	t.Run("call close no error", func(t *testing.T) {
		var (
			component = &closerMock{}
			sut       = decorate.New(component)
		)

		component.StartFunc = func(ctx context.Context) error {
			return nil
		}
		component.CloseFunc = func() error {
			return nil
		}

		// act
		sut = decorate.New(component)

		// assert
		assert.NoError(t, sut.Setup(t.Context()))
		assert.NoError(t, sut.Start(t.Context()))
		assert.NoError(t, sut.Close())
		assert.Equal(t, "*decorate_test.closerMock", sut.Name())
	})

	t.Run("call name", func(t *testing.T) {
		var (
			component = &namerMock{}
			sut       = decorate.New(component)
		)

		component.StartFunc = func(ctx context.Context) error {
			return nil
		}
		component.NameFunc = func() string {
			return "TEST NAME"
		}

		// act
		sut = decorate.New(component)

		// assert
		assert.NoError(t, sut.Setup(t.Context()))
		assert.NoError(t, sut.Start(t.Context()))
		assert.NoError(t, sut.Close())
		assert.Equal(t, "TEST NAME", sut.Name())
	})

}
