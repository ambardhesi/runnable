.PHONY: binaries clean test vendor

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
ROOT=$(shell pwd)
COMMANDS=$(shell go list ./... | grep cmd/)

binaries:
	$(GOBUILD) -o "." $(COMMANDS) 	

clean:
	rm -rf vendor
	$(GOCLEAN)

test: vendor
	$(GOTEST) -v -race ./...

vendor:
	$(GOCMD) mod vendor
