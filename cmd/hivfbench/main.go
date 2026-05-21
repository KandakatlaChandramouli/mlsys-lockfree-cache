package main

import (
    "fmt"
    "math/rand"
    "time"

    "fluxruntime/internal/hivf"
    "fluxruntime/internal/kmeans"
)

const (
    Dim     = 384
    Coarse  = 16
    PerLeaf = 16
    Queries = 100000
)

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

func main() {

    root := hivf.Build(
        Coarse,
        PerLeaf,
    )

    coarseRaw := make(
        [][]float32,
        Coarse,
    )

    for i := range coarseRaw {
        coarseRaw[i] = randomVector()
    }

    root.Coarse = kmeans.Build(
        coarseRaw,
    )

    for i := range root.Leaves {

        raw := make(
            [][]float32,
            PerLeaf,
        )

        for j := range raw {
            raw[j] = randomVector()
        }

        root.Leaves[i].Centroids =
            kmeans.Build(raw)
    }

    start := time.Now()

    for i := 0; i < Queries; i++ {

        q := randomVector()

        _, _ = hivf.Route(
            q,
            root,
        )
    }

    fmt.Println(
        "queries:",
        Queries,
    )

    fmt.Println(
        "coarse:",
        Coarse,
    )

    fmt.Println(
        "per leaf:",
        PerLeaf,
    )

    fmt.Println(
        "hierarchical routing latency:",
        time.Since(start),
    )
}
