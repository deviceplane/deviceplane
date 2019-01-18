package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/deviceplane/deviceplane/pkg/client"
	"github.com/urfave/cli"
)

var release = cli.Command{
	Name:    "release",
	Aliases: []string{"r"},
	Subcommands: []cli.Command{
		{
			Name: "create",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "project",
				},
				cli.StringFlag{
					Name: "application",
				},
				cli.StringFlag{
					Name: "config",
				},
			},
			Action: func(c *cli.Context) error {
				url := c.GlobalString("url")
				client := client.NewClient(url, nil)

				projectID := c.String("project")
				applicationID := c.String("application")
				configLocation := c.String("config")

				configFile, err := os.Open(configLocation)
				if err != nil {
					return err
				}

				configBytes, err := ioutil.ReadAll(configFile)
				if err != nil {
					return err
				}

				release, err := client.CreateRelease(context.TODO(), projectID, applicationID, string(configBytes))
				if err != nil {
					return err
				}

				fmt.Println(release.ID)

				return nil
			},
		},
	},
}
