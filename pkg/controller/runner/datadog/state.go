package datadog

import (
	"github.com/deviceplane/deviceplane/pkg/metrics/datadog"
	"github.com/deviceplane/deviceplane/pkg/models"
)

func (r *Runner) getProjectMetrics(
	project *models.Project,
	device *models.Device,
) models.DatadogSeries {
	return []models.DatadogMetric{
		models.DatadogMetric{
			Metric: "devices",
			Points: [][2]interface{}{
				datadog.NewPoint(1),
			},
			Type: "count",
			Tags: []string{"deviceplane.status:" + string(device.Status)},
		},
	}
}
