package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/deviceplane/deviceplane/pkg/client"
	"github.com/urfave/cli"
)

var execute = cli.Command{
	Name: "execute",
	Flags: []cli.Flag{
		projectFlag,
		deviceFlag,
	},
	Action: func(c *cli.Context) error {
		return withClient(c, func(client *client.Client) error {
			project := c.String("project")
			device := c.String("device")

			executeResponse, err := client.Execute(
				context.TODO(), project, device,
				strings.Join(c.Args(), " "),
			)
			if err != nil {
				return err
			}

			fmt.Println("Exit code:", executeResponse.ExitCode)

			return nil
		})
	},
}
