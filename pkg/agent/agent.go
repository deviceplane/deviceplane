package agent

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/engine"
	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/deviceplane/deviceplane/pkg/spec"
)

type Agent struct {
	client   *Client // TODO: interface
	engine   engine.Engine
	stateDir string
}

func NewAgent(client *Client, engine engine.Engine) *Agent {
	return &Agent{
		client: client,
		engine: engine,
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

			var applicationConfig spec.Application
			if err := yaml.Unmarshal([]byte(application.LatestRelease.Config), &applicationConfig); err != nil {
				log.WithError(err).Error("unmarshal config")
				continue
			}

			if err := a.Reconcile(context.TODO(), applicationConfig); err != nil {
				log.WithError(err).Error("reconcile")
				continue
			}
		}
	}

	return nil
}

func (a *Agent) Reconcile(ctx context.Context, application spec.Application) error {
	// TODO
	for name, service := range application.Services {
		service.Name = name
		application.Services[name] = service
	}

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		if err := a.removeDanglingServices(ctx, application); err != nil {
			log.WithError(err).Error("remove dangling containers")
		}
		wg.Done()
	}()

	for _, container := range application.Services {
		wg.Add(1)
		go func(service spec.Service) {
			if err := a.reconcileContainer(ctx, application, service); err != nil {
				log.WithError(err).Error("reconcile container")
			}
			wg.Done()
		}(container)
	}

	wg.Wait()

	return nil
}

func (a *Agent) reconcileContainer(ctx context.Context, application spec.Application, container spec.Service) error {
	instances, err := a.engine.List(ctx, nil, map[string]string{
		models.HashLabel: container.Hash(),
	})
	if err != nil {
		log.WithError(err).Error("list containers")
	}

	if len(instances) > 1 {
		return nil
	}

	containerWithHash := container.WithHash()
	// TODO
	containerWithHash.Name = fmt.Sprintf("%s-%s", container.Name, container.Hash()[:6])
	instance, err := a.engine.Create(ctx, containerWithHash)
	if err != nil {
		return err
	}
	if err = a.engine.Start(ctx, instance.ID); err != nil {
		return err
	}

	return nil
}

func (a *Agent) removeDanglingServices(ctx context.Context, application spec.Application) error {
	expectedHashValues := make(map[string]bool)
	for _, container := range application.Services {
		expectedHashValues[container.Hash()] = true
	}

	instances, err := a.engine.List(ctx, map[string]bool{
		models.HashLabel: true,
	}, nil)
	if err != nil {
		log.WithError(err).Error("list containers")
	}

	containersToBeRemoved := make(map[string]bool)
	for _, instance := range instances {
		hashValue, ok := instance.Labels[models.HashLabel]
		if !ok {
			continue
		}

		_, ok = expectedHashValues[hashValue]
		if !ok {
			containersToBeRemoved[instance.ID] = true
		}
	}

	wg := &sync.WaitGroup{}

	for id := range containersToBeRemoved {
		wg.Add(1)
		go func(id string) {
			if err := a.engine.Stop(ctx, id); err != nil {
				log.WithError(err).Error("stop container")
				return
			}

			if err := a.engine.Remove(ctx, id); err != nil {
				log.WithError(err).Error("remove container")
				return
			}
		}(id)
	}

	wg.Wait()

	return nil
}
