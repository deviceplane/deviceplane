package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/deviceplane/deviceplane/pkg/client"
	"github.com/deviceplane/deviceplane/pkg/interpolation"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

var deploy = cli.Command{
	Name: "deploy",
	Flags: []cli.Flag{
		projectFlag,
		applicationFlag,
	},
	Action: func(c *cli.Context) error {
		return withClient(c, func(client *client.Client) error {
			project := c.String("project")
			application := c.String("application")

			configBytes, err := ioutil.ReadFile(c.Args().First())
			if err != nil {
				return err
			}

			var configMap map[string]interface{}
			if err := yaml.Unmarshal(configBytes, &configMap); err != nil {
				return err
			}

			if err := interpolation.Interpolate(configMap, os.Getenv); err != nil {
				return err
			}

			configBytes, err = yaml.Marshal(configMap)
			if err != nil {
				return err
			}

			release, err := client.CreateRelease(context.TODO(), project, application, string(configBytes))
			if err != nil {
				return err
			}

			fmt.Println(release.ID)

			return nil
		})
	},
}
