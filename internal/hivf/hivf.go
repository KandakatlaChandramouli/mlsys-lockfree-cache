package hivf

import (
    "fluxruntime/internal/kmeans"
)

type Leaf struct {
    Centroids []kmeans.Centroid
}

type Root struct {
    Coarse []kmeans.Centroid
    Leaves []Leaf
}

func Build(
    coarse int,
    perLeaf int,
) *Root {

    leaves := make(
        []Leaf,
        coarse,
    )

    coarseVecs := make(
        []kmeans.Centroid,
        coarse,
    )

    for i := range leaves {

        leaves[i] = Leaf{
            Centroids: make(
                []kmeans.Centroid,
                perLeaf,
            ),
        }
    }

    return &Root{
        Coarse: coarseVecs,
        Leaves: leaves,
    }
}

func Route(
    query []float32,
    root *Root,
) (int, int) {

    coarseIdx := kmeans.Nearest(
        query,
        root.Coarse,
    )

    leaf := root.Leaves[coarseIdx]

    leafIdx := kmeans.Nearest(
        query,
        leaf.Centroids,
    )

    return coarseIdx, leafIdx
}
