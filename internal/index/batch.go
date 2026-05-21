package index

import (
    "fluxruntime/internal/ingest"
)

func (h *HNSWIndex) InsertBatch(
    jobs []ingest.Job,
) {

    h.mu.Lock()
    defer h.mu.Unlock()

    for _, job := range jobs {

        h.addUnsafe(
            job.ID,
            job.Embedding,
        )
    }
}
