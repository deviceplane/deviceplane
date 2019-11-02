package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

var version = "dev"
var name = "deviceplane"

func main() {
	app := cli.NewApp()

	app.EnableBashCompletion = true
	app.Name = name
	app.Version = version
	app.Usage = "Device Plane CLI"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Hidden: true,
			Name:   "url",
			Value:  "https://cloud.deviceplane.com:443/api",
		},
		cli.StringFlag{
			Name:   "access-key",
			EnvVar: "DEVICEPLANE_ACCESS_KEY",
		},
	}

	app.Commands = []cli.Command{
		project,
		application,
		edit,
		deploy,
		ssh,
		execute,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
