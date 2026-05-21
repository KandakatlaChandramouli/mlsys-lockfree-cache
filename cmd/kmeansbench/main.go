package main

import (
    "fmt"
    "math/rand"
    "time"

    "fluxruntime/internal/kmeans"
)

const (
    Dim       = 384
    Centroids = 256
    Queries   = 100000
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

    raw := make(
        [][]float32,
        Centroids,
    )

    for i := range raw {
        raw[i] = randomVector()
    }

    centroids := kmeans.Build(
        raw,
    )

    start := time.Now()

    for i := 0; i < Queries; i++ {

        q := randomVector()

        _ = kmeans.Nearest(
            q,
            centroids,
        )
    }

    fmt.Println(
        "queries:",
        Queries,
    )

    fmt.Println(
        "centroids:",
        Centroids,
    )

    fmt.Println(
        "routing latency:",
        time.Since(start),
    )
}
