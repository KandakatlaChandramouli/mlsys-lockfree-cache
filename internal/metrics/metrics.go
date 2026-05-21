package metrics

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
)

var (
	Requests atomic.Uint64
	Failures atomic.Uint64
	ActiveConnections atomic.Uint64
	QueuedRequests atomic.Uint64
)

type Snapshot struct {
	Requests         uint64 `json:"requests"`
	Failures         uint64 `json:"failures"`
	ActiveConnections uint64 `json:"active_connections"`
	QueuedRequests   uint64 `json:"queued_requests"`
	Inflight uint64 `json:"inflight"`
}

func Stats() Snapshot {

	return Snapshot{
		Requests: Requests.Load(),
		Failures: Failures.Load(),
		ActiveConnections: ActiveConnections.Load(),
		QueuedRequests: QueuedRequests.Load(),
		Inflight: ActiveConnections.Load(),
	}
}

func Handler(
	w http.ResponseWriter,
	r *http.Request,
) {

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		Stats(),
	)
}
