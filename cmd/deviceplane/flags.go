package main

import (
	"net/url"

	"github.com/deviceplane/deviceplane/pkg/client"
	"github.com/urfave/cli"
)

var (
	projectFlag = cli.StringFlag{
		Name:   "project",
		EnvVar: "DEVICE_PLANE_PROJECT",
	}
	applicationFlag = cli.StringFlag{
		Name:   "application",
		EnvVar: "DEVICE_PLANE_APPLICATION",
	}
	deviceFlag = cli.StringFlag{
		Name:   "device",
		EnvVar: "DEVICE_PLANE_DEVICE",
	}
)

func withClient(c *cli.Context, f func(*client.Client) error) error {
	u, err := url.Parse(c.GlobalString("url"))
	if err != nil {
		return err
	}
	u2, err := url.Parse(c.GlobalString("url2"))
	if err != nil {
		return err
	}
	accessKey := c.GlobalString("access-key")
	return f(client.NewClient(u, u2, accessKey, nil))
}
