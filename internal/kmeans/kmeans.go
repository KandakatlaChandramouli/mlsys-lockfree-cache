package kmeans

import (
    "fluxruntime/internal/avx"
)

type Centroid struct {
    Vec  []float32
    Norm float32
}

func norm(
    v []float32,
) float32 {

    var out float32

    for i := range v {
        out += v[i] * v[i]
    }

    return out
}

func L2(
    a []float32,
    b []float32,
) float32 {

    aa := norm(a)
    bb := norm(b)

    dot := avx.DotProduct(
        a,
        b,
    )

    return aa + bb - 2*dot
}

func Build(
    centroids [][]float32,
) []Centroid {

    out := make(
        []Centroid,
        len(centroids),
    )

    for i := range centroids {

        out[i] = Centroid{
            Vec: centroids[i],
            Norm: norm(
                centroids[i],
            ),
        }
    }

    return out
}

func Nearest(
    vec []float32,
    centroids []Centroid,
) int {

    queryNorm := norm(
        vec,
    )

    best := 0

    bestDist :=
        queryNorm +
        centroids[0].Norm -
        2*avx.DotProduct(
            vec,
            centroids[0].Vec,
        )

    for i := 1; i < len(centroids); i++ {

        d :=
            queryNorm +
            centroids[i].Norm -
            2*avx.DotProduct(
                vec,
                centroids[i].Vec,
            )

        if d < bestDist {
            bestDist = d
            best = i
        }
    }

    return best
}
