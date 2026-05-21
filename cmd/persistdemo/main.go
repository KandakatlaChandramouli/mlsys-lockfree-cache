
package main

import (
    "fmt"

    "fluxruntime/internal/storage"
    "fluxruntime/internal/vectorstore"
)

func main() {

    store := vectorstore.New()

    store.Add(
        "doc-1",
        []float32{1, 2, 3},
    )

    store.Add(
        "doc-2",
        []float32{4, 5, 6},
    )

    err := storage.Save(
        "vectors.bin",
        store.Records(),
    )

    if err != nil {
        panic(err)
    }

    loaded, err := storage.Load(
        "vectors.bin",
    )

    if err != nil {
        panic(err)
    }

    restored := vectorstore.New()

    restored.Restore(
        loaded,
    )

    results := restored.Search(
        []float32{1, 2, 3},
        2,
    )

    fmt.Println(results)
}
