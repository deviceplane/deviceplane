package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/deviceplane/deviceplane/pkg/client"
	"github.com/urfave/cli"
)

var edit = cli.Command{
	Name: "edit",
	Flags: []cli.Flag{
		projectFlag,
		applicationFlag,
	},
	Action: func(c *cli.Context) error {
		return withClient(c, func(client *client.Client) error {
			project := c.String("project")
			application := c.String("application")

			release, err := client.GetLatestRelease(context.TODO(), project, application)
			if err != nil {
				return err
			}

			var yamlConfig string
			if release != nil {
				yamlConfig = release.RawConfig
			}

			tmpfile, err := ioutil.TempFile("", "")
			if err != nil {
				return err
			}
			defer os.Remove(tmpfile.Name())

			if _, err := tmpfile.Write([]byte(yamlConfig)); err != nil {
				return err
			}

			editor := os.Getenv("EDITOR")
			if editor == "" {
				editor = "vi"
			}

			cmd := exec.Command(editor, tmpfile.Name())
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err = cmd.Run(); err != nil {
				fmt.Println("Edit cancelled, no changes made.")
				return nil
			}

			if err := tmpfile.Close(); err != nil {
				return err
			}

			yamlConfigFile, err := os.Open(tmpfile.Name())
			if err != nil {
				return err
			}

			yamlConfigBytes, err := ioutil.ReadAll(yamlConfigFile)
			if err != nil {
				return err
			}

			release, err = client.CreateRelease(context.TODO(), project, application, string(yamlConfigBytes))
			if err != nil {
				return err
			}

			fmt.Println(release.ID)

			return nil
		})
	},
}
