version= $(shell git describe --tags --always --dirty="-dev")
agent_version= 1.0.0

db-reset:
	docker-compose down
	docker-compose build
	docker-compose up -d
	sleep 20
	./scripts/seed

controller:
	docker build -t deviceplane/deviceplane:${version} -f dockerfiles/controller/Dockerfile --build-arg version=${version} .

push-controller: controller
	docker push deviceplane/deviceplane:${version}

agent:
	VERSION=${agent_version} ./scripts/build-agent

push-agent: agent
	docker manifest push deviceplane/agent:${agent_version}

cli-ci:
	docker build -t deviceplane/cli-ci:${version} -f dockerfiles/cli-ci/Dockerfile --build-arg version=${version} .

push-cli-ci: cli-ci
	docker push deviceplane/cli-ci:${version}

cli-binaries:
	VERSION=${version} ./scripts/build-cli-binaries

upload-cli-binary-redirects:
	VERSION=${version} ./scripts/upload-cli-binary-redirects
