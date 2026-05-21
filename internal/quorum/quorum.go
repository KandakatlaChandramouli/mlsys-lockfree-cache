package quorum

import (
    "sync"
    "sync/atomic"
)

type Replica interface {
    Apply(
        id string,
        embedding []float32,
    )
}

type Cluster struct {
    replicas []Replica
    quorum   int
}

func New(
    quorum int,
    replicas ...Replica,
) *Cluster {

    return &Cluster{
        replicas: replicas,
        quorum: quorum,
    }
}

func (c *Cluster) Write(
    id string,
    embedding []float32,
) bool {

    var success atomic.Int32

    var wg sync.WaitGroup

    wg.Add(
        len(c.replicas),
    )

    for _, replica := range c.replicas {

        rep := replica

        go func() {

            defer wg.Done()

            rep.Apply(
                id,
                embedding,
            )

            success.Add(
                1,
            )
        }()
    }

    wg.Wait()

    return int(
        success.Load(),
    ) >= c.quorum
}
