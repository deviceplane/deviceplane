package datadog

import (
	"fmt"

	"github.com/deviceplane/deviceplane/pkg/metrics/datadog"
	"github.com/deviceplane/deviceplane/pkg/models"
)

const WildcardMetric = string("*")

func FilterMetrics(
	project *models.Project,
	app *models.Application,
	device *models.Device,
	targetType models.MetricTargetType,
	config models.MetricConfig,
	metrics []datadog.Metric,
) (passedMetrics []datadog.Metric) {
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

	allowedMetrics := make(map[string]*models.Metric, len(config.Metrics))

	for i, m := range config.Metrics {
		allowedMetrics[m.Metric] = &config.Metrics[i]
	}

	for _, m := range metrics {
		correspondingConfig := allowedMetrics[m.Metric]
		if correspondingConfig == nil {
			correspondingConfig = allowedMetrics[WildcardMetric]
		}

		if correspondingConfig == nil {
			continue
		}

		// Prefix metric name
		m.Metric = fmt.Sprintf("%s.%s", metricPrefix, m.Metric)

		// Helper
		addTag := func(tag, value string) {
			m.Tags = append(
				m.Tags,
				fmt.Sprintf("%s:%s", tag, value),
			)
		}

		// Optional labels
		for _, label := range correspondingConfig.Labels {
			labelValue, ok := device.Labels[label]
			if ok {
				addTag("label"+"."+label, labelValue)
			}
		}

		// Optional tags
		// implementation could be less manual
		for _, tag := range correspondingConfig.Tags {
			switch tag {
			case "device":
				addTag(tag, device.Name)
			case "application":
				if app == nil {
					continue
				}
				addTag(tag, app.Name)
			}
		}

		// Guaranteed tags
		addTag("project", project.Name)

		passedMetrics = append(passedMetrics, m)
	}

	return passedMetrics
}
