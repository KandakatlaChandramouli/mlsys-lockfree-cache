package main

import (
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"fluxruntime/internal/cluster"
)

var (
	concurrency = flag.Int("concurrency", 8, "")
	duration    = flag.Duration("duration", 10*time.Second, "")
	tokens      = flag.Int("tokens", 64, "")
)

var (
	latencies []int64
	latMu     sync.Mutex
)

var scheduler = cluster.NewScheduler(
	[]string{
		"localhost:7001",
		"localhost:7002",
		"localhost:7003",
	},
)

func worker(
	ctx context.Context,
	reqs *atomic.Uint64,
	fail *atomic.Uint64,
	tokenCount int,
) {

	node := scheduler.Pick()

	defer scheduler.Complete(
		node,
	)

	conn, err := net.Dial(
		"tcp",
		node.Addr,
	)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	toks := make(
		[]uint32,
		tokenCount,
	)

	for {

		select {

		case <-ctx.Done():
			return

		default:
		}

		start := time.Now()

		hash := uint64(
			time.Now().UnixNano(),
		)

		err = binary.Write(
			conn,
			binary.LittleEndian,
			hash,
		)

		if err != nil {

			fail.Add(1)

			continue
		}

		err = binary.Write(
			conn,
			binary.LittleEndian,
			uint32(tokenCount),
		)

		if err != nil {

			fail.Add(1)

			continue
		}

		err = binary.Write(
			conn,
			binary.LittleEndian,
			toks,
		)

		if err != nil {

			fail.Add(1)

			continue
		}

		var vecLen uint32

		err = binary.Read(
			conn,
			binary.LittleEndian,
			&vecLen,
		)

		if err != nil {

			fail.Add(1)

			continue
		}

		vec := make(
			[]float32,
			vecLen,
		)

		err = binary.Read(
			conn,
			binary.LittleEndian,
			&vec,
		)

		if err != nil {

			fail.Add(1)

			continue
		}

		lat := time.Since(start)

		latMu.Lock()

		latencies = append(
			latencies,
			lat.Nanoseconds(),
		)

		latMu.Unlock()

		reqs.Add(1)
	}
}

func percentile(
	p float64,
) time.Duration {

	if len(latencies) == 0 {
		return 0
	}

	idx := int(
		float64(len(latencies)) * p,
	)

	if idx >= len(latencies) {
		idx = len(latencies) - 1
	}

	return time.Duration(
		latencies[idx],
	)
}

func main() {

	flag.Parse()

	ctx, cancel := context.WithTimeout(
		context.Background(),
		*duration,
	)

	defer cancel()

	var reqs atomic.Uint64
	var fail atomic.Uint64

	for i := 0; i < *concurrency; i++ {

		go worker(
			ctx,
			&reqs,
			&fail,
			*tokens,
		)
	}

	<-ctx.Done()

	sort.Slice(
		latencies,
		func(i, j int) bool {
			return latencies[i] < latencies[j]
		},
	)

	rps := float64(
		reqs.Load(),
	) / duration.Seconds()

	fmt.Println("")
	fmt.Println("==== ADAPTIVE DISTRIBUTED BENCHMARK ====")
	fmt.Println("")

	fmt.Printf(
		"Requests/sec %.2f\n",
		rps,
	)

	fmt.Printf(
		"Total Requests %d\n",
		reqs.Load(),
	)

	fmt.Printf(
		"Failures %d\n",
		fail.Load(),
	)

	fmt.Println("")

	fmt.Printf(
		"p50 %v\n",
		percentile(0.50),
	)

	fmt.Printf(
		"p95 %v\n",
		percentile(0.95),
	)

	fmt.Printf(
		"p99 %v\n",
		percentile(0.99),
	)

	fmt.Println("")
}
