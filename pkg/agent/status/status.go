package status

import (
	"context"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/deviceplane/deviceplane/pkg/spec"
	"gopkg.in/yaml.v2"
)

type GarbageCollector struct {
	deleteApplicationStatus func(ctx context.Context, applicationID string) error
	deleteServiceStatus     func(ctx context.Context, applicationID, service string) error

	bundle                                models.Bundle
	applicationStatusGarbageCollectorDone chan struct{}
	serviceStatusGarbageCollectorDone     chan struct{}

	once   sync.Once
	lock   sync.RWMutex
	ctx    context.Context
	cancel func()
}

func NewGarbageCollector(
	deleteApplicationStatus func(ctx context.Context, applicationID string) error,
	deleteServiceStatus func(ctx context.Context, applicationID, service string) error,
) *GarbageCollector {
	ctx, cancel := context.WithCancel(context.Background())
	return &GarbageCollector{
		deleteApplicationStatus: deleteApplicationStatus,
		deleteServiceStatus:     deleteServiceStatus,

		applicationStatusGarbageCollectorDone: make(chan struct{}),
		serviceStatusGarbageCollectorDone:     make(chan struct{}),

		ctx:    ctx,
		cancel: cancel,
	}
}

func (gc *GarbageCollector) SetBundle(bundle models.Bundle) {
	gc.lock.Lock()
	gc.bundle = bundle
	gc.lock.Unlock()

	gc.once.Do(func() {
		go gc.applicationStatusGarbageCollector()
		go gc.serviceStatusGarbageCollector()
	})
}

func (gc *GarbageCollector) Stop() {
	gc.cancel()
	// TODO: don't do this if SetBundle was never called
	<-gc.applicationStatusGarbageCollectorDone
	<-gc.serviceStatusGarbageCollectorDone
}

func (gc *GarbageCollector) applicationStatusGarbageCollector() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		gc.lock.RLock()
		bundle := gc.bundle
		gc.lock.RUnlock()

		applications := make(map[string]struct{})
		for _, application := range bundle.Applications {
			applications[application.Application.ID] = struct{}{}
		}

		for _, applicationStatus := range bundle.ApplicationStatuses {
			if _, ok := applications[applicationStatus.ApplicationID]; !ok {
				if err := gc.deleteApplicationStatus(gc.ctx, applicationStatus.ApplicationID); err != nil {
					log.WithField("application", applicationStatus.ApplicationID).
						WithError(err).
						Error("delete application status")
				}
			}
		}

		select {
		case <-gc.ctx.Done():
			gc.applicationStatusGarbageCollectorDone <- struct{}{}
			return
		case <-ticker.C:
			continue
		}
	}
}

func (gc *GarbageCollector) serviceStatusGarbageCollector() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		gc.lock.RLock()
		bundle := gc.bundle
		gc.lock.RUnlock()

		services := make(map[string]map[string]struct{})
		for _, application := range bundle.Applications {
			var applicationConfig map[string]spec.Service
			if err := yaml.Unmarshal([]byte(application.LatestRelease.Config), &applicationConfig); err != nil {
				log.WithError(err).Error("unmarshal")
				continue
			}

			services[application.Application.ID] = make(map[string]struct{})
			for serviceName := range applicationConfig {
				services[application.Application.ID][serviceName] = struct{}{}
			}
		}

		serviceStatuses := make(map[string]map[string]struct{})
		for _, serviceStatus := range bundle.ServiceStatuses {
			if _, ok := serviceStatuses[serviceStatus.ApplicationID]; !ok {
				serviceStatuses[serviceStatus.ApplicationID] = make(map[string]struct{})
			}
			serviceStatuses[serviceStatus.ApplicationID][serviceStatus.Service] = struct{}{}
		}

		deleteServiceStatus := func(applicationID, service string) {
			if err := gc.deleteServiceStatus(gc.ctx, applicationID, service); err != nil {
				log.WithField("application", applicationID).
					WithField("service", service).
					WithError(err).
					Error("delete service status")
			}
		}

		for applicationID, serviceStatuses := range serviceStatuses {
			if services, ok := services[applicationID]; ok {
				for service := range serviceStatuses {
					if _, ok = services[service]; !ok {
						deleteServiceStatus(applicationID, service)
					}
				}
			} else {
				for service := range serviceStatuses {
					deleteServiceStatus(applicationID, service)
				}
			}
		}

		select {
		case <-gc.ctx.Done():
			gc.serviceStatusGarbageCollectorDone <- struct{}{}
			return
		case <-ticker.C:
			continue
		}
	}
}
