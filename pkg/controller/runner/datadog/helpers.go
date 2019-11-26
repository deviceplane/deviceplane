package datadog

import (
	"github.com/deviceplane/deviceplane/pkg/models"
)

func addedInternalTags(project *models.Project, device *models.Device) []string {
	return []string{
		"project:" + project.Name,
		"device:" + device.Name,
	}
}
