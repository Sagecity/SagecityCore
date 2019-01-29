.PHONY: build cluster test

default: build

build:
	go build -o bin/swarm    github.com/SagecityCore/sagecity/cmd/swarm
	go build -o bin/geth     github.com/SagecityCore/sagecity/cmd/geth
	go build -o bin/bootnode github.com/SagecityCore/sagecity/cmd/bootnode

cluster: build
	scripts/boot-cluster.sh

test:
	go test -v github.com/SagecityCore/sagecity/swarm/...
