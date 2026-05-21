package hnsw

import (
	"container/heap"
	"math"
	"math/rand"
	"sort"
)

type Node struct {
	ID        string
	Vector    []float32
	Level     int
	Neighbors map[int][]int
}

type Graph struct {
	Nodes          []Node
	EntryPoint     int
	MaxLevel       int
	M              int
	EfConstruction int
}

func New(
	m int,
	ef int,
) *Graph {

	return &Graph{
		EntryPoint:     -1,
		M:              m,
		EfConstruction: ef,
	}
}

func randomLevel() int {

	lvl := 0

	for rand.Float64() < 0.5 &&
		lvl < 16 {

		lvl++
	}

	return lvl
}

func cosine(
	a []float32,
	b []float32,
) float32 {

	var dot float32
	var na float32
	var nb float32

	for i := range a {

		dot += a[i] * b[i]
		na += a[i] * a[i]
		nb += b[i] * b[i]
	}

	if na == 0 || nb == 0 {
		return 0
	}

	return dot / float32(
		math.Sqrt(float64(na*nb)),
	)
}

type Candidate struct {
	Index int
	Score float32
}

type MaxHeap []Candidate

func (h MaxHeap) Len() int {
	return len(h)
}

func (h MaxHeap) Less(
	i int,
	j int,
) bool {
	return h[i].Score > h[j].Score
}

func (h MaxHeap) Swap(
	i int,
	j int,
) {
	h[i], h[j] = h[j], h[i]
}

func (h *MaxHeap) Push(
	x interface{},
) {
	*h = append(
		*h,
		x.(Candidate),
	)
}

func (h *MaxHeap) Pop() interface{} {

	old := *h
	n := len(old)

	x := old[n-1]

	*h = old[:n-1]

	return x
}

func (g *Graph) Insert(
	id string,
	vec []float32,
) {

	level := randomLevel()

	node := Node{
		ID:     id,
		Vector: vec,
		Level:  level,
		Neighbors: make(
			map[int][]int,
		),
	}

	idx := len(g.Nodes)

	g.Nodes = append(
		g.Nodes,
		node,
	)

	if g.EntryPoint == -1 {

		g.EntryPoint = idx
		g.MaxLevel = level

		return
	}

	ep := g.EntryPoint

	for l := g.MaxLevel; l > level; l-- {

		ep = g.greedySearch(
			ep,
			vec,
			l,
		)
	}

	for l := min(
		level,
		g.MaxLevel,
	); l >= 0; l-- {

		neighbors := g.searchLayer(
			ep,
			vec,
			l,
			g.EfConstruction,
		)

		if len(neighbors) > g.M {
			neighbors = neighbors[:g.M]
		}

		for _, n := range neighbors {

			g.Nodes[idx].Neighbors[l] = append(
				g.Nodes[idx].Neighbors[l],
				n.Index,
			)

			g.Nodes[n.Index].Neighbors[l] = append(
				g.Nodes[n.Index].Neighbors[l],
				idx,
			)
		}
	}

	if level > g.MaxLevel {

		g.EntryPoint = idx
		g.MaxLevel = level
	}
}

func (g *Graph) greedySearch(
	ep int,
	q []float32,
	level int,
) int {

	best := ep

	improved := true

	for improved {

		improved = false

		bestScore := cosine(
			q,
			g.Nodes[best].Vector,
		)

		for _, n := range g.Nodes[best].Neighbors[level] {

			s := cosine(
				q,
				g.Nodes[n].Vector,
			)

			if s > bestScore {

				best = n
				bestScore = s
				improved = true
			}
		}
	}

	return best
}

func (g *Graph) searchLayer(
	ep int,
	q []float32,
	level int,
	ef int,
) []Candidate {

	visited := make(
		map[int]struct{},
	)

	pq := &MaxHeap{}

	heap.Init(pq)

	heap.Push(
		pq,
		Candidate{
			Index: ep,
			Score: cosine(
				q,
				g.Nodes[ep].Vector,
			),
		},
	)

	visited[ep] = struct{}{}

	out := make(
		[]Candidate,
		0,
		ef,
	)

	for pq.Len() > 0 &&
		len(out) < ef {

		cur := heap.Pop(
			pq,
		).(Candidate)

		out = append(
			out,
			cur,
		)

		for _, n := range g.Nodes[cur.Index].Neighbors[level] {

			if _, ok := visited[n]; ok {
				continue
			}

			visited[n] = struct{}{}

			heap.Push(
				pq,
				Candidate{
					Index: n,
					Score: cosine(
						q,
						g.Nodes[n].Vector,
					),
				},
			)
		}
	}

	sort.Slice(
		out,
		func(i, j int) bool {
			return out[i].Score >
				out[j].Score
		},
	)

	return out
}

func (g *Graph) Search(
	q []float32,
	k int,
	ef int,
) []Candidate {

	ep := g.EntryPoint

	for l := g.MaxLevel; l > 0; l-- {

		ep = g.greedySearch(
			ep,
			q,
			l,
		)
	}

	out := g.searchLayer(
		ep,
		q,
		0,
		ef,
	)

	if len(out) > k {
		out = out[:k]
	}

	return out
}

func min(
	a int,
	b int,
) int {

	if a < b {
		return a
	}

	return b
}
