package index

func (h *HNSWIndex) Insert(
    id string,
    embedding []float32,
) {

    h.mu.Lock()
    defer h.mu.Unlock()

    h.addUnsafe(
        id,
        embedding,
    )
}
