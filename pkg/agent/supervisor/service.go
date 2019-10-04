package supervisor

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/deviceplane/deviceplane/pkg/agent/utils"
	"github.com/deviceplane/deviceplane/pkg/engine"
	"github.com/deviceplane/deviceplane/pkg/hash"
	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/deviceplane/deviceplane/pkg/spec"
)

type ServiceSupervisor struct {
	applicationID string
	serviceName   string
	engine        engine.Engine
	reporter      *Reporter

	release             string
	service             spec.Service
	keepAliveRelease    chan string
	keepAliveService    chan spec.Service
	keepAliveDeactivate chan struct{}
	reconcileLoopDone   chan struct{}
	keepAliveDone       chan struct{}

	once   sync.Once
	lock   sync.RWMutex
	ctx    context.Context
	cancel func()
}

func NewServiceSupervisor(
	applicationID string,
	serviceName string,
	engine engine.Engine,
	reporter *Reporter,
) *ServiceSupervisor {
	ctx, cancel := context.WithCancel(context.Background())
	return &ServiceSupervisor{
		applicationID: applicationID,
		serviceName:   serviceName,
		engine:        engine,
		reporter:      reporter,

		keepAliveRelease:    make(chan string),
		keepAliveService:    make(chan spec.Service),
		keepAliveDeactivate: make(chan struct{}),
		reconcileLoopDone:   make(chan struct{}),
		keepAliveDone:       make(chan struct{}),

		ctx:    ctx,
		cancel: cancel,
	}
}

func (s *ServiceSupervisor) SetService(release string, service spec.Service) {
	s.lock.Lock()
	s.release = release
	s.service = service
	s.lock.Unlock()

	s.once.Do(func() {
		go s.reconcileLoop()
		go s.keepAlive()
	})
}

func (s *ServiceSupervisor) Stop() {
	s.cancel()
	// TODO: don't do this if SetService was never called
	<-s.reconcileLoopDone
	<-s.keepAliveDone
}

func (s *ServiceSupervisor) reconcileLoop() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		s.lock.RLock()
		release := s.release
		service := s.service
		s.lock.RUnlock()

		ctx, cancel := context.WithCancel(s.ctx)

		startCanceler := func() {
			go func() {
				ticker := time.NewTicker(time.Second)
				defer ticker.Stop()

				for {
					select {
					case <-ctx.Done():
						return
					case <-ticker.C:
						s.lock.RLock()
						if s.service.Hash(s.serviceName) != service.Hash(s.serviceName) {
							cancel()
						}
						s.lock.RUnlock()
					}
				}
			}()
		}

		instances := utils.ContainerList(ctx, s.engine, nil, map[string]string{
			models.ApplicationLabel: s.applicationID,
			models.ServiceLabel:     s.serviceName,
		}, true)

		if len(instances) > 0 {
			// TODO: filter down to just one instance if we find more
			instance := instances[0]

			if hashLabel, ok := instance.Labels[models.HashLabel]; ok && hashLabel == service.Hash(s.serviceName) {
				s.sendKeepAliveService(service)
				s.sendKeepAliveRelease(release)
				goto cont
			}

			startCanceler()
			utils.ImagePull(ctx, s.engine, service.Image)

			s.sendKeepAliveDeactivate()

			utils.ContainerStop(ctx, s.engine, instance.ID)
			utils.ContainerRemove(ctx, s.engine, instance.ID)
		} else {
			startCanceler()
			utils.ImagePull(ctx, s.engine, service.Image)
		}

		s.sendKeepAliveDeactivate()

		utils.ContainerCreate(
			ctx,
			s.engine,
			strings.Join([]string{s.serviceName, hash.ShortHash(s.applicationID), service.ShortHash(s.serviceName)}, "-"),
			service.WithStandardLabels(s.applicationID, s.serviceName),
		)

		s.sendKeepAliveService(service)
		s.sendKeepAliveRelease(release)

		cancel()

	cont:
		select {
		case <-s.ctx.Done():
			s.reconcileLoopDone <- struct{}{}
			return
		case <-ticker.C:
			continue
		}
	}
}

func (s *ServiceSupervisor) sendKeepAliveRelease(release string) {
	select {
	case <-s.ctx.Done():
		break
	default:
		s.keepAliveRelease <- release
	}
}

func (s *ServiceSupervisor) sendKeepAliveService(service spec.Service) {
	select {
	case <-s.ctx.Done():
		break
	default:
		s.keepAliveService <- service
	}
}

func (s *ServiceSupervisor) sendKeepAliveDeactivate() {
	select {
	case <-s.ctx.Done():
		break
	default:
		s.keepAliveDeactivate <- struct{}{}
	}
}

func (s *ServiceSupervisor) keepAlive() {
	active := false
	var release string
	var service spec.Service

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			s.keepAliveDone <- struct{}{}
			return
		case release = <-s.keepAliveRelease:
			break
		case service = <-s.keepAliveService:
			active = true
		case <-s.keepAliveDeactivate:
			active = false
		case <-ticker.C:
			if !active {
				continue
			}

			instances := utils.ContainerList(s.ctx, s.engine, nil, map[string]string{
				models.ApplicationLabel: s.applicationID,
				models.ServiceLabel:     s.serviceName,
				models.HashLabel:        service.Hash(s.serviceName),
			}, true)

			if len(instances) == 0 {
				active = false
				continue
			}

			s.reporter.SetServiceRelease(s.serviceName, release)

			// TODO: filter down to just one instance if we find more
			instance := instances[0]

			if !instance.Running {
				utils.ContainerStart(s.ctx, s.engine, instance.ID)
			}
		}
	}
}
