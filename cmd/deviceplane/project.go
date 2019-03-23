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
				url := c.GlobalString("url")
				accessKey := c.GlobalString("access-key")
				client := client.NewClient(url, accessKey, nil)

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
