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
	serviceName *string,
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
		metricPrefix = fmt.Sprintf("deviceplane.user_defined.%s.%s", app.Name, *serviceName)
	case models.MetricStateTargetType:
		metricPrefix = "deviceplane"
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
				addTag("deviceplane.labels."+label, labelValue)
			}
		}

		// Optional properties
		for _, tag := range correspondingConfig.Tags {
			switch tag {
			case "device":
				addTag("deviceplane.device", device.Name)
			}
		}

		// Guaranteed tags
		addTag("deviceplane.project", project.Name)

		passedMetrics = append(passedMetrics, m)
	}

	return passedMetrics
}
