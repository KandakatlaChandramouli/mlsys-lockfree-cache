
package tests

import (
    "fmt"
    "testing"

    "fluxruntime/internal/index"
)

func randomVector(
    dim int,
) []float32 {

    out := make(
        []float32,
        dim,
    )

    for i := range out {
        out[i] = float32(i) * 0.01
    }

    return out
}

func TestHNSWInsertSearch(
    t *testing.T,
) {

    idx := index.NewHNSW()

    for i := 0; i < 1000; i++ {

        idx.Insert(
            fmt.Sprintf(
                "vec-%d",
                i,
            ),
            randomVector(384),
        )
    }

    out := idx.Search(
        randomVector(384),
        10,
    )

    if len(out) == 0 {
        t.Fatalf(
            "search returned zero results",
        )
    }
}
