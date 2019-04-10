package engine

import (
	"github.com/deviceplane/deviceplane/pkg/spec"
)

type Instance struct {
	ID string
	spec.Service
}
