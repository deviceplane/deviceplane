# Local Development

Run the UI with `npm start`.

Run `make db-reset` to setup the database.

Run the controller with `go run cmd/controller/main.go`. By default it runs on port 8080.

Run the agent with `go run cmd/agent/main.go --controller=http://localhost:8080 --controller2=http://localhost:8080 --project=prj_xxx --state-dir=./cmd/agent/state  --log-level=debug --registration-token=drt_1Lgz2FGL4jSvjqdZB3Bd7Z2ZdGn`
