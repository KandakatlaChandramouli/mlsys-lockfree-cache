package main

import (
	"log"
	"net"
	"net/http"

	_ "net/http/pprof"

	"google.golang.org/grpc"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"fluxruntime/internal/metrics"
	"fluxruntime/internal/router"
	routerv1 "fluxruntime/proto/v1"
)

func main() {

	metrics.Register()

	go func() {

		http.Handle("/metrics", promhttp.Handler())

		log.Println("📊 metrics server listening on :6060")
		log.Println("🔥 pprof enabled on :6060")

		log.Println(http.ListenAndServe(":6060", nil))

	}()

	rtr := router.New()

	server := grpc.NewServer()

	routerv1.RegisterRouterServiceServer(server, rtr)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	log.Println("🚀 gRPC runtime listening on :50051")

	if err := server.Serve(lis); err != nil {
		panic(err)
	}
}
