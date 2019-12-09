package main

import (
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"

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
	// Start zombie reaper
	go reaper()

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

	client := agent_client.NewClient(controllerURL, config.Project, http.DefaultClient)
	agent, err := agent.NewAgent(client, engine, config.Project, config.RegistrationToken,
		config.ConfDir, config.StateDir, version, os.Args[0], config.ServerPort)
	if err != nil {
		log.WithError(err).Fatal("failure creating agent")
	}

	if err := agent.Initialize(); err != nil {
		log.WithError(err).Fatal("failure while initializing agent")
	}

	agent.Run()
}

// reaper reaps zombie processes
// This is needed because SSH spawns child processes that can die
// https://en.wikipedia.org/wiki/Zombie_process
func reaper() {
	c := make(chan os.Signal, 2048)
	signal.Notify(c, syscall.SIGCHLD)

	for range c {
		for {
			if pid, err := syscall.Wait4(
				-1, nil, syscall.WNOHANG, nil,
			); err != nil || pid <= 0 {
				break
			}
		}
	}
}
