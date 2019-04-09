package scheduler

import (
	"errors"
	"strings"

	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/deviceplane/deviceplane/pkg/spec"
)

var (
	errInvalidScheduling = errors.New("invalid scheduling")
)

func TransformSpec(application spec.Application, deviceLabels []models.DeviceLabel) (spec.Application, error) {
	deviceLabelsMap := make(map[string]string)
	for _, deviceLabel := range deviceLabels {
		deviceLabelsMap[deviceLabel.Key] = deviceLabel.Value
	}

	var transformedApplication spec.Application
	transformedApplication.Services = make(map[string]spec.Service)
	for serviceName, service := range application.Services {
		if service.Scheduling == "" {
			transformedApplication.Services[serviceName] = service
			continue
		}

		schedulingSplit := strings.Split(service.Scheduling, "=")
		if len(schedulingSplit) != 2 {
			return spec.Application{}, errInvalidScheduling
		}

		key := schedulingSplit[0]
		value := schedulingSplit[1]

		if deviceLabelsMap[key] == value {
			transformedApplication.Services[serviceName] = service
		}
	}

	return transformedApplication, nil
}
