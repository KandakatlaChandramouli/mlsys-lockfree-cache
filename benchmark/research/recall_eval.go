package main

import (
    "fmt"
    "math/rand"
    "sort"
    "time"

    "fluxruntime/internal/index"
    "fluxruntime/internal/telemetry"
    "fluxruntime/internal/vectorstore"
)

const (
    Dim     = 384
    Vectors = 10000
    Queries = 100
    TopK    = 10
)

type Pair struct {
    idx int
    sim float32
}

func randomVector() []float32 {

    out := make(
        []float32,
        Dim,
    )

    for i := range out {
        out[i] = rand.Float32()
    }

    return out
}

func dot(
    a []float32,
    b []float32,
) float32 {

    var s float32

    for i := range a {
        s += a[i] * b[i]
    }

    return s
}

func bruteForce(
    db [][]float32,
    q []float32,
) []int {

    scores := make(
        []Pair,
        0,
        len(db),
    )

    for i, v := range db {

        s := dot(
            q,
            v,
        )

        scores = append(
            scores,
            Pair{
                idx: i,
                sim: s,
            },
        )
    }

    sort.Slice(
        scores,
        func(i, j int) bool {
            return scores[i].sim >
                scores[j].sim
        },
    )

    out := make(
        []int,
        TopK,
    )

    for i := 0; i < TopK; i++ {
        out[i] = scores[i].idx
    }

    return out
}

func recall(
    gt []int,
    pred []vectorstore.SearchResult,
) float64 {

    gtset := make(
        map[string]struct{},
    )

    for _, g := range gt {

        key := fmt.Sprintf(
            "vec-%d",
            g,
        )

        gtset[key] = struct{}{}
    }

    hit := 0

    for _, p := range pred {

        if _, ok := gtset[p.ID]; ok {
            hit++
        }
    }

    return float64(hit) /
        float64(len(gt))
}

func main() {

    db := make(
        [][]float32,
        0,
        Vectors,
    )

    idx := index.NewHNSW()

    for i := 0; i < Vectors; i++ {

        v := randomVector()

        db = append(
            db,
            v,
        )

        idx.Insert(
            fmt.Sprintf(
                "vec-%d",
                i,
            ),
            v,
        )
    }

    hist := telemetry.NewHistogram()

    totalRecall := 0.0

    for i := 0; i < Queries; i++ {

        q := randomVector()

        gt := bruteForce(
            db,
            q,
        )

        start := time.Now()

        pred := idx.Search(
            q,
            TopK,
        )

        hist.Record(
            start,
        )

        totalRecall += recall(
            gt,
            pred,
        )
    }

    fmt.Println(
        "queries:",
        Queries,
    )

    fmt.Println(
        "recall@10:",
        totalRecall/
            float64(Queries),
    )

    fmt.Println(
        "p50(ns):",
        hist.Percentile(50),
    )

    fmt.Println(
        "p95(ns):",
        hist.Percentile(95),
    )

    fmt.Println(
        "p99(ns):",
        hist.Percentile(99),
    )
}
