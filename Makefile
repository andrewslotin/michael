container: build compress
	docker build -t slack-deploy-command .

build: slack-deploy-command

compress:
ifeq (, $(shell command -v upx 2>/dev/null))
	@echo "upx not found in PATH, proceeding with unpacked binary"
else
	upx -q slack-deploy-command
endif

slack-deploy-command:
	GOOS=linux GOARCH=amd64 go build

clean:
	go clean
