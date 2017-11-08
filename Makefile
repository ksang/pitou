# Define parameters
BINARY=pitou
SHELL := /bin/bash
GOPACKAGES = $(shell go list ./... | grep -v vendor)
ROOTDIR = $(pwd)

.PHONY: build install test linux

default: build

build: pitou.go
	go build -v -o ./build/${BINARY} pitou.go

install:
	go install  ./...

test:
	go test -race -cover ${GOPACKAGES}

clean:
	rm -rf build

linux: pitou.go
	GOOS=linux GOARCH=amd64 go build -o ./build/linux/${BINARY} pitou.go
