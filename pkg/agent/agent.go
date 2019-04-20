package agent

import (
	"context"
	"time"

	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/agent/supervisor"
	"github.com/deviceplane/deviceplane/pkg/engine"
)

type Agent struct {
	client     *Client // TODO: interface
	engine     engine.Engine
	stateDir   string
	supervisor *supervisor.Supervisor
}

func NewAgent(client *Client, engine engine.Engine) *Agent {
	return &Agent{
		client:     client,
		engine:     engine,
		supervisor: supervisor.NewSupervisor(engine),
	}
}

func (a *Agent) Run() error {
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		bundle, err := a.client.getBundle(context.TODO())
		if err != nil {
			log.WithError(err).Error("get bundle")
			continue
		}

		for _, application := range bundle.Applications {
			if application.LatestRelease == nil {
				continue
			}

			if err := a.supervisor.SetApplication(application); err != nil {
				log.WithError(err).Error("set application")
				continue
			}
		}
	}

	return nil
}
