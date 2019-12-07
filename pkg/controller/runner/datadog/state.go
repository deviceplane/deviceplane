package datadog

import (
	"context"

	"github.com/deviceplane/deviceplane/pkg/metrics/datadog"
	"github.com/deviceplane/deviceplane/pkg/models"
)

func (r *Runner) getStateMetrics(ctx context.Context, project *models.Project, device *models.Device, metricConfig *models.MetricTargetConfig) datadog.Series {
	stateMetrics := []datadog.Metric{
		datadog.Metric{
			Metric: "devices",
			Points: [][2]interface{}{
				datadog.NewPoint(1),
			},
			Type: "count",
			Tags: []string{"status:" + string(device.Status)},
		},
	}

	config := metricConfig.Configs[0]
	return FilterMetrics(project, nil, device, metricConfig.Type, config, stateMetrics)
}
