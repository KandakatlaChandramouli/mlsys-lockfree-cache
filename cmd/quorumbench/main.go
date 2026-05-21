package main

import (
    "fmt"
    "math/rand"
    "time"

    "fluxruntime/internal/index"
    "fluxruntime/internal/quorum"
)

const (
    Dim      = 384
    Writes   = 100000
    Replicas = 3
    Quorum   = 2
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

    replicas := make(
        []quorum.Replica,
        0,
        Replicas,
    )

    for i := 0; i < Replicas; i++ {

        replicas = append(
            replicas,
            index.NewHNSW(),
        )
    }

    cluster := quorum.New(
        Quorum,
        replicas...,
    )

    start := time.Now()

    committed := 0

    for i := 0; i < Writes; i++ {

        ok := cluster.Write(
            fmt.Sprintf(
                "vec-%d",
                i,
            ),
            randomVector(),
        )

        if ok {
            committed++
        }
    }

    fmt.Println(
        "writes:",
        Writes,
    )

    fmt.Println(
        "committed:",
        committed,
    )

    fmt.Println(
        "quorum:",
        Quorum,
    )

    fmt.Println(
        "quorum latency:",
        time.Since(start),
    )
}
