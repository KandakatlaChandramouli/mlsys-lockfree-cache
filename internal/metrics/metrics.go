
package metrics

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
	"time"
)

type Metrics struct {
	Requests    atomic.Uint64
	Failures    atomic.Uint64
	BatchCount  atomic.Uint64
	TotalLatency atomic.Uint64
}

type Snapshot struct {
	Requests     uint64  `json:"requests"`
	Failures     uint64  `json:"failures"`
	Batches      uint64  `json:"batches"`
	AvgLatencyMS float64 `json:"avg_latency_ms"`
}

func New() *Metrics {
	return &Metrics{}
}

func (m *Metrics) RecordRequest(
	latency time.Duration,
) {

	m.Requests.Add(1)

	m.TotalLatency.Add(
		uint64(latency.Milliseconds()),
	)
}

func (m *Metrics) RecordFailure() {
	m.Failures.Add(1)
}

func (m *Metrics) RecordBatch() {
	m.BatchCount.Add(1)
}

func (m *Metrics) Snapshot() Snapshot {

	reqs := m.Requests.Load()

	var avg float64

	if reqs > 0 {

		avg = float64(
			m.TotalLatency.Load(),
		) / float64(reqs)
	}

	return Snapshot{
		Requests:     reqs,
		Failures:     m.Failures.Load(),
		Batches:      m.BatchCount.Load(),
		AvgLatencyMS: avg,
	}
}

func Handler(
	m *Metrics,
) http.HandlerFunc {

	return func(
		w http.ResponseWriter,
		r *http.Request,
	) {

		w.Header().Set(
			"Content-Type",
			"application/json",
		)

		json.NewEncoder(w).Encode(
			m.Snapshot(),
		)
	}
}
