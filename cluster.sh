#!/bin/bash

pkill -f rawtcp
pkill -f loadbalancer

sleep 1

PORT=7001 go run ./cmd/rawtcp &
PORT=7002 go run ./cmd/rawtcp &
PORT=7003 go run ./cmd/rawtcp &

sleep 2

go run ./cmd/loadbalancer
