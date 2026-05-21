package cluster

import (
	"sync/atomic"
)

type Node struct {
	Addr string

	Queued atomic.Uint64
	Active atomic.Uint64
}

func NewNode(
	addr string,
) *Node {

	return &Node{
		Addr: addr,
	}
}

func (n *Node) Score() uint64 {

	return n.Queued.Load()*2 +
		n.Active.Load()
}
