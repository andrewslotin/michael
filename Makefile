NAME := slack-deploy-command
SOURCES := $(shell find . -name '*.go')

PACKAGES := $$(go list ./... | grep -v /vendor/ | grep -v /cmd/)
test:
	go test $(PACKAGES)

build: slack-deploy-command

slack-deploy-command: $(SOURCES)
	GOGC=off GOOS=linux GOARCH=amd64 go build -o $(NAME)

compress:
ifeq (, $(shell command -v upx 2>/dev/null))
	@echo "upx not found in PATH, proceeding with unpacked binary"
else
	upx -q $(NAME)
endif

all: test build
.DEFAULT_GOAL := all

container: all compress
	docker build -t $(NAME) .

clean:
	go clean

.PHONY: 
