package main

import (
	"net/url"

	"github.com/deviceplane/deviceplane/pkg/client"
	"github.com/urfave/cli"
)

var (
	// Global flags
	urlFlag = cli.StringFlag{
		Hidden: true,
		Name:   "url",
		Value:  "https://cloud.deviceplane.com:443/api",
	}
	accessKeyFlag = cli.StringFlag{
		Name:   "access-key",
		EnvVar: "DEVICEPLANE_ACCESS_KEY",
	}

	// Common flags
	projectFlag = cli.StringFlag{
		Name:   "project",
		EnvVar: "DEVICEPLANE_PROJECT",
	}
	applicationFlag = cli.StringFlag{
		Name:   "application",
		EnvVar: "DEVICEPLANE_APPLICATION",
	}
	deviceFlag = cli.StringFlag{
		Name:   "device",
		EnvVar: "DEVICEPLANE_DEVICE",
	}

	// SSH command flags
	sshConnectTimeoutFlag = cli.StringFlag{
		Name:   "connect-timeout",
		EnvVar: "DEVICEPLANE_SSH_CONNECT_TIMEOUT",
		Value:  "60",
	}
)

func withClient(c *cli.Context, f func(*client.Client) error) error {
	apiURL, err := url.Parse(c.GlobalString(urlFlag.Name))
	if err != nil {
		return err
	}
	accessKey := c.GlobalString(accessKeyFlag.Name)
	return f(client.NewClient(apiURL, accessKey, nil))
}
