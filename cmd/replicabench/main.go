
package main

import (
    "fmt"
    "math/rand"
    "time"

    "fluxruntime/internal/index"
    "fluxruntime/internal/replication"
)

const (
    Dim     = 384
    Writes  = 100000
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
        []replication.Replica,
        0,
        Replicas,
    )

    for i := 0; i < Replicas; i++ {

        replicas = append(
            replicas,
            index.NewHNSW(),
        )
    }

    repl := replication.New(
        replicas...,
    )

    start := time.Now()

    for i := 0; i < Writes; i++ {

        repl.Broadcast(
            fmt.Sprintf(
                "vec-%d",
                i,
            ),
            randomVector(),
        )
    }

    fmt.Println(
        "writes:",
        Writes,
    )

    fmt.Println(
        "replicas:",
        Replicas,
    )

    fmt.Println(
        "replication latency:",
        time.Since(start),
    )
}
