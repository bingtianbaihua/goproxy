all:
	GOBIN="`pwd`/bin" go install -v ./cmd/...

install: all
	@echo

clean:
	rm -rf bin

gofmt:
	find . -name '*.go' | xargs -l1 go fmt
