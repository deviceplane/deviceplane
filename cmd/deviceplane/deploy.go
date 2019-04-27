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
		cli.StringFlag{
			Name: "project",
		},
		cli.StringFlag{
			Name: "application",
		},
	},
	Action: func(c *cli.Context) error {
		url := c.GlobalString("url")
		accessKey := c.GlobalString("access-key")
		client := client.NewClient(url, accessKey, nil)

		projectID := c.String("project")
		applicationID := c.String("application")

		release, err := client.GetLatestRelease(context.TODO(), projectID, applicationID)
		if err != nil {
			return err
		}

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

		release, err = client.CreateRelease(context.TODO(), projectID, applicationID, string(configBytes))
		if err != nil {
			return err
		}

		fmt.Println(release.ID)

		return nil
	},
}
