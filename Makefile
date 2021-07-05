.PHONY: clean test vendor

clean:
	rm -rf vendor
	go clean

test: vendor
	go test -v -race ./...

vendor:
	go mod vendor
