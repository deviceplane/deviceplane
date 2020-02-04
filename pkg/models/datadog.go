package models

type DatadogPostMetricsRequest struct {
	Series DatadogSeries `json:"series"`
}

type IntermediateServiceMetricsRequest map[string](map[string]DatadogSeries)

type DatadogSeries []DatadogMetric

type DatadogMetric struct {
	Metric   string           `json:"metric"`
	Points   [][2]interface{} `json:"points"`
	Type     string           `json:"type"`
	Interval *int64           `json:"interval,omitempty"`
	Host     string           `json:"host,omitempty"`
	Tags     []string         `json:"tags"`
}
