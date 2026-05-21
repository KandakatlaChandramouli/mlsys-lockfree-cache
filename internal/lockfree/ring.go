package lockfree

import "sync/atomic"

const RingSize = 2048

type Slot struct {
	req any
}

type Ring struct {
	head atomic.Uint64
	tail atomic.Uint64

	buffer [RingSize]Slot
}

func NewRing() *Ring {

	return &Ring{}
}

func (r *Ring) Push(
	v any,
) bool {

	head := r.head.Load()
	tail := r.tail.Load()

	if head-tail >= RingSize {
		return false
	}

	idx := head % RingSize

	r.buffer[idx].req = v

	r.head.Add(1)

	return true
}

func (r *Ring) Pop() (any, bool) {

	tail := r.tail.Load()

	if tail >= r.head.Load() {
		return nil, false
	}

	idx := tail % RingSize

	v := r.buffer[idx].req

	r.tail.Add(1)

	return v, true
}
