#!/bin/bash

apt update -y
apt install -y protobuf-compiler tree

wget https://go.dev/dl/go1.24.3.linux-amd64.tar.gz

rm -rf /usr/local/go
tar -C /usr/local -xzf go1.24.3.linux-amd64.tar.gz

export PATH=$PATH:/usr/local/go/bin

/usr/local/go/bin/go version

/usr/local/go/bin/go mod tidy
/usr/local/go/bin/go build ./...
