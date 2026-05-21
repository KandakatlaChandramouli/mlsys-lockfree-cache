package core

import (
	"sync"
	"time"
)

const (
	ProcessingDelay = 100 * time.Microsecond
	MaxBatchSize    = 64
)

type RawRequest struct {
	QueryHash uint64
	Tokens    []uint32
	ConnID    uint64
}

type RawResponse struct {
	Vector []float32
}

type Worker struct {
	id int
}

var vectorPool = sync.Pool{
	New: func() any {

		v := make([]float32, 512)

		return &v
	},
}

func AcquireVector(
	size int,
) *[]float32 {

	v := vectorPool.Get().(*[]float32)

	if cap(*v) < size {

		n := make([]float32, size)

		return &n
	}

	tmp := (*v)[:size]

	return &tmp
}

func ReleaseVector(
	v *[]float32,
) {

	vectorPool.Put(v)
}

func NewWorker(
	id int,
) *Worker {

	return &Worker{
		id: id,
	}
}

func (w *Worker) InferBatch(
	reqs []*RawRequest,
) []*RawResponse {

	time.Sleep(ProcessingDelay)

	out := make(
		[]*RawResponse,
		len(reqs),
	)

	for i, req := range reqs {

		vec := AcquireVector(
			len(req.Tokens),
		)

		for j := range *vec {
			(*vec)[j] = float32(j) * 0.1
		}

		out[i] = &RawResponse{
			Vector: *vec,
		}
	}

	return out
}
