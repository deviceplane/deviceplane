# Deviceplane

[![CircleCI](https://circleci.com/gh/deviceplane/deviceplane.svg?style=svg)](https://circleci.com/gh/deviceplane/deviceplane)

Deviceplane is an open source platform for managing IoT devices and edge servers with a modern developer experience. Beyond remote updates, Deviceplane handles the entire lifecycle of managing remote devices - that’s everything from provisioning to access management and monitoring.

A hosted version of Deviceplane is available at [https://cloud.deviceplane.com/](https://cloud.deviceplane.com/).

We also have documentation for self-hosting Deviceplane available at [https://deviceplane.com/docs/self-hosted/](https://deviceplane.com/docs/self-hosted/).

This repository contains all Deviceplane code including the backend (controller), agent, and CLI.

## Local Development

Run `make db-reset` to setup the database. This command will reset the database to an empty state and then seed it with some basic data.

Run the controller with `go run cmd/controller/main.go --allowed-origin "http://localhost:3000"`. By default it runs on port 8080.

Run the UI with `npm start` in the `ui/` folder. The login is `email@example.com` / `password`.

To run the agent navigate to the "Add Device" button in the UI. A command to run the agent locally will be generated.

## Releasing

Release the controller, agent, or CLI by pushing git tags.

```
git tag controller-x.x.x
git push origin controller-x.x.x
```

```
git tag agent-x.x.x
git push origin agent-x.x.x
```

```
git tag cli-x.x.x
git push origin cli-x.x.x
```
