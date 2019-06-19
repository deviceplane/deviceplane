version= $(shell git describe --tags --always --dirty="-dev")

db-reset:
	docker-compose down
	docker-compose build
	docker-compose up -d
	sleep 20
	./scripts/seed

controller:
	docker build -t deviceplane/deviceplane:${version} -f Dockerfile.controller --build-arg version=${version} .

push-controller: controller
	docker push deviceplane/deviceplane:${version}

agent:
	VERSION=${version} ./scripts/build-agent

push-agent: agent
	docker manifest push deviceplane/agent:${version}

cli-ci:
	docker build -t deviceplane/cli-ci:${version} -f Dockerfile.cli-ci --build-arg version=${version} .

push-cli-ci: cli-ci
	docker push deviceplane/cli-ci:${version}

cli-binaries:
	VERSION=${version} ./scripts/build-cli-binaries

upload-cli-binary-redirects:
	VERSION=${version} ./scripts/upload-cli-binary-redirects
