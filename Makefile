APP_NAME := mime-api
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
DOCKER_IMAGE := $(APP_NAME):$(VERSION)

.PHONY: build clean test lint fmt run docker-build docker-run

build:
	@mkdir -p dist
	go build -ldflags="-w -s" -o dist/$(APP_NAME) .

clean:
	rm -rf dist
	go clean

test:
	go test -v ./...

lint:
	golangci-lint run ./...

fmt:
	go fmt ./...

run:
	go run .

docker-build:
	docker build -t $(DOCKER_IMAGE) .
	docker tag $(DOCKER_IMAGE) $(APP_NAME):latest

docker-run:
	docker run --rm -p 8080:8080 $(APP_NAME):latest