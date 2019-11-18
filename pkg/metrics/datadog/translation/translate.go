package translation

import (
	"io"
	"strings"
	"time"

	"github.com/deviceplane/deviceplane/pkg/metrics/datadog"

	prometheus "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

func ConvertOpenMetricsToDataDog(in io.Reader) ([]datadog.Metric, error) {
	parser := expfmt.TextParser{}
	promMetrics, err := parser.TextToMetricFamilies(in)
	if err != nil {
		return nil, err
	}

	ddMetrics := make([]datadog.Metric, 0)
	for _, promMetric := range promMetrics {
		if promMetric.Type == nil {
			continue
		}
		if *promMetric.Type != prometheus.MetricType_GAUGE {
			continue
		}

		if promMetric.Metric == nil {
			continue
		}

		points := make([][2]float32, 0)
		tags := make([]string, 0)

		promValues := promMetric.GetMetric()
		for _, v := range promValues {
			gauge := v.GetGauge()
			if gauge != nil {
				points = append(points, [2]float32{float32(time.Now().Unix()), float32(gauge.GetValue())})
			}

			labels := v.GetLabel()
			if len(labels) != 0 {
				for _, l := range labels {
					if l == nil {
						continue
					}
					tag := strings.Join([]string{l.GetName(), l.GetValue()}, ":")
					tags = append(tags, tag)
				}
			}
		}

		m := datadog.Metric{
			Metric: promMetric.GetName(),
			Points: points,
			Type:   "gauge",
			Tags:   tags,
		}
		ddMetrics = append(ddMetrics, m)
	}

	return ddMetrics, nil
}
