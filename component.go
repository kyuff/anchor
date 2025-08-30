package anchor

import (
	"context"
	"io"

	"github.com/kyuff/anchor/internal/decorate"
)

// Component is a central part of an application that needs to have it's lifetime managed.
//
// A Component can optionally Setup and Close methods..
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

// Close creates a component that have an empty Start() and Setup() method, but
// have Close. It is a convenience to run cleanup code
func Close(name string, closer io.Closer) Component {
	return decorate.Close(name, closer)
}

// MakeProbe a component by it's setup and probe functions. A convenience when the Component is not needed
// as a reference by other parts of the application, but just needs it's lifecycle handled.
// The probe function is called after Start to check if the Component is ready.
func MakeProbe[T Component](name string, setup func() (T, error), probe func(ctx context.Context) error) Component {
	return decorate.MakeProbe(name, setup, probe)
}

// contextSetupComponent allows a Component to create resources before Start
// The context gives the Deadline in which Setup must be complete.
type contextSetupComponent interface {
	Setup(ctx context.Context) error
}

// contextProbeComponent allows a Component to probe its readiness.
// Probe is called after Start and before Close.
// It will continue to be called until it returns nil or the Context is done.
type contextProbeComponent interface {
	Probe(ctx context.Context) error
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
	contextProbeComponent
	contextCloseComponent
	namedComponent
}
