package supervisor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/engine"
	canonical_image "github.com/deviceplane/deviceplane/pkg/image"
	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/deviceplane/deviceplane/pkg/spec"
	"gopkg.in/yaml.v2"
)

type Supervisor struct {
	engine                  engine.Engine
	reportApplicationStatus func(ctx context.Context, applicationID string, currentReleaseID string) error
	reportServiceStatus     func(ctx context.Context, applicationID, service, currentReleaseID string) error

	latestDesiredApplicationReleases map[string]string
	reportedApplicationReleases      map[string]string

	serviceReleases         map[string]map[string]string
	reportedServiceReleases map[string]map[string]string

	reconcilingServices map[string]map[string]struct{}
	keepAliveShutdowns  map[string]map[string]chan struct{}
	keepAliveAcks       map[string]map[string]chan struct{}

	lock sync.Mutex
}

func NewSupervisor(
	engine engine.Engine,
	reportApplicationStatus func(ctx context.Context, applicationID, currentReleaseID string) error,
	reportServiceStatus func(ctx context.Context, applicationID, service, currentReleaseID string) error,
) *Supervisor {
	supervisor := &Supervisor{
		engine:                  engine,
		reportApplicationStatus: reportApplicationStatus,
		reportServiceStatus:     reportServiceStatus,

		latestDesiredApplicationReleases: make(map[string]string),
		reportedApplicationReleases:      make(map[string]string),

		serviceReleases:         make(map[string]map[string]string),
		reportedServiceReleases: make(map[string]map[string]string),

		reconcilingServices: make(map[string]map[string]struct{}),
		keepAliveShutdowns:  make(map[string]map[string]chan struct{}),
		keepAliveAcks:       make(map[string]map[string]chan struct{}),
	}
	go supervisor.applicationStatusReporter()
	go supervisor.serviceStatusReporter()
	return supervisor
}

func (s *Supervisor) applicationStatusReporter() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		s.lock.Lock()

		completedApplicationReleases := make(map[string]string)
		for applicationID, serviceReleases := range s.serviceReleases {
			desiredApplicationRelease := s.latestDesiredApplicationReleases[applicationID]

			applicationCompleted := true
			for _, releaseID := range serviceReleases {
				if releaseID != desiredApplicationRelease {
					applicationCompleted = false
					break
				}
			}

			if applicationCompleted {
				completedApplicationReleases[applicationID] = desiredApplicationRelease
			}
		}

		diff := make(map[string]string)
		copy := make(map[string]string)
		for applicationID, releaseID := range completedApplicationReleases {
			reportedReleaseID, ok := s.reportedApplicationReleases[applicationID]
			if !ok || reportedReleaseID != releaseID {
				diff[applicationID] = releaseID
			}
			copy[applicationID] = releaseID
		}

		s.lock.Unlock()

		for applicationID, releaseID := range diff {
			retry(func(ctx context.Context) error {
				if err := s.reportApplicationStatus(ctx, applicationID, releaseID); err != nil {
					log.WithError(err).Error("report application status")
					return err
				}
				return nil
			}, time.Minute)
		}

		s.reportedApplicationReleases = copy

		select {
		case <-ticker.C:
			continue
		}
	}
}

func (s *Supervisor) serviceStatusReporter() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		s.lock.Lock()
		diff := make(map[string]map[string]string)
		copy := make(map[string]map[string]string)
		for applicationID, serviceReleases := range s.serviceReleases {
			diff[applicationID] = make(map[string]string)
			copy[applicationID] = make(map[string]string)

			if _, ok := s.reportedServiceReleases[applicationID]; !ok {
				s.reportedServiceReleases[applicationID] = make(map[string]string)
			}

			for service, releaseID := range serviceReleases {
				reportedReleaseID, ok := s.reportedServiceReleases[applicationID][service]
				if !ok || reportedReleaseID != releaseID {
					diff[applicationID][service] = releaseID
				}
				copy[applicationID][service] = releaseID
			}
		}
		s.lock.Unlock()

		for applicationID, applicationDiff := range diff {
			for service, releaseID := range applicationDiff {
				retry(func(ctx context.Context) error {
					if err := s.reportServiceStatus(ctx, applicationID, service, releaseID); err != nil {
						log.WithError(err).Error("report service status")
						return err
					}
					return nil
				}, time.Minute)
			}
		}

		s.reportedServiceReleases = copy

		select {
		case <-ticker.C:
			continue
		}
	}
}

