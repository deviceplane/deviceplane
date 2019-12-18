package models

import "time"

const (
	DefaultServiceMetricPort uint = 2112
	DefaultServiceMetricPath      = "/metrics"
)

// ServiceMetricConfig is used for metrics scraping
type ServiceMetricConfig struct {
	Port uint   `json:"port" yaml:"port"`
	Path string `json:"path" yaml:"path"`
}

// The following are used for metrics forwarding

type ExposedMetricConfigHolder struct {
	ID        string                `json:"id" yaml:"id"`
	CreatedAt time.Time             `json:"createdAt" yaml:"createdAt"`
	ProjectID string                `json:"projectId" yaml:"projectId"`
	Type      ExposedMetricType     `json:"type" yaml:"type"`
	Configs   []ExposedMetricConfig `json:"configs" yaml:"configs"`
}

type ExposedMetricType string

const (
	ExposedServiceMetric ExposedMetricType = "service"
	ExposedHostMetric    ExposedMetricType = "host"
	ExposedStateMetric   ExposedMetricType = "state"
)

type ExposedMetricConfig struct {
	Params  *ExposedServiceMetricParams `json:"params,omitempty" json:"yaml,omitempty"`
	Metrics []ExposedMetric             `json:"metrics" yaml:"metrics"`
}

type ExposedServiceMetricParams struct {
	ApplicationID string `json:"applicationId" yaml:"applicationId"`
	Service       string `json:"service" yaml:"service"`
}

type ExposedMetric struct {
	Metric string   `json:"metric" yaml:"metric"`
	Labels []string `json:"labels" yaml:"labels"`
	Tags   []string `json:"tags" yaml:"tags"`
}
