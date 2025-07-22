package decorate

import (
	"sync/atomic"
)

func New(component starter) *Component {
	var (
		startTime = &atomic.Int64{}
		c         = &Component{inner: component}
	)

	c.setup = makeSetupCall(c)
	c.name = makeNameCall(c)
	c.close = makeClose(c)
	c.start = makeStart(c, startTime)
	c.probe = makeProbe(startTime, probeInner(c))

	return c

}
