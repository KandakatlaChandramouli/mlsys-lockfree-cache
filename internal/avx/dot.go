package avx

import "fluxruntime/internal/asm"

func DotProduct(
    a []float32,
    b []float32,
) float32 {

    n := len(a)

    if n == 0 {
        return 0
    }

    return asm.DotProductAVX2(
        a,
        b,
    )
}
