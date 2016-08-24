PROJECT := slack-deploy-command
SOURCES := $(shell find . -name '*.go')

VERSION := 1.1.1
GIT_REVISION := $$(git rev-parse HEAD | cut -c -6)
GOVERSION := $(shell go version)
BUILDDATE := $(shell date -u +"%B %d, %Y")
BUILDER := $(shell echo "`git config user.name` <`git config user.email`>")
LDFLAGS := -X 'main.version=$(VERSION)' \
           -X 'main.buildDate=$(BUILDDATE)' \
           -X 'main.builder=$(BUILDER)' \
           -X 'main.buildRev=$(GIT_REVISION)' \
           -X 'main.buildGoVersion=$(GOVERSION)'

OS ?= $(shell go env GOOS)
ARCH ?= $(shell go env GOARCH)

PACKAGES := $$(go list ./... | grep -v /vendor/ | grep -v /cmd/)
test:
	go test $(PACKAGES)

build: slack-deploy-command

slack-deploy-command: $(SOURCES)
	GOGC=off GOOS=$(OS) GOARCH=$(ARCH) go build -ldflags "$(LDFLAGS)" -o $(PROJECT)

compress:
ifeq (, $(shell command -v upx 2>/dev/null))
	@echo "upx not found in PATH, proceeding with unpacked binary"
else
	upx -q $(PROJECT)
endif

all: test build
.DEFAULT_GOAL := all

container: OS = linux
container: ARCH = amd64
container: all
	docker build -t $(PROJECT) .

clean:
	go clean

.PHONY: test build slack-deploy-command compress all container clean
