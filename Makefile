BINARY=star-wars

VERSION=$(shell cat VERSION)
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
BUILD_TS=$(shell date '+%Y-%m-%dT%H:%M:%S')

BUILD_DIR=build/${BINARY}
LDFLAGS = -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.branch=${BRANCH} -X main.buildTimestamp=${BUILD_TS}"


build:
	go build -trimpath ${LDFLAGS} -o ${BUILD_DIR}/${BINARY} cmd/${BINARY}/${BINARY}.go

clean:
	rm -rf build

fmt:
	gofmt -s -w .

run:
	go run cmd/star-wars/star-wars.go

test:
	go test

.PHONY: build fmt run test 