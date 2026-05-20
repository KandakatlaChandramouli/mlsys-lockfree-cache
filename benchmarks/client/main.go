package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"sync"
	"sync/atomic"
	"text/tabwriter"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	routerv1 "fluxruntime/proto/v1"
)

var (
	flagAddr        = flag.String("addr", "localhost:50051", "gRPC server address")
	flagConcurrency = flag.Int("concurrency", 32, "number of concurrent workers")
	flagDuration    = flag.Duration("duration", 30*time.Second, "benchmark duration")
	flagTokens      = flag.Int("tokens", 64, "tokens per request")
	flagTraffic     = flag.String("traffic", "uniform", "uniform | zipf")
	flagRampUp      = flag.Duration("ramp-up", 2*time.Second, "warmup duration")
	flagTimeout     = flag.Duration("request-timeout", 5*time.Second, "request timeout")
	flagConnections = flag.Int("connections", 4, "grpc connection pool size")
)

type histogram struct {
	samples []int64
}

func mergeHistograms(locals [][]int64) histogram {
	total := 0
	for _, s := range locals {
		total += len(s)
	}

	merged := make([]int64, 0, total)

	for _, s := range locals {
		merged = append(merged, s...)
	}

	sort.Slice(merged, func(i, j int) bool {
		return merged[i] < merged[j]
	})

	return histogram{samples: merged}
}

func (h *histogram) percentile(p float64) time.Duration {
	if len(h.samples) == 0 {
		return 0
	}

	idx := int(math.Ceil((p / 100.0) * float64(len(h.samples))))
	idx--

	if idx < 0 {
		idx = 0
	}

	if idx >= len(h.samples) {
		idx = len(h.samples) - 1
	}

	return time.Duration(h.samples[idx]) * time.Microsecond
}

func buildConnPool(addr string, n int) ([]*grpc.ClientConn, error) {
	pool := make([]*grpc.ClientConn, n)

	for i := range pool {
		conn, err := grpc.Dial(
			addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                20 * time.Second,
				Timeout:             10 * time.Second,
				PermitWithoutStream: true,
			}),
		)

		if err != nil {
			return nil, err
		}

		pool[i] = conn
	}

	return pool, nil
}

type workerResult struct {
	latencies []int64
	requests  int64
	failures  int64
}

func runWorker(
	ctx context.Context,
	client routerv1.RouterServiceClient,
	tokenCount int,
	measuring *atomic.Bool,
	reqTimeout time.Duration,
) workerResult {

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var res workerResult

	res.latencies = make([]int64, 0, 4096)

	for {

		select {
		case <-ctx.Done():
			return res
		default:
		}

		tokens := make([]uint32, tokenCount)

		for i := range tokens {
			tokens[i] = rng.Uint32()
		}

		req := &routerv1.RouteRequest{
			QueryHash: rng.Uint64(),
			Tokens:    tokens,
		}

		rctx, cancel := context.WithTimeout(ctx, reqTimeout)

		start := time.Now()

		_, err := client.Route(rctx, req)

		elapsed := time.Since(start)

		cancel()

		if !measuring.Load() {
			continue
		}

		res.requests++

		if err != nil {
			res.failures++
			continue
		}

		res.latencies = append(
			res.latencies,
			elapsed.Microseconds(),
		)
	}
}

func main() {

	flag.Parse()

	pool, err := buildConnPool(
		*flagAddr,
		*flagConnections,
	)

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		for _, c := range pool {
			c.Close()
		}
	}()

	clients := make([]routerv1.RouterServiceClient, len(pool))

	for i, conn := range pool {
		clients[i] = routerv1.NewRouterServiceClient(conn)
	}

	totalDuration := *flagRampUp + *flagDuration

	ctx, cancel := context.WithTimeout(
		context.Background(),
		totalDuration+5*time.Second,
	)

	defer cancel()

	var measuring atomic.Bool

	results := make([]workerResult, *flagConcurrency)

	var wg sync.WaitGroup

	log.Printf(
		"starting benchmark: concurrency=%d duration=%s",
		*flagConcurrency,
		*flagDuration,
	)

	for i := 0; i < *flagConcurrency; i++ {

		wg.Add(1)

		go func(id int) {

			defer wg.Done()

			results[id] = runWorker(
				ctx,
				clients[id%len(clients)],
				*flagTokens,
				&measuring,
				*flagTimeout,
			)

		}(i)
	}

	time.Sleep(*flagRampUp)

	measuring.Store(true)

	start := time.Now()

	time.Sleep(*flagDuration)

	elapsed := time.Since(start)

	cancel()

	wg.Wait()

	var totalReqs int64
	var totalFail int64

	localLatencies := make([][]int64, *flagConcurrency)

	for i, r := range results {
		totalReqs += r.requests
		totalFail += r.failures
		localLatencies[i] = r.latencies
	}

	hist := mergeHistograms(localLatencies)

	rps := float64(totalReqs) / elapsed.Seconds()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "==== MLSYS RUNTIME BENCHMARK ====")
	fmt.Fprintln(w, "")

	fmt.Fprintf(w, "Requests/sec\t%.2f\n", rps)
	fmt.Fprintf(w, "Total Requests\t%d\n", totalReqs)
	fmt.Fprintf(w, "Failures\t%d\n", totalFail)

	fmt.Fprintln(w, "")

	fmt.Fprintf(w, "p50\t%s\n", hist.percentile(50))
	fmt.Fprintf(w, "p95\t%s\n", hist.percentile(95))
	fmt.Fprintf(w, "p99\t%s\n", hist.percentile(99))

	fmt.Fprintln(w, "")

	w.Flush()
}
