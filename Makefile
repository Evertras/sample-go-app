################################################################################
# Common commands
test:
	ginkgo ./...

run: bin/server
	@./bin/server

skaffold-dev: ./bin/skaffold
	./bin/skaffold dev --port-forward

docker:
	docker build -t "evertras/sample-go-app" .

clean:
	rm -rf bin

################################################################################
# Dependencies
GO_FILES = $(shell find . -type f -name '*.go')
UNAME = $(shell uname)

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

