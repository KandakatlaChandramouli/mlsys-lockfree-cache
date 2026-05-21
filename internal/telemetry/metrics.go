
package telemetry

import (
    "sort"
    "sync/atomic"
    "time"
)

const MaxSamples = 1_000_000

type Histogram struct {
    idx     atomic.Uint64
    samples []int64
}

func NewHistogram() *Histogram {

    return &Histogram{
        samples: make(
            []int64,
            MaxSamples,
        ),
    }
}

func (h *Histogram) Record(
    start time.Time,
) {

    i := h.idx.Add(
        1,
    ) - 1

    if i >= MaxSamples {
        return
    }

    h.samples[i] = time.Since(
        start,
    ).Nanoseconds()
}

func (h *Histogram) Percentile(
    p float64,
) int64 {

    n := int(
        h.idx.Load(),
    )

    if n == 0 {
        return 0
    }

    cp := make(
        []int64,
        n,
    )

    copy(
        cp,
        h.samples[:n],
    )

    sort.Slice(
        cp,
        func(i, j int) bool {
            return cp[i] < cp[j]
        },
    )

    idx := int(
        (p / 100.0) *
            float64(n-1),
    )

    return cp[idx]
}
