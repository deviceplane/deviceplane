package supervisor

type Lookup interface {
	GetContainerID(applicationID string, service string) *string
}

var _ Lookup = &Supervisor{}

func (s *Supervisor) GetContainerID(applicationID, service string) *string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	applicationSupervisor, ok := s.applicationSupervisors[applicationID]
	if !ok {
		return nil
	}

	return applicationSupervisor.GetContainerID(service)
}

func (s *ApplicationSupervisor) GetContainerID(service string) *string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	serviceSupervisor, ok := s.serviceSupervisors[service]
	if !ok {
		return nil
	}

	return serviceSupervisor.GetContainerID()
}

func (s *ServiceSupervisor) GetContainerID() *string {
	value := s.containerID.Load()
	if value == nil {
		return nil
	}
	containerID, ok := value.(string)
	if !ok || containerID == "" {
		return nil
	}
	return &containerID
}
