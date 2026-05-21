
package main

import (
    "fmt"
    "math/rand"
    "time"

    "fluxruntime/internal/index"
    "fluxruntime/internal/snapshot"
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

    for i := 0; i < Vectors; i++ {

        idx.Insert(
            fmt.Sprintf(
                "vec-%d",
                i,
            ),
            randomVector(),
        )
    }

    start := time.Now()

    err := snapshot.Save(
        "snapshot.bin",
        idx,
    )

    if err != nil {
        panic(err)
    }

    fmt.Println(
        "snapshot save latency:",
        time.Since(start),
    )

    start = time.Now()

    restored, err := snapshot.Load(
        "snapshot.bin",
    )

    if err != nil {
        panic(err)
    }

    fmt.Println(
        "restored nodes:",
        len(restored.Search(
            randomVector(),
            10,
        )),
    )

    fmt.Println(
        "snapshot restore latency:",
        time.Since(start),
    )
}
