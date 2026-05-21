package core

import (
	"sync"
)

const (
	MaxBatchSize = 64
	MaxTokens    = 512
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

		v := make([]float32, MaxTokens)

		return &v
	},
}

var tokenPool = sync.Pool{
	New: func() any {

		t := make([]uint32, MaxTokens)

		return &t
	},
}

func AcquireTokens() *[]uint32 {

	return tokenPool.Get().(*[]uint32)
}

func ReleaseTokens(
	t *[]uint32,
) {

	tokenPool.Put(t)
}

func AcquireVector() *[]float32 {

	return vectorPool.Get().(*[]float32)
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

	out := make(
		[]*RawResponse,
		len(reqs),
	)

	for i, req := range reqs {

		vec := AcquireVector()

		v := (*vec)[:len(req.Tokens)]

		for j := range v {
			v[j] = float32(j) * 0.1
		}

		out[i] = &RawResponse{
			Vector: v,
		}
	}

	return out
}
