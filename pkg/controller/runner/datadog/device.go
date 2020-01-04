package datadog

import (
	"net"

	"github.com/apex/log"

	"github.com/deviceplane/deviceplane/pkg/agent/service/client"
	"github.com/deviceplane/deviceplane/pkg/metrics/datadog"
	"github.com/deviceplane/deviceplane/pkg/metrics/datadog/translation"
	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/deviceplane/deviceplane/pkg/utils"
)

func (r *Runner) getDeviceMetrics(
	deviceConn net.Conn,
	project *models.Project,
	device *models.Device,
) []datadog.Metric {
	// Get metrics from device
	deviceMetricsResp, err := client.GetDeviceMetrics(deviceConn)
	if err != nil || deviceMetricsResp.StatusCode != 200 {
		r.st.Incr("runner.datadog.device_metrics_pull", append([]string{"status:failure"}, utils.InternalTags(project.ID)...), 1)
		return nil
	}
	r.st.Incr("runner.datadog.device_metrics_pull", append([]string{"status:success"}, utils.InternalTags(project.ID)...), 1)

	// Convert request to DataDog format
	metrics, err := translation.ConvertOpenMetricsToDataDog(
		deviceMetricsResp.Body,
		r.statsCache,
		translation.GetMetricsPrefix(project, device, models.DeviceMetricsConfigKey),
	)
	if err != nil {
		log.WithField("project_id", project.ID).
			WithField("device_id", device.ID).
			WithError(err).Error("parsing openmetrics")
		return nil
	}

	return metrics
}
