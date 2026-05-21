package leader

import (
    "sync/atomic"
)

type Replica interface {
    Apply(
        id string,
        embedding []float32,
    )
}

type Node struct {
    id int
    replica Replica
}

type Cluster struct {
    nodes []Node
    leader atomic.Int32
}

func New(
    replicas ...Replica,
) *Cluster {

    nodes := make(
        []Node,
        0,
        len(replicas),
    )

    for i, r := range replicas {

        nodes = append(
            nodes,
            Node{
                id: i,
                replica: r,
            },
        )
    }

    c := &Cluster{
        nodes: nodes,
    }

    c.leader.Store(0)

    return c
}

func (c *Cluster) Leader() int {

    return int(
        c.leader.Load(),
    )
}

func (c *Cluster) Elect(
    id int,
) {

    c.leader.Store(
        int32(id),
    )
}

func (c *Cluster) Write(
    id string,
    embedding []float32,
) {

    leader := c.Leader()

    c.nodes[leader].replica.Apply(
        id,
        embedding,
    )
}

func (c *Cluster) Failover() {

    next := (c.Leader() + 1) % len(c.nodes)

    c.Elect(next)
}
