
package kmeans

func L2(
    a []float32,
    b []float32,
) float32 {

    var out float32

    for i := range a {

        d := a[i] - b[i]
        out += d * d
    }

    return out
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
