
package avx

func DotProduct(
    a []float32,
    b []float32,
) float32 {

    var out float32

    n := len(a)

    for i := 0; i < n; i++ {
        out += a[i] * b[i]
    }

    return out
}
