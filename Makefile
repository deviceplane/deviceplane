test:
	go test -v ./... -mod vendor

get-releases:
	@git log | grep "Release agent" | head -n 1
	@git log | grep "Release controller" | head -n 1
	@git log | grep "Release CLI" | head -n 1

db-reset: state-reset
	docker-compose down
	docker-compose build
	docker-compose up -d
	sleep 30
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

statik:
	npm run build --prefix ./ui
	statik -src=./ui/build -dest=./pkg
