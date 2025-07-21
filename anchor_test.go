package anchor_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/kyuff/anchor"
	"github.com/kyuff/anchor/internal/assert"
)

func TestAnchor(t *testing.T) {
	type call string
	const (
		setupCalled  call = "setupCalled"
		setupSkipped call = "setupSkipped"
		startCalled  call = "startCalled"
		startSkipped call = "startSkipped"
		closeCalled  call = "closeCalled"
		closeSkipped call = "closeSkipped"
	)
	var (
		assertCalls = func(t *testing.T, component *fullComponentMock, calls ...call) {
			t.Helper()
			for _, c := range calls {
				switch c {
				case setupCalled:
					assert.Equalf(t, 1, len(component.SetupCalls()), "Setup Called: %s", component.Name())
				case setupSkipped:
					assert.Equalf(t, 0, len(component.SetupCalls()), "Setup Skipped: %s", component.Name())
				case startCalled:
					assert.Equalf(t, 1, len(component.StartCalls()), "Start Called: %s", component.Name())
				case startSkipped:
					assert.Equalf(t, 0, len(component.StartCalls()), "Start Skipped: %s", component.Name())
				case closeCalled:
					assert.Equalf(t, 1, len(component.CloseCalls()), "Close Called: %s", component.Name())
				case closeSkipped:
					assert.Equalf(t, 0, len(component.CloseCalls()), "Close Skipped: %s", component.Name())
				default:
					t.Fatalf("Unexpected call %q", c)
				}
			}
		}
		doneOnStart = func(wg *sync.WaitGroup) func(c *fullComponentMock) {
			return func(c *fullComponentMock) {
				start := c.StartFunc
				c.StartFunc = func(ctx context.Context) error {
					// slight delay to mimic a process running
					wg.Done()
					return start(ctx)
				}
			}
		}
		errorOnStart = func(err error) func(c *fullComponentMock) {
			return func(c *fullComponentMock) {
				c.StartFunc = func(ctx context.Context) error {
					return err
				}
			}
		}
		panicOnStart = func(msg any) func(c *fullComponentMock) {
			return func(c *fullComponentMock) {
				c.StartFunc = func(ctx context.Context) error {
					panic(msg)
				}
			}
		}
		errorOnSetup = func(err error) func(c *fullComponentMock) {
			return func(c *fullComponentMock) {
				c.SetupFunc = func(ctx context.Context) error {
					return err
				}
			}
		}
		panicOnSetup = func(msg any) func(c *fullComponentMock) {
			return func(c *fullComponentMock) {
				c.SetupFunc = func(ctx context.Context) error {
					panic(msg)
				}
			}
		}
		sleepOnSetup = func(duration time.Duration) func(c *fullComponentMock) {
			return func(c *fullComponentMock) {
				c.SetupFunc = func(ctx context.Context) error {
					time.Sleep(duration)
					return nil
				}
			}
		}
		blockOnStart = func(duration time.Duration) func(c *fullComponentMock) {
			return func(c *fullComponentMock) {
				c.StartFunc = func(ctx context.Context) error {
					select {
					case <-ctx.Done():
						return ctx.Err()
					case <-time.After(duration):
						return nil
					}
				}
			}
		}
		errorOnClose = func(err error) func(c *fullComponentMock) {
			return func(c *fullComponentMock) {
				c.CloseFunc = func(ctx context.Context) error {
					return err
				}
			}
		}
		panicOnClose = func(msg any) func(c *fullComponentMock) {
			return func(c *fullComponentMock) {
				c.CloseFunc = func(ctx context.Context) error {
					panic(msg)
				}
			}
		}
		newComponent = func(name string, mods ...func(c *fullComponentMock)) *fullComponentMock {
			c := &fullComponentMock{
				CloseFunc: func(ctx context.Context) error {
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
				ProbeFunc: func(ctx context.Context) error { return nil },
			}

			for _, mod := range mods {
				mod(c)
			}

			return c
		}

		newWire = func(t *testing.T, wg *sync.WaitGroup) *WireMock {
			return &WireMock{
				WireFunc: func(ctx context.Context) (context.Context, context.CancelFunc) {
					ctx, cancel := context.WithCancel(ctx)
					t.Cleanup(cancel)
					go func() {
						wg.Wait()
						cancel()
					}()
					return ctx, cancel
				},
			}
		}
	)
	t.Run("use all components", func(t *testing.T) {
		// arrange
		var (
			wg         = &sync.WaitGroup{}
			components = []*fullComponentMock{
				newComponent("c-0", doneOnStart(wg)),
				newComponent("c-1", doneOnStart(wg)),
				newComponent("c-2", doneOnStart(wg)),
			}
			wire = newWire(t, wg)
			sut  = anchor.New(wire)
		)

		for _, component := range components {
			wg.Add(1)
			sut.Add(component)
		}

		// act
		code := sut.Run()

		// assert
		assert.Equal(t, anchor.OK, code)
		for _, component := range components {
			assertCalls(t, component, setupCalled, startCalled, closeCalled)
		}
	})

	t.Run("panic on nil components", func(t *testing.T) {
		// arrange
		var (
			wg   = &sync.WaitGroup{}
			wire = newWire(t, wg)
			sut  = anchor.New(wire)
		)

		// assert
		assert.Panic(t, func() {
			// act
			_ = sut.Add(nil)
		})
	})

	t.Run("panic on add to running anchor", func(t *testing.T) {
		// arrange
		var (
			wireWg  = &sync.WaitGroup{}
			startWg = &sync.WaitGroup{}
			wire    = newWire(t, wireWg)
			sut     = anchor.New(wire)
		)

		wireWg.Add(1)
		startWg.Add(1)
		t.Cleanup(wireWg.Done)

		sut.Add(newComponent("c-0", doneOnStart(startWg)))

		go func() {
			_ = sut.Run()
		}()

		startWg.Wait()

		// assert
		assert.Panic(t, func() {
			// act
			_ = sut.Add(newComponent("c-0"))
		})
	})

	t.Run("panic on second run", func(t *testing.T) {
		// arrange
		var (
			wireWg  = &sync.WaitGroup{}
			startWg = &sync.WaitGroup{}
			wire    = newWire(t, wireWg)
			sut     = anchor.New(wire)
		)

		wireWg.Add(1)
		startWg.Add(1)
		t.Cleanup(wireWg.Done)

		sut.Add(newComponent("c-0", doneOnStart(startWg)))

		go func() {
			_ = sut.Run()
		}()

		startWg.Wait()

		// assert
		assert.Panic(t, func() {
			// act
			_ = sut.Run()
		})
	})

	t.Run("exit with no components", func(t *testing.T) {
		// arrange
		var (
			wg   = &sync.WaitGroup{}
			wire = newWire(t, wg)
			sut  = anchor.New(wire)
		)

		// act
		code := sut.Run()

		// assert
		assert.Equal(t, anchor.OK, code)
	})

	t.Run("break on setup error", func(t *testing.T) {
		// arrange
		var (
			wg         = &sync.WaitGroup{}
			components = []*fullComponentMock{
				newComponent("c-0"),
				newComponent("c-1", errorOnSetup(errors.New("FAIL"))),
				newComponent("c-2"),
			}
			wire = newWire(t, wg)
			sut  = anchor.New(wire)
		)

		for _, component := range components {
			wg.Add(1)
			sut.Add(component)
		}

		// act
		code := sut.Run()

		// assert
		assert.Equal(t, anchor.SetupFailed, code)
		assertCalls(t, components[0], setupCalled, startSkipped, closeCalled)
		assertCalls(t, components[1], setupCalled, startSkipped, closeCalled)
		assertCalls(t, components[2], setupSkipped, startSkipped, closeSkipped)
	})

	t.Run("break on setup panic", func(t *testing.T) {
		// arrange
		var (
			wg         = &sync.WaitGroup{}
			components = []*fullComponentMock{
				newComponent("c-0"),
				newComponent("c-1", panicOnSetup("TEST")),
				newComponent("c-2"),
			}
			wire = newWire(t, wg)
			sut  = anchor.New(wire)
		)

		for _, component := range components {
			wg.Add(1)
			sut.Add(component)
		}

		// act
		code := sut.Run()

		// assert
		assert.Equal(t, anchor.SetupFailed, code)
		assertCalls(t, components[0], setupCalled, startSkipped, closeCalled)
		assertCalls(t, components[1], setupCalled, startSkipped, closeCalled)
		assertCalls(t, components[2], setupSkipped, startSkipped, closeSkipped)
	})

	t.Run("break on setup timeout with short runtime", func(t *testing.T) {
		// arrange
		var (
			wg         = &sync.WaitGroup{}
			components = []*fullComponentMock{
				newComponent("c-0", doneOnStart(wg)),
				newComponent("c-1", doneOnStart(wg), sleepOnSetup(time.Second)),
				newComponent("c-2", doneOnStart(wg)),
			}
			wire = newWire(t, wg)
			sut  = anchor.New(wire, anchor.WithSetupTimeout(time.Millisecond*50))
		)

		for _, component := range components {
			wg.Add(1)
			sut.Add(component)
		}

		// act
		code := sut.Run()

		// assert
		assert.Equal(t, anchor.Interrupted, code)
		assertCalls(t, components[0], setupCalled, startSkipped, closeCalled)
		assertCalls(t, components[1], setupCalled, startSkipped, closeCalled)
		assertCalls(t, components[2], setupSkipped, startSkipped, closeSkipped)
	})

	t.Run("break on setup timeout with long runtime", func(t *testing.T) {
		// arrange
		var (
			wg         = &sync.WaitGroup{}
			components = []*fullComponentMock{
				newComponent("c-0", doneOnStart(wg)),
				newComponent("c-1", doneOnStart(wg), sleepOnSetup(time.Minute)),
				newComponent("c-2", doneOnStart(wg)),
			}
			wire = newWire(t, wg)
			sut  = anchor.New(wire, anchor.WithSetupTimeout(time.Millisecond*50))
		)

		for _, component := range components {
			wg.Add(1)
			sut.Add(component)
		}

		// act
		code := sut.Run()

		// assert
		assert.Equal(t, anchor.Interrupted, code)
		assertCalls(t, components[0], setupCalled, startSkipped, closeCalled)
		assertCalls(t, components[1], setupCalled, startSkipped, closeCalled)
		assertCalls(t, components[2], setupSkipped, startSkipped, closeSkipped)
	})

	t.Run("break on ready check notification error", func(t *testing.T) {
		// arrange
		var (
			wg         = &sync.WaitGroup{}
			components = []*fullComponentMock{
				newComponent("c-0", blockOnStart(time.Second)),
				newComponent("c-1", blockOnStart(time.Second)),
				newComponent("c-2", blockOnStart(time.Second)),
			}
			wire = newWire(t, wg)
			sut  = anchor.New(wire,
				anchor.WithReady(func(ctx context.Context) error {
					return errors.New("FAIL")
				}),
				anchor.WithDefaultSlog(),
			)
		)

		for _, component := range components {
			wg.Add(1)
			sut.Add(component)
		}

		// act
		code := sut.Run()

		// assert
		assert.Equal(t, anchor.Internal, code)
		assertCalls(t, components[0], setupCalled, startCalled, closeCalled)
		assertCalls(t, components[1], setupCalled, startCalled, closeCalled)
		assertCalls(t, components[2], setupCalled, startCalled, closeCalled)
	})

	t.Run("no break on close error", func(t *testing.T) {
		// arrange
		var (
			wg         = &sync.WaitGroup{}
			components = []*fullComponentMock{
				newComponent("c-0", doneOnStart(wg)),
				newComponent("c-1", doneOnStart(wg), errorOnClose(errors.New("FAIL"))),
				newComponent("c-2", doneOnStart(wg)),
			}
			wire = newWire(t, wg)
			sut  = anchor.New(wire)
		)

		for _, component := range components {
			wg.Add(1)
			sut.Add(component)
		}

		// act
		code := sut.Run()

		// assert
		assert.Equal(t, anchor.OK, code)
		assertCalls(t, components[0], setupCalled, startCalled, closeCalled)
		assertCalls(t, components[1], setupCalled, startCalled, closeCalled)
		assertCalls(t, components[2], setupCalled, startCalled, closeCalled)
	})

	t.Run("no break on close panic", func(t *testing.T) {
		// arrange
		var (
			wg         = &sync.WaitGroup{}
			components = []*fullComponentMock{
				newComponent("c-0", doneOnStart(wg)),
				newComponent("c-1", doneOnStart(wg), panicOnClose("TEST")),
				newComponent("c-2", doneOnStart(wg)),
			}
			wire = newWire(t, wg)
			sut  = anchor.New(wire)
		)

		for _, component := range components {
			wg.Add(1)
			sut.Add(component)
		}

		// act
		code := sut.Run()

		// assert
		assert.Equal(t, anchor.OK, code)
		assertCalls(t, components[0], setupCalled, startCalled, closeCalled)
		assertCalls(t, components[1], setupCalled, startCalled, closeCalled)
		assertCalls(t, components[2], setupCalled, startCalled, closeCalled)
	})

	t.Run("return Internal when start errors", func(t *testing.T) {
		// arrange
		var (
			wg         = &sync.WaitGroup{}
			components = []*fullComponentMock{
				newComponent("c-0", doneOnStart(wg)),
				newComponent("c-1", doneOnStart(wg), errorOnStart(errors.New("FAIL"))),
				newComponent("c-2", doneOnStart(wg)),
			}
			wire = newWire(t, wg)
			sut  = anchor.New(wire)
		)

		for _, component := range components {
			wg.Add(1)
			sut.Add(component)
		}

		// act
		code := sut.Run()

		// assert
		assert.Equal(t, anchor.Internal, code)
		for _, component := range components {
			assertCalls(t, component, setupCalled, startCalled, closeCalled)
		}
	})

	t.Run("return Internal when start panics", func(t *testing.T) {
		// arrange
		var (
			wg         = &sync.WaitGroup{}
			components = []*fullComponentMock{
				newComponent("c-0", doneOnStart(wg)),
				newComponent("c-1", doneOnStart(wg), panicOnStart("TEST")),
				newComponent("c-2", doneOnStart(wg)),
			}
			wire = newWire(t, wg)
			sut  = anchor.New(wire, anchor.WithDefaultSlog())
		)

		for _, component := range components {
			wg.Add(1)
			sut.Add(component)
		}

		// act
		code := sut.Run()

		// assert
		assert.Equal(t, anchor.Internal, code)
		for _, component := range components {
			assertCalls(t, component, setupCalled, startCalled, closeCalled)
		}
	})

	t.Run("setup in order of registration", func(t *testing.T) {
		// arrange
		var (
			wg         = &sync.WaitGroup{}
			names      []string
			recordName = func() func(c *fullComponentMock) {
				return func(c *fullComponentMock) {
					c.SetupFunc = func(ctx context.Context) error {
						names = append(names, c.Name())
						return nil
					}
				}
			}
			components = []*fullComponentMock{
				newComponent("c-0", doneOnStart(wg), recordName()),
				newComponent("c-1", doneOnStart(wg), recordName()),
				newComponent("c-2", doneOnStart(wg), recordName()),
				newComponent("c-3", doneOnStart(wg), recordName()),
			}
			wire = newWire(t, wg)
			sut  = anchor.New(wire)
		)

		for _, component := range components {
			wg.Add(1)
			sut.Add(component)
		}

		// act
		code := sut.Run()

		// assert
		assert.Equal(t, anchor.OK, code)
		assert.EqualSlice(t, []string{"c-0", "c-1", "c-2", "c-3"}, names)
	})

	t.Run("close in reverse order of registration", func(t *testing.T) {
		// arrange
		var (
			wg         = &sync.WaitGroup{}
			names      []string
			recordName = func() func(c *fullComponentMock) {
				return func(c *fullComponentMock) {
					c.CloseFunc = func(ctx context.Context) error {
						names = append(names, c.Name())
						return nil
					}
				}
			}
			components = []*fullComponentMock{
				newComponent("c-0", doneOnStart(wg), recordName()),
				newComponent("c-1", doneOnStart(wg), recordName()),
				newComponent("c-2", doneOnStart(wg), recordName()),
				newComponent("c-3", doneOnStart(wg), recordName()),
			}
			wire = newWire(t, wg)
			sut  = anchor.New(wire)
		)

		for _, component := range components {
			wg.Add(1)
			sut.Add(component)
		}

		// act
		code := sut.Run()

		// assert
		assert.Equal(t, anchor.OK, code)
		assert.EqualSlice(t, []string{"c-3", "c-2", "c-1", "c-0"}, names)
	})

	t.Run("interrupt on external context", func(t *testing.T) {
		// arrange
		var (
			testCtx, cancel = context.WithCancel(t.Context())
			wg              = &sync.WaitGroup{}
			components      = []*fullComponentMock{
				newComponent("c-0", doneOnStart(wg)),
				newComponent("c-1", doneOnStart(wg)),
				newComponent("c-2", doneOnStart(wg)),
			}
			wire = &WireMock{ // noop Wire
				WireFunc: func(ctx context.Context) (context.Context, context.CancelFunc) {
					return context.WithCancel(ctx)
				},
			}
			sut = anchor.New(wire, anchor.WithContext(testCtx))
		)

		for _, component := range components {
			wg.Add(1)
			sut.Add(component)
		}

		go func() {
			wg.Wait()
			cancel()
		}()

		// act
		code := sut.Run()

		// assert
		assert.Equal(t, anchor.OK, code)
		for _, component := range components {
			assertCalls(t, component, setupCalled, startCalled, closeCalled)
		}
	})

	t.Run("call CancelFunc on shutdown", func(t *testing.T) {
		// arrange
		var (
			testCtx, cancel = context.WithCancel(t.Context())
			wg              = &sync.WaitGroup{}
			components      = []*fullComponentMock{
				newComponent("c-0", doneOnStart(wg)),
				newComponent("c-1", doneOnStart(wg)),
				newComponent("c-2", doneOnStart(wg)),
			}
			called = false
			wire   = &WireMock{ // noop Wire
				WireFunc: func(ctx context.Context) (context.Context, context.CancelFunc) {
					wireCtx, cancel := context.WithCancel(ctx)
					return wireCtx, func() {
						called = true
						cancel()
					}
				},
			}
			sut = anchor.New(wire, anchor.WithContext(testCtx))
		)

		for _, component := range components {
			wg.Add(1)
			sut.Add(component)
		}

		go func() {
			wg.Wait()
			cancel()
		}()

		// act
		code := sut.Run()

		// assert
		assert.Equal(t, anchor.OK, code)
		assert.Truef(t, called, "CancelFunc called")
	})

	t.Run("close when start blocks", func(t *testing.T) {
		// arrange
		var (
			wg        = &sync.WaitGroup{}
			component = newComponent("blocking component", func(c *fullComponentMock) {
				c.StartFunc = func(ctx context.Context) error {
					wg.Done()
					// block eternal
					<-t.Context().Done()
					return nil
				}
			})
			wire = newWire(t, wg)
			sut  = anchor.New(wire,
				anchor.WithCloseTimeout(time.Millisecond*200),
			)
		)

		wg.Add(1)
		sut.Add(component)

		// act
		code := sut.Run()

		// assert
		assert.Equal(t, anchor.OK, code)
		assertCalls(t, component, setupCalled, startCalled, closeCalled)
	})

	t.Run("interrupt when close blocks", func(t *testing.T) {
		// arrange
		var (
			wg        = &sync.WaitGroup{}
			component = newComponent("blocking component", func(c *fullComponentMock) {
				c.StartFunc = func(ctx context.Context) error {
					wg.Done()
					// block eternal
					<-t.Context().Done()
					return nil
				}
				c.CloseFunc = func(ctx context.Context) error {
					// block eternal
					<-t.Context().Done()
					return nil
				}
			})
			wire = newWire(t, wg)
			sut  = anchor.New(wire,
				anchor.WithCloseTimeout(time.Millisecond*200),
			)
		)

		wg.Add(1)
		sut.Add(component)

		// act
		code := sut.Run()

		// assert
		assert.Equal(t, anchor.Interrupted, code)
		assertCalls(t, component, setupCalled, startCalled, closeCalled)
	})

	t.Run("call ready after probes", func(t *testing.T) {
		// arrange
		var (
			wg         = &sync.WaitGroup{}
			mu         sync.Mutex
			calls      []string
			wire       = newWire(t, wg)
			components = []*fullComponentMock{
				newComponent("c-0", doneOnStart(wg)),
			}
			record = func(action string) func(_ context.Context) error {
				return func(_ context.Context) error {
					mu.Lock()
					defer mu.Unlock()
					calls = append(calls, action)
					return nil
				}
			}
			sut = anchor.New(wire,
				anchor.WithDefaultSlog(),
				anchor.WithReady(record("ready")),
			)
		)

		for _, component := range components {
			wg.Add(1)
			sut.Add(component)
			component.SetupFunc = record("setup")
			component.StartFunc = record("start")
			component.CloseFunc = record("close")
			component.ProbeFunc = func(ctx context.Context) error {
				return nil
			}
		}

		// act
		code := sut.Run()

		// assert
		assert.Equal(t, anchor.OK, code)
		fmt.Printf("%#v\n", calls)
		assert.EqualSlice(t, []string{"setup", "start", "ready", "close"}, calls)
	})

}
