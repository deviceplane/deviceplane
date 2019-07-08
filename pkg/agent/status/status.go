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
	deleteServiceStatus func(ctx context.Context, applicationID, service string) error

	bundle                            models.Bundle
	serviceStatusGarbageCollectorDone chan struct{}

	once   sync.Once
	lock   sync.RWMutex
	ctx    context.Context
	cancel func()
}

func NewGarbageCollector(
	deleteServiceStatus func(ctx context.Context, applicationID, service string) error,
) *GarbageCollector {
	ctx, cancel := context.WithCancel(context.Background())
	return &GarbageCollector{
		deleteServiceStatus: deleteServiceStatus,

		serviceStatusGarbageCollectorDone: make(chan struct{}),

		ctx:    ctx,
		cancel: cancel,
	}
}

func (gc *GarbageCollector) SetBundle(bundle models.Bundle) {
	gc.lock.Lock()
	gc.bundle = bundle
	gc.lock.Unlock()

	gc.once.Do(func() {
		go gc.serviceStatusGarbageCollector()
	})
}

func (gc *GarbageCollector) Stop() {
	gc.cancel()
	// TODO: don't do this if SetBundle was never called
	<-gc.serviceStatusGarbageCollectorDone
}

func (gc *GarbageCollector) serviceStatusGarbageCollector() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		gc.lock.RLock()
		bundle := gc.bundle
		gc.lock.RUnlock()

		services := make(map[string]map[string]struct{})
		serviceStatuses := make(map[string]map[string]struct{})
		for _, application := range bundle.Applications {
			var applicationConfig map[string]spec.Service
			if err := yaml.Unmarshal([]byte(application.LatestRelease.Config), &applicationConfig); err != nil {
				log.WithError(err).Error("unmarshal")
				continue
			}

			services[application.Application.ID] = make(map[string]struct{})
			serviceStatuses[application.Application.ID] = make(map[string]struct{})

			for serviceName := range applicationConfig {
				services[application.Application.ID][serviceName] = struct{}{}
			}

			for _, serviceStatus := range application.ServiceStatuses {
				serviceStatuses[application.Application.ID][serviceStatus.Service] = struct{}{}
			}
		}

		for applicationID, serviceStatuses := range serviceStatuses {
			for service := range serviceStatuses {
				if _, ok := services[applicationID][service]; !ok {
					if err := gc.deleteServiceStatus(gc.ctx, applicationID, service); err != nil {
						log.WithField("service", service).
							WithField("application", applicationID).
							WithError(err).
							Error("delete service status")
					}
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
