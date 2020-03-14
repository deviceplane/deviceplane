package supervisor

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/deviceplane/deviceplane/pkg/spec"

	"github.com/apex/log"

	"github.com/deviceplane/deviceplane/pkg/agent/validator"
	"github.com/deviceplane/deviceplane/pkg/agent/variables"
	"github.com/deviceplane/deviceplane/pkg/engine"
	"github.com/deviceplane/deviceplane/pkg/hash"
	"github.com/deviceplane/deviceplane/pkg/models"
)

const (
	deviceIDEnvironmentVariableKey   = "DEVICEPLANE_DEVICE_ID"
	deviceNameEnvironmentVariableKey = "DEVICEPLANE_DEVICE_NAME"
)

type ServiceSupervisor struct {
	applicationID string
	serviceName   string
	engine        engine.Engine
	reporter      *Reporter
	validators    []validator.Validator

	imagePuller *imagePuller

	bundle              models.Bundle
	release             string
	service             models.Service
	keepAliveRelease    chan string
	keepAliveService    chan models.Service
	keepAliveDeactivate chan struct{}
	reconcileLoopDone   chan struct{}
	keepAliveDone       chan struct{}

	containerID atomic.Value

	once   sync.Once
	lock   sync.RWMutex
	ctx    context.Context
	cancel func()
}

func NewServiceSupervisor(
	applicationID string,
	serviceName string,
	engine engine.Engine,
	variables variables.Interface,
	reporter *Reporter,
	validators []validator.Validator,
) *ServiceSupervisor {
	ctx, cancel := context.WithCancel(context.Background())
	return &ServiceSupervisor{
		applicationID: applicationID,
		serviceName:   serviceName,
		engine:        engine,
		reporter:      reporter,
		validators:    validators,

		imagePuller: newImagePuller(applicationID, serviceName, engine, variables),

		keepAliveRelease:    make(chan string),
		keepAliveService:    make(chan models.Service),
		keepAliveDeactivate: make(chan struct{}),
		reconcileLoopDone:   make(chan struct{}),
		keepAliveDone:       make(chan struct{}),

		ctx:    ctx,
		cancel: cancel,
	}
}

func (s *ServiceSupervisor) Set(bundle models.Bundle, release string, service models.Service) {
	s.lock.Lock()
	s.bundle = bundle
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
	ticker := time.NewTicker(defaultTickerFrequency)
	defer ticker.Stop()

	for {
		s.reconcile()

		select {
		case <-s.ctx.Done():
			s.reconcileLoopDone <- struct{}{}
			return
		case <-ticker.C:
			continue
		}
	}
}

func (s *ServiceSupervisor) reconcile() {
	s.lock.RLock()
	release := s.release
	service := s.service
	s.lock.RUnlock()

	ctx, cancel := context.WithCancel(s.ctx)
	defer cancel()

	startCanceler := func() {
		go func() {
			ticker := time.NewTicker(defaultTickerFrequency)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					s.lock.RLock()
					if spec.Hash(s.service, s.serviceName) != spec.Hash(service, s.serviceName) {
						cancel()
					}
					s.lock.RUnlock()
				}
			}
		}()
	}

	instances, err := containerList(ctx, s.engine, nil, map[string]string{
		models.ApplicationLabel: s.applicationID,
		models.ServiceLabel:     s.serviceName,
	}, true)
	if err != nil {
		return
	}

	if len(instances) > 0 {
		// TODO: filter down to just one instance if we find more
		instance := instances[0]

		if hashLabel, ok := instance.Labels[models.HashLabel]; ok && hashLabel == spec.Hash(service, s.serviceName) {
			s.sendKeepAliveService(service)
			s.sendKeepAliveRelease(release)
			return
		}

		startCanceler()

		s.reporter.SetServiceState(s.serviceName, models.SetDeviceServiceStateRequest{
			State:        models.ServiceStatePullingImage,
			ErrorMessage: "",
		})
		if err = s.imagePuller.Pull(ctx, service.Image); err != nil {
			s.reporter.SetServiceState(s.serviceName, models.SetDeviceServiceStateRequest{
				State:        models.ServiceStatePullingImage,
				ErrorMessage: err.Error(),
			})
			return
		}

		s.sendKeepAliveDeactivate()

		s.reporter.SetServiceState(s.serviceName, models.SetDeviceServiceStateRequest{
			State:        models.ServiceStateStoppingPreviousContainer,
			ErrorMessage: "",
		})
		if err = containerStop(ctx, s.engine, instance.ID); err != nil {
			s.reporter.SetServiceState(s.serviceName, models.SetDeviceServiceStateRequest{
				State:        models.ServiceStateStoppingPreviousContainer,
				ErrorMessage: err.Error(),
			})
			return
		}

		s.reporter.SetServiceState(s.serviceName, models.SetDeviceServiceStateRequest{
			State:        models.ServiceStateRemovingPreviousContainer,
			ErrorMessage: "",
		})
		if err = containerRemove(ctx, s.engine, instance.ID); err != nil {
			s.reporter.SetServiceState(s.serviceName, models.SetDeviceServiceStateRequest{
				State:        models.ServiceStateRemovingPreviousContainer,
				ErrorMessage: err.Error(),
			})
			return
		}
	} else {
		startCanceler()

		s.reporter.SetServiceState(s.serviceName, models.SetDeviceServiceStateRequest{
			State:        models.ServiceStatePullingImage,
			ErrorMessage: "",
		})
		if err = s.imagePuller.Pull(ctx, service.Image); err != nil {
			s.reporter.SetServiceState(s.serviceName, models.SetDeviceServiceStateRequest{
				State:        models.ServiceStatePullingImage,
				ErrorMessage: err.Error(),
			})
			return
		}
	}

	s.sendKeepAliveDeactivate()

	for _, v := range s.validators {
		err := v.Validate(s.service)
		if err != nil {
			log.WithField("service", s.serviceName).
				WithField("validator", v.Name()).
				WithError(err).
				Error("validation failed")
			return
		}
	}

	s.reporter.SetServiceState(s.serviceName, models.SetDeviceServiceStateRequest{
		State:        models.ServiceStateCreatingContainer,
		ErrorMessage: "",
	})
	if _, err = containerCreate(
		ctx,
		s.engine,
		strings.Join([]string{s.serviceName, hash.ShortHash(s.applicationID), spec.ShortHash(service, s.serviceName)}, "-"),
		s.transformService(spec.WithStandardLabels(service, s.applicationID, s.serviceName)),
	); err != nil {
		s.reporter.SetServiceState(s.serviceName, models.SetDeviceServiceStateRequest{
			State:        models.ServiceStateCreatingContainer,
			ErrorMessage: err.Error(),
		})
		return
	}

	s.sendKeepAliveService(service)
	s.sendKeepAliveRelease(release)
}

