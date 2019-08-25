db-reset: state-reset
	docker-compose down
	docker-compose build
	docker-compose up -d
	sleep 20
	./scripts/seed

state-reset:
	rm -rf ./cmd/controller/state

controller:
	./scripts/build-controller

push-controller: controller
	docker push deviceplane/deviceplane:${CONTROLLER_VERSION}

agent:
	./scripts/build-agent

push-agent: agent
	docker manifest push deviceplane/agent:${AGENT_VERSION}

cli:
	./scripts/build-cli

push-cli: cli
	docker push deviceplane/cli:${CLI_VERSION}

cli-binaries:
	./scripts/build-cli-binaries

upload-cli-binary-redirects:
	./scripts/upload-cli-binary-redirects
