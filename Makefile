GO_FILES = $(shell find . -type f -name '*.go')

bin/server: $(GO_FILES)
	CGO_ENABLED=0 go build -o bin/server ./cmd/server/main.go

test:
	ginkgo ./...

