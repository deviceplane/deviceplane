package validator

import (
	"github.com/deviceplane/deviceplane/pkg/models"
)

type Validator interface {
	Validate(models.Service) error
	Name() string
}
