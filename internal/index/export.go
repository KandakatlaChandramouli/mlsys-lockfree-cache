package index

func (h *HNSWIndex) ExportNodes() []Node {

    h.mu.RLock()
    defer h.mu.RUnlock()

    out := make(
        []Node,
        len(h.nodes),
    )

    copy(
        out,
        h.nodes,
    )

    return out
}

func (h *HNSWIndex) ImportNodes(
    nodes []Node,
) {

    h.mu.Lock()
    defer h.mu.Unlock()

    h.nodes = nodes
}
