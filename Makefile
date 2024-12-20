# Makefile for common Go operations

.PHONY: fmt lint test build clean run

fmt:
	go fmt ./...

lint: fmt
	go vet -v ./...

test:
	go test -v -race ./...

build: clean
	mkdir -p bin
	go build -v -o bin/eth-tx-parser cmd/eth-tx-parser/main.go

run: build
	./bin/eth-tx-parser

clean:
	rm -rf bin/*