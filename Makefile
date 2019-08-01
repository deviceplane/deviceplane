version= $(shell git describe --tags --always --dirty="-dev")

db-reset:
	docker-compose down
	docker-compose build
	docker-compose up -d
	sleep 20
	./scripts/seed

controller:
	./scripts/build-controller

push-controller: controller
	docker push deviceplane/deviceplane:${CONTROLLER_VERSION}

agent:
	./scripts/build-agent

push-agent: agent
	docker manifest push deviceplane/agent:${AGENT_VERSION}

cli-ci:
	docker build -t deviceplane/cli-ci:${version} -f dockerfiles/cli-ci/Dockerfile --build-arg version=${version} .

push-cli-ci: cli-ci
	docker push deviceplane/cli-ci:${version}

cli-binaries:
	VERSION=${version} ./scripts/build-cli-binaries

upload-cli-binary-redirects:
	VERSION=${version} ./scripts/upload-cli-binary-redirects
