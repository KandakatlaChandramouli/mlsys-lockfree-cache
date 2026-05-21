
package simd

import "math"


func Cosine(
    a []float32,
    b []float32,
) float32 {

    var dot float64
    var na float64
    var nb float64

    n := len(a)

    i := 0

    for i+7 < n {

        av0 := float64(a[i+0])
        bv0 := float64(b[i+0])

        av1 := float64(a[i+1])
        bv1 := float64(b[i+1])

        av2 := float64(a[i+2])
        bv2 := float64(b[i+2])

        av3 := float64(a[i+3])
        bv3 := float64(b[i+3])

        av4 := float64(a[i+4])
        bv4 := float64(b[i+4])

        av5 := float64(a[i+5])
        bv5 := float64(b[i+5])

        av6 := float64(a[i+6])
        bv6 := float64(b[i+6])

        av7 := float64(a[i+7])
        bv7 := float64(b[i+7])

        dot +=
            av0*bv0 +
            av1*bv1 +
            av2*bv2 +
            av3*bv3 +
            av4*bv4 +
            av5*bv5 +
            av6*bv6 +
            av7*bv7

        na +=
            av0*av0 +
            av1*av1 +
            av2*av2 +
            av3*av3 +
            av4*av4 +
            av5*av5 +
            av6*av6 +
            av7*av7

        nb +=
            bv0*bv0 +
            bv1*bv1 +
            bv2*bv2 +
            bv3*bv3 +
            bv4*bv4 +
            bv5*bv5 +
            bv6*bv6 +
            bv7*bv7

        i += 8
    }

    for ; i < n; i++ {

        av := float64(a[i])
        bv := float64(b[i])

        dot += av * bv
        na += av * av
        nb += bv * bv
    }

    if na == 0 || nb == 0 {
        return 0
    }

    return float32(dot / (math.Sqrt(na) * math.Sqrt(nb)))
}
