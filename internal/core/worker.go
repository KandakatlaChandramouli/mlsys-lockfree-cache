package core

import (
	"context"
	"time"

	routerv1 "fluxruntime/proto/v1"
)

const (
	workerQueueDepth = 256
	processingDelay  = 100 * time.Microsecond
)

type Job struct {
	Req    *routerv1.RouteRequest
	Result chan<- Result
	Ctx    context.Context
}

type Result struct {
	Resp *routerv1.RouteResponse
	Err  error
}

type Worker struct {
	id    int
	queue chan Job
	done  chan struct{}
}

func NewWorker(id int) *Worker {
	w := &Worker{
		id:    id,
		queue: make(chan Job, workerQueueDepth),
		done:  make(chan struct{}),
	}

	go w.run()

	return w
}

func (w *Worker) Enqueue(ctx context.Context, job Job) bool {
	select {
	case w.queue <- job:
		return true
	case <-ctx.Done():
		return false
	default:
		return false
	}
}

func (w *Worker) Shutdown() {
	close(w.queue)
	<-w.done
}

func (w *Worker) run() {
	defer close(w.done)

	for job := range w.queue {
		w.process(job)
	}
}

func (w *Worker) process(job Job) {
	resp, err := w.infer(job.Ctx, job.Req)

	job.Result <- Result{
		Resp: resp,
		Err:  err,
	}
}

func (w *Worker) infer(
	ctx context.Context,
	req *routerv1.RouteRequest,
) (*routerv1.RouteResponse, error) {

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(processingDelay):
	}

	vec := make([]float32, len(req.Tokens))

	for i := range vec {
		vec[i] = float32(i) * 0.1
	}

	return &routerv1.RouteResponse{
		Vector: vec,
	}, nil
}
