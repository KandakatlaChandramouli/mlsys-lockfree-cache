
package main

import (
    "fmt"
    "math/rand"
    "time"

    "fluxruntime/internal/index"
)

const (
    Dim = 384
    N   = 10000
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

    idx := index.NewHNSW()

    start := time.Now()

    for i := 0; i < N; i++ {

        idx.Add(
            fmt.Sprintf("vec-%d", i),
            randomVector(),
        )
    }

    fmt.Println(
        "index build:",
        time.Since(start),
    )

    q := randomVector()

    start = time.Now()

    results := idx.Search(
        q,
        10,
    )

    fmt.Println(
        "search latency:",
        time.Since(start),
    )

    fmt.Println(
        "results:",
        len(results),
    )
}
