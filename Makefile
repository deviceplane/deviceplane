test:
	go test -v ./... -mod vendor

db-reset: state-reset
	docker-compose down
	docker-compose up -d
	sleep 30
	./scripts/seed

state-reset:
	rm -rf ./cmd/controller/state

controller:
	./scripts/build-controller

push-controller: controller
	docker push deviceplane/deviceplane:${CONTROLLER_VERSION}

controller-with-db:
	./scripts/build-controller-with-db

push-controller-with-db: controller-with-db
	docker push deviceplane/deviceplane:${CONTROLLER_WITH_DB_VERSION}-with-db

agent-binaries:
	./scripts/build-agent-binaries

cli:
	./scripts/build-cli

push-cli: cli
	docker push deviceplane/cli:${CLI_VERSION}

cli-binaries:
	./scripts/build-cli-binaries
