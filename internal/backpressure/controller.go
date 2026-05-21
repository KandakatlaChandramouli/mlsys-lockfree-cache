package backpressure

import (
	"sync/atomic"
)

type Controller struct {
	limit uint64

	inflight atomic.Uint64
	rejected atomic.Uint64
}

func New(
	limit uint64,
) *Controller {

	return &Controller{
		limit: limit,
	}
}

func (c *Controller) Acquire() bool {

	for {

		cur := c.inflight.Load()

		if cur >= c.limit {

			c.rejected.Add(1)

			return false
		}

		if c.inflight.CompareAndSwap(
			cur,
			cur+1,
		) {

			return true
		}
	}
}

func (c *Controller) Release() {

	c.inflight.Add(^uint64(0))
}

func (c *Controller) Inflight() uint64 {

	return c.inflight.Load()
}

func (c *Controller) Rejected() uint64 {

	return c.rejected.Load()
}
