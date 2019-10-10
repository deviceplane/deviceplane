#!make
-include .env
export

DB_COMMAND=./scripts/db-command.sh
WAIT_FOR_DB=./scripts/wait-for-db.sh

db-reset: state-reset
	docker-compose down
	docker-compose build
	docker-compose up -d
	$$WAIT_FOR_DB $$(env ENV=local $$DB_COMMAND)
	$$(env ENV=local $$DB_COMMAND) < pkg/controller/store/mysql/schema.sql
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

dump-remote-db:
	mkdir -p localdump
	if [[ -z "$$DB_PASS" ]]; then \
		echo "DB_PASS is not set"; \
		exit 1; \
	fi
	$$(env ENV=prod DUMP=true $$DB_COMMAND) > localdump/db.sql

load-local-db-from-dump: state-reset
	docker-compose down
	docker-compose build
	docker-compose up -d
	$$WAIT_FOR_DB $$(env ENV=local $$DB_COMMAND)
	$$(env ENV=local $$DB_COMMAND) < localdump/db.sql

apply-migration:
	if [[ -z "$$MIGRATION_FILE" ]]; then \
		echo "MIGRATION_FILE is not set"; \
		exit 1; \
	fi
	if [[ -z "$$ENV" ]]; then \
		echo "ENV is not set"; \
		exit 1; \
	fi
	$$($$DB_COMMAND) < $$MIGRATION_FILE

try-last-migration-locally: load-local-db-from-dump
	lastFile=$$(ls ./migrations | tail -n 1); env ENV=local MIGRATION_FILE=./migrations/$$lastFile make apply-migration