################################################################################
# This Makefile builds/generates/downloads any dependencies/tools in the project.
# If a machine with Go and Docker cannot run any of these commands, that's a bug
# in the Makefile and should be fixed ASAP.
#
# Except Windows.  GLHF, PRs accepted <3

################################################################################
# Common commands
#
# These are the commands a dev should generally have to be aware of.  We should
# keep this list reasonably limited.

# Runs all unit tests
#
# This is the default command and is a good sanity check.
test: ./vendor
	go run ./vendor/github.com/onsi/ginkgo/ginkgo -r --progress --succinct=false ./internal

# Views test coverage as a pretty HTML document
coverage: coverage.out ./vendor
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out

# Runs Skaffold
#
# Starts Skaffold in dev mode, which will deploy ./deploy/k8s-local-dev.yaml
# as a dev stack and watch all files for changes.  Whenever a change is made,
# it will rebuild our service's Docker image and redeploy it to the dev stack.
# This lets us run everything with external dependencies in k8s.
skaffold-dev: ./bin/skaffold ./vendor
	./bin/skaffold dev --port-forward

# Builds the Docker image for our service
docker: ./vendor
	docker build -t "evertras/sample-go-app" .

# Cleans generated files/directories
#
# Generally shouldn't need to use this, but can be helpful for sanity.
clean:
	rm -rf bin
	rm -rf vendor

# Runs go fmt on all of our code
#
# This will ONLY output to stdout if anything changed, which means CI can use it
# as a check; if 'make fmt' outputs anything, then someone didn't run this command
# before checking into git and the CI check should fail.
fmt: $(GO_FILES)
	@go fmt ./internal/...
	@go fmt ./cmd/...

# Runs the Swagger editor at http://localhost:8005
#
# This is run in terminal as opposed to detached so we don't get zombies
swagger:
	@echo "Running swagger editor at http://localhost:8005"
	@docker run --rm -it -v ${PWD}/specs/:/specs -p 8005:8080 -e SWAGGER_FILE=/specs/api.yaml swaggerapi/swagger-editor

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

SKAFFOLD_VERSION = v1.13.2
./bin/skaffold: ./bin
ifeq ($(UNAME), Darwin)
	@echo Downloading Skaffold for Mac
	curl -L -o ./bin/skaffold https://storage.googleapis.com/skaffold/releases/$(SKAFFOLD_VERSION)/skaffold-darwin-amd64
else ifeq ($(UNAME), Linux)
	@echo Downloading Skaffold for Linux
	curl -L -o ./bin/skaffold https://storage.googleapis.com/skaffold/releases/$(SKAFFOLD_VERSION)/skaffold-linux-amd64
else
	# Windows gets rekt
	@echo Download the proper skaffold binary for your OS to ./bin/skaffold to continue
	@exit 1
endif
	chmod +x ./bin/skaffold
	./bin/skaffold version

./vendor: go.mod go.sum
	rm -rf vendor
	go mod vendor

