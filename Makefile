# Makefile for resty
# 
# Use `git tag` to tag version
#

BINARY=resty

HTTPTESTFILE=_dev/jsonplaceholder.http

VERSION=`git describe --tags`
BUILD=`date +%FT%T%z`
LDFLAGS=-ldflags "-w -s -X main.Version=${VERSION} -X main.Build=${BUILD}"

info:
	@echo "Targets:"
	@echo "  build     - build for current OS"
	@echo "  build-win - build for Windows AMD64"
	@echo "  test      - run tests"
	@echo "  run       - run with test file"
	@echo "  clean     - remove build artifacts"

build:
	go build ${LDFLAGS} -o ${BINARY}

build-win:
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY}.exe

test:
	go test -v ./...

run: clean build
	./resty ${HTTPTESTFILE}

clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY}; fi
	@if [ -f ${BINARY}.exe ] ; then rm ${BINARY}.exe ; fi

.PHONY: clean build

