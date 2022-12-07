
dep: ; go mod tidy

.PHONY: run
run: dep ; go run *.go

.PHONY: all
all: build test lint

.PHONY: build
build: dep ; go build -o bin/exmo-fix-client ./...

.PHONY: test
test: dep ; go test ./...

.PHONY: lint
lint: ; golangci-lint run
