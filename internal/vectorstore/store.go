
package vectorstore

import (
    "math"
    "sort"
    "sync"
)

type VectorRecord struct {
    ID        string
    Embedding []float32
}

type SearchResult struct {
    ID    string
    Score float32
}

type Store struct {
    mu      sync.RWMutex
    vectors []VectorRecord
}

func New() *Store {

    return &Store{}
}

func (s *Store) Add(
    id string,
    embedding []float32,
) {

    s.mu.Lock()
    defer s.mu.Unlock()

    cp := make(
            []float32,
            len(embedding),
    )

    copy(cp, embedding)

    s.vectors = append(
            s.vectors,
            VectorRecord{
                    ID: id,
                    Embedding: cp,
            },
    )
}

func (s *Store) Search(
    query []float32,
    topK int,
) []SearchResult {

    s.mu.RLock()
    defer s.mu.RUnlock()

    out := make(
            []SearchResult,
            0,
            len(s.vectors),
    )

    for _, v := range s.vectors {

            score := cosine(
                    query,
                    v.Embedding,
            )

            out = append(
                    out,
                    SearchResult{
                            ID: v.ID,
                            Score: score,
                    },
            )
    }

    sort.Slice(
            out,
            func(i, j int) bool {
                    return out[i].Score > out[j].Score
            },
    )

    if len(out) > topK {
            out = out[:topK]
    }

    return out
}




func cosine(
    a []float32,
    b []float32,
) float32 {

    var dot float64
    var na float64
    var nb float64

    for i := range a {

        av := float64(a[i])
        bv := float64(b[i])

        dot += av * bv
        na += av * av
        nb += bv * bv
    }

    if na == 0 || nb == 0 {
        return 0
    }

    return float32(dot / (math.Sqrt(na) * math.Sqrt(nb)))
}
