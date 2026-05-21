package index

func (h *HNSWIndex) addUnsafe(
    id string,
    embedding []float32,
) {

    node := Node{
        ID: id,
        Embedding: embedding,
    }

    idx := len(h.nodes)

    if len(h.nodes) > 0 {

        limit := MaxNeighbors

        if len(h.nodes) < limit {
            limit = len(h.nodes)
        }

        for i := 0; i < limit; i++ {

            node.Neighbors = append(
                node.Neighbors,
                i,
            )

            h.nodes[i].Neighbors = append(
                h.nodes[i].Neighbors,
                idx,
            )
        }
    }

    h.nodes = append(
        h.nodes,
        node,
    )
}
