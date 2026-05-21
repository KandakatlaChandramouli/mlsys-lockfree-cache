package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"

	"fluxruntime/internal/hnsw"
)

const (
	Dim      = 384
	Vectors  = 10000
	Queries  = 100
	TopK     = 10
	EfSearch = 128
)

func randomVector() []float32 {

	out := make(
		[]float32,
		Dim,
	)

	var norm float32

	for i := range out {

		v := rand.Float32()

		out[i] = v

		norm += v * v
	}

	inv := float32(
		1.0 / math.Sqrt(
			float64(norm),
		),
	)

	for i := range out {
		out[i] *= inv
	}

	return out
}

func cosine(
	a []float32,
	b []float32,
) float32 {

	var s float32

	for i := range a {
		s += a[i] * b[i]
	}

	return s
}

func bruteForce(
	db [][]float32,
	q []float32,
) []int {

	type pair struct {
		idx int
		sim float32
	}

	scores := make(
		[]pair,
		0,
		len(db),
	)

	for i, v := range db {

		scores = append(
			scores,
			pair{
				idx: i,
				sim: cosine(
					q,
					v,
				),
			},
		)
	}

	sort.Slice(
		scores,
		func(i, j int) bool {
			return scores[i].sim >
				scores[j].sim
		},
	)

	out := make(
		[]int,
		0,
		TopK,
	)

	for i := 0; i < TopK; i++ {
		out = append(
			out,
			scores[i].idx,
		)
	}

	return out
}

func recall(
	gt []int,
	pred []hnsw.Candidate,
) float64 {

	m := make(
		map[int]struct{},
	)

	for _, g := range gt {
		m[g] = struct{}{}
	}

	hit := 0

	for _, p := range pred {

		if _, ok := m[p.Index]; ok {

			hit++
		}
	}

	return float64(hit) /
		float64(len(gt))
}

func main() {

	idx := hnsw.New(
		16,
		200,
	)

	db := make(
		[][]float32,
		0,
		Vectors,
	)

	for i := 0; i < Vectors; i++ {

		v := randomVector()

		db = append(
			db,
			v,
		)

		idx.Insert(
			fmt.Sprintf(
				"vec-%d",
				i,
			),
			v,
		)
	}

	totalRecall := 0.0

	latencies := make(
		[]int64,
		0,
		Queries,
	)

	for i := 0; i < Queries; i++ {

		q := randomVector()

		gt := bruteForce(
			db,
			q,
		)

		start := time.Now()

		pred := idx.Search(
			q,
			TopK,
			EfSearch,
		)

		latencies = append(
			latencies,
			time.Since(start).Nanoseconds(),
		)

		totalRecall += recall(
			gt,
			pred,
		)
	}

	sort.Slice(
		latencies,
		func(i, j int) bool {
			return latencies[i] <
				latencies[j]
		},
	)

	fmt.Println(
		"queries:",
		Queries,
	)

	fmt.Println(
		"recall@10:",
		totalRecall/Queries,
	)

	fmt.Println(
		"p50(ns):",
		latencies[len(latencies)*50/100],
	)

	fmt.Println(
		"p95(ns):",
		latencies[len(latencies)*95/100],
	)

	fmt.Println(
		"p99(ns):",
		latencies[len(latencies)*99/100],
	)
}
