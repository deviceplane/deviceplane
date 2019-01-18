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
			Name:  "url",
			Value: "http://0.0.0.0:8080",
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
