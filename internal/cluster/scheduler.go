package cluster

type Scheduler struct {
	nodes []*Node
}

func NewScheduler(
	addrs []string,
) *Scheduler {

	nodes := make(
		[]*Node,
		len(addrs),
	)

	for i, addr := range addrs {

		nodes[i] = NewNode(
			addr,
		)
	}

	return &Scheduler{
		nodes: nodes,
	}
}

func (s *Scheduler) Pick() *Node {

	best := s.nodes[0]

	bestScore := best.Score()

	for i := 1; i < len(s.nodes); i++ {

		score := s.nodes[i].Score()

		if score < bestScore {

			best = s.nodes[i]

			bestScore = score
		}
	}

	best.Queued.Add(1)

	return best
}

func (s *Scheduler) Complete(
	n *Node,
) {

	n.Queued.Add(^uint64(0))
}
