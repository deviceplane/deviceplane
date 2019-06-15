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
			Value:  "https://api.deviceplane.io:443",
		},
		cli.StringFlag{
			Hidden: true,
			Name:   "url2",
			Value:  "https://api2.deviceplane.io:443",
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
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
