.PHONY: build test fmt vet lint clean

build:
	go build ./...

test:
	go test ./...

fmt:
	gofmt -w .

vet:
	go vet ./...
check: fmt vet test build