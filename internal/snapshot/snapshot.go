package snapshot

import (
    "encoding/gob"
    "os"

    "fluxruntime/internal/index"
)

type Snapshot struct {
    Nodes []index.Node
}

func Save(
    path string,
    idx *index.HNSWIndex,
) error {

    f, err := os.Create(
        path,
    )

    if err != nil {
        return err
    }

    defer f.Close()

    enc := gob.NewEncoder(
        f,
    )

    snap := Snapshot{
        Nodes: idx.ExportNodes(),
    }

    return enc.Encode(
        snap,
    )
}

func Load(
    path string,
) (*index.HNSWIndex, error) {

    f, err := os.Open(
        path,
    )

    if err != nil {
        return nil, err
    }

    defer f.Close()

    dec := gob.NewDecoder(
        f,
    )

    var snap Snapshot

    err = dec.Decode(
        &snap,
    )

    if err != nil {
        return nil, err
    }

    idx := index.NewHNSW()

    idx.ImportNodes(
        snap.Nodes,
    )

    return idx, nil
}
