package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"
)

type Request struct {
	QueryHash uint64
	TokenCount uint32
}

func handle(
	conn net.Conn,
) {

	defer conn.Close()

	for {

		var req Request

		err := binary.Read(
			conn,
			binary.LittleEndian,
			&req,
		)

		if err != nil {

			if err != io.EOF {
				log.Println(err)
			}

			return
		}

		tokens := make(
			[]uint32,
			req.TokenCount,
		)

		err = binary.Read(
			conn,
			binary.LittleEndian,
			&tokens,
		)

		if err != nil {
			return
		}

		vecLen := uint32(len(tokens))

		err = binary.Write(
			conn,
			binary.LittleEndian,
			vecLen,
		)

		if err != nil {
			return
		}

		resp := make(
			[]float32,
			vecLen,
		)

		for i := range resp {
			resp[i] = float32(i) * 0.1
		}

		err = binary.Write(
			conn,
			binary.LittleEndian,
			resp,
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
