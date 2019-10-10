# Device Plane

This repository contains the backend (controller), agent, and CLI code.

## Local Development

Run `make db-reset` to setup the database. This command will reset the database to an empty state and then seed it with some basic data.

Run the controller with `go run cmd/controller/main.go`. By default it runs on port 8080.

Run the [UI](https://github.com/deviceplane/app) with `npm start`. The login is `email@example.com` / `password`.

Run the agent with `go run cmd/agent/main.go --controller=http://localhost:8080/api --project=prj_xxx --conf-dir=./cmd/agent/conf --state-dir=./cmd/agent/state --log-level=debug --registration-token=drt_1Lgz2FGL4jSvjqdZB3Bd7Z2ZdGn`. You will need to change the registration token for this command. To create a new one click on "Add Device" on the devices page.

## Releasing

While on an up-to-date master branch run one of the following commands and then push to master.

```
git commit --allow-empty -m "Release controller x.x.x"
git commit --allow-empty -m "Release agent x.x.x"
git commit --allow-empty -m "Release CLI x.x.x"
```

## Testing

There are no tests. Yet. To manually test things like db migrations, run `make clone-db-locally` to get a local copy of the prod database. You can run db migrations, the ui, the controller, and local agents as needed.