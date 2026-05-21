
package main

import (
    "fmt"
    "time"

    "fluxruntime/internal/wal"
)

const Writes = 100000

func main() {

    log, err := wal.Open(
        "wal",
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

    fmt.Println(
        "writes:",
        Writes,
    )

    fmt.Println(
        "segments:",
        Writes / 10000,
    )

    fmt.Println(
        "segmented wal latency:",
        time.Since(start),
    )
}
