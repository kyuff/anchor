package anchor

import "sync/atomic"

func New(opts ...Option) *Anchor {
	return &Anchor{
		cfg: applyOptions(defaultOptions(), opts...),
	}
}

type Anchor struct {
	cfg        *Config
	components []Component
	running    atomic.Bool
}

// Add will manage the Component list by the Anchor.
//
// When Run is called, all Components will be started in the order they were
// given to the Anchor.
func (a *Anchor) Add(components ...Component) *Anchor {
	if a.running.Load() {
		// even though panic is frowned upon, this is one of the few places
		// it makes sense to do, as it's part of application setup nad should fail fast.
		panic("cannot add components after Run is called")
	}

	a.components = append(a.components, components...)
	return a
}

func (a *Anchor) Run() int {
	a.running.Store(true)

	return 0
}
