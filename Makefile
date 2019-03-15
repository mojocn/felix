LDFLAGS := "-s -w -X main.buildTime=$(shell date -u '+%Y-%m-%dT%I:%M:%S%p') -X main.gitHash=$(shell git rev-parse HEAD)"
GO ?= go
GOFMT ?= gofmt "-s"
PACKAGES ?= $(shell $(GO) list ./... | grep -v /vendor/)
VETPACKAGES ?= $(shell $(GO) list ./... | grep -v /vendor/ | grep -v /examples/)
GOFILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")


run: install
	./build/felix -V
install:
	go install -ldflags $(LDFLAGS)
vuejs:
	felix ginbin -s dist -p felixbin

build:vuejs
	go build -race -ldflags $(LDFLAGS)  -o build/felix *.go

release:vuejs
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64  CXX_FOR_TARGET=i686-w64-mingw32-g++ CC_FOR_TARGET=i686-w64-mingw32-gcc go build -ldflags $(LDFLAGS) -o _release/felix-amd64-win.exe *.go
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags $(LDFLAGS) -o _release/felix-amd64-linux *.go
	CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -ldflags $(LDFLAGS) -o _release/felix-amd64-linux-arm *.go
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -ldflags $(LDFLAGS) -o _release/felix-amd64-darwin *.go


.PHONY: release

