
package main

import (
    "fmt"
    "math/rand"
    "time"

    "fluxruntime/internal/index"
    "fluxruntime/internal/ingest"
)

const (
    Dim     = 384
    Vectors = 100000
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

    pipe := ingest.New(
        8,
        4096,
        idx,
    )

    start := time.Now()

    for i := 0; i < Vectors; i++ {

        pipe.Submit(
            ingest.Job{
                ID: fmt.Sprintf(
                    "vec-%d",
                    i,
                ),
                Embedding: randomVector(),
            },
        )
    }

    pipe.Close()

    fmt.Println(
        "vectors:",
        Vectors,
    )

    fmt.Println(
        "workers:",
        8,
    )

    fmt.Println(
        "async ingestion latency:",
        time.Since(start),
    )
}
