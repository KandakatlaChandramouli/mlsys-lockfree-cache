
package index

import (
    "fluxruntime/internal/vectorstore"
)

type HNSWIndex struct {
    backend *vectorstore.Store
}

func NewHNSW() *HNSWIndex {

    return &HNSWIndex{
        backend: vectorstore.New(),
    }
}

func (h *HNSWIndex) Add(
    id string,
    embedding []float32,
) {

    h.backend.Add(
        id,
        embedding,
    )
}

func (h *HNSWIndex) Search(
    query []float32,
    topK int,
) []vectorstore.SearchResult {

    return h.backend.Search(
        query,
        topK,
    )
}
