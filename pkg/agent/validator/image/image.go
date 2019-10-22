package image

import (
	"errors"
	"strings"

	"github.com/deviceplane/deviceplane/pkg/agent/variables"
	"github.com/deviceplane/deviceplane/pkg/spec"
)

var (
	ErrNonWhitelistedImage = errors.New("image is not found in the device's non-empty whitelist")
)

type Validator struct {
	variables variables.Interface
}

func NewValidator(variables variables.Interface) *Validator {
	return &Validator{
		variables: variables,
	}
}

func (i *Validator) Validate(s spec.Service) error {
	whitelistedImages := i.variables.GetWhitelistedImages()

	if isValid(s.Image, whitelistedImages) {
		return nil
	}

	return ErrNonWhitelistedImage
}

func (i *Validator) Name() string { return "ImageValidator" }

func isValid(image string, whitelistedImages []string) bool {
	// If the file doesn't exist, or there are no whitelisted, we allow
	// everything
	if len(whitelistedImages) == 0 {
		return true
	}

	for _, wlImage := range whitelistedImages {
		if image == wlImage {
			return true
		}
		if strings.HasPrefix(image, wlImage) {
			return true
		}
	}
	return false
}
