package main

import (
    "fmt"
    "math/rand"
    "time"

    "fluxruntime/internal/index"
    "fluxruntime/internal/leader"
)

const (
    Dim      = 384
    Writes   = 100000
    Replicas = 3
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
        []leader.Replica,
        0,
        Replicas,
    )

    for i := 0; i < Replicas; i++ {

        replicas = append(
            replicas,
            index.NewHNSW(),
        )
    }

    cluster := leader.New(
        replicas...,
    )

    start := time.Now()

    for i := 0; i < Writes; i++ {

        cluster.Write(
            fmt.Sprintf(
                "vec-%d",
                i,
            ),
            randomVector(),
        )
    }

    fmt.Println(
        "leader:",
        cluster.Leader(),
    )

    fmt.Println(
        "writes:",
        Writes,
    )

    fmt.Println(
        "leader write latency:",
        time.Since(start),
    )

    cluster.Elect(
        1,
    )

    fmt.Println(
        "new leader:",
        cluster.Leader(),
    )
}
