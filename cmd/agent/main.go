package main

import (
	"context"
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
}

func init() {
	config.Controller = "https://api.deviceplane.io"
}

func main() {
	conf.Load(&config)

	engine, err := docker.NewEngine()
	if err != nil {
		log.WithError(err).Fatal("create docker client")
	}

	client := agent.NewClient(config.Controller, config.Project, http.DefaultClient)

	// TODO: check for existing access key

	registerDeviceResponse, err := client.RegisterDevice(context.Background(), config.RegistrationToken)
	if err != nil {
		log.WithError(err).Fatal("register device")
	}

	client.SetDeviceID(registerDeviceResponse.DeviceID)
	client.SetDeviceID(registerDeviceResponse.DeviceAccessKeyValue)

	agent := agent.NewAgent(client, engine)
	if err := agent.Run(); err != nil {
		log.WithError(err).Fatal("agent error")
	}
}
