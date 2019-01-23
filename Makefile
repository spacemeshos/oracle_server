BINARY := oracle_server
VERSION := 0.0.1
COMMIT = $(shell git rev-parse HEAD)
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
BIN_DIR = $(shell pwd)/build
CURR_DIR = $(shell pwd)
CURR_DIR_WIN = $(shell cd)
export GO111MODULE = on

# Setup the -ldflags option to pass vars defined here to app vars
LDFLAGS = -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.branch=${BRANCH}"

PKGS = $(shell go list ./...)

build:
	./pb/compile.sh
	go build ${LDFLAGS} -o $(CURR_DIR)/$(BINARY)
.PHONY: build

tidy:
	go mod tidy
.PHONY: tidy

test:
	go test -short -p 1 ./...
.PHONY: test

test-tidy:
	# We expect `go mod tidy` not to change anything, the test should fail otherwise
	make tidy
	git diff --exit-code
.PHONY: test-tidy

cover:
	@echo "mode: count" > cover-all.out
	@$(foreach pkg,$(PKGS),\
		go test -coverprofile=cover.out -covermode=count $(pkg);\
		tail -n +2 cover.out >> cover-all.out;)
	go tool cover -html=cover-all.out
.PHONY: cover
