package lockfree

import (
	"time"

	"fluxruntime/internal/core"
)

type Request struct {
	Req      *core.RawRequest
	Callback func(*core.RawResponse)
}

type Engine struct {
	ring   *Ring
	worker *core.Worker
}

func NewEngine(id int) *Engine {

	e := &Engine{
		ring:   NewRing(),
		worker: core.NewWorker(id),
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

	batch := make(
		[]*Request,
		0,
		core.MaxBatchSize,
	)

	for {

		v, ok := e.ring.Pop()

		if !ok {

			time.Sleep(time.Microsecond)

			continue
		}

		batch = batch[:0]

		batch = append(
			batch,
			v.(*Request),
		)

		for len(batch) < core.MaxBatchSize {

			v, ok := e.ring.Pop()

			if !ok {
				break
			}

			batch = append(
				batch,
				v.(*Request),
			)
		}

		reqs := make(
			[]*core.RawRequest,
			len(batch),
		)

		for i := range batch {
			reqs[i] = batch[i].Req
		}

		resps := e.worker.InferBatch(
			reqs,
		)

		for i := range batch {
			batch[i].Callback(resps[i])
		}
	}
}
