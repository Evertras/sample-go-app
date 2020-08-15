################################################################################
# This Makefile builds/generates/downloads any dependencies/tools in the project.
# If a machine with Go and Docker cannot run any of these commands, that's a bug
# in the Makefile and should be fixed ASAP.
#
# Except Windows.  GLHF, PRs accepted <3

################################################################################
# Common commands
test: ./vendor
	go run ./vendor/github.com/onsi/ginkgo/ginkgo -r ./internal

coverage: coverage.out ./vendor
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out

run: bin/server ./vendor
	@./bin/server

skaffold-dev: ./bin/skaffold ./vendor
	./bin/skaffold dev --port-forward

docker: ./vendor
	docker build -t "evertras/sample-go-app" .

clean:
	rm -rf bin
	rm -rf vendor

fmt: $(GO_FILES)
	@go fmt ./internal/...
	@go fmt ./cmd/...

################################################################################
# Dependencies
GO_FILES = $(shell find . -type f -name '*.go')
UNAME = $(shell uname)

coverage.out: $(GO_FILES)
	go test -coverprofile=coverage.out ./...

# Bin contains any built binaries or tools
./bin:
	@mkdir bin

./bin/server: $(GO_FILES) ./bin
	CGO_ENABLED=0 go build -o bin/server ./cmd/server/main.go

./bin/skaffold: ./bin
ifeq ($(UNAME), Darwin)
	@echo Downloading Skaffold for Mac
	curl -L -o ./bin/skaffold https://storage.googleapis.com/skaffold/releases/latest/skaffold-darwin-amd64
else ifeq ($(UNAME), Linux)
	@echo Downloading Skaffold for Linux
	curl -L -o ./bin/skaffold https://storage.googleapis.com/skaffold/releases/latest/skaffold-linux-amd64
else
	# Windows gets rekt
	@echo Download the proper skaffold binary for your OS to ./bin/skaffold to continue
	@exit 1
endif
	chmod +x ./bin/skaffold

./vendor: go.mod go.sum
	rm -rf vendor
	go mod vendor

