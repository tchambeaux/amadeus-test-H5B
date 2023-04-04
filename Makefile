.PHONY: test
test:
	staticcheck ./...
	go vet ./...
	go test ./... -cover

.PHONY: build
build:
	go build -ldflags="-s -w" -o bin/test_H5B

.PHONY: clean
clean:
	rm -rf bin/

.PHONY: all
all: clean test build