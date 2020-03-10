package main

import (
	"net/url"
	"os"

	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/agent"
	agent_client "github.com/deviceplane/deviceplane/pkg/agent/client"
	dphttp "github.com/deviceplane/deviceplane/pkg/http"
	"github.com/segmentio/conf"
)

var version = "dev"
var name = "deviceplane-agent"

var config struct {
	Controller        string `conf:"controller"`
	Project           string `conf:"project"`
	RegistrationToken string `conf:"registration-token"`
	ConfDir           string `conf:"conf-dir"`
	StateDir          string `conf:"state-dir"`
	ServerPort        int    `conf:"server-port"`
	LogLevel          string `conf:"log-level"`
}

func init() {
	config.Controller = "https://cloud.deviceplane.com:443/api"
	config.ConfDir = "/etc/deviceplane"
	config.StateDir = "/var/lib/deviceplane"
	config.ServerPort = 4444
	config.LogLevel = "info"
}

func main() {
	conf.Load(&config)

	lvl, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		log.WithError(err).Fatal("--log-level")
	}
	log.SetLevel(lvl)

	controllerURL, err := url.Parse(config.Controller)
	if err != nil {
		log.WithError(err).Fatal("parse controller URL")
	}

	client := agent_client.NewClient(controllerURL, config.Project, dphttp.DefaultClient)
	agent, err := agent.NewAgent(client, config.Project, config.RegistrationToken,
		config.ConfDir, config.StateDir, version, os.Args[0], config.ServerPort)
	if err != nil {
		log.WithError(err).Fatal("failure creating agent")
	}

	if err := agent.Initialize(); err != nil {
		log.WithError(err).Fatal("failure while initializing agent")
	}

	agent.Run()
}
