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

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Hidden: true,
			Name:   "url",
			Value:  "https://api.deviceplane.io",
		},
		cli.StringFlag{
			Name:   "access-key",
			EnvVar: "DEVICE_PLANE_ACCESS_KEY",
		},
	}

	app.Commands = []cli.Command{
		project,
		application,
		edit,
		deploy,
		release,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
