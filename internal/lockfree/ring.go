package lockfree

type Ring struct {
	ch chan any
}

func NewRing() *Ring {

	return &Ring{
		ch: make(chan any, 4096),
	}
}

func (r *Ring) Push(
	v any,
) bool {

	select {

	case r.ch <- v:
		return true

	default:
		return false
	}
}

func (r *Ring) Pop() (any, bool) {

	select {

	case v := <-r.ch:
		return v, true

	default:
		return nil, false
	}
}
