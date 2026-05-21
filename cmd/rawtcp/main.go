package main

import (
	"bufio"
	"math"
	"encoding/binary"
	"io"
	"log"
	"net"
	"sync"
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

	respPool = sync.Pool{
		New: func() any {

			b := make([]byte, 8192)

			return &b
		},
	}
)

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

		bufPtr := respPool.Get().(*[]byte)

		buf := (*bufPtr)[:0]

		vecLen := uint32(
			len(resp.Vector),
		)

		tmp := make([]byte, 4)

		binary.LittleEndian.PutUint32(
			tmp,
			vecLen,
		)

		buf = append(buf, tmp...)

		for _, f := range resp.Vector {

			bits := binary.LittleEndian.AppendUint32(
				nil,
				math.Float32bits(f),
			)

			buf = append(buf, bits...)
		}

		_, err = writer.Write(buf)

		if err != nil {
			return
		}

		err = writer.Flush()

		if err != nil {
			return
		}

		respPool.Put(bufPtr)

		core.ReleaseVector(
			&resp.Vector,
		)

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
