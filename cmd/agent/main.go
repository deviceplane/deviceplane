package main

import (
	"net/http"
	"net/url"

	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/agent"
	agent_client "github.com/deviceplane/deviceplane/pkg/agent/client"
	"github.com/deviceplane/deviceplane/pkg/engine/docker"
	"github.com/segmentio/conf"
)

var version = "dev"
var name = "deviceplane-agent"

var config struct {
	Controller        string `conf:"controller"`
	Controller2       string `conf:"controller2"`
	Project           string `conf:"project"`
	RegistrationToken string `conf:"registration-token"`
	StateDir          string `conf:"state-dir"`
	LogLevel          string `conf:"log-level"`
}

func init() {
	config.Controller = "https://api.deviceplane.io:443"
	config.Controller2 = "https://api2.deviceplane.io:443"
	config.StateDir = "/var/lib/deviceplane"
	config.LogLevel = "info"
}

func main() {
	conf.Load(&config)

	lvl, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		log.WithError(err).Fatal("--log-level")
	}
	log.SetLevel(lvl)

	engine, err := docker.NewEngine()
	if err != nil {
		log.WithError(err).Fatal("create docker client")
	}

	controllerURL, err := url.Parse(config.Controller)
	if err != nil {
		log.WithError(err).Fatal("parse controller URL")
	}

	controller2URL, err := url.Parse(config.Controller2)
	if err != nil {
		log.WithError(err).Fatal("parse controller2 URL")
	}

	client := agent_client.NewClient(controllerURL, controller2URL, config.Project, http.DefaultClient)
	agent := agent.NewAgent(client, engine, config.Project, config.RegistrationToken, config.StateDir)

	if err := agent.Initialize(); err != nil {
		log.WithError(err).Fatal("failure while initializing agent")
	}

	agent.Run()
}
