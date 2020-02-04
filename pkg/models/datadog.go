package models

type DatadogPostMetricsRequest struct {
	Series DatadogSeries `json:"series"`
}

type IntermediateServiceMetricsRequest map[string](map[string]DatadogSeries)

func (i IntermediateServiceMetricsRequest) flatten() []DatadogMetric {
	var series []DatadogMetric

	for _, v := range i {
		for _, v2 := range v {
			series = append(series, v2...)
		}
	}

	return series
}

type DatadogSeries []DatadogMetric

type DatadogMetric struct {
	Metric   string           `json:"metric"`
	Points   [][2]interface{} `json:"points"`
	Type     string           `json:"type"`
	Interval *int64           `json:"interval,omitempty"`
	Host     string           `json:"host,omitempty"`
	Tags     []string         `json:"tags"`
}
