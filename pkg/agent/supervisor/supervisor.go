package supervisor

import (
	"context"
	"sync"
	"time"

	"github.com/deviceplane/deviceplane/pkg/agent/utils"
	"github.com/deviceplane/deviceplane/pkg/agent/validator"
	"github.com/deviceplane/deviceplane/pkg/engine"
	"github.com/deviceplane/deviceplane/pkg/models"
)

type Supervisor struct {
	engine                  engine.Engine
	reportApplicationStatus func(ctx context.Context, applicationID string, currentReleaseID string) error
	reportServiceStatus     func(ctx context.Context, applicationID, service, currentReleaseID string) error
	validators              []validator.Validator

	applicationIDs         map[string]struct{}
	applicationSupervisors map[string]*ApplicationSupervisor
	once                   sync.Once

	lock   sync.RWMutex
	ctx    context.Context
	cancel func()
}

func NewSupervisor(
	engine engine.Engine,
	reportApplicationStatus func(ctx context.Context, applicationID, currentReleaseID string) error,
	reportServiceStatus func(ctx context.Context, applicationID, service, currentReleaseID string) error,
	validators []validator.Validator,
) *Supervisor {
	ctx, cancel := context.WithCancel(context.Background())
	return &Supervisor{
		engine:                  engine,
		reportApplicationStatus: reportApplicationStatus,
		reportServiceStatus:     reportServiceStatus,
		validators:              validators,

		applicationIDs:         make(map[string]struct{}),
		applicationSupervisors: make(map[string]*ApplicationSupervisor),

		ctx:    ctx,
		cancel: cancel,
	}
}

func (s *Supervisor) SetApplications(applications []models.ApplicationFull2) {
	applicationIDs := make(map[string]struct{})
	for _, application := range applications {
		s.lock.Lock()
		applicationSupervisor, ok := s.applicationSupervisors[application.Application.ID]
		if !ok {
			applicationSupervisor = NewApplicationSupervisor(
				application.Application.ID,
				s.engine,
				NewReporter(application.Application.ID, s.reportApplicationStatus, s.reportServiceStatus),
				s.validators,
			)
			s.applicationSupervisors[application.Application.ID] = applicationSupervisor
		}
		s.lock.Unlock()

		applicationSupervisor.SetApplication(application)

		applicationIDs[application.Application.ID] = struct{}{}
	}

	s.lock.Lock()
	s.applicationIDs = applicationIDs
	s.lock.Unlock()

	s.once.Do(func() {
		go s.applicationSupervisorGC()
		go s.containerGC()
	})
}

func (s *Supervisor) applicationSupervisorGC() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		s.lock.RLock()
		danglingApplicationSupervisors := make(map[string]*ApplicationSupervisor)
		for applicationID, applicationSupervisor := range s.applicationSupervisors {
			if _, ok := s.applicationIDs[applicationID]; !ok {
				danglingApplicationSupervisors[applicationID] = applicationSupervisor
			}
		}
		s.lock.RUnlock()

		for applicationID, applicationSupervisor := range danglingApplicationSupervisors {
			applicationSupervisor.Stop()
			s.lock.Lock()
			delete(s.applicationSupervisors, applicationID)
			s.lock.Unlock()
		}

		select {
		case <-ticker.C:
			continue
		}
	}
}

func (s *Supervisor) containerGC() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		instances := utils.ContainerList(s.ctx, s.engine, map[string]struct{}{
			models.ApplicationLabel: struct{}{},
		}, nil, true)

		s.lock.RLock()
		for _, instance := range instances {
			applicationID := instance.Labels[models.ApplicationLabel]
			if _, ok := s.applicationSupervisors[applicationID]; !ok {
				// TODO: this could start many goroutines
				go func(instanceID string) {
					utils.ContainerStop(s.ctx, s.engine, instanceID)
					utils.ContainerRemove(s.ctx, s.engine, instanceID)
				}(instance.ID)
			}
		}
		s.lock.RUnlock()

		select {
		case <-ticker.C:
			continue
		}
	}
}
