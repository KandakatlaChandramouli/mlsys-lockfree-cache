package loadbalancer

import (
	"io"
	"log"
	"net"
	"sync/atomic"
)

type Balancer struct {
	nodes []string
	next  atomic.Uint64
}

func New(nodes []string) *Balancer {
	return &Balancer{
		nodes: nodes,
	}
}

func (b *Balancer) pick() string {
	idx := b.next.Add(1)

	return b.nodes[
		idx%uint64(len(b.nodes)),
	]
}

func pipe(
	dst net.Conn,
	src net.Conn,
) {
	defer dst.Close()
	defer src.Close()

	io.Copy(dst, src)
}

func (b *Balancer) handle(
	client net.Conn,
) {
	target := b.pick()

	backend, err := net.Dial(
		"tcp",
		target,
	)

	if err != nil {
		client.Close()
		return
	}

	go pipe(
		backend,
		client,
	)

	go pipe(
		client,
		backend,
	)
}

func (b *Balancer) Listen(
	addr string,
) error {

	ln, err := net.Listen(
		"tcp",
		addr,
	)

	if err != nil {
		return err
	}

	log.Println(
		"🔥 load balancer listening on",
		addr,
	)

	for {
		conn, err := ln.Accept()

		if err != nil {
			continue
		}

		go b.handle(conn)
	}
}
