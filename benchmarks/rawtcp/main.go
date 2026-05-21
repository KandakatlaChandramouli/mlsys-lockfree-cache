package main

import (
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"sync/atomic"
	"time"
)

var (
	concurrency = flag.Int("concurrency", 8, "")
	duration = flag.Duration("duration", 10*time.Second, "")
	tokens = flag.Int("tokens", 64, "")
)

func worker(
	ctx context.Context,
	reqs *atomic.Uint64,
	fail *atomic.Uint64,
	tokenCount int,
) {

	conn, err := net.Dial(
		"tcp",
		"localhost:7001",
	)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	toks := make([]uint32, tokenCount)

	for {

		select {

		case <-ctx.Done():
			return

		default:
		}

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

		reqs.Add(1)
	}
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

	rps := float64(
		reqs.Load(),
	) / duration.Seconds()

	fmt.Println("")
	fmt.Println("==== RAW TCP BENCHMARK ====")
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
}
