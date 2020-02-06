package datadog

// import (
// 	"context"
// 	"fmt"
// 	"net"

// 	"github.com/apex/log"

// 	"github.com/deviceplane/deviceplane/pkg/agent/service/client"
// 	"github.com/deviceplane/deviceplane/pkg/controller/scheduling"
// 	"github.com/deviceplane/deviceplane/pkg/metrics/datadog"
// 	"github.com/deviceplane/deviceplane/pkg/metrics/datadog/translation"
// 	"github.com/deviceplane/deviceplane/pkg/models"
// 	"github.com/deviceplane/deviceplane/pkg/utils"
// )

// func (r *Runner) getServiceMetrics(
// 	ctx context.Context,
// 	deviceConn net.Conn,
// 	project *models.Project,
// 	device *models.Device,
// 	apps []models.Application,
// 	appsByID map[string]*models.Application,
// 	getReleaseByID func(releaseID string, appID string) *models.Release,
// 	serviceMetricsConfigs []models.ServiceMetricsConfig,
// ) (metrics datadog.Series) {
// 	scheduledAppReleaseOnDevice := map[string]string{}
// 	for _, serviceMetricsConfig := range serviceMetricsConfigs {
// 		app, exists := appsByID[serviceMetricsConfig.ApplicationID]
// 		if !exists {
// 			continue
// 		}

// 		releaseID, exists := scheduledAppReleaseOnDevice[app.ID]
// 		if !exists {
// 			var err error
// 			scheduled, scheduledDevice, err := scheduling.IsApplicationScheduled(*device, app.SchedulingRule)
// 			if err != nil {
// 				log.WithField("application", app.ID).
// 					WithField("device", device.ID).
// 					WithError(err).Error("evaluate application scheduling rule")
// 			} else if scheduled {
// 				releaseID = scheduledDevice.ReleaseID
// 			}

// 			scheduledAppReleaseOnDevice[app.ID] = releaseID
// 		}
// 		if releaseID == "" {
// 			continue
// 		}

// 		// Don't hit if there is no release
// 		release := getReleaseByID(releaseID, app.ID)
// 		if release == nil {
// 			continue
// 		}

// 		// Don't hit if service doesn't exist in config
// 		_, exists = release.Config[serviceMetricsConfig.Service]
// 		if !exists {
// 			continue
// 		}

// 		// Don't hit if user hasn't configured any metrics to export
// 		if len(serviceMetricsConfig.ExposedMetrics) == 0 {
// 			continue
// 		}

// 		serviceMetricEndpointConfig, exists := app.MetricEndpointConfigs[serviceMetricsConfig.Service]
// 		if !exists {
// 			serviceMetricEndpointConfig.Port = models.DefaultMetricPort
// 			serviceMetricEndpointConfig.Path = models.DefaultMetricPath
// 		}

// 		// Get metrics from services
// 		serviceMetricsResp, err := client.GetServiceMetrics(ctx, deviceConn, app.ID, serviceMetricsConfig.Service, serviceMetricEndpointConfig.Path, serviceMetricEndpointConfig.Port)
// 		serviceStatTags := append([]string{"service:" + serviceMetricsConfig.Service, "application:" + app.Name, "device:" + device.Name}, utils.InternalTags(project.Name)...)
// 		if err != nil || serviceMetricsResp.StatusCode != 200 {
// 			r.st.Incr("runner.datadog.service_metrics_pull", append(serviceStatTags, "status:failure"), 1)
// 			// TODO: we want to present to the user a list
// 			// of applications that don't have functioning
// 			// endpoints
// 			continue
// 		}
// 		r.st.Incr("runner.datadog.service_metrics_pull", append(serviceStatTags, "status:success"), 1)

// 		// Convert request to DataDog format
// 		serviceMetrics, err := translation.ConvertOpenMetricsToDataDog(
// 			serviceMetricsResp.Body,
// 			r.statsCache,
// 			translation.GetMetricsPrefix(project, device, fmt.Sprintf("service-(%s)(%s)", app.ID, serviceMetricsConfig.Service)),
// 		)
// 		if err != nil {
// 			log.WithField("project_id", project.ID).
// 				WithField("device_id", device.ID).
// 				WithError(err).Error("parsing openmetrics")
// 			continue
// 		}

// 		metrics = append(
// 			metrics,
// 			FilterMetrics(serviceMetrics, project, device, models.ServiceMetricsConfigKey, serviceMetricsConfig.ExposedMetrics, app, &serviceMetricsConfig.Service)...,
// 		)
// 	}

// 	return metrics
// }
