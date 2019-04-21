package engine

import (
	"context"
	"errors"

	"github.com/deviceplane/deviceplane/pkg/spec"
)

var (
	ErrInstanceNotFound = errors.New("instance not found")
)

type Engine interface {
	Create(context.Context, string, spec.Service) (string, error)
	Start(context.Context, string) error
	List(context.Context, map[string]bool, map[string]string, bool) ([]Instance, error)
	Stop(context.Context, string) error
	Remove(context.Context, string) error
}

type Instance struct {
	ID      string
	Labels  map[string]string
	Running bool
}
