package main

import (
	"bufio"
	"encoding/binary"
	"io"
	"log"
	"net"
	"sync/atomic"
	"syscall"

	"fluxruntime/internal/core"
	"fluxruntime/internal/lockfree"
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
			return
		}

		err = binary.Write(
			writer,
			binary.LittleEndian,
			resp.Vector,
		)

		if err != nil {
			return
		}

		err = writer.Flush()

		if err != nil {
			return
		}

		core.ReleaseVector(
			&resp.Vector,
		)
	}
}

func handle(
	conn net.Conn,
) {

	defer conn.Close()

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
			return
		}

		respCh := make(
			chan *core.RawResponse,
			1,
		)

		reqID := reqCounter.Add(1)

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

func main() {

	ln, err := net.Listen(
		"tcp",
		":7001",
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("🔥 rawtcp runtime listening on :7001")

	for {

		conn, err := ln.Accept()

		if err != nil {
			continue
		}

		go handle(conn)
	}
}
