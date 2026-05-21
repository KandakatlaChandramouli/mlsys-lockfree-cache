package compress

import (
    "bytes"

    "github.com/klauspost/compress/zstd"
)

var encoder, _ = zstd.NewWriter(
    nil,
    zstd.WithEncoderLevel(
        zstd.SpeedFastest,
    ),
)

var decoder, _ = zstd.NewReader(
    nil,
)

func Encode(
    b []byte,
) []byte {

    return encoder.EncodeAll(
        b,
        make(
            []byte,
            0,
            len(b)/4,
        ),
    )
}

func Decode(
    b []byte,
) ([]byte, error) {

    out, err := decoder.DecodeAll(
        b,
        nil,
    )

    if err != nil {
        return nil, err
    }

    return bytes.Clone(
        out,
    ), nil
}
