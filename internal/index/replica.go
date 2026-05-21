
package index

func (h *HNSWIndex) Apply(
    id string,
    embedding []float32,
) {

    h.Insert(
        id,
        embedding,
    )
}