func (s *Supervisor) SetApplication(application models.ApplicationAndLatestRelease) error {
	var applicationConfig map[string]spec.Service
	if err := yaml.Unmarshal([]byte(application.LatestRelease.Config), &applicationConfig); err != nil {
		return err
	}

	s.lock.Lock()

	if _, ok := s.serviceReleases[application.Application.ID]; !ok {
		s.serviceReleases[application.Application.ID] = make(map[string]string)
	}
	if _, ok := s.reportedServiceReleases[application.Application.ID]; !ok {
		s.reportedServiceReleases[application.Application.ID] = make(map[string]string)
	}
	if _, ok := s.reconcilingServices[application.Application.ID]; !ok {
		s.reconcilingServices[application.Application.ID] = make(map[string]struct{})
	}
	if _, ok := s.keepAliveShutdowns[application.Application.ID]; !ok {
		s.keepAliveShutdowns[application.Application.ID] = make(map[string]chan struct{})
	}
	if _, ok := s.keepAliveAcks[application.Application.ID]; !ok {
		s.keepAliveAcks[application.Application.ID] = make(map[string]chan struct{})
	}

	s.latestDesiredApplicationReleases[application.Application.ID] = application.LatestRelease.ID

	s.lock.Unlock()

	for serviceName, service := range applicationConfig {
		go s.reconcile(application.Application.ID, application.LatestRelease.ID, serviceName, service)
	}

	return nil
}

func (s *Supervisor) reconcile(applicationID, releaseID, serviceName string, service spec.Service) {
	s.lock.Lock()
	// If this service is already in reconciling state then exit
	if _, ok := s.reconcilingServices[applicationID][serviceName]; ok {
		s.lock.Unlock()
		return
	}
	// Set this service to reconciling state
	s.reconcilingServices[applicationID][serviceName] = struct{}{}
	s.lock.Unlock()

	defer func() {
		s.lock.Lock()
		delete(s.reconcilingServices[applicationID], serviceName)
		s.lock.Unlock()
	}()

	stopKeepAlive := func() {
		s.lock.Lock()

		shutdown, ok := s.keepAliveShutdowns[applicationID][serviceName]
		if !ok {
			s.lock.Unlock()
			return
		}

		ack := make(chan struct{})
		s.keepAliveAcks[applicationID][serviceName] = ack
		shutdown <- struct{}{}
		<-ack
		delete(s.keepAliveAcks[applicationID], serviceName)

		s.lock.Unlock()
	}

	instances := s.containerList(nil, map[string]string{
		models.ServiceLabel: serviceName,
	}, true)

	if len(instances) > 0 {
		// TODO: filter down to just one instance if we find more
		instance := instances[0]

		if hashLabel, ok := instance.Labels[models.HashLabel]; ok && hashLabel == service.Hash() {
			stopKeepAlive()
			go s.keepAlive(applicationID, releaseID, serviceName, service)
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

	go s.keepAlive(applicationID, releaseID, serviceName, service)
}

func (s *Supervisor) keepAlive(applicationID, releaseID, serviceName string, service spec.Service) {
	s.lock.Lock()

	if _, ok := s.keepAliveShutdowns[applicationID][serviceName]; ok {
		s.lock.Unlock()
		return
	}
	shutdown := make(chan struct{}, 1)
	s.keepAliveShutdowns[applicationID][serviceName] = shutdown

	s.serviceReleases[applicationID][serviceName] = releaseID

	s.lock.Unlock()

	defer func() {
		s.lock.Lock()
		delete(s.keepAliveShutdowns[applicationID], serviceName)
		s.lock.Unlock()
	}()

	dead := false
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-shutdown:
			// TODO: need a lock here?
			s.keepAliveAcks[applicationID][serviceName] <- struct{}{}
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
				go s.reconcile(applicationID, releaseID, serviceName, service)
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
	}, 2*time.Minute)

	return id
}

func (s *Supervisor) containerStart(id string) {
	retry(func(ctx context.Context) error {
		if err := s.engine.StartContainer(ctx, id); err != nil && err != engine.ErrInstanceNotFound {
			log.WithError(err).Error("start container")
			return err
		}
		return nil
	}, 2*time.Minute)
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
	}, 2*time.Minute)

	return instances
}

func (s *Supervisor) containerStop(id string) {
	retry(func(ctx context.Context) error {
		if err := s.engine.StopContainer(ctx, id); err != nil && err != engine.ErrInstanceNotFound {
			log.WithError(err).Error("stop container")
			return err
		}
		return nil
	}, 2*time.Minute)
}

func (s *Supervisor) containerRemove(id string) {
	retry(func(ctx context.Context) error {
		if err := s.engine.RemoveContainer(ctx, id); err != nil && err != engine.ErrInstanceNotFound {
			log.WithError(err).Error("remove container")
			return err
		}
		return nil
	}, 2*time.Minute)
}

func (s *Supervisor) imagePull(image string) {
	retry(func(ctx context.Context) error {
		if err := s.engine.PullImage(ctx, canonical_image.ToCanonical(image)); err != nil {
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
