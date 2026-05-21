
package tests

import (
    "fmt"
    "testing"

    "fluxruntime/internal/wal"
)

func TestWALReplay(
    t *testing.T,
) {

    log, err := wal.Open(
        "testwal",
    )

    if err != nil {
        t.Fatal(err)
    }

    defer log.Close()

    for i := 0; i < 1000; i++ {

        err := log.Append(
            fmt.Sprintf(
                "vec-%d",
                i,
            ),
        )

        if err != nil {
            t.Fatal(err)
        }
    }

    err = log.Sync()

    if err != nil {
        t.Fatal(err)
    }

    count := 0

    err = wal.Replay(
        "testwal/000000.log",
        func(id string) {
            count++
        },
    )

    if err != nil {
        t.Fatal(err)
    }

    if count != 1000 {
        t.Fatalf(
            "expected 1000 entries got %d",
            count,
        )
    }
}
