package datadog

import (
	"fmt"

	"github.com/deviceplane/deviceplane/pkg/metrics/datadog"
	"github.com/deviceplane/deviceplane/pkg/models"
)

func addedInternalTags(project *models.Project, device *models.Device) []string {
	return []string{
		"project:" + project.Name,
		"device:" + device.Name,
	}
}

func getFilteredMetrics(
	project *models.Project,
	app *models.Application,
	device *models.Device,
	targetType models.MetricTargetType,
	config models.MetricConfig,
	metrics []datadog.Metric,
) []datadog.Metric {
	var metricPrefix string
	switch targetType {
	case models.MetricHostTargetType:
		metricPrefix = "deviceplane.host"
	case models.MetricServiceTargetType:
		metricPrefix = "deviceplane.service"
	case models.MetricStateTargetType:
		metricPrefix = "deviceplane.state"
	default:
		return nil
	}

	returnedMetrics := make([]datadog.Metric, 0)
	returnedMetricsLookup := make(map[string]bool, len(metrics))

	for _, metricConfig := range config.Metrics {
		for _, m := range metrics {
			if returnedMetricsLookup[m.Metric] {
				continue
			}
			if m.Metric == metricConfig.Metric {
				returnedMetricsLookup[m.Metric] = true

				m.Metric = fmt.Sprintf("%s.%s", metricPrefix, m.Metric)

				// Optional labels
				for _, label := range metricConfig.Labels {
					labelValue, ok := device.Labels[label]
					if ok {
						m.Tags = append(
							m.Tags,
							fmt.Sprintf("%s.%s:%s", "label", label, labelValue),
						)
					}
				}

				// Optional tags
				// implementation could be less manual
				for _, tag := range metricConfig.Tags {
					switch tag {
					case "device":
						m.Tags = append(
							m.Tags,
							fmt.Sprintf("%s:%s", tag, device.Name),
						)
					case "application":
						if app == nil {
							continue
						}
						m.Tags = append(
							m.Tags,
							fmt.Sprintf("%s:%s", tag, app.Name),
						)
					}
				}

				// Guaranteed tags
				m.Tags = append(
					m.Tags,
					fmt.Sprintf("%s:%s", "project", project.Name),
				)

				returnedMetrics = append(returnedMetrics, m)
			}
		}
	}

	return returnedMetrics
}
