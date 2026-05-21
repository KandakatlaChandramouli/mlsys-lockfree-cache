package core

import "time"

const processingDelay = 200 * time.Microsecond

type RawRequest struct {
	QueryHash uint64
	Tokens []uint32
}

type RawResponse struct {
	Vector []float32
}

type Worker struct {
	id int
}

func NewWorker(
	id int,
) *Worker {

	return &Worker{
		id: id,
	}
}

func (w *Worker) InferRaw(
	req *RawRequest,
) *RawResponse {

	time.Sleep(processingDelay)

	vec := make(
		[]float32,
		len(req.Tokens),
	)

	for i := range vec {
		vec[i] = float32(i) * 0.1
	}

	return &RawResponse{
		Vector: vec,
	}
}
