package main

import (
    "fmt"
    "time"

    "fluxruntime/internal/compaction"
)

func main() {

    start := time.Now()

    err := compaction.Merge(
        "wal",
        "wal.compacted",
    )

    if err != nil {
        panic(err)
    }

    fmt.Println(
        "compaction latency:",
        time.Since(start),
    )
}
