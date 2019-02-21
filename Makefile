all: fmt build

build:
	go build -o ./bin/goproxy ./cmd

clean:
	rm -rf bin

fmt:
	go fmt ./...
