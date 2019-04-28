.PHONY: fmt test lint cilint

all: fmt test lint cilint

fmt:
	go fmt $(go list ./... | grep -v /vendor/)

test:
	go test -cover $(go list ./... | grep -v /vendor/ | grep -v /_examples/)

lint:
	golint ./... | sed '/^vendor/ d' | sed '/^_examples/ d'

cilint:
	golangci-lint run ./... | sed '/^vendor/ d' | sed '/^_examples/ d'

generate:
	go generate ./...