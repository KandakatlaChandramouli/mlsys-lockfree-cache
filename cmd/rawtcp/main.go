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

var (
	engine = lockfree.NewSharded()

	connCounter atomic.Uint64
)

func handle(
	conn net.Conn,
) {

	defer conn.Close()

	reader := bufio.NewReaderSize(
		conn,
		64*1024,
	)

	writer := bufio.NewWriterSize(
		conn,
		64*1024,
	)

	connID := connCounter.Add(1)

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

		tokens := make(
			[]uint32,
			hdr.TokenCount,
		)

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

		ok := engine.Submit(
			&lockfree.Request{
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

		resp := <-respCh

		vecLen := uint32(
			len(resp.Vector),
		)

		err = binary.Write(
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
