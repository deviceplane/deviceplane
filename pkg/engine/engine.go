package engine

import (
	"context"

	"github.com/deviceplane/deviceplane/pkg/spec"
)

type Engine interface {
	Create(context.Context, spec.Service) (*Instance, error)
	Start(context.Context, string) error
	List(context.Context, map[string]bool, map[string]string) ([]Instance, error)
	Stop(context.Context, string) error
	Remove(context.Context, string) error
}
