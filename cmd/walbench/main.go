
package main

import (
    "fmt"
    "time"

    "fluxruntime/internal/wal"
)

const Writes = 100000

func main() {

    log, err := wal.Open(
        "vectors.wal",
    )

    if err != nil {
        panic(err)
    }

    defer log.Close()

    start := time.Now()

    for i := 0; i < Writes; i++ {

        err := log.Append(
            fmt.Sprintf(
                "vec-%d",
                i,
            ),
        )

        if err != nil {
            panic(err)
        }
    }

    if err := log.Sync(); err != nil {
        panic(err)
    }

    fmt.Println(
        "writes:",
        Writes,
    )

    fmt.Println(
        "wal latency:",
        time.Since(start),
    )
}
