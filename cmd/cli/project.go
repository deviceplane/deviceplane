package main

import (
	"context"
	"fmt"

	"github.com/deviceplane/deviceplane/pkg/client"
	"github.com/urfave/cli"
)

var project = cli.Command{
	Name:    "project",
	Aliases: []string{"p"},
	Subcommands: []cli.Command{
		{
			Name: "create",
			Action: func(c *cli.Context) error {
				accessToken := c.GlobalString("access-token")

				client := client.NewClient(accessToken, nil)

				project, err := client.CreateProject(context.TODO())
				if err != nil {
					return err
				}

				fmt.Println(project.ID)

				return nil
			},
		},
	},
}
