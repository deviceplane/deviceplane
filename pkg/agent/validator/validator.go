package validator

import (
	"github.com/deviceplane/deviceplane/pkg/spec"
)

type Validator interface {
	Validate(spec.Service) error
	Name() string
}
