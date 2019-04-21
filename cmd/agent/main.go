package main

import (
	"net/http"

	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/agent"
	"github.com/deviceplane/deviceplane/pkg/engine/docker"
	"github.com/segmentio/conf"
)

var version = "dev"
var name = "deviceplane-agent"

var config struct {
	Controller        string `conf:"controller"`
	Project           string `conf:"project"`
	RegistrationToken string `conf:"registration-token"`
	StateDir          string `conf:"state-dir"`
}

func init() {
	config.Controller = "https://api.deviceplane.io"
	config.StateDir = "/var/lib/deviceplane"
}

func main() {
	conf.Load(&config)

	engine, err := docker.NewEngine()
	if err != nil {
		log.WithError(err).Fatal("create docker client")
	}

	client := agent.NewClient(config.Controller, config.Project, http.DefaultClient)
	agent := agent.NewAgent(client, engine, config.Project, config.RegistrationToken, config.StateDir)

	if err := agent.Initialize(); err != nil {
		log.WithError(err).Fatal("failure while initializing agent")
	}

	if err := agent.Run(); err != nil {
		log.WithError(err).Fatal("failure while running agent")
	}
}
