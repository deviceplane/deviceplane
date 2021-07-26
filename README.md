
> :warning: This repository has been deprecated. The new repositories can be found here:
> [Agent](https://github.com/deviceplane/agent) | [CLI](https://github.com/deviceplane/cli) | [Controller](https://github.com/deviceplane/controller).
<p align="center">
    <img src="./logo/name-black.png" alt="Build Status" width="500">
</p>

<p align="center">
    <a aria-label="Build Status" href="https://circleci.com/gh/deviceplane/deviceplane" target="_blank">
        <img src="https://img.shields.io/circleci/build/github/deviceplane/deviceplane?style=flat-square" alt="Build Status">
    </a>
    <a aria-label="License" href="https://github.com/deviceplane/deviceplane/LICENSE.md" target="_blank">
        <img src="https://img.shields.io/github/license/deviceplane/deviceplane?color=%23000&style=flat-square" alt="Build Status">
    </a>
</p>

Deviceplane is an open source device management tool for embedded systems and edge computing. It solves various infrastructure problems related to remote device management such as:

- Network connectivity and SSH access
- Orchestration and deployment of remote updates
- Host and application monitoring
- Device organization: naming, labeling, searching, and filtering of devices
- Access and security controls

Deviceplane integrates with your device by running a lightweight static binary via your system supervisor. It can be used with nearly any Linux distro, which means you can continue using Ubuntu, Raspbian, a Yocto build, or whatever else fits your needs.

A hosted version of Deviceplane is available at [https://cloud.deviceplane.com/](https://cloud.deviceplane.com/).

## Documentation

Visit <a aria-label="next.js learn" href="https://deviceplane.com/docs">https://deviceplane.com/docs</a> to view the full documentation.

## Getting Started

The architecture of Deviceplane is simple and designed to be easy to run and manage. The backend can be run with a single Docker command:

```
docker run -d --restart=unless-stopped -p 8080:8080 deviceplane/deviceplane
```

For more information on hosting Deviceplane yourself, check out the [self-hosted docs](https://deviceplane.com/docs/self-hosted/).

## Local Development

#### Setup the database

This command will reset the database to an empty state and then seed it with some basic data.

```
make db-reset
```

#### Run the controller

This command starts the controller running on port 8080 by default.

```
go run cmd/controller/main.go --allowed-origin "http://localhost:3000"
```

#### Run the web application

Run the following commands in the `ui/` directory

```
npm install
npm start
```

The login is `email@example.com` / `password`.

#### Run the agent

Navigate to the `/devices/register` route in the web application. A command to run the agent locally will be generated and printed to the browser console.

## Support

For bugs, issues, and feature requests please [submit](//github.com/deviceplane/deviceplane/issues/new) a GitHub issue.

For security issues, please email security@deviceplane.com instead of posting a public issue on GitHub.

For support, please email support@deviceplane.com.

## License

Copyright (c) Deviceplane, Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
