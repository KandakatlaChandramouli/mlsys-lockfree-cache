package main

import (
	"log"
	"net"
)

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

		conn.Close()
	}
}
