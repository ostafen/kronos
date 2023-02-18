BIN_FOLDER=./bin
EXEC_NAME=kronos

BUILD_TIME := $(shell date)
COMMIT := $(shell git rev-parse HEAD)

build: vendor
	@mkdir -p $(BIN_FOLDER)
	go build -mod vendor -ldflags '-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X "main.buildTime=$(BUILD_TIME)"' -o $(BIN_FOLDER)/$(EXEC_NAME) cmd/main.go

vendor:
	go mod vendor
