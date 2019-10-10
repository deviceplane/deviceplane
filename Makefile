db-reset: state-reset
	docker-compose down
	docker-compose build
	docker-compose up -d
	$$WAIT_FOR_DB
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

clone-db-locally: dump-remote-db load-local-db-from-dump

dump-remote-db:
	mkdir -p localdump
	echo $$DB_PASS
	if [[ -z "$$DB_PASS" ]]; then \
		echo "DB_PASS is not set"; \
		exit 1; \
	fi
	ssh ubuntu@54.200.126.157 "mysqldump -h deviceplane.coed7waagekn.us-west-2.rds.amazonaws.com -u deviceplane --password=$$DB_PASS -P 3306 --databases deviceplane" > localdump/db.sql

load-local-db-from-dump: state-reset
	docker-compose down
	docker-compose build
	docker-compose up -d
	$$WAIT_FOR_DB
	mysql -h 127.0.0.1 -u user --password=pass -P 3306 -D deviceplane < localdump/db.sql
