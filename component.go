package anchor

import (
	"context"

	"github.com/kyuff/anchor/internal/decorate"
)

// Component is a central part of an application that needs to have it's lifetime managed.
//
// A Component can optionally implement the setupComponent or closeComponent interfaces.
// By doing so, Anchor will guarantee to call the methods in order: Setup, Start, Close.
// This allows applications to prepare and gracefully clean up.
//
// An example of a Component could be a database connection.
// It can read it's configuration in the Setup phase, connect at Start and disconnect at Close.
// If the database connection disconnects due to a network outage, it can return and error from Start,
// which closes down the entire application.
type Component interface {
	Start(ctx context.Context) error
}

// Setup creates a component that have an empty Start() and Close() method, but
// have Setup. It is a convenience to run code before full application start.
func Setup(name string, fn func() error) Component {
	return decorate.Setup(name, fn)
}

// Make a component by it's setup func. A convenience when the Component is not needed
// as a reference by other parts of the application, but just needs it's lifecycle handled.
func Make[T Component](name string, setup func() (T, error)) Component {
	return decorate.Make(name, setup)
}

// setupComponent allows a Component to create resources before Start
type setupComponent interface {
	Setup() error
}

// contextSetupComponent allows a Component to create resources before Start
// The context gives the Deadline in which Setup must be complete.
type contextSetupComponent interface {
	Setup(ctx context.Context) error
}

// closeComponent is a standard io.Closer to free up resources on a graceful shutdown.
type closeComponent interface {
	Close() error
}

// contextCloseComponent is a component that close within the Deadline of the Context.
type contextCloseComponent interface {
	Close(ctx context.Context) error
}

type namedComponent interface {
	Name() string
}

type fullComponent interface {
	contextSetupComponent
	Component
	contextCloseComponent
	namedComponent
}
