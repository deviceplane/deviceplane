version= $(shell git describe --tags --always --dirty="-dev")

cli:
	go build -o ./dist/cli ./cmd/cli

controller:
	GOOS=linux go build -o ./dist/controller ./cmd/controller

controller-image: controller
	docker build -t deviceplane/deviceplane:${version} -f Dockerfile.controller .

agent:
	GOOS=linux GOARCH=arm go build -o ./dist/agent ./cmd/agent

agent-image: agent
	docker build -t deviceplane/agent:${version} -f Dockerfile.agent .
