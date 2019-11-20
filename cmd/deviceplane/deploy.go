package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/deviceplane/deviceplane/pkg/client"
	"github.com/deviceplane/deviceplane/pkg/interpolation"
	"github.com/urfave/cli"
)

var deploy = cli.Command{
	Name: "deploy",
	Flags: []cli.Flag{
		projectFlag,
		applicationFlag,
	},
	Action: func(c *cli.Context) error {
		return withClient(c, func(client *client.Client) error {
			project := c.String(projectFlag.Name)
			application := c.String(applicationFlag.Name)

			yamlConfigBytes, err := ioutil.ReadFile(c.Args().First())
			if err != nil {
				return err
			}

			finalYamlConfig, err := interpolation.Interpolate(string(yamlConfigBytes), os.Getenv)
			if err != nil {
				return err
			}

			release, err := client.CreateRelease(context.TODO(), project, application, finalYamlConfig)
			if err != nil {
				return err
			}

			fmt.Println(release.ID)

			return nil
		})
	},
}
