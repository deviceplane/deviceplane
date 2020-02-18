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
			addTag := func(key, value string) {
				hasTag := false
				for _, tag := range m.Tags { // TODO: Refactor m.Tags to be key:value pair
					if strings.HasPrefix(tag, key) {
						hasTag = true
					}
				}
				if !hasTag {
					m.Tags = append(
						m.Tags,
						fmt.Sprintf("%s:%s", key, value),
					)
				}
			}

			// Only keep tags on whitelist
			if len(exposedMetric.WhitelistedTags) != 0 { // Trim only if whitelist exists
				var allowedTags []string
				for _, tagPair := range m.Tags {
					var tagKey string
					var tagValue string
					parsedTagPair := strings.Split(tagPair, ":")
					if len(parsedTagPair) > 0 {
						tagKey = parsedTagPair[0]
					}
					if len(parsedTagPair) > 1 {
						tagValue = parsedTagPair[1]
					}

					if strings.HasPrefix(tagKey, "deviceplane.") {
						allowedTags = append(allowedTags, tagPair)
						continue
					}
					for _, wTag := range exposedMetric.WhitelistedTags {
						if tagKey == wTag.Key {
							if len(wTag.Values) == 0 {
								allowedTags = append(allowedTags, tagPair)
							} else {
								for _, value := range wTag.Values {
									if tagValue == value {
										allowedTags = append(allowedTags, tagPair)
										break
									}
								}
							}
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
