GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=datetime-client
DOCKER_IMAGE_NAME=datetime-client-image

all: deps fmt lint build test docker-build docker-run

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

test:
	$(GOTEST) -v ./client...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

deps:
	$(GOMOD) download

fmt:
	gofmt -s -w .

lint:
	golangci-lint run

docker-build:
	docker build -t $(DOCKER_IMAGE_NAME) .

docker-run:
	docker run --network host -e SERVER_URL=$(SERVER_URL) -e SERVER_PORT=$(SERVER_PORT) -e SERVER_PORT_GIN=$(SERVER_PORT_GIN) $(DOCKER_IMAGE_NAME)

.PHONY: all build test clean run deps fmt lint docker-build docker-run