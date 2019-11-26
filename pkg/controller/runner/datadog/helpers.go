package datadog

import (
	"github.com/deviceplane/deviceplane/pkg/models"
)

func addedInternalTags(project *models.Project) []string {
	return []string{
		"project:" + project.Name,
	}
}
