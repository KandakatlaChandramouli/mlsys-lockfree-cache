
package main

import (
    "fmt"

    "fluxruntime/internal/rag"
    "fluxruntime/internal/vectorstore"
)

func main() {

    store := vectorstore.New()

    engine := rag.New(
        store,
    )

    engine.Insert(
        "doc-1",
        []float32{1, 0, 0},
    )

    engine.Insert(
        "doc-2",
        []float32{0, 1, 0},
    )

    engine.Insert(
        "doc-3",
        []float32{0.9, 0.1, 0},
    )

    results := engine.Search(
        []float32{1, 0, 0},
        2,
    )

    for _, r := range results {

        fmt.Printf(
            "%s score=%.4f\n",
            r.ID,
            r.Score,
        )
    }
}
