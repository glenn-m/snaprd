BINARY_NAME=snaprd

all: fmt vet test build

build:
	@CGO_ENABLED=0 go build -o "bin/$(BINARY_NAME)" cmd/main.go

run: build
	@./bin/${BINARY_NAME}

fmt:
	@go fmt ./...

vet:
	@go vet ./...

test:
	@go test -cover -v ./...
