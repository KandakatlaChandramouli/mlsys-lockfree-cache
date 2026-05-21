package lockfree

import (
	"time"

	"fluxruntime/internal/core"
)

type request struct {
	req *core.RawRequest
	resp chan *core.RawResponse
}

type Engine struct {
	ring *Ring

	worker *core.Worker
}

func NewEngine() *Engine {

	e := &Engine{
		ring: NewRing(),
		worker: core.NewWorker(0),
	}

	go e.loop()

	return e
}

func (e *Engine) Submit(
	req *core.RawRequest,
) *core.RawResponse {

	slot := &request{
		req: req,
		resp: make(chan *core.RawResponse),
	}

	deadline := time.Now().Add(
		5 * time.Millisecond,
	)

	for !e.ring.Push(slot) {

		if time.Now().After(deadline) {
			return &core.RawResponse{}
		}

		time.Sleep(time.Microsecond)
	}

	return <-slot.resp
}

func (e *Engine) loop() {

	for {

		v, ok := e.ring.Pop()

		if !ok {
			time.Sleep(time.Microsecond)
			continue
		}

		req := v.(*request)

		resp := e.worker.InferRaw(
			req.req,
		)

		req.resp <- resp
	}
}
