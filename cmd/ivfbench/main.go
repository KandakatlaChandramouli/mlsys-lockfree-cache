
package main

import (
    "fmt"
    "math/rand"
    "time"

    "fluxruntime/internal/ivf"
)

const (
    Dim = 384
    N   = 100000
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

    idx := ivf.New(256)

    start := time.Now()

    for i := 0; i < N; i++ {

        idx.Add(
            randomVector(),
        )
    }

    fmt.Println(
        "vectors:",
        N,
    )

    fmt.Println(
        "clusters:",
        idx.NList,
    )

    fmt.Println(
        "ivf build latency:",
        time.Since(start),
    )
}
