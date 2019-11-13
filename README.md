# Deviceplane

This repository contains the backend (controller), agent, and CLI code.

## Local Development

Run `make db-reset` to setup the database. This command will reset the database to an empty state and then seed it with some basic data.

Run the controller with `go run cmd/controller/main.go`. By default it runs on port 8080.

Run the UI with `npm start` in the `ui/` folder. The login is `email@example.com` / `password`.

To run the agent navigate to the "Add Device" button in the UI. A command to run the agent locally will be generated.

## Releasing

While on an up-to-date master branch run one of the following commands and then push to master.

```
git commit --allow-empty -m "Release controller x.x.x"
git commit --allow-empty -m "Release agent x.x.x"
git commit --allow-empty -m "Release CLI x.x.x"
```
