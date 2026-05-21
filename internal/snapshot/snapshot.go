package snapshot

import (
    "bufio"
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

    w := bufio.NewWriterSize(
        f,
        1<<20,
    )

    enc := gob.NewEncoder(
        w,
    )

    snap := Snapshot{
        Nodes: idx.ExportNodes(),
    }

    if err := enc.Encode(
        snap,
    ); err != nil {
        return err
    }

    return w.Flush()
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

    r := bufio.NewReaderSize(
        f,
        1<<20,
    )

    dec := gob.NewDecoder(
        r,
    )

    var snap Snapshot

    if err := dec.Decode(
        &snap,
    ); err != nil {
        return nil, err
    }

    idx := index.NewHNSW()

    idx.ImportNodes(
        snap.Nodes,
    )

    return idx, nil
}
