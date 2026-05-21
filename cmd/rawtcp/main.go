package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"sync/atomic"

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

	connID := connCounter.Add(1)

	for {

		var hdr Header

		err := binary.Read(
			conn,
			binary.LittleEndian,
			&hdr,
		)

		if err != nil {

			if err != io.EOF {
				log.Println(err)
			}

			return
		}

		tokens := make(
			[]uint32,
			hdr.TokenCount,
		)

		err = binary.Read(
			conn,
			binary.LittleEndian,
			&tokens,
		)

		if err != nil {
			return
		}

		done := make(
			chan struct{},
		)

		ok := engine.Submit(
			&lockfree.Request{
				Req: &core.RawRequest{
					QueryHash: hdr.QueryHash,
					Tokens: tokens,
					ConnID: connID,
				},

				Callback: func(
					resp *core.RawResponse,
				) {

					vecLen := uint32(
						len(resp.Vector),
					)

					err := binary.Write(
						conn,
						binary.LittleEndian,
						vecLen,
					)

					if err == nil {

						err = binary.Write(
							conn,
							binary.LittleEndian,
							resp.Vector,
						)
					}

					core.ReleaseVector(
						&resp.Vector,
					)

					close(done)
				},
			},
		)

		if !ok {
			return
		}

		<-done
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
