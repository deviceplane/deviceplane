package processing

import (
	"fmt"

	"github.com/deviceplane/deviceplane/pkg/models"
)

const WildcardMetric = string("*")

func mapMetrics(exposedMetrics []models.ExposedMetric) map[string]*models.ExposedMetric {
	allowedMetricsByName := make(map[string]*models.ExposedMetric, len(exposedMetrics))
	for i, m := range exposedMetrics {
		allowedMetricsByName[m.Name] = &exposedMetrics[i]
	}
	return allowedMetricsByName
}

func addNoTags(func(tag, value string)) {}

var ProcessProjectMetrics = metricProcessorFunc("deviceplane.device", addNoTags)
var ProcessDeviceMetrics = metricProcessorFunc("deviceplane.device", addNoTags)
var ProcessServiceMetrics = func(applicationName, serviceName string) metricProcessor {
	return metricProcessorFunc("deviceplane.device",
		func(addTag func(tag, value string)) {
			if applicationName != "" {
				addTag("deviceplane.application", applicationName)
			}
			if serviceName != "" {
				addTag("deviceplane.service", serviceName)
			}
		},
	)
}

type metricProcessor func(
	controllerSide bool,
	metrics []models.DatadogMetric,
	exposedMetrics []models.ExposedMetric,
	project *models.Project,
	device *models.Device,
) (filteredMetrics []models.DatadogMetric)

func metricProcessorFunc(
	metricPrefix string,
	addTags func(func(tag, value string)),
) metricProcessor {
	return func(
		controllerSide bool,
		metrics []models.DatadogMetric,
		exposedMetrics []models.ExposedMetric,
		project *models.Project,
		device *models.Device,
	) (filteredMetrics []models.DatadogMetric) {

		exposedMetricsKV := mapMetrics(exposedMetrics)

		for _, m := range metrics {
			// Get exposed metric settings
			exposedMetric := exposedMetricsKV[m.Metric]
			if exposedMetric == nil {
				exposedMetric = exposedMetricsKV[WildcardMetric]
				if exposedMetric == nil {
					continue
				}
			}

			// Prefix metric name
			if !controllerSide {
				m.Metric = fmt.Sprintf("%s.%s", metricPrefix, m.Metric)
			}

			// Helper
			addTag := func(tag, value string) {
				m.Tags = append(
					m.Tags,
					fmt.Sprintf("%s:%s", tag, value),
				)
			}

			// Optional labels
			if device != nil {
				for _, label := range exposedMetric.Labels {
					labelValue, ok := device.Labels[label]
					if ok {
						addTag("deviceplane.labels."+label, labelValue)
					}
				}
			}

			// Optional properties
			for _, tag := range exposedMetric.Properties {
				switch tag {
				case "device":
					if device != nil {
						addTag("deviceplane.device", device.Name)
					}
				}
			}

			// Guaranteed tags
			if project != nil {
				addTag("deviceplane.project", project.Name)
			}

			// Func-defined tags
			addTags(addTag)

			filteredMetrics = append(filteredMetrics, m)
		}

		return filteredMetrics
	}
}
