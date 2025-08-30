package anchor_test

import (
	"context"
	"testing"
	"time"

	"github.com/kyuff/anchor"
	"github.com/kyuff/anchor/internal/assert"
)

type closeFunc func() error

func (fn closeFunc) Close() error {
	return fn()
}

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

	t.Run("Close", func(t *testing.T) {
		// arrange
		var (
			called = false
		)

		// act
		sut := anchor.Close("TEST NAME", closeFunc(func() error {
			called = true
			return nil
		}))

		// assert
		assert.NoError(t, sut.Start(t.Context()))
		component, ok := sut.(interface {
			Close(ctx context.Context) error
			Name() string
		})
		if assert.Truef(t, ok, "expected Close() method") {
			assert.NoError(t, component.Close(t.Context()))
			assert.Truef(t, called, "not called")
			assert.Equal(t, "TEST NAME", component.Name())
		}
	})

	t.Run("Make", func(t *testing.T) {
		// arrange
		var (
			called = false
		)

		// act
		sut := anchor.Make("TEST NAME", func() (*ComponentMock, error) {
			called = true
			return &ComponentMock{
				StartFunc: func(ctx context.Context) error {
					return nil
				},
			}, nil
		})

		// assert
		component, ok := sut.(interface {
			Setup(ctx context.Context) error
			Start(ctx context.Context) error
			Close(ctx context.Context) error
			Name() string
		})
		if assert.Truef(t, ok, "expected a full component") {
			assert.NoError(t, component.Setup(t.Context()))
			assert.NoError(t, sut.Start(t.Context()))
			assert.Truef(t, called, "not called")
			assert.Equal(t, "TEST NAME", component.Name())
		}
	})

	t.Run("MakeProbe", func(t *testing.T) {
		// arrange
		var (
			called = false
			probed = false
		)

		// act
		sut := anchor.MakeProbe("TEST NAME",
			func() (*ComponentMock, error) {
				called = true
				return &ComponentMock{
					StartFunc: func(ctx context.Context) error {
						return nil
					},
				}, nil
			},
			func(ctx context.Context) error {
				probed = true
				return nil
			})

		// assert
		component, ok := sut.(interface {
			Setup(ctx context.Context) error
			Start(ctx context.Context) error
			Probe(ctx context.Context) error
			Close(ctx context.Context) error
			Name() string
		})
		if assert.Truef(t, ok, "expected a full component") {
			assert.NoError(t, component.Setup(t.Context()))
			assert.NoError(t, sut.Start(t.Context()))
			assert.NoErrorEventually(t, time.Millisecond*100, func() error {
				return component.Probe(t.Context())
			})
			assert.Truef(t, called, "not called")
			assert.Truef(t, probed, "not probed")
			assert.Equal(t, "TEST NAME", component.Name())
		}
	})
}
