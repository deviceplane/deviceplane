package handoff

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/deviceplane/deviceplane/pkg/agent/utils"
	"github.com/deviceplane/deviceplane/pkg/engine"
	"github.com/deviceplane/deviceplane/pkg/models"
)

type Coordinator struct {
	engine  engine.Engine
	version string
	port    int
}

func NewCoordinator(engine engine.Engine, version string, port int) *Coordinator {
	return &Coordinator{
		engine:  engine,
		version: version,
		port:    port,
	}
}

func (c *Coordinator) Takeover() net.Listener {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", c.port))
		if err == nil {
			return listener
		}

		instances, err := utils.ContainerList(context.TODO(), c.engine, map[string]struct{}{
			models.AgentVersionLabel: struct{}{},
		}, nil, true)
		if err != nil {
			goto cont
		}

		for _, instance := range instances {
			if instance.Labels[models.AgentVersionLabel] != c.version {
				if err = utils.ContainerStop(context.TODO(), c.engine, instance.ID); err != nil {
					goto cont
				}
				if err = utils.ContainerRemove(context.TODO(), c.engine, instance.ID); err != nil {
					goto cont
				}
			}
		}

	cont:
		select {
		case <-ticker.C:
			continue
		}
	}
}
