package anchor_test

import (
	"context"
	"fmt"
	"time"

	"github.com/kyuff/anchor"
)

type ExampleService struct {
	name string
}

func (svc *ExampleService) Setup(ctx context.Context) error {
	fmt.Printf("Setup of %s\n", svc.name)
	return nil
}

func (svc *ExampleService) Start(ctx context.Context) error {
	fmt.Printf("Starting component\n")
	<-ctx.Done()
	fmt.Printf("Context was cancelled\n")
	return nil
}

func (svc *ExampleService) Close() error {
	fmt.Printf("Closing %s\n", svc.name)
	return nil
}

func ExampleAnchor() {

	code := anchor.New(newAutoClosingWire(time.Millisecond*100)).
		Add(
			&ExampleService{name: "component-a"},
			&ExampleService{name: "component-b"},
		).
		Run()

	fmt.Printf("Exit code: %d\n", code)

	// Output:
	// Setup of component-a
	// Setup of component-b
	// Starting component
	// Starting component
	// Closing down the Anchor
	// Context was cancelled
	// Context was cancelled
	// Closing component-b
	// Closing component-a
	// Exit code: 0
}

func newAutoClosingWire(duration time.Duration) anchor.WireFunc {
	return func(ctx context.Context) (context.Context, context.CancelFunc) {
		ctx, cancel := context.WithCancel(ctx)
		go func() {
			defer cancel()
			time.Sleep(duration)
			fmt.Printf("Closing down the Anchor\n")
		}()
		return ctx, cancel
	}
}
