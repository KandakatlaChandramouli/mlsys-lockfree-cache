
package tests

import (
    "testing"

    "fluxruntime/internal/avx"
)

func TestSIMDDot(
    t *testing.T,
) {

    a := make(
        []float32,
        1024,
    )

    b := make(
        []float32,
        1024,
    )

    for i := range a {

        a[i] = 1
        b[i] = 2
    }

    out := avx.DotProduct(
        a,
        b,
    )

    if out <= 0 {
        t.Fatalf(
            "invalid simd output",
        )
    }
}
