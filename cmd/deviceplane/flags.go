package main

import (
	"github.com/urfave/cli"
)

var (
	projectFlag = cli.StringFlag{
		Name:   "project",
		EnvVar: "DEVICE_PLANE_PROJECT",
	}
	applicationFlag = cli.StringFlag{
		Name:   "application",
		EnvVar: "DEVICE_PLANE_APPLICATION",
	}
)
