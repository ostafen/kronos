BIN_FOLDER = ./bin
EXEC_NAME = kronos

BUILD_TIME := $(shell date)
COMMIT := $(shell git rev-parse HEAD)
VERSION ?= latest

IMG_NAME ?= ghcr.io/ostafen/kronos
IMG_TAG ?= latest

build: vendor
	@mkdir -p $(BIN_FOLDER)
	go build -mod vendor -ldflags '-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X "main.buildTime=$(BUILD_TIME)"' -o $(BIN_FOLDER)/$(EXEC_NAME) cmd/main.go

vendor:
	go mod vendor

docker-build:
	docker build -t ${IMG_NAME}:${IMG_TAG} --build-arg VERSION=$(VERSION) --build-arg COMMIT=$(COMMIT) --build-arg BUILD_TIME="$(BUILD_TIME)" . 

docker-push:
	docker push ${IMG_NAME}:${IMG_TAG}
