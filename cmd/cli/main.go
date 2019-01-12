package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/deviceplane/deviceplane/pkg/client"
	"github.com/urfave/cli"
)

var version = "dev"
var name = "deviceplane"

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "url",
			Value: "http://0.0.0.0:8080",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "project",
			Aliases: []string{"p"},
			Subcommands: []cli.Command{
				{
					Name: "create",
					Action: func(c *cli.Context) error {
						url := c.GlobalString("url")

						client := client.NewClient(url, nil)

						project, err := client.CreateProject(context.TODO())
						if err != nil {
							return err
						}

						fmt.Println(project.ID)

						return nil
					},
				},
			},
		},
		{
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
		},
		{
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
		},
		{
			Name:    "bundle",
			Aliases: []string{"b"},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "project",
				},
			},
			Action: func(c *cli.Context) error {
				url := c.GlobalString("url")
				client := client.NewClient(url, nil)

				projectID := c.String("project")

				bundle, err := client.GetBundle(context.TODO(), projectID)
				if err != nil {
					return err
				}

				bundleBytes, err := json.Marshal(bundle)
				if err != nil {
					return err
				}

				fmt.Println(string(bundleBytes))

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
