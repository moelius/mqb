.PHONY: fmt lint cilint test

fmt:
	go fmt $(go list ./... | grep -v /vendor/)

test:
	go test -cover $(go list ./... | grep -v /vendor/ | grep -v /_examples/)

lint:
	golint ./... | sed '/^vendor/ d' | sed '/^_examples/ d'

cilint:
	golangci-lint run ./... | sed '/^vendor/ d' | sed '/^_examples/ d'