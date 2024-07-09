BINARY=star-wars

VERSION=$(shell cat VERSION)
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
BUILD_TS=$(shell date '+%Y-%m-%dT%H:%M:%S')

BUILD_DIR=build/bin
LDFLAGS= -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.branch=${BRANCH} -X main.buildTimestamp=${BUILD_TS}"

IMAGE_NAME=mattslater.io/${BINARY}:${VERSION}

build:
	GOOS=darwin GOARCH=arm64 go build -trimpath ${LDFLAGS} -o ${BUILD_DIR}/darwin/${BINARY}-darwin cmd/${BINARY}/${BINARY}.go
	GOOS=linux GOARH=amd64 go build -trimpath ${LDFLAGS} -o ${BUILD_DIR}/linux/${BINARY}-linux cmd/${BINARY}/${BINARY}.go
	GOOS=windows GOARCH=amd64 go build -trimpath ${LDFLAGS} -o ${BUILD_DIR}/windows/${BINARY}-windows.exe cmd/${BINARY}/${BINARY}.go

build-container:
	cp build/bin/linux/star-wars-linux build/package/star-wars-linux
	docker image build -t ${IMAGE_NAME} -f build/package/Dockerfile build/package
	rm -f build/package/star-wars-linux

clean:
	rm -rf build/${BINARY}

fmt:
	goimports -w .
	gofmt -s -w .

run:
	go run cmd/${BINARY}/${BINARY}.go

.PHONY: build build-container fmt run 