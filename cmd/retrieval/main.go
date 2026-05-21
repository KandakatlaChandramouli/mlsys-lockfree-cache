
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

    docs := map[string][]float32{
        "transformers": {1, 0, 0},
        "databases":    {0, 1, 0},
        "mlsys":        {0.9, 0.1, 0},
    }

    for id, emb := range docs {

        engine.Insert(
            id,
            emb,
        )
    }

    query := []float32{
        1,
        0,
        0,
    }

    results := engine.Search(
        query,
        3,
    )

    fmt.Println("Top Results")

    for _, r := range results {

        fmt.Printf(
            "%s -> %.4f\n",
            r.ID,
            r.Score,
        )
    }
}
