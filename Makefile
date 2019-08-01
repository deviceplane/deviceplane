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
	./scripts/build-cli-ci

push-cli-ci: cli-ci
	docker push deviceplane/cli-ci:${CLI_VERSION}

cli-binaries:
	./scripts/build-cli-binaries

upload-cli-binary-redirects:
	./scripts/upload-cli-binary-redirects
