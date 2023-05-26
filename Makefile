.PHONY: start build

NOW = $(shell date -u '+%Y%m%d%I%M%S')

RELEASE_VERSION = v1.0

APP 		= nightwatcher
SERVER_BIN  	= ./cmd/bin/${APP}
RELEASE_ROOT 	= release
RELEASE_SERVER 	= release/${APP}
GIT_COUNT 	= $(shell git rev-list --all --count)
GIT_HASH        = $(shell git rev-parse --short HEAD)
RELEASE_TAG     = $(RELEASE_VERSION).$(GIT_COUNT).$(GIT_HASH)

all: start

build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s -X main.VERSION=$(RELEASE_TAG)" -o $(SERVER_BIN) -a -installsuffix cgo ./

swagger:
	@swag init --parseDependency --generalInfo ./main.go --output ./docs

clean:
	rm -rf data release $(SERVER_BIN) internal/app/test/data cmd/bin
