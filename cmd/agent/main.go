package main

import (
	"fmt"

	"github.com/deviceplane/deviceplane/pkg/agent"
	"github.com/deviceplane/deviceplane/pkg/client"
	"github.com/deviceplane/deviceplane/pkg/engine/docker"
	"github.com/segmentio/conf"
)

var version = "dev"
var name = "deviceplane-agent"

var config struct {
	Controller string `conf:"controller"`
	Project    string `conf:"project"`
}

func init() {
	config.Controller = "http://0.0.0.0:8080"
}

func main() {
	conf.Load(&config)

	client := client.NewClient(config.Controller, nil)

	engine, err := docker.NewEngine()
	if err != nil {
		panic(err)
	}

	agent := agent.NewAgent(client, engine, config.Project)

	fmt.Println(agent.Run())

	/*err = agent.Reconcile(context.Background(), spec.Application{
		Containers: map[string]spec.Container{
			"josh": {
				Image:   "ubuntu",
				Command: []string{"sleep", "infinity"},
			},
		},
	})
	fmt.Println(err)*/
}
