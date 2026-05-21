
package metrics

import (
    "sync/atomic"
)

type Metrics struct {
    Requests uint64
    Failures uint64
    CacheHits uint64
}

func New() *Metrics {

    return &Metrics{}
}

func (m *Metrics) RecordRequest() {

    atomic.AddUint64(
        &m.Requests,
        1,
    )
}

func (m *Metrics) RecordFailure() {

    atomic.AddUint64(
        &m.Failures,
        1,
    )
}

func (m *Metrics) RecordCacheHit() {

    atomic.AddUint64(
        &m.CacheHits,
        1,
    )
}

func (m *Metrics) Snapshot() (
    uint64,
    uint64,
    uint64,
) {

    return atomic.LoadUint64(
            &m.Requests,
    ),
    atomic.LoadUint64(
            &m.Failures,
    ),
    atomic.LoadUint64(
            &m.CacheHits,
    )
}
