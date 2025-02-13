package anchor_test

import (
	"context"
	"sync"
	"testing"

	"github.com/kyuff/anchor"
	"github.com/kyuff/anchor/internal/assert"
)

func TestAnchor(t *testing.T) {
	var (
		newComponent = func(name string, mods ...func(c *fullComponentMock)) *fullComponentMock {
			c := &fullComponentMock{
				CloseFunc: func() error {
					return nil
				},
				NameFunc: func() string {
					return name
				},
				SetupFunc: func(ctx context.Context) error {
					return nil
				},
				StartFunc: func(ctx context.Context) error {
					return nil
				},
			}

			for _, mod := range mods {
				mod(c)
			}

			return c
		}
	)
	t.Run("should setup, start and close all components", func(t *testing.T) {
		// arrange
		var (
			components = []*fullComponentMock{
				newComponent("c-0"),
				newComponent("c-1"),
				newComponent("c-2"),
			}
			sut = anchor.New()
		)

		for _, component := range components {
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

	t.Run("start components in go routine", func(t *testing.T) {
		// arrange
		var (
			withWaitGroupOnStart = func(wg *sync.WaitGroup) func(c *fullComponentMock) {
				return func(c *fullComponentMock) {
					c.StartFunc = func(ctx context.Context) error {
						wg.Done()
						return nil
					}
				}
			}
			ctx, cancel = context.WithCancel(t.Context())
			started     = &sync.WaitGroup{}
			components  = []*fullComponentMock{
				newComponent("c-0", withWaitGroupOnStart(started)),
				newComponent("c-1", withWaitGroupOnStart(started)),
				newComponent("c-2", withWaitGroupOnStart(started)),
			}

			sut = anchor.New(anchor.WithContext(ctx))
		)

		for _, component := range components {
			sut.Add(component)
			started.Add(1)
		}

		// act
		go func() {
			_ = sut.Run()
		}()

		// assert
		started.Wait()
		cancel()
		for _, component := range components {
			assert.Equal(t, 1, len(component.StartCalls()))
		}
	})
}
