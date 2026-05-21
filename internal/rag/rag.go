
package rag

import (
    "fluxruntime/internal/vectorstore"
)

type Engine struct {
    store *vectorstore.Store
}

func New(
    store *vectorstore.Store,
) *Engine {

    return &Engine{
        store: store,
    }
}

func (e *Engine) Insert(
    id string,
    embedding []float32,
) {

    e.store.Add(
        id,
        embedding,
    )
}

func (e *Engine) Search(
    embedding []float32,
    topK int,
) []vectorstore.SearchResult {

    return e.store.Search(
        embedding,
        topK,
    )
}
