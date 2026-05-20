package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "runtime_requests_total",
			Help: "Total requests processed",
		},
	)

	RequestFailures = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "runtime_request_failures_total",
			Help: "Total failed requests",
		},
	)

	RequestLatency = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "runtime_request_latency_seconds",
			Help:    "Request latency",
			Buckets: prometheus.DefBuckets,
		},
	)

	WorkerQueueDepth = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "runtime_worker_queue_depth",
			Help: "Current queue depth per worker",
		},
		[]string{"worker"},
	)
)

func Register() {

	prometheus.MustRegister(RequestsTotal)
	prometheus.MustRegister(RequestFailures)
	prometheus.MustRegister(RequestLatency)
	prometheus.MustRegister(WorkerQueueDepth)
}
