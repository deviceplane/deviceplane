package datadog

import (
	"net"

	"github.com/apex/log"

	"github.com/deviceplane/deviceplane/pkg/agent/service/client"
	"github.com/deviceplane/deviceplane/pkg/metrics/datadog"
	"github.com/deviceplane/deviceplane/pkg/metrics/datadog/translation"
	"github.com/deviceplane/deviceplane/pkg/models"
)

func (r *Runner) getHostMetrics(deviceConn net.Conn, project *models.Project, device *models.Device, metricConfig *models.MetricTargetConfig) datadog.Series {
	// Get metrics from host
	deviceMetricsResp, err := client.GetHostMetrics(deviceConn)
	if err != nil || deviceMetricsResp.StatusCode != 200 {
		r.st.Incr("runner.datadog.host_metrics_pull", append([]string{"status:failure"}, addedInternalTags(project, device)...), 1)
		return nil
	}
	r.st.Incr("runner.datadog.host_metrics_pull", append([]string{"status:success"}, addedInternalTags(project, device)...), 1)

	// Convert request to DataDog format
	metrics, err := translation.ConvertOpenMetricsToDataDog(deviceMetricsResp.Body)
	if err != nil {
		log.WithField("project_id", project.ID).
			WithField("device_id", device.ID).
			WithError(err).Error("parsing openmetrics")
		return nil
	}

	config := metricConfig.Configs[0]
	return FilterMetrics(project, nil, device, metricConfig.Type, config, metrics)
}
