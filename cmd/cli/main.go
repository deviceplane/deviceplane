package main

import (
	"os"

	"github.com/urfave/cli"
)

var version = "dev"
var name = "deviceplane"

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "access-token",
		},
	}

	app.Commands = []cli.Command{
		project,
		application,
		edit,
		register,
		release,
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
