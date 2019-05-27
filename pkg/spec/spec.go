package spec

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/deviceplane/deviceplane/pkg/models"
)

type Application struct {
	Services map[string]Service `yaml:"services,omitempty"`
}

type Service struct {
	Name       string            `yaml:"name,omitempty"`
	Image      string            `yaml:"image,omitempty"`
	Entrypoint []string          `yaml:"entrypoint,omitempty"`
	Command    []string          `yaml:"command,omitempty"`
	Labels     map[string]string `yaml:"labels,omitempty"`
	Scheduling string            `yaml:"scheduling,omitempty"`
}

func (s Service) WithStandardLabels(serviceName string) Service {
	// TODO
	if s.Labels == nil {
		s.Labels = make(map[string]string)
	}
	s.Labels[models.ServiceLabel] = serviceName
	s.Labels[models.HashLabel] = s.Hash()
	return s
}

func (s Service) Hash() string {
	var parts []string

	parts = append(parts, s.Name)
	parts = append(parts, s.Image)
	parts = append(parts, s.Entrypoint...)
	parts = append(parts, s.Command...)
	// TODO: labels

	return hash(parts)
}

func hash(parts []string) string {
	sum := sha256.Sum256([]byte(strings.Join(parts, "")))
	return fmt.Sprintf("%x", sum)
}
