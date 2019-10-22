package supervisor

import (
	"context"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/agent/utils"
	"github.com/deviceplane/deviceplane/pkg/agent/validator"
	"github.com/deviceplane/deviceplane/pkg/engine"
	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/deviceplane/deviceplane/pkg/spec"
	"gopkg.in/yaml.v2"
)

type ApplicationSupervisor struct {
	applicationID string
	engine        engine.Engine
	reporter      *Reporter
	validators    []validator.Validator

	serviceNames            map[string]struct{}
	serviceSupervisors      map[string]*ServiceSupervisor
	serviceSupervisorGCDone chan struct{}
	containerGCDone         chan struct{}

	once     sync.Once
	lock     sync.RWMutex
	stopLock sync.Mutex
	ctx      context.Context
	cancel   func()
}

func NewApplicationSupervisor(
	applicationID string,
	engine engine.Engine,
	reporter *Reporter,
	validators []validator.Validator,
) *ApplicationSupervisor {
	ctx, cancel := context.WithCancel(context.Background())
	return &ApplicationSupervisor{
		applicationID: applicationID,
		engine:        engine,
		reporter:      reporter,
		validators:    validators,

		serviceNames:            make(map[string]struct{}),
		serviceSupervisors:      make(map[string]*ServiceSupervisor),
		serviceSupervisorGCDone: make(chan struct{}),
		containerGCDone:         make(chan struct{}),

		ctx:    ctx,
		cancel: cancel,
	}
}

func (s *ApplicationSupervisor) SetApplication(application models.ApplicationFull2) {
	s.stopLock.Lock()
	defer s.stopLock.Unlock()

	select {
	case <-s.ctx.Done():
		return
	default:
		break
	}

	var applicationConfig map[string]spec.Service
	if err := yaml.Unmarshal([]byte(application.LatestRelease.Config), &applicationConfig); err != nil {
		log.WithError(err).Error("unmarshal")
		return
	}

	s.reporter.SetDesiredApplication(application.LatestRelease.ID, applicationConfig)

	serviceNames := make(map[string]struct{})
	for serviceName, service := range applicationConfig {
		s.lock.Lock()
		serviceSupervisor, ok := s.serviceSupervisors[serviceName]
		if !ok {
			serviceSupervisor = NewServiceSupervisor(
				application.Application.ID,
				serviceName,
				s.engine,
				s.reporter,
				s.validators,
			)
			s.serviceSupervisors[serviceName] = serviceSupervisor
		}
		s.lock.Unlock()

		serviceSupervisor.SetService(application.LatestRelease.ID, service)

		serviceNames[serviceName] = struct{}{}
	}

	s.lock.Lock()
	s.serviceNames = serviceNames
	s.lock.Unlock()

	s.once.Do(func() {
		go s.serviceSupervisorGC()
		go s.containerGC()
	})
}

func (s *ApplicationSupervisor) Stop() {
	s.stopLock.Lock()
	defer s.stopLock.Unlock()

	s.cancel()

	wg := &sync.WaitGroup{}
	wg.Add(len(s.serviceSupervisors) + 3)

	go func() {
		s.reporter.Stop()
		wg.Done()
	}()
	go func() {
		<-s.serviceSupervisorGCDone
		wg.Done()
	}()
	go func() {
		<-s.containerGCDone
		wg.Done()
	}()
	for _, serviceSupervisor := range s.serviceSupervisors {
		go func(serviceSupervisor *ServiceSupervisor) {
			serviceSupervisor.Stop()
			wg.Done()
		}(serviceSupervisor)
	}

	wg.Wait()
}

func (s *ApplicationSupervisor) serviceSupervisorGC() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		s.lock.RLock()
		danglingServiceSupervisors := make(map[string]*ServiceSupervisor)
		for serviceName, serviceSupervisor := range s.serviceSupervisors {
			if _, ok := s.serviceNames[serviceName]; !ok {
				danglingServiceSupervisors[serviceName] = serviceSupervisor
			}
		}
		s.lock.RUnlock()

		for serviceName, serviceSupervisor := range danglingServiceSupervisors {
			serviceSupervisor.Stop()
			s.lock.Lock()
			delete(s.serviceSupervisors, serviceName)
			s.lock.Unlock()
		}

		select {
		case <-s.ctx.Done():
			s.serviceSupervisorGCDone <- struct{}{}
			return
		case <-ticker.C:
			continue
		}
	}
}

func (s *ApplicationSupervisor) containerGC() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		instances := utils.ContainerList(s.ctx, s.engine, map[string]struct{}{
			models.ServiceLabel: struct{}{},
		}, map[string]string{
			models.ApplicationLabel: s.applicationID,
		}, true)

		s.lock.RLock()
		for _, instance := range instances {
			serviceName := instance.Labels[models.ServiceLabel]
			if _, ok := s.serviceSupervisors[serviceName]; !ok {
				// TODO: this could start many goroutines
				go func(instanceID string) {
					utils.ContainerStop(s.ctx, s.engine, instanceID)
					utils.ContainerRemove(s.ctx, s.engine, instanceID)
				}(instance.ID)
			}
		}
		s.lock.RUnlock()

		select {
		case <-s.ctx.Done():
			s.containerGCDone <- struct{}{}
			return
		case <-ticker.C:
			continue
		}
	}
}
