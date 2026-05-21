
package index

import (
    "fluxruntime/internal/vectorstore"
)

type Index interface {

    Add(
        id string,
        embedding []float32,
    )

    Search(
        query []float32,
        topK int,
    ) []vectorstore.SearchResult
}
