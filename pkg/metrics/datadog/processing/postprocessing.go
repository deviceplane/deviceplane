package processing

import (
	"fmt"
	"strings"

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

var ProcessProjectMetrics = metricProcessorFunc("deviceplane.", addNoTags)
var ProcessDeviceMetrics = metricProcessorFunc("deviceplane.device.", addNoTags)
var ProcessServiceMetrics = func(applicationName, serviceName string) metricProcessor {
	return metricProcessorFunc("deviceplane.service.",
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
		metrics []models.DatadogMetric,
		exposedMetrics []models.ExposedMetric,
		project *models.Project,
		device *models.Device,
	) (filteredMetrics []models.DatadogMetric) {

		exposedMetricsKV := mapMetrics(exposedMetrics)

		for _, m := range metrics {
			// Prefix metric name
			if !strings.HasPrefix(m.Metric, metricPrefix) {
				m.Metric = metricPrefix + m.Metric
			}

			// Get unprefixed metric name
			start := strings.Index(m.Metric, metricPrefix) + len(metricPrefix)
			unprefixedMetricName := m.Metric[start:]

			// Get exposed metric settings, or wildcard settings
			exposedMetric := exposedMetricsKV[unprefixedMetricName]
			if exposedMetric == nil {
				exposedMetric = exposedMetricsKV[WildcardMetric]
				if exposedMetric == nil {
					continue
				}
			}

			// Helper
			addTag := func(tag, value string) {
				m.Tags = append(
					m.Tags,
					fmt.Sprintf("%s:%s", tag, value),
				)
			}

			// Only keep tags on whitelist
			if len(exposedMetric.WhitelistedTags) != 0 { // Trim only if whitelist exists
				var allowedTags []string
				for _, tagPair := range m.Tags {
					tag := strings.Split(tagPair, ":")[0]

					if strings.HasPrefix(tag, "deviceplane.") {
						allowedTags = append(allowedTags, tagPair)
						continue
					}
					for _, wTagName := range exposedMetric.WhitelistedTags {
						if tag == wTagName {
							allowedTags = append(allowedTags, tagPair)
							break
						}
					}
				}
				m.Tags = allowedTags
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

			// Add modified metric
			filteredMetrics = append(filteredMetrics, m)
		}

		return filteredMetrics
	}
}
