
package ivf

import (
    "math/rand"
)

type Cluster struct {
    Vectors [][]float32
}

type Index struct {
    Clusters []Cluster
    NList    int
}

func New(
    nlist int,
) *Index {

    clusters := make(
        []Cluster,
        nlist,
    )

    return &Index{
        Clusters: clusters,
        NList: nlist,
    }
}

func (i *Index) Add(
    vec []float32,
) {

    bucket := rand.Intn(
        i.NList,
    )

    i.Clusters[bucket].Vectors = append(
        i.Clusters[bucket].Vectors,
        vec,
    )
}
