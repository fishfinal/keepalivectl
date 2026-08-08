.PHONY: build clean test lint fmt install

BINARY_NAME=keepalivectl
VERSION=1.0.0
BUILD_TIME=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS=-ldflags "-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME}"

build:
	go build ${LDFLAGS} -o bin/${BINARY_NAME} cmd/keepalivectl/main.go

build-all:
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}_linux_amd64 cmd/keepalivectl/main.go
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}_darwin_amd64 cmd/keepalivectl/main.go
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o bin/${BINARY_NAME}_darwin_arm64 cmd/keepalivectl/main.go
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY_NAME}_windows_amd64.exe cmd/keepalivectl/main.go

test:
	go test -v -race ./...

lint:
	golangci-lint run

fmt:
	go fmt ./...

clean:
	rm -rf bin/
	go clean

install:
	go install ${LDFLAGS} ./cmd/keepalivectl

docker-build:
	docker build -t ${BINARY_NAME}:${VERSION} .

docker-run:
	docker run --rm ${BINARY_NAME}:${VERSION}

release: build-all
	gh release create v${VERSION} bin/${BINARY_NAME}_* --title "Release ${VERSION}"
