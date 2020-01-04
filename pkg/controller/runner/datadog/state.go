package datadog

import (
	"context"

	"github.com/deviceplane/deviceplane/pkg/metrics/datadog"
	"github.com/deviceplane/deviceplane/pkg/models"
)

func (r *Runner) getProjectMetrics(
	ctx context.Context,
	project *models.Project,
	device *models.Device,
) datadog.Series {
	return []datadog.Metric{
		datadog.Metric{
			Metric: "devices",
			Points: [][2]interface{}{
				datadog.NewPoint(1),
			},
			Type: "count",
			Tags: []string{"deviceplane.status:" + string(device.Status)},
		},
	}
}
