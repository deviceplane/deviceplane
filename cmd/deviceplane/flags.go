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
)

func withClient(c *cli.Context, f func(*client.Client) error) error {
	u := c.GlobalString("url")
	url, err := url.Parse(u)
	if err != nil {
		return err
	}
	accessKey := c.GlobalString("access-key")
	return f(client.NewClient(url, accessKey, nil))
}
