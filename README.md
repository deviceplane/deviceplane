# Deviceplane

[![CircleCI](https://circleci.com/gh/deviceplane/deviceplane.svg?style=svg)](https://circleci.com/gh/deviceplane/deviceplane)

Deviceplane is an open source device management tool for embedded systems and edge computing. It helps solve various infrastructure problems relating to remote device management such as:

- Network connectivity and SSH access
- Orchestrating and deployment of remote updates
- Host and application monitoring
- Device organization: naming, labeling, searching, and filtering of devices
- Access and security controls

Deviceplane integrates with your device by running a lightweight static binary via your system supervisor. It can be used with nearly any Linux distro, which means you can continue using Ubuntu, Raspbian, a Yocto build, or whatever else fits your needs.

Deviceplane is completely open source and all of our code, including the code for the backend and UI, can be found in this repo. The architecture of Deviceplane is simple and designed to be very easy to run and manage. The backend can even be run with one simple Docker command:

```
docker run -d --restart=unless-stopped -p 8080:8080 deviceplane/deviceplane
```

For more information on hosting Deviceplane yourself, check out our [self-hosted docs](https://deviceplane.com/docs/self-hosted/).

If you'd rather jump right into managing devices there's a hosted version of Deviceplane available at [https://cloud.deviceplane.com/](https://cloud.deviceplane.com/).

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
