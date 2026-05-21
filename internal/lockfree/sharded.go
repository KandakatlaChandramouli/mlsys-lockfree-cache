package lockfree

import (
	"hash/fnv"
	"runtime"
)

type ShardedEngine struct {
	shards []*Engine
	count  uint64
}

func NewSharded() *ShardedEngine {

	n := runtime.NumCPU()

	shards := make(
		[]*Engine,
		n,
	)

	for i := range shards {
		shards[i] = NewEngine(i)
	}

	return &ShardedEngine{
		shards: shards,
		count:  uint64(n),
	}
}

func (s *ShardedEngine) route(
	hash uint64,
) *Engine {

	h := fnv.New64a()

	var b [8]byte

	b[0] = byte(hash)
	b[1] = byte(hash >> 8)
	b[2] = byte(hash >> 16)
	b[3] = byte(hash >> 24)

	h.Write(b[:])

	idx := h.Sum64() % s.count

	return s.shards[idx]
}

func (s *ShardedEngine) Submit(
	req *Request,
) bool {

	engine := s.route(
		req.Req.QueryHash,
	)

	return engine.Submit(req)
}
