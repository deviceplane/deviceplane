package spec

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/deviceplane/deviceplane/pkg/models"
)

type Application struct {
	Containers map[string]Container `yaml:"containers"`
}

type Container struct {
	Name       string            `yaml:"name"`
	Image      string            `yaml:"image"`
	Entrypoint []string          `yaml:"entrypoint"`
	Command    []string          `yaml:"command"`
	Labels     map[string]string `yaml:"labels"`
}

func (c Container) WithHash() Container {
	// TODO
	if c.Labels == nil {
		c.Labels = make(map[string]string)
	}
	c.Labels[models.HashLabel] = c.Hash()
	return c
}

func (c Container) Hash() string {
	var parts []string

	parts = append(parts, c.Name)
	parts = append(parts, c.Image)
	parts = append(parts, c.Entrypoint...)
	parts = append(parts, c.Command...)
	// TODO: labels

	return hash(parts)
}

func hash(parts []string) string {
	sum := sha256.Sum256([]byte(strings.Join(parts, "")))
	return fmt.Sprintf("%x", sum)
}
