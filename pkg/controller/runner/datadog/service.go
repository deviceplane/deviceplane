package datadog

import (
	"net"

	"github.com/apex/log"

	"github.com/deviceplane/deviceplane/pkg/agent/service/client"
	"github.com/deviceplane/deviceplane/pkg/controller/query"
	"github.com/deviceplane/deviceplane/pkg/metrics/datadog"
	"github.com/deviceplane/deviceplane/pkg/metrics/datadog/translation"
	"github.com/deviceplane/deviceplane/pkg/models"
)

func (r *Runner) getServiceMetrics(
	deviceConn net.Conn,
	project *models.Project,
	device *models.Device,
	metricConfig *models.MetricTargetConfig,
	apps []models.Application,
	appsByID map[string]*models.Application,
	latestAppReleaseByAppID map[string]*models.Release,
) (metrics datadog.Series) {
	appIsScheduled := map[string]bool{} // we have denormalized (app, serv), (app, serv2) tuples in metricConfig.Configs
	for _, config := range metricConfig.Configs {
		if config.Params == nil {
			return nil
		}

		app, exists := appsByID[config.Params.ApplicationID]
		if !exists {
			return nil
		}

		scheduled, exists := appIsScheduled[app.ID]
		if !exists {
			var err error
			scheduled, err = query.DeviceMatchesQuery(*device, app.SchedulingRule)
			if err != nil {
				log.WithField("application", app.ID).
					WithField("device", device.ID).
					WithError(err).Error("evaluate application scheduling rule")
				scheduled = false
			}
			appIsScheduled[app.ID] = scheduled
		}
		if !scheduled {
			continue
		}

		release, exists := latestAppReleaseByAppID[app.ID]
		if !exists {
			continue
		}

		_, exists = release.Config[config.Params.Service]
		if !exists {
			continue
		}

		// Get metrics from services
		deviceMetricsResp, err := client.GetServiceMetrics(deviceConn, app.ID, config.Params.Service)
		if err != nil || deviceMetricsResp.StatusCode != 200 {
			r.st.Incr("runner.datadog.service_metrics_pull", append([]string{"status:failure"}, addedInternalTags(project)...), 1)
			// TODO: we want to present to the user a list
			// of applications that don't have functioning
			// endpoints
			continue
		}
		r.st.Incr("runner.datadog.service_metrics_pull", append([]string{"status:success"}, addedInternalTags(project)...), 1)

		// Convert request to DataDog format
		serviceMetrics, err := translation.ConvertOpenMetricsToDataDog(deviceMetricsResp.Body)
		if err != nil {
			log.WithField("project_id", project.ID).
				WithField("device_id", device.ID).
				WithError(err).Error("parsing openmetrics")
			continue
		}

		metrics = append(
			metrics,
			FilterMetrics(project, app, device, metricConfig.Type, config, serviceMetrics)...,
		)
	}

	return metrics
}
