
package index

import (
    "fluxruntime/internal/simd"
    "sort"
    "sync"

    "fluxruntime/internal/vectorstore"
)

const (
    MaxNeighbors = 8
    BeamWidth    = 16
)

type HNSWIndex struct {
    mu    sync.RWMutex
    nodes []Node
}

func NewHNSW() *HNSWIndex {

    return &HNSWIndex{}
}

func (h *HNSWIndex) Add(
    id string,
    embedding []float32,
) {

    h.mu.Lock()
    defer h.mu.Unlock()

    node := Node{
        ID: id,
        Embedding: embedding,
    }

    idx := len(h.nodes)

    if idx > 0 {

        type pair struct {
            idx   int
            score float32
        }

        candidates := make(
            []pair,
            0,
            len(h.nodes),
        )

        start := 0

    if len(h.nodes) > 128 {
        start = len(h.nodes) - 128
    }

    for i := start; i < len(h.nodes); i++ {

        n := h.nodes[i]

            score := cosine(
                embedding,
                n.Embedding,
            )

            candidates = append(
                candidates,
                pair{
                    idx: i,
                    score: score,
                },
            )
        }

        sort.Slice(
            candidates,
            func(i, j int) bool {
                return candidates[i].score >
                       candidates[j].score
            },
        )

        limit := MaxNeighbors

        if len(candidates) < limit {
            limit = len(candidates)
        }

        for i := 0; i < limit; i++ {

            neighbor := candidates[i].idx

            node.Neighbors = append(
                node.Neighbors,
                neighbor,
            )

            h.nodes[neighbor].Neighbors = append(
                h.nodes[neighbor].Neighbors,
                idx,
            )
        }
    }

    h.nodes = append(
        h.nodes,
        node,
    )
}

func (h *HNSWIndex) Search(
    query []float32,
    topK int,
) []vectorstore.SearchResult {

    h.mu.RLock()
    defer h.mu.RUnlock()

    if len(h.nodes) == 0 {
        return nil
    }

    visited := make(
        map[int]bool,
        BeamWidth * 4,
    )

    frontier := make([]int, 0, BeamWidth)
    frontier = append(frontier, 0)

    best := make(
        []vectorstore.SearchResult,
        0,
        BeamWidth,
    )

    for len(frontier) > 0 {

        if len(frontier) > BeamWidth {
            frontier = frontier[:BeamWidth]
        }

        current := frontier[0]
        frontier = frontier[1:]

        if visited[current] {
            continue
        }

        visited[current] = true

        node := h.nodes[current]

        score := cosine(
            query,
            node.Embedding,
        )

        best = append(
            best,
            vectorstore.SearchResult{
                ID: node.ID,
                Score: score,
            },
        )

        for _, n := range node.Neighbors {

            if !visited[n] {
                frontier = append(
                    frontier,
                    n,
                )
            }
        }
    }

    sort.Slice(
        best,
        func(i, j int) bool {
            return best[i].Score >
                   best[j].Score
        },
    )

    if len(best) > topK {
        best = best[:topK]
    }

    return best
}





func cosine(
    a []float32,
    b []float32,
) float32 {

    return simd.Cosine(
        a,
        b,
    )
}
