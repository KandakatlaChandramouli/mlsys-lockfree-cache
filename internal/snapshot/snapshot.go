package snapshot

import (
    "bufio"
    "bytes"
    "encoding/gob"
    "os"

    "fluxruntime/internal/compress"
    "fluxruntime/internal/index"
)

type EncodedNode struct {
    ID        string
    Neighbors []int
}

type Snapshot struct {
    Nodes      []EncodedNode
    Embeddings []float32
    Dim        int
}

func Save(
    path string,
    idx *index.HNSWIndex,
) error {

    rawNodes := idx.ExportNodes()

    snap := Snapshot{
        Nodes: make(
            []EncodedNode,
            0,
            len(rawNodes),
        ),
    }

    if len(rawNodes) > 0 {
        snap.Dim = len(
            rawNodes[0].Embedding,
        )
    }

    snap.Embeddings = make(
        []float32,
        0,
        len(rawNodes)*snap.Dim,
    )

    for _, n := range rawNodes {

        snap.Nodes = append(
            snap.Nodes,
            EncodedNode{
                ID: n.ID,
                Neighbors: n.Neighbors,
            },
        )

        snap.Embeddings = append(
            snap.Embeddings,
            n.Embedding...,
        )
    }

    var raw bytes.Buffer

    enc := gob.NewEncoder(
        &raw,
    )

    if err := enc.Encode(
        snap,
    ); err != nil {
        return err
    }

    compressed := compress.Encode(
        raw.Bytes(),
    )

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

    if _, err := w.Write(
        compressed,
    ); err != nil {
        return err
    }

    return w.Flush()
}

func Load(
    path string,
) (*index.HNSWIndex, error) {

    raw, err := os.ReadFile(
        path,
    )

    if err != nil {
        return nil, err
    }

    decoded, err := compress.Decode(
        raw,
    )

    if err != nil {
        return nil, err
    }

    dec := gob.NewDecoder(
        bytes.NewReader(
            decoded,
        ),
    )

    var snap Snapshot

    if err := dec.Decode(
        &snap,
    ); err != nil {
        return nil, err
    }

    nodes := make(
        []index.Node,
        len(snap.Nodes),
    )

    offset := 0

    for i, n := range snap.Nodes {

        emb := make(
            []float32,
            snap.Dim,
        )

        copy(
            emb,
            snap.Embeddings[offset:offset+snap.Dim],
        )

        offset += snap.Dim


        nodes[i] = index.Node{
            ID: n.ID,
            Embedding: emb,
            Neighbors: n.Neighbors,
        }
    }

    idx := index.NewHNSW()

    idx.ImportNodes(
        nodes,
    )

    return idx, nil
}
