package translation

import (
	"io"

	"github.com/deviceplane/deviceplane/pkg/metrics/datadog"

	prometheus "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

func ConvertOpenMetricsToDataDog(in io.Reader, statsCache *StatsCache, statsPrefix string) ([]datadog.Metric, error) {
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
		if promMetric.Metric == nil {
			continue
		}

		promValues := promMetric.GetMetric()
		for _, v := range promValues {
			points := make([][2]interface{}, 0)
			tags := make([]string, 0)

			labels := v.GetLabel()
			if len(labels) != 0 {
				for _, l := range labels {
					if l == nil {
						continue
					}
					tag := l.GetName() + ":" + l.GetValue()
					tags = append(tags, tag)
				}
			}

			switch *promMetric.Type {
			case prometheus.MetricType_GAUGE:
				gauge := v.GetGauge()
				if gauge == nil {
					continue
				}

				points = append(points, datadog.NewPoint(float32(gauge.GetValue())))
				m := datadog.Metric{
					Metric: promMetric.GetName(),
					Points: points,
					Type:   "gauge",
					Tags:   tags,
				}
				ddMetrics = append(ddMetrics, m)

			case prometheus.MetricType_COUNTER:
				counter := v.GetCounter()
				if counter == nil {
					continue
				}

				delta, ok := statsCache.UpdateCount(statsPrefix, promMetric.GetName(), tags, counter.GetValue())
				if !ok {
					continue
				}

				points = append(points, datadog.NewPoint(float32(delta)))
				m := datadog.Metric{
					Metric: promMetric.GetName(),
					Points: points,
					Type:   "count",
					Tags:   tags,
				}
				ddMetrics = append(ddMetrics, m)
			}
		}
	}

	return ddMetrics, nil
}
