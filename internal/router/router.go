package router

import (
	"context"
	"runtime"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"fluxruntime/internal/core"
	"fluxruntime/internal/metrics"
	routerv1 "fluxruntime/proto/v1"
)

const requestTimeout = 5 * time.Second

var resultPool = sync.Pool{
	New: func() any {
		return make(chan core.Result, 1)
	},
}

type Router struct {
	routerv1.UnimplementedRouterServiceServer

	workers    []*core.Worker
	numWorkers int
}

func New() *Router {

	numWorkers := runtime.NumCPU()

	workers := make([]*core.Worker, numWorkers)

	for i := range workers {
		workers[i] = core.NewWorker(i)
	}

	return &Router{
		workers:    workers,
		numWorkers: numWorkers,
	}
}

func (r *Router) Route(
	ctx context.Context,
	req *routerv1.RouteRequest,
) (*routerv1.RouteResponse, error) {

	metrics.RequestsTotal.Inc()

	start := time.Now()

	if _, ok := ctx.Deadline(); !ok {

		var cancel context.CancelFunc

		ctx, cancel = context.WithTimeout(
			ctx,
			requestTimeout,
		)

		defer cancel()
	}

	idx := int(req.QueryHash % uint64(r.numWorkers))

	resultCh := resultPool.Get().(chan core.Result)

	job := core.Job{
		Req:    req,
		Result: resultCh,
		Ctx:    ctx,
	}

	if !r.workers[idx].Enqueue(ctx, job) {

		metrics.RequestFailures.Inc()

		resultPool.Put(resultCh)

		return nil, status.Error(
			codes.ResourceExhausted,
			"worker queue full",
		)
	}

	select {

	case result := <-resultCh:

		metrics.RequestLatency.Observe(
			time.Since(start).Seconds(),
		)

		resultPool.Put(resultCh)

		return result.Resp, result.Err

	case <-ctx.Done():

		metrics.RequestFailures.Inc()

		resultPool.Put(resultCh)

		return nil, status.Error(
			codes.DeadlineExceeded,
			ctx.Err().Error(),
		)
	}
}
