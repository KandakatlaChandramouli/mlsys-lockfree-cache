
package transport

import (
	"bufio"
	"context"
	"net"
	"sync/atomic"
	"time"

	"fluxruntime/internal/core"
	"fluxruntime/internal/metrics"
	"fluxruntime/internal/protocol"
	"fluxruntime/internal/tokenizer"
)

type Server struct {
	addr string
	pool *core.ShardedPool
	tok  *tokenizer.Tokenizer
	m    *metrics.Metrics
	reqs atomic.Uint64
}

func NewServer(
	addr string,
	pool *core.ShardedPool,
) *Server {

	return &Server{
		addr: addr,
		pool: pool,
		tok:  tokenizer.New(),
		m:    metrics.New(),
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

		start := time.Now()

		text, err := protocol.ReadString(
			reader,
		)

		if err != nil {
			return
		}

		tokens := s.tok.Encode(
			text,
		)

		embedding, err := s.pool.Embed(
			tokens,
		)

		if err != nil {

			s.m.RecordFailure()

			return
		}

		err = protocol.WriteEmbedding(
			writer,
			embedding,
		)

		if err != nil {
			return
		}

		err = writer.Flush()

		if err != nil {
			return
		}

		s.m.RecordRequest(
			time.Since(start),
		)

		s.reqs.Add(1)
	}
}
