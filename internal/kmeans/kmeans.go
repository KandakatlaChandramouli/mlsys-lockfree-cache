package kmeans

import (
    "fluxruntime/internal/avx"
)

func L2(
    a []float32,
    b []float32,
) float32 {

    var aa float32
    var bb float32

    for i := range a {
        aa += a[i] * a[i]
        bb += b[i] * b[i]
    }

    dot := avx.DotProduct(
        a,
        b,
    )

    return aa + bb - 2*dot
}

func Nearest(
    vec []float32,
    centroids [][]float32,
) int {

    best := 0

    bestDist := L2(
        vec,
        centroids[0],
    )

    for i := 1; i < len(centroids); i++ {

        d := L2(
            vec,
            centroids[i],
        )

        if d < bestDist {
            bestDist = d
            best = i
        }
    }

    return best
}
