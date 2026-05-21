
package cache

import (
    "sync"
)

type EmbeddingCache struct {
    mu sync.RWMutex

    data map[string][]float32
}

func New() *EmbeddingCache {

    return &EmbeddingCache{
        data: make(
            map[string][]float32,
        ),
    }
}

func (c *EmbeddingCache) Get(
    key string,
) ([]float32, bool) {

    c.mu.RLock()
    defer c.mu.RUnlock()

    v, ok := c.data[key]

    return v, ok
}

func (c *EmbeddingCache) Put(
    key string,
    val []float32,
) {

    cp := make(
        []float32,
        len(val),
    )

    copy(cp, val)

    c.mu.Lock()
    c.data[key] = cp
    c.mu.Unlock()
}
