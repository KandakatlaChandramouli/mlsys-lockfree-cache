
package main

import (
    "fmt"
    "math/rand"
    "time"

    "fluxruntime/internal/avx"
)

const (
    Dim = 384
    N   = 100000
)

func randVec() []float32 {

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

    a := randVec()
    b := randVec()

    start := time.Now()

    var out float32

    for i := 0; i < N; i++ {
        out += avx.DotProduct(
            a,
            b,
        )
    }

    fmt.Println(
        "result:",
        out,
    )

    fmt.Println(
        "latency:",
        time.Since(start),
    )
}
