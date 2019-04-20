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

// TODO: we might need state to not constsntly retry failing updates
type Supervisor struct {
	engine              engine.Engine
	reconcilingServices map[string]struct{}
	keepAliveShutdowns  map[string]chan struct{}
	keepAliveAcks       map[string]chan struct{}
	lock                sync.RWMutex
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

// this is a little confusing because ports of this can run at the same time as keepalive
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

	instances, err := s.engine.List(context.TODO(), nil, map[string]string{
		models.ServiceLabel: serviceName,
	})
	if err != nil {
		// TODO
		log.WithError(err).Error("list containers")
		panic(err)
	}

	if len(instances) > 0 {
		instance := instances[0]

		if hashLabel, ok := instance.Labels[models.HashLabel]; ok && hashLabel == service.Hash() {
			go s.keepAlive(serviceName, service)
			return
		}

		stopKeepAlive()

		if err := s.engine.Stop(context.TODO(), instance.ID); err != nil {
			// TODO
			log.WithError(err).Error("stop container")
			panic(err)
		}

		if err := s.engine.Remove(context.TODO(), instance.ID); err != nil {
			// TODO
			log.WithError(err).Error("remove container")
			panic(err)
		}
	}

	// TODO: is this really the right thing to do here?
	stopKeepAlive()

	serviceWithHash := service.WithStandardLabels(serviceName)
	// TODO
	serviceWithHash.Name = fmt.Sprintf("%s-%s", serviceName, service.Hash()[:6])
	instance, err := s.engine.Create(context.TODO(), serviceWithHash)
	if err != nil {
		// TODO
		log.WithError(err).Error("create container 1")
		panic(err)
	}
	if err = s.engine.Start(context.TODO(), instance.ID); err != nil {
		// TODO
		log.WithError(err).Error("start container")
		panic(err)
	}

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

	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-shutdown:
			s.keepAliveAcks[serviceName] <- struct{}{}
			return

		case <-ticker.C:
			instances, err := s.engine.List(context.TODO(), nil, map[string]string{
				// TODO: probably want AND hash here
				models.ServiceLabel: serviceName,
			})
			if err != nil {
				// TODO
				log.WithError(err).Error("list containers")
			}

			if len(instances) > 0 {
				continue
			}

			serviceWithHash := service.WithStandardLabels(serviceName)
			// TODO
			serviceWithHash.Name = fmt.Sprintf("%s-%s", serviceName, service.Hash()[:6])
			instance, err := s.engine.Create(context.TODO(), serviceWithHash)
			if err != nil {
				// TODO
				log.WithError(err).Error("create container 2")
				panic(err)
			}
			if err = s.engine.Start(context.TODO(), instance.ID); err != nil {
				// TODO
				log.WithError(err).Error("start container")
				panic(err)
			}
		}
	}
}
