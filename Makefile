BINARY=star-wars

BUILD_DIR=build/${BINARY}
LDFLAGS = -ldflags "-X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH}"

VERSION=$(shell cat VERSION)
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

build:
	go build -trimpath ${LDFLAGS} -o ${BUILD_DIR}/${BINARY}

clean:
	rm -rf build

fmt:
	gofmt -w .

run:
	go run cmd/star-wars/star-wars.go

test:

.PHONY: build fmt run test 