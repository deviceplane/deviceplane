package spec

import (
	"strings"

	"github.com/deviceplane/deviceplane/pkg/hash"
	"github.com/deviceplane/deviceplane/pkg/models"
)

type Service struct {
	Name       string            `yaml:"name,omitempty"`
	Image      string            `yaml:"image,omitempty"`
	Entrypoint []string          `yaml:"entrypoint,omitempty"`
	Command    []string          `yaml:"command,omitempty"`
	Labels     map[string]string `yaml:"labels,omitempty"`
	Scheduling string            `yaml:"scheduling,omitempty"`
}

func (s Service) WithStandardLabels(applicationID, serviceName string) Service {
	// TODO
	if s.Labels == nil {
		s.Labels = make(map[string]string)
	}
	s.Labels[models.ApplicationLabel] = applicationID
	s.Labels[models.ServiceLabel] = serviceName
	s.Labels[models.HashLabel] = s.Hash()
	return s
}

func (s Service) Hash() string {
	return s.hash(hash.Hash)
}

func (s Service) ShortHash() string {
	return s.hash(hash.ShortHash)
}

func (s Service) hash(hash func(string) string) string {
	var parts []string

	parts = append(parts, s.Name)
	parts = append(parts, s.Image)
	parts = append(parts, s.Entrypoint...)
	parts = append(parts, s.Command...)
	// TODO: labels

	return hash(strings.Join(parts, ""))
}
