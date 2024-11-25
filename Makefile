.PHONY: all build test lint clean run docker-build docker-run

# Variables
BINARY_NAME=api-server
DOCKER_IMAGE=yourapp

all: lint test build

build:
	go build -o ${BINARY_NAME} ./cmd/api

test:
	go test -v -race ./...

lint:
	golangci-lint run

clean:
	go clean
	rm -f ${BINARY_NAME}

run:
	go run ./cmd/api

docker-build:
	docker build -t ${DOCKER_IMAGE} .

docker-run:
	docker run -p 8080:8080 ${DOCKER_IMAGE}