package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"

	"fluxruntime/internal/core"
	"fluxruntime/internal/lockfree"
)

type Header struct {
	QueryHash uint64
	TokenCount uint32
}

var engine = lockfree.NewSharded()

func handle(
	conn net.Conn,
) {

	defer conn.Close()

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

		req := &core.RawRequest{
			QueryHash: hdr.QueryHash,
			Tokens: tokens,
		}

		resp := engine.Submit(req)

		vecLen := uint32(
			len(resp.Vector),
		)

		err = binary.Write(
			conn,
			binary.LittleEndian,
			vecLen,
		)

		if err != nil {
			return
		}

		err = binary.Write(
			conn,
			binary.LittleEndian,
			resp.Vector,
		)

		if err != nil {
			return
		}
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
