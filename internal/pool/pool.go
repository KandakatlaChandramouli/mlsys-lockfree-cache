package pool

import (
	"net"
	"sync"
)

type Pool struct {
	addr string

	conns chan net.Conn
}

func New(
	addr string,
	size int,
) *Pool {

	p := &Pool{
		addr: addr,

		conns: make(
			chan net.Conn,
			size,
		),
	}

	for i := 0; i < size; i++ {

		conn, err := net.Dial(
			"tcp",
			addr,
		)

		if err != nil {
			continue
		}

		p.conns <- conn
	}

	return p
}

func (p *Pool) Get() net.Conn {

	select {

	case conn := <-p.conns:
		return conn

	default:

		conn, err := net.Dial(
			"tcp",
			p.addr,
		)

		if err != nil {
			return nil
		}

		return conn
	}
}

func (p *Pool) Put(
	conn net.Conn,
) {

	if conn == nil {
		return
	}

	select {

	case p.conns <- conn:

	default:
		conn.Close()
	}
}

var (
	pools sync.Map
)

func GetPool(
	addr string,
) *Pool {

	v, ok := pools.Load(addr)

	if ok {
		return v.(*Pool)
	}

	p := New(
		addr,
		128,
	)

	actual, _ := pools.LoadOrStore(
		addr,
		p,
	)

	return actual.(*Pool)
}
