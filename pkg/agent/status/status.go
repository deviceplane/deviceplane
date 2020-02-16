package status

import (
	"context"
	"sync"
	"time"

	"github.com/apex/log"
	dpcontext "github.com/deviceplane/deviceplane/pkg/context"
	"github.com/deviceplane/deviceplane/pkg/models"
)

type GarbageCollector struct {
	deleteApplicationStatus func(ctx *dpcontext.Context, applicationID string) error
	deleteServiceStatus     func(ctx *dpcontext.Context, applicationID, service string) error
	deleteServiceState      func(ctx *dpcontext.Context, applicationID, service string) error

	bundle                                models.Bundle
	applicationStatusGarbageCollectorDone chan struct{}
	serviceStatusGarbageCollectorDone     chan struct{}
	serviceStateGarbageCollectorDone      chan struct{}

	once   sync.Once
	lock   sync.RWMutex
	ctx    context.Context
	cancel func()
}

func NewGarbageCollector(
	deleteApplicationStatus func(ctx *dpcontext.Context, applicationID string) error,
	deleteServiceStatus func(ctx *dpcontext.Context, applicationID, service string) error,
	deleteServiceState func(ctx *dpcontext.Context, applicationID, service string) error,
) *GarbageCollector {
	ctx, cancel := context.WithCancel(context.Background())
	return &GarbageCollector{
		deleteApplicationStatus: deleteApplicationStatus,
		deleteServiceStatus:     deleteServiceStatus,
		deleteServiceState:      deleteServiceState,

		applicationStatusGarbageCollectorDone: make(chan struct{}),
		serviceStatusGarbageCollectorDone:     make(chan struct{}),
		serviceStateGarbageCollectorDone:      make(chan struct{}),

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
		go gc.serviceStateGarbageCollector()
	})
}

func (gc *GarbageCollector) Stop() {
	gc.cancel()
	// TODO: don't do this if SetBundle was never called
	<-gc.applicationStatusGarbageCollectorDone
	<-gc.serviceStatusGarbageCollectorDone
	<-gc.serviceStateGarbageCollectorDone
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
				ctx, cancel := dpcontext.New(gc.ctx, time.Minute)

				if err := gc.deleteApplicationStatus(ctx, applicationStatus.ApplicationID); err != nil {
					log.WithField("application", applicationStatus.ApplicationID).
						WithError(err).
						Error("delete application status")
				}

				cancel()
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
			services[application.Application.ID] = make(map[string]struct{})
			for serviceName := range application.LatestRelease.Config {
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
			ctx, cancel := dpcontext.New(gc.ctx, time.Minute)

			if err := gc.deleteServiceStatus(ctx, applicationID, service); err != nil {
				log.WithField("application", applicationID).
					WithField("service", service).
					WithError(err).
					Error("delete service status")
			}

			cancel()
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

func (gc *GarbageCollector) serviceStateGarbageCollector() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		gc.lock.RLock()
		bundle := gc.bundle
		gc.lock.RUnlock()

		services := make(map[string]map[string]struct{})
		for _, application := range bundle.Applications {
			services[application.Application.ID] = make(map[string]struct{})
			for serviceName := range application.LatestRelease.Config {
				services[application.Application.ID][serviceName] = struct{}{}
			}
		}

		serviceStates := make(map[string]map[string]struct{})
		for _, serviceState := range bundle.ServiceStates {
			if _, ok := serviceStates[serviceState.ApplicationID]; !ok {
				serviceStates[serviceState.ApplicationID] = make(map[string]struct{})
			}
			serviceStates[serviceState.ApplicationID][serviceState.Service] = struct{}{}
		}

		deleteServiceState := func(applicationID, service string) {
			ctx, cancel := dpcontext.New(gc.ctx, time.Minute)

			if err := gc.deleteServiceState(ctx, applicationID, service); err != nil {
				log.WithField("application", applicationID).
					WithField("service", service).
					WithError(err).
					Error("delete service state")
			}

			cancel()
		}

		for applicationID, serviceStates := range serviceStates {
			if services, ok := services[applicationID]; ok {
				for service := range serviceStates {
					if _, ok = services[service]; !ok {
						deleteServiceState(applicationID, service)
					}
				}
			} else {
				for service := range serviceStates {
					deleteServiceState(applicationID, service)
				}
			}
		}

		select {
		case <-gc.ctx.Done():
			gc.serviceStateGarbageCollectorDone <- struct{}{}
			return
		case <-ticker.C:
			continue
		}
	}
}
