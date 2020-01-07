package datadog

import (
	"fmt"

	"github.com/deviceplane/deviceplane/pkg/metrics/datadog"
	"github.com/deviceplane/deviceplane/pkg/models"
)

const WildcardMetric = string("*")

func FilterMetrics(
	metrics []datadog.Metric,
	project *models.Project,
	device *models.Device,
	metricType string,
	exposedMetrics []models.ExposedMetric,
	app *models.Application,
	serviceName *string,
) (filteredMetrics []datadog.Metric) {
	var metricPrefix string
	switch metricType {
	case models.DeviceMetricsConfigKey:
		metricPrefix = "deviceplane.device"
	case models.ServiceMetricsConfigKey:
		metricPrefix = "deviceplane.service"
	case models.ProjectMetricsConfigKey:
		metricPrefix = "deviceplane"
	default:
		return nil
	}

	// Build kv pair for efficiency
	allowedMetricsByName := make(map[string]*models.ExposedMetric, len(exposedMetrics))
	for i, m := range exposedMetrics {
		allowedMetricsByName[m.Name] = &exposedMetrics[i]
	}

	for _, m := range metrics {
		exposedMetric := allowedMetricsByName[m.Metric]
		if exposedMetric == nil {
			exposedMetric = allowedMetricsByName[WildcardMetric]
		}
		if exposedMetric == nil {
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
		for _, label := range exposedMetric.Labels {
			labelValue, ok := device.Labels[label]
			if ok {
				addTag("deviceplane.labels."+label, labelValue)
			}
		}

		// Optional properties
		for _, tag := range exposedMetric.Properties {
			switch tag {
			case "device":
				addTag("deviceplane.device", device.Name)
			}
		}

		// Guaranteed tags
		addTag("deviceplane.project", project.Name)
		if metricType == models.ServiceMetricsConfigKey {
			addTag("deviceplane.application", app.Name)
			addTag("deviceplane.service", *serviceName)
		}

		filteredMetrics = append(filteredMetrics, m)
	}

	return filteredMetrics
}
