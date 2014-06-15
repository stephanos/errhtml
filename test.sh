#!/bin/bash

export GOPATH="/Users/stephan/Dev/workspace/golang"

gofmt -w .
go vet .
#$GOPATH/bin/golint .
go test ./...