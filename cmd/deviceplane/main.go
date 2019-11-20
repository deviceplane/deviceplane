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
	app.Usage = "Deviceplane CLI"

	app.Flags = []cli.Flag{
		urlFlag,
		accessKeyFlag,
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
