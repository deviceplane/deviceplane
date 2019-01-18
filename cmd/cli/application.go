package main

import (
	"context"
	"fmt"

	"github.com/deviceplane/deviceplane/pkg/client"
	"github.com/urfave/cli"
)

var application = cli.Command{
	Name:    "application",
	Aliases: []string{"a"},
	Subcommands: []cli.Command{
		{
			Name: "create",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "project",
				},
			},
			Action: func(c *cli.Context) error {
				url := c.GlobalString("url")
				client := client.NewClient(url, nil)

				projectID := c.String("project")

				application, err := client.CreateApplication(context.TODO(), projectID)
				if err != nil {
					return err
				}

				fmt.Println(application.ID)

				return nil
			},
		},
	},
}
