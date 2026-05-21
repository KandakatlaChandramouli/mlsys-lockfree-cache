
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"fluxruntime/internal/core"
	"fluxruntime/internal/transport"
)

func main() {

	ctx, cancel := context.WithCancel(
		context.Background(),
	)

	defer cancel()

	modelPath := "models/model.onnx"

	pool, err := core.NewShardedPool(
		ctx,
		modelPath,
	)

	if err != nil {
		log.Fatal(err)
	}

	defer pool.Shutdown()

	server := transport.NewServer(
		":9000",
		pool,
	)

	go func() {

		sig := make(
			chan os.Signal,
			1,
		)

		signal.Notify(
			sig,
			syscall.SIGINT,
			syscall.SIGTERM,
		)

		<-sig

		cancel()
	}()

	log.Println(
		"runtime listening on :9000",
	)

	err = server.Listen(ctx)

	if err != nil {
		log.Fatal(err)
	}
}
