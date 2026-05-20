package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	routerv1 "fluxruntime/proto/v1"

	"fluxruntime/internal/router"
)

func main() {

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
