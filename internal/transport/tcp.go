
package transport

import (
	"bufio"
	"context"
	"encoding/binary"
	"io"
	"net"
	"sync/atomic"

	"fluxruntime/internal/core"
)

type Server struct {
	addr string
	pool *core.ShardedPool
	reqs atomic.Uint64
}

func NewServer(
	addr string,
	pool *core.ShardedPool,
) *Server {

	return &Server{
		addr: addr,
		pool: pool,
	}
}

func (s *Server) Listen(
	ctx context.Context,
) error {

	ln, err := net.Listen(
		"tcp",
		s.addr,
	)

	if err != nil {
		return err
	}

	go func() {

		<-ctx.Done()

		ln.Close()
	}()

	for {

		conn, err := ln.Accept()

		if err != nil {

			select {

			case <-ctx.Done():
				return nil

			default:
			}

			continue
		}

		go s.handle(
			ctx,
			conn,
		)
	}
}

func (s *Server) handle(
	ctx context.Context,
	conn net.Conn,
) {

	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {

		select {

		case <-ctx.Done():
			return

		default:
		}

		var size uint32

		err := binary.Read(
			reader,
			binary.LittleEndian,
			&size,
		)

		if err != nil {

			if err == io.EOF {
				return
			}

			return
		}

		buf := make(
			[]byte,
			size,
		)

		_, err = io.ReadFull(
			reader,
			buf,
		)

		if err != nil {
			return
		}

		text := string(buf)

		tokens := tokenize(
			text,
		)

		embedding, err := s.pool.Embed(
			tokens,
		)

		if err != nil {
			return
		}

		outSize := uint32(
			len(embedding),
		)

		err = binary.Write(
			writer,
			binary.LittleEndian,
			outSize,
		)

		if err != nil {
			return
		}

		for _, v := range embedding {

			err = binary.Write(
				writer,
				binary.LittleEndian,
				v,
			)

			if err != nil {
				return
			}
		}

		err = writer.Flush()

		if err != nil {
			return
		}

		s.reqs.Add(1)
	}
}

func tokenize(
	text string,
) []int32 {

	out := []int32{
		101,
	}

	for _, b := range []byte(text) {

		out = append(
			out,
			int32(b)%30000,
		)
	}

	out = append(
		out,
		102,
	)

	return out
}
