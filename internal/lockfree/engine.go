package lockfree

import (
	"time"

	"fluxruntime/internal/core"
)

type Request struct {
	Req  *core.RawRequest
	Resp chan *core.RawResponse
}

type Engine struct {
	ring   *Ring
	worker *core.Worker

	batch []*Request
	reqs  []*core.RawRequest
}

func NewEngine(
	id int,
) *Engine {

	e := &Engine{
		ring:   NewRing(),
		worker: core.NewWorker(id),

		batch: make(
			[]*Request,
			core.MaxBatchSize,
		),

		reqs: make(
			[]*core.RawRequest,
			core.MaxBatchSize,
		),
	}

	go e.loop()

	return e
}

func (e *Engine) Submit(
	r *Request,
) bool {

	deadline := time.Now().Add(
		5 * time.Millisecond,
	)

	for !e.ring.Push(r) {

		if time.Now().After(deadline) {
			return false
		}

		time.Sleep(time.Microsecond)
	}

	return true
}

func (e *Engine) loop() {

	for {

		v, ok := e.ring.Pop()

		if !ok {

			time.Sleep(time.Microsecond)

			continue
		}

		n := 0

		e.batch[n] = v.(*Request)

		n++

		for n < core.MaxBatchSize {

			v, ok := e.ring.Pop()

			if !ok {
				break
			}

			e.batch[n] = v.(*Request)

			n++
		}

		for i := 0; i < n; i++ {
			e.reqs[i] = e.batch[i].Req
		}

		resps := e.worker.InferBatch(
			e.reqs[:n],
		)

		for i := 0; i < n; i++ {

			e.batch[i].Resp <- resps[i]
		}
	}
}
