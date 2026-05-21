package main

import (
	"log"

	"fluxruntime/internal/loadbalancer"
)

func main() {

	nodes := []string{
		"localhost:7001",
		"localhost:7002",
		"localhost:7003",
	}

	lb := loadbalancer.New(
		nodes,
	)

	log.Fatal(
		lb.Listen(":8000"),
	)
}
