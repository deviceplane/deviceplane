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
				projectFlag,
			},
			Action: func(c *cli.Context) error {
				return withClient(c, func(client *client.Client) error {
					project := c.String(projectFlag.Name)

					application, err := client.CreateApplication(context.TODO(), project)
					if err != nil {
						return err
					}

					fmt.Println(application.ID)

					return nil
				})
			},
		},
	},
}
