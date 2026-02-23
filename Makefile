.PHONY: all build format test lint

all: build

build:
	go build -ldflags="-s -w"

format:
	gofmt -w .

test:
	go test ./...

lint:
	go vet ./...
	golangci-lint run