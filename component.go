package anchor

import "context"

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

type setupComponent interface {
	Setup(ctx context.Context) error
}

type closeComponent interface {
	Close() error
}

type namedComponent interface {
	Name() string
}

type fullComponent interface {
	setupComponent
	Component
	closeComponent
	namedComponent
}
