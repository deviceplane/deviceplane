package models

type ProjectConfig struct {
	ProjectID string `json:"projectId" yaml:"projectId"`
	Key       string `json:"key" yaml:"key"`
	Value     string `json:"value" yaml:"value"`
}

const (
	ServiceMetricsConfigKey = "service-metrics-config"
	ProjectMetricsConfigKey = "project-metrics-config"
	DeviceMetricsConfigKey  = "device-metrics-config"
)

type ServiceMetricsConfig struct {
	ApplicationID  string          `json:"applicationId" yaml:"applicationId"`
	Service        string          `json:"service" yaml:"service"`
	ExposedMetrics []ExposedMetric `json:"exposedMetrics" yaml:"exposedMetrics"`
}

type ProjectMetricsConfig struct {
	ExposedMetrics []ExposedMetric `json:"exposedMetrics" yaml:"exposedMetrics"`
}

type DeviceMetricsConfig struct {
	ExposedMetrics []ExposedMetric `json:"exposedMetrics" yaml:"exposedMetrics"`
}

type ExposedMetric struct {
	Name            string           `json:"name" yaml:"name"`
	Labels          []string         `json:"labels" yaml:"labels"`
	Properties      []string         `json:"properties" yaml:"properties"`
	WhitelistedTags []WhitelistedTag `json:"whitelistedTags" yaml:"whitelistedTags"`
}

type WhitelistedTag struct {
	Key    string   `json:"key" yaml:"key"`
	Values []string `json:"values" yaml:"values"`
}
