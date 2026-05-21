
package replication

import (
    "sync"
)

type Replica interface {
    Apply(
        id string,
        embedding []float32,
    )
}

type Replicator struct {
    replicas []Replica
}

func New(
    replicas ...Replica,
) *Replicator {

    return &Replicator{
        replicas: replicas,
    }
}

func (r *Replicator) Broadcast(
    id string,
    embedding []float32,
) {

    var wg sync.WaitGroup

    wg.Add(
        len(r.replicas),
    )

    for _, replica := range r.replicas {

        rep := replica

        go func() {

            defer wg.Done()

            rep.Apply(
                id,
                embedding,
            )
        }()
    }

    wg.Wait()
}
