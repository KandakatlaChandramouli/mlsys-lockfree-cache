
package storage

import (
    "encoding/gob"
    "os"

    "fluxruntime/internal/vectorstore"
)

func Save(
    path string,
    data []vectorstore.VectorRecord,
) error {

    f, err := os.Create(path)

    if err != nil {
        return err
    }

    defer f.Close()

    enc := gob.NewEncoder(f)

    return enc.Encode(data)
}

func Load(
    path string,
) ([]vectorstore.VectorRecord, error) {

    f, err := os.Open(path)

    if err != nil {
        return nil, err
    }

    defer f.Close()

    var out []vectorstore.VectorRecord

    dec := gob.NewDecoder(f)

    err = dec.Decode(&out)

    return out, err
}
