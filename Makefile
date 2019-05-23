version= $(shell git describe --tags --always --dirty="-dev")

db-reset:
	docker-compose down
	docker-compose build
	docker-compose up -d
	sleep 20
	./scripts/seed

cli:
	go build -o ./dist/cli ./cmd/cli

controller:
	docker build -t deviceplane/deviceplane:${version} -f Dockerfile.controller .

push-controller: controller
	docker push deviceplane/deviceplane:${version}

agent:
	docker build -t deviceplane/agent:${version} -f Dockerfile.agent .

push-agent: agent
	docker push deviceplane/agent:${version}

cli-ci:
	docker build -t deviceplane/cli-ci:${version} -f Dockerfile.cli-ci .

push-cli-ci: cli-ci
	docker push deviceplane/cli-ci:${version}
