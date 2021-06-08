OS = $(shell uname -s | tr '[:upper:]' '[:lower:]')
export CGO_ENABLED ?= 0

# Debug symbols
ifeq (${DEBUG},)
else
GOARGS=-gcflags="all=-N -l"
endif
LDFLAGS="-w -s"

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
	docker push deviceplane/controller:${CONTROLLER_VERSION}

controller-with-db:
	./scripts/build-controller-with-db

push-controller-with-db: controller-with-db
	docker push deviceplane/deviceplane:${CONTROLLER_WITH_DB_VERSION}

agent-binaries:
	./scripts/build-agent-binaries

cli:
	./scripts/build-cli

push-cli: cli
	docker push deviceplane/cli:${CLI_VERSION}

cli-binaries:
	./scripts/build-cli-binaries

.PHONY: build
build: GOARGS += -mod=vendor -ldflags $(LDFLAGS)
build: ## Build binary
	go build $(GOARGS) -o bin/controller cmd/controller/main.go
