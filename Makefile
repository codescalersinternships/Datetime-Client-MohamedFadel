GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=datetime-client
DOCKER_IMAGE_NAME=datetime-client-image

SERVER_URL ?= http://localhost:8000
CONTENT_TYPE ?= application/json

all: deps fmt lint build test docker-build docker-run

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME) -url=$(SERVER_URL) -type=$(CONTENT_TYPE)

deps:
	$(GOMOD) download

fmt:
	gofmt -s -w .

lint:
	golangci-lint run

docker-build:
	docker build -t $(DOCKER_IMAGE_NAME) .

docker-run:
	docker run --network host $(DOCKER_IMAGE_NAME) -url=$(SERVER_URL) -type=$(CONTENT_TYPE)

.PHONY: all build test clean run deps fmt lint docker-build docker-run