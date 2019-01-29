.PHONY: build cluster test

default: build

build:
	go build -o bin/swarm    github.com/zdbrig/sagecity/cmd/swarm
	go build -o bin/geth     github.com/zdbrig/sagecity/cmd/geth
	go build -o bin/bootnode github.com/zdbrig/sagecity/cmd/bootnode

cluster: build
	scripts/boot-cluster.sh

test:
	go test -v github.com/zdbrig/sagecity/swarm/...
