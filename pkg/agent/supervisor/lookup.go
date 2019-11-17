package supervisor

type Lookup interface {
	GetContainerID(applicationID string, service string) (string, bool)
	GetImagePullProgress(applicationID string, service string) (map[string]PullEvent, bool)
}

var _ Lookup = &Supervisor{}

func (s *Supervisor) GetContainerID(applicationID, service string) (string, bool) {
	var containerID string
	var ok bool
	s.withServiceSupervisor(applicationID, service, func(s *ServiceSupervisor) {
		value := s.containerID.Load()
		if value == nil {
			ok = false
			return
		}
		containerID, ok = value.(string)
		if containerID == "" {
			ok = false
		}
	})
	return containerID, ok
}

func (s *Supervisor) GetImagePullProgress(applicationID, service string) (map[string]PullEvent, bool) {
	var progress map[string]PullEvent
	var ok bool
	s.withServiceSupervisor(applicationID, service, func(s *ServiceSupervisor) {
		progress, ok = s.imagePuller.Progress()
	})
	return progress, ok
}

func (s *Supervisor) withServiceSupervisor(
	applicationID, service string,
	f func(*ServiceSupervisor),
) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	applicationSupervisor, ok := s.applicationSupervisors[applicationID]
	if !ok {
		return
	}

	applicationSupervisor.withServiceSupervisor(service, f)
}

func (s *ApplicationSupervisor) withServiceSupervisor(
	service string,
	f func(*ServiceSupervisor),
) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	serviceSupervisor, ok := s.serviceSupervisors[service]
	if !ok {
		return
	}

	f(serviceSupervisor)
}
