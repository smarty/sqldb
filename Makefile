#!/usr/bin/make -f

test: fmt
	go test -timeout=1s -race -covermode=atomic ./...
	go test github.com/smarty/sqldb/integration

fmt:
	go mod tidy && go fmt ./...

compile:
	go build ./...

build: test compile

.PHONY: test fmt compile build
