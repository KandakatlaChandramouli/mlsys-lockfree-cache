
package pq

type QuantizedVector struct {
    Codes []uint8
}

type ProductQuantizer struct {
    Dim       int
    Subspaces int
    Ks        int
}

func New(
    dim int,
    subspaces int,
    ks int,
) *ProductQuantizer {

    return &ProductQuantizer{
        Dim: dim,
        Subspaces: subspaces,
        Ks: ks,
    }
}

func (p *ProductQuantizer) Encode(
    vec []float32,
) QuantizedVector {

    codes := make(
        []uint8,
        p.Subspaces,
    )

    chunk := p.Dim / p.Subspaces

    for i := 0; i < p.Subspaces; i++ {

        start := i * chunk
        end := start + chunk

        var sum float32

        for j := start; j < end; j++ {
            sum += vec[j]
        }

        codes[i] = uint8(
            int(sum * 10) % 256,
        )
    }

    return QuantizedVector{
        Codes: codes,
    }
}
