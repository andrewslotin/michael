container: slack-deploy-command
	docker build -t slack-deploy-command .

build: slack-deploy-command

slack-deploy-command:
	GOOS=linux GOARCH=amd64 go build

clean:
	go clean
