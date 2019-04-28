package supervisor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/engine"
	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/deviceplane/deviceplane/pkg/spec"
	"gopkg.in/yaml.v2"
)

type Supervisor struct {
	engine              engine.Engine
	reconcilingServices map[string]struct{}
	keepAliveShutdowns  map[string]chan struct{}
	keepAliveAcks       map[string]chan struct{}
	lock                sync.Mutex
}

func NewSupervisor(engine engine.Engine) *Supervisor {
	return &Supervisor{
		engine:              engine,
		reconcilingServices: make(map[string]struct{}),
		keepAliveShutdowns:  make(map[string]chan struct{}),
		keepAliveAcks:       make(map[string]chan struct{}),
	}
}

// TODO: removes applications and services
func (s *Supervisor) SetApplication(application models.ApplicationAndLatestRelease) error {
	var applicationConfig spec.Application
	if err := yaml.Unmarshal([]byte(application.LatestRelease.Config), &applicationConfig); err != nil {
		return err
	}

	for serviceName, service := range applicationConfig.Services {
		go s.reconcile(serviceName, service)
	}

	return nil
}

func (s *Supervisor) reconcile(serviceName string, service spec.Service) {
	s.lock.Lock()
	// If this service is already in reconciling state then exit
	_, ok := s.reconcilingServices[serviceName]
	if ok {
		s.lock.Unlock()
		return
	}
	// Set this service to reconciling state
	s.reconcilingServices[serviceName] = struct{}{}
	s.lock.Unlock()

	defer func() {
		s.lock.Lock()
		delete(s.reconcilingServices, serviceName)
		s.lock.Unlock()
	}()

	stopKeepAlive := func() {
		s.lock.Lock()

		shutdown, ok := s.keepAliveShutdowns[serviceName]
		if !ok {
			s.lock.Unlock()
			return
		}

		ack := make(chan struct{})
		s.keepAliveAcks[serviceName] = ack
		shutdown <- struct{}{}
		<-ack
		delete(s.keepAliveAcks, serviceName)

		s.lock.Unlock()
	}

	instances := s.containerList(nil, map[string]string{
		models.ServiceLabel: serviceName,
	}, true)

	if len(instances) > 0 {
		// TODO: filter down to just one instance if we find more
		instance := instances[0]

		if hashLabel, ok := instance.Labels[models.HashLabel]; ok && hashLabel == service.Hash() {
			go s.keepAlive(serviceName, service)
			return
		}

		s.imagePull(service.Image)
		stopKeepAlive()

		s.containerStop(instance.ID)
		s.containerRemove(instance.ID)
	}

	s.imagePull(service.Image)
	stopKeepAlive()

	name := fmt.Sprintf("%s-%s", serviceName, service.Hash()[:6])
	serviceWithHash := service.WithStandardLabels(serviceName)
	s.containerCreate(name, serviceWithHash)

	go s.keepAlive(serviceName, service)
}

func (s *Supervisor) keepAlive(serviceName string, service spec.Service) {
	s.lock.Lock()
	if _, ok := s.keepAliveShutdowns[serviceName]; ok {
		s.lock.Unlock()
		return
	}
	shutdown := make(chan struct{}, 1)
	s.keepAliveShutdowns[serviceName] = shutdown
	s.lock.Unlock()

	defer func() {
		s.lock.Lock()
		delete(s.keepAliveShutdowns, serviceName)
		s.lock.Unlock()
	}()

	dead := false
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-shutdown:
			s.keepAliveAcks[serviceName] <- struct{}{}
			return

		case <-ticker.C:
			if dead {
				continue
			}
			instances := s.containerList(nil, map[string]string{
				models.ServiceLabel: serviceName,
				models.HashLabel:    service.Hash(),
			}, true)

			if len(instances) == 0 {
				go s.reconcile(serviceName, service)
				dead = true
				continue
			}

			// TODO: filter down to just one instance if we find more
			instance := instances[0]

			if !instance.Running {
				s.containerStart(instance.ID)
			}
		}
	}
}

func (s *Supervisor) containerCreate(name string, service spec.Service) string {
	var id string

	retry(func(ctx context.Context) error {
		var err error
		id, err = s.engine.CreateContainer(ctx, name, service)
		if err != nil {
			log.WithError(err).Error("create container")
			return err
		}
		return nil
	}, 5*time.Minute)

	return id
}

func (s *Supervisor) containerStart(id string) {
	retry(func(ctx context.Context) error {
		if err := s.engine.StartContainer(ctx, id); err != nil && err != engine.ErrInstanceNotFound {
			log.WithError(err).Error("start container")
			return err
		}
		return nil
	}, 5*time.Minute)
}

func (s *Supervisor) containerList(keyFilters map[string]bool, keyAndValueFilters map[string]string, all bool) []engine.Instance {
	var instances []engine.Instance

	retry(func(ctx context.Context) error {
		var err error
		instances, err = s.engine.ListContainers(context.TODO(), keyFilters, keyAndValueFilters, all)
		if err != nil {
			log.WithError(err).Error("list containers")
			return err
		}
		return nil
	}, 5*time.Minute)

	return instances
}

func (s *Supervisor) containerStop(id string) {
	retry(func(ctx context.Context) error {
		if err := s.engine.StopContainer(ctx, id); err != nil && err != engine.ErrInstanceNotFound {
			log.WithError(err).Error("stop container")
			return err
		}
		return nil
	}, 5*time.Minute)
}

func (s *Supervisor) containerRemove(id string) {
	retry(func(ctx context.Context) error {
		if err := s.engine.RemoveContainer(ctx, id); err != nil && err != engine.ErrInstanceNotFound {
			log.WithError(err).Error("remove container")
			return err
		}
		return nil
	}, 5*time.Minute)
}

func (s *Supervisor) imagePull(image string) {
	retry(func(ctx context.Context) error {
		if err := s.engine.PullImage(ctx, image); err != nil {
			log.WithError(err).Error("pull image")
			return err
		}
		return nil
	}, 30*time.Minute)
}

func retry(f func(context.Context) error, timeout time.Duration) {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			if err := f(ctx); err == nil {
				return
			}
		}
	}
}
