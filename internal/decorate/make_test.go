package decorate_test

import (
	"context"
	"errors"
	"testing"

	"github.com/kyuff/anchor/internal/assert"
	"github.com/kyuff/anchor/internal/decorate"
)

func TestMake(t *testing.T) {

	t.Run("fail on nil setup", func(t *testing.T) {
		// act
		sut := decorate.Make[*starterMock]("TEST NAME", nil)

		// assert
		assert.Error(t, sut.Setup(t.Context()))
		assert.Error(t, sut.Start(t.Context()))
		assert.Error(t, sut.Close(t.Context()))
	})

	t.Run("fail on nil component", func(t *testing.T) {
		// arrange
		var (
			setup = false
		)

		// act
		sut := decorate.Make("TEST NAME", func() (*starterMock, error) {
			setup = true
			return nil, nil
		})

		// assert
		assert.Error(t, sut.Setup(t.Context()))
		assert.Error(t, sut.Start(t.Context()))
		assert.Error(t, sut.Close(t.Context()))
		assert.Equal(t, "TEST NAME", sut.Name())
		assert.Truef(t, setup, "setup")
	})

	t.Run("call setup success", func(t *testing.T) {
		// arrange
		var (
			setup   = false
			started = false
		)

		// act
		sut := decorate.Make("TEST NAME", func() (*starterMock, error) {
			setup = true
			return &starterMock{
				StartFunc: func(ctx context.Context) error {
					started = true
					return nil
				},
			}, nil
		})

		// assert
		assert.NoError(t, sut.Setup(t.Context()))
		assert.NoError(t, sut.Start(t.Context()))
		assert.NoError(t, sut.Close(t.Context()))
		assert.Equal(t, "TEST NAME", sut.Name())
		assert.Truef(t, setup, "setup")
		assert.Truef(t, started, "started")
	})

	t.Run("call setup error", func(t *testing.T) {
		// arrange
		var (
			setup = false
		)

		// act
		sut := decorate.Make("TEST NAME", func() (*starterMock, error) {
			setup = true
			return nil, errors.New("some error")
		})

		// assert
		assert.Error(t, sut.Setup(t.Context()))
		assert.Error(t, sut.Start(t.Context()))
		assert.Error(t, sut.Close(t.Context()))
		assert.Equal(t, "TEST NAME", sut.Name())
		assert.Truef(t, setup, "setup")
	})

	t.Run("call setup error and inner", func(t *testing.T) {
		// arrange
		var (
			setup = false
		)

		// act
		sut := decorate.Make("TEST NAME", func() (*starterMock, error) {
			setup = true
			return &starterMock{
				StartFunc: func(ctx context.Context) error {
					return nil
				},
			}, errors.New("some error")
		})

		// assert
		assert.Error(t, sut.Setup(t.Context()))
		// anchor.Anchor should not call Start, but for good measure we test it.
		assert.NoError(t, sut.Start(t.Context()))
		assert.NoError(t, sut.Close(t.Context()))
		assert.Equal(t, "TEST NAME", sut.Name())
		assert.Truef(t, setup, "setup")
	})

	t.Run("call setup on setupper", func(t *testing.T) {
		// arrange
		var (
			setup      = false
			setupInner = false
		)

		// act
		sut := decorate.Make("TEST NAME", func() (*setupperMock, error) {
			setup = true
			return &setupperMock{
				StartFunc: func(ctx context.Context) error {
					return nil
				},
				SetupFunc: func() error {
					setupInner = true
					return nil
				},
			}, nil
		})

		// assert
		assert.NoError(t, sut.Setup(t.Context()))
		assert.NoError(t, sut.Start(t.Context()))
		assert.NoError(t, sut.Close(t.Context()))
		assert.Equal(t, "TEST NAME", sut.Name())
		assert.Truef(t, setup, "setup")
		assert.Truef(t, setupInner, "setupInner")
	})

	t.Run("call setup on contextSetupper", func(t *testing.T) {
		// arrange
		var (
			setup      = false
			setupInner = false
		)

		// act
		sut := decorate.Make("TEST NAME", func() (*contextSetupperMock, error) {
			setup = true
			return &contextSetupperMock{
				StartFunc: func(ctx context.Context) error {
					return nil
				},
				SetupFunc: func(ctx context.Context) error {
					setupInner = true
					return nil
				},
			}, nil
		})

		// assert
		assert.NoError(t, sut.Setup(t.Context()))
		assert.NoError(t, sut.Start(t.Context()))
		assert.NoError(t, sut.Close(t.Context()))
		assert.Equal(t, "TEST NAME", sut.Name())
		assert.Truef(t, setup, "setup")
		assert.Truef(t, setupInner, "setupInner")
	})

	t.Run("call close on closer", func(t *testing.T) {
		// arrange
		var (
			setup      = false
			closeInner = false
		)

		// act
		sut := decorate.Make("TEST NAME", func() (*closerMock, error) {
			setup = true
			return &closerMock{
				StartFunc: func(ctx context.Context) error {
					return nil
				},
				CloseFunc: func() error {
					closeInner = true
					return nil
				},
			}, nil
		})

		// assert
		assert.NoError(t, sut.Setup(t.Context()))
		assert.NoError(t, sut.Start(t.Context()))
		assert.NoError(t, sut.Close(t.Context()))
		assert.Equal(t, "TEST NAME", sut.Name())
		assert.Truef(t, setup, "setup")
		assert.Truef(t, closeInner, "closeInner")
	})

	t.Run("call close on contextCloser", func(t *testing.T) {
		// arrange
		var (
			setup      = false
			closeInner = false
		)

		// act
		sut := decorate.Make("TEST NAME", func() (*contextCloserMock, error) {
			setup = true
			return &contextCloserMock{
				StartFunc: func(ctx context.Context) error {
					return nil
				},
				CloseFunc: func(ctx context.Context) error {
					closeInner = true
					return nil
				},
			}, nil
		})

		// assert
		assert.NoError(t, sut.Setup(t.Context()))
		assert.NoError(t, sut.Start(t.Context()))
		assert.NoError(t, sut.Close(t.Context()))
		assert.Equal(t, "TEST NAME", sut.Name())
		assert.Truef(t, setup, "setup")
		assert.Truef(t, closeInner, "closeInner")
	})

}