func (s *ServiceSupervisor) transformService(service models.Service) models.Service {
	service.Environment = append(
		service.Environment,
		fmt.Sprintf("%s=%s", deviceIDEnvironmentVariableKey, s.bundle.DeviceID),
		fmt.Sprintf("%s=%s", deviceNameEnvironmentVariableKey, s.bundle.DeviceName),
	)
	for key, val := range s.bundle.EnvironmentVariables {
		service.Environment = append(
			service.Environment,
			fmt.Sprintf("%s=%s", key, val),
		)
	}
	return service
}

func (s *ServiceSupervisor) sendKeepAliveRelease(release string) {
	select {
	case <-s.ctx.Done():
		break
	default:
		s.keepAliveRelease <- release
	}
}

func (s *ServiceSupervisor) sendKeepAliveService(service models.Service) {
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
	var service models.Service

	ticker := time.NewTicker(defaultTickerFrequency)
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
				s.containerID.Store("")
				continue
			}

			instances, err := containerList(s.ctx, s.engine, nil, map[string]string{
				models.ApplicationLabel: s.applicationID,
				models.ServiceLabel:     s.serviceName,
				models.HashLabel:        spec.Hash(service, s.serviceName),
			}, true)
			if err != nil {
				continue
			}

			if len(instances) == 0 {
				active = false
				continue
			}

			// TODO: filter down to just one instance if we find more
			instance := instances[0]

			if instance.State == models.ServiceStateRunning {
				s.reporter.SetServiceState(s.serviceName, models.SetDeviceServiceStateRequest{
					State:        models.ServiceStateRunning,
					ErrorMessage: "",
				})
				s.reporter.SetServiceStatus(s.serviceName, models.SetDeviceServiceStatusRequest{
					CurrentReleaseID: release,
				})
				s.containerID.Store(instance.ID)
			} else {
				inspectResponse, err := s.engine.InspectContainer(s.ctx, instance.ID)
				s.reporter.SetServiceState(s.serviceName, models.SetDeviceServiceStateRequest{
					State: instance.State,
					ErrorMessage: func() string {
						if err != nil {
							return "unknown error, cannot inspect container"
						}
						if inspectResponse.ExitCode != nil {
							message := fmt.Sprintf(
								"container exited with exit code %d",
								*inspectResponse.ExitCode,
							)
							if inspectResponse.Error == "" {
								return message
							}
							return fmt.Sprintf("%s (error: %s)",
								message,
								inspectResponse.Error,
							)
						}
						return ""
					}(),
				})

				containerStart(s.ctx, s.engine, instance.ID)
			}
		}
	}
}
