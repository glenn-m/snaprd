BINARY_NAME=snaprd

all: fmt vet test build

build:
	@CGO_ENABLED=0 go build -o "bin/$(BINARY_NAME)" cmd/${BINARY_NAME}/main.go

run: build
	@./bin/${BINARY_NAME}

fmt:
	@go fmt ./...

vet:
	@go vet ./...

lint:
	@golangci-lint run

test:
	@go test -cover -v ./...
