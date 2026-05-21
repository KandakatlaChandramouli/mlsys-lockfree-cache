package main

import (
    "fmt"
    "time"

    "fluxruntime/internal/mmap"
)

func main() {

    start := time.Now()

    mapping, err := mmap.Open(
        "snapshot.bin",
    )

    if err != nil {
        panic(err)
    }

    defer mapping.Close()

    fmt.Println(
        "mapped bytes:",
        len(mapping.Data),
    )

    fmt.Println(
        "mmap latency:",
        time.Since(start),
    )
}
