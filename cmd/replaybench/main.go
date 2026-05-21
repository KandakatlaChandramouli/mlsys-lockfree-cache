package main

import (
    "fmt"
    "time"

    "fluxruntime/internal/wal"
)

func main() {

    count := 0

    start := time.Now()

    err := wal.Replay(
        "vectors.wal",
        func(id string) {
            count++
        },
    )

    if err != nil {
        panic(err)
    }

    fmt.Println(
        "replayed:",
        count,
    )

    fmt.Println(
        "replay latency:",
        time.Since(start),
    )
}
