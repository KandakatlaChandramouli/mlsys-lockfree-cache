
package mempool

import (
	"sync"
)

type Float32Pool struct {
	pool sync.Pool
	size int
}

func NewFloat32Pool(
	size int,
) *Float32Pool {

	return &Float32Pool{
		size: size,
		pool: sync.Pool{
			New: func() any {
				return make(
					[]float32,
					size,
				)
			},
		},
	}
}

func (p *Float32Pool) Get() []float32 {

	return p.pool.Get().([]float32)
}

func (p *Float32Pool) Put(
	buf []float32,
) {

	if len(buf) != p.size {
		return
	}

	clear(buf)

	p.pool.Put(buf)
}
