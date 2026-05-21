
package main

import (
    "fmt"
    "math/rand"
    "time"

    "fluxruntime/internal/pq"
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

    quantizer := pq.New(
        Dim,
        48,
        256,
    )

    vectors := make(
        [][]float32,
        N,
    )

    for i := range vectors {
        vectors[i] = randomVector()
    }

    start := time.Now()

    encoded := make(
        []pq.QuantizedVector,
        0,
        N,
    )

    for _, v := range vectors {

        encoded = append(
            encoded,
            quantizer.Encode(v),
        )
    }

    fmt.Println(
        "encoded:",
        len(encoded),
    )

    fmt.Println(
        "pq encode latency:",
        time.Since(start),
    )

    fmt.Println(
        "original vector bytes:",
        Dim * 4,
    )

    fmt.Println(
        "pq vector bytes:",
        48,
    )
}
