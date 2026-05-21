
package tests

import (
    "testing"

    "fluxruntime/internal/snapshot"
    "fluxruntime/internal/index"
)

func TestSnapshot(
    t *testing.T,
) {

    idx := index.NewHNSW()

    for i := 0; i < 1000; i++ {

        idx.Insert(
            "x",
            []float32{
                1,
                2,
                3,
            },
        )
    }

    err := snapshot.Save(
        "test.snapshot",
        idx,
    )

    if err != nil {
        t.Fatal(err)
    }

    restored, err := snapshot.Load(
        "test.snapshot",
    )

    if err != nil {
        t.Fatal(err)
    }

    if restored == nil {
        t.Fatalf(
            "snapshot restore failed",
        )
    }
}
