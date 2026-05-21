package main

import (
	"bufio"
	"encoding/binary"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"syscall"

	"fluxruntime/internal/core"
	"fluxruntime/internal/health"
	"fluxruntime/internal/lockfree"
	"fluxruntime/internal/metrics"
)

type Header struct {
	QueryHash uint64
	TokenCount uint32
}

type Pending struct {
	ID   uint64
	Resp chan *core.RawResponse
}

var (
	engine = lockfree.NewSharded()

	connCounter atomic.Uint64
	reqCounter  atomic.Uint64
)

func writerLoop(
	writer *bufio.Writer,
	pending <-chan Pending,
) {

	for p := range pending {

		resp := <-p.Resp

		vecLen := uint32(
			len(resp.Vector),
		)

		err := binary.Write(
			writer,
			binary.LittleEndian,
			vecLen,
		)

		if err != nil {

			metrics.Failures.Add(1)

			return
		}

		err = binary.Write(
			writer,
			binary.LittleEndian,
			resp.Vector,
		)

		if err != nil {

			metrics.Failures.Add(1)

			return
		}

		err = writer.Flush()

		if err != nil {

			metrics.Failures.Add(1)

			return
		}

		core.ReleaseVector(
			&resp.Vector,
		)

		metrics.Requests.Add(1)
		metrics.QueuedRequests.Add(^uint64(0))
	}
}

func handle(
	conn net.Conn,
) {

	metrics.ActiveConnections.Add(1)

	defer func() {
		metrics.ActiveConnections.Add(^uint64(0))
		conn.Close()
	}()

	reader := bufio.NewReaderSize(
		conn,
		128*1024,
	)

	writer := bufio.NewWriterSize(
		conn,
		128*1024,
	)

	connID := connCounter.Add(1)

	pending := make(
		chan Pending,
		4096,
	)

	go writerLoop(
		writer,
		pending,
	)

	for {

		var hdr Header

		err := binary.Read(
			reader,
			binary.LittleEndian,
			&hdr,
		)

		if err != nil {

			if err != io.EOF &&
				err != syscall.ECONNRESET {

				log.Println(err)

				metrics.Failures.Add(1)
			}

			return
		}

		tokenBuf := core.AcquireTokens()

		tokens := (*tokenBuf)[:hdr.TokenCount]

		err = binary.Read(
			reader,
			binary.LittleEndian,
			&tokens,
		)

		if err != nil {

			metrics.Failures.Add(1)

			return
		}

		respCh := make(
			chan *core.RawResponse,
			1,
		)

		reqID := reqCounter.Add(1)

		metrics.QueuedRequests.Add(1)

		ok := engine.Submit(
			&lockfree.Request{
				ID: reqID,

				Req: &core.RawRequest{
					QueryHash: hdr.QueryHash,
					Tokens:    tokens,
					ConnID:    connID,
				},

				Resp: respCh,
			},
		)

		if !ok {

			metrics.Failures.Add(1)

			return
		}

		pending <- Pending{
			ID:   reqID,
			Resp: respCh,
		}

		core.ReleaseTokens(
			tokenBuf,
		)
	}
}

func metricsServer() {

	mux := http.NewServeMux()

	mux.HandleFunc(
		"/metrics",
		metrics.Handler,
	)

	mux.HandleFunc(
		"/health",
		health.Handler,
	)

	log.Println(
		"🔥 metrics server listening on :9000",
	)

	http.ListenAndServe(
		":9000",
		mux,
	)
}

func main() {

	go metricsServer()

	port := os.Getenv("PORT")

	if port == "" {
		port = "7001"
	}

	_, err := strconv.Atoi(port)

	if err != nil {
		log.Fatal(err)
	}

	addr := ":" + port

	ln, err := net.Listen(
		"tcp",
		addr,
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(
		"🔥 rawtcp runtime listening on",
		addr,
	)

	for {

		conn, err := ln.Accept()

		if err != nil {
			continue
		}

		go handle(conn)
	}
}
