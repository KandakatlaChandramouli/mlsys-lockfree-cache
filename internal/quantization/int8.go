
package quantization

func QuantizeF32ToI8(
    in []float32,
) []int8 {

    out := make(
        []int8,
        len(in),
    )

    for i, v := range in {

        if v > 1 {
            v = 1
        }

        if v < -1 {
            v = -1
        }

        out[i] = int8(v * 127)
    }

    return out
}
