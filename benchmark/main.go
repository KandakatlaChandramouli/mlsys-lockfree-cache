
package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"net"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

var (
	concurrency = flag.Int("concurrency", 32, "")
	duration    = flag.Duration("duration", 10*time.Second, "")
)

var latencies []int64
var latMu sync.Mutex

func worker(
	reqs *atomic.Uint64,
	fail *atomic.Uint64,
	stop <-chan struct{},
	wg *sync.WaitGroup,
) {

	defer wg.Done()

	conn, err := net.Dial(
		"tcp",
		"localhost:7001",
	)

	if err != nil {
		fail.Add(1)
		return
	}

	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	payload := []byte(
		"hello inference runtime",
	)

	for {

		select {

		case <-stop:
			return

		default:
		}

		start := time.Now()

		size := uint32(len(payload))

		err = binary.Write(
			writer,
			binary.LittleEndian,
			size,
		)

		if err != nil {
			fail.Add(1)
			return
		}

		_, err = writer.Write(payload)

		if err != nil {
			fail.Add(1)
			return
		}

		err = writer.Flush()

		if err != nil {
			fail.Add(1)
			return
		}

		var outSize uint32

		err = binary.Read(
			reader,
			binary.LittleEndian,
			&outSize,
		)

		if err != nil {
			fail.Add(1)
			return
		}

		buf := make(
			[]byte,
			outSize*4,
		)

		_, err = io.ReadFull(
			reader,
			buf,
		)

		if err != nil {
			fail.Add(1)
			return
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
		float64(len(latencies)-1) * p,
	)

	return time.Duration(
		latencies[idx],
	)
}

func main() {

	flag.Parse()

	var reqs atomic.Uint64
	var fail atomic.Uint64

	stop := make(
		chan struct{},
	)

	var wg sync.WaitGroup

	for i := 0; i < *concurrency; i++ {

		wg.Add(1)

		go worker(
			&reqs,
			&fail,
			stop,
			&wg,
		)
	}

	time.Sleep(*duration)

	close(stop)

	wg.Wait()

	sort.Slice(
		latencies,
		func(i, j int) bool {
			return latencies[i] < latencies[j]
		},
	)

	rps := float64(
		reqs.Load(),
	) / duration.Seconds()

	fmt.Println()
	fmt.Println("==== INFERENCE BENCHMARK ====")
	fmt.Println()
	fmt.Printf("Requests/sec %.2f\n", rps)
	fmt.Printf("Total Requests %d\n", reqs.Load())
	fmt.Printf("Failures %d\n", fail.Load())
	fmt.Println()
	fmt.Printf("p50 %v\n", percentile(0.50))
	fmt.Printf("p95 %v\n", percentile(0.95))
	fmt.Printf("p99 %v\n", percentile(0.99))
}
