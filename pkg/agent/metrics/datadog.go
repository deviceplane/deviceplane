package metrics

import (
	"context"

	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/agent/client"
	"github.com/deviceplane/deviceplane/pkg/metrics/datadog/processing"
	"github.com/deviceplane/deviceplane/pkg/metrics/datadog/translation"
	"github.com/deviceplane/deviceplane/pkg/models"
)

// func (a *Agent) beginMetricsCollection(elem ...string) string {
// 	// load bundle
// 	// // get metrics
// 	// wait a minute
// }

type MetricsPusher struct {
	client                *client.Client
	statsCache            *translation.StatsCache
	serviceMetricsFetcher *ServiceMetricsFetcher
}

func NewMetricsPusher(
	client *client.Client,
	serviceMetricsFetcher *ServiceMetricsFetcher,
) *MetricsPusher {
	return &MetricsPusher{
		client:                client,
		serviceMetricsFetcher: serviceMetricsFetcher,

		statsCache: translation.NewStatsCache(),
	}
}

func (m *MetricsPusher) PushDeviceMetrics(ctx context.Context, bundle *models.Bundle) {
	if len(bundle.DeviceMetricsConfig.ExposedMetrics) == 0 {
		return
	}

	deviceMetrics := GetFilteredHostMetrics()
	convertedMetrics, err := translation.ConvertOpenMetricsToDataDog(
		deviceMetrics,
		m.statsCache,
		"service-metrics",
	)
	processedMetrics := processing.ProcessDeviceMetrics(
		convertedMetrics,
		bundle.DeviceMetricsConfig.ExposedMetrics,
		nil,
		nil,
	)

	if len(processedMetrics) == 0 {
		return
	}

	err = m.client.SendDeviceMetrics(ctx, models.DatadogPostMetricsRequest{
		Series: processedMetrics,
	})
	if err != nil {
		log.WithError(err).Error("could not POST device metrics")
	}
}

func (m *MetricsPusher) PushServiceMetrics(ctx context.Context, bundle *models.Bundle) {
	if len(bundle.ServiceMetricsConfigs) == 0 {
		return
	}

	var datadogMetrics = make(models.IntermediateServiceMetricsRequest)

	// Faster accessing
	appsByID := make(map[string]*models.FullBundledApplication, len(bundle.Applications))
	for i, app := range bundle.Applications {
		appsByID[app.Application.ID] = &bundle.Applications[i]
	}

	// Faster accessing
	serviceConfigsByID := make(map[string]*models.ServiceMetricsConfig, len(bundle.ServiceMetricsConfigs))
	for i, config := range bundle.ServiceMetricsConfigs {
		serviceConfigsByID[config.ApplicationID] = &bundle.ServiceMetricsConfigs[i]
	}

	for _, service := range bundle.ServiceStatuses {
		app, exists := appsByID[service.ApplicationID]
		if !exists {
			log.WithField("application_id", service.ApplicationID).
				WithField("service", service.Service).Error("could not get application for metrics")
			continue
		}

		serviceConfig, exists := serviceConfigsByID[app.Application.ID]
		if !exists {
			log.WithField("application_id", service.ApplicationID).
				WithField("service", service.Service).Error("could not get service metrics config for metrics")
			continue
		}

		config, exists := app.Application.MetricEndpointConfigs[service.Service]
		if !exists {
			config.Port = models.DefaultMetricPort
			config.Path = models.DefaultMetricPath
		}

		metricResponse, err := m.serviceMetricsFetcher.ContainerServiceMetrics(
			ctx,
			service.ApplicationID,
			service.Service,
			int(config.Port),
			config.Path,
		)
		if err != nil {
			log.WithField("application_id", service.ApplicationID).
				WithField("service", service.Service).Error("could not fetch service metrics")
		}
		defer metricResponse.Body.Close()

		convertedMetrics, err := translation.ConvertOpenMetricsToDataDog(
			metricResponse.Body,
			m.statsCache,
			"service-metrics",
		)

		processedMetrics := processing.ProcessServiceMetrics(app.Application.Name, service.Service)(
			convertedMetrics,
			serviceConfig.ExposedMetrics,
			nil,
			nil,
		)

		_, exists = datadogMetrics[app.Application.ID]
		if !exists {
			datadogMetrics[app.Application.ID] = make(map[string]models.DatadogSeries)
		}
		datadogMetrics[app.Application.ID][service.Service] = processedMetrics
	}

	if len(datadogMetrics) == 0 {
		return
	}

	err := m.client.SendServiceMetrics(ctx, datadogMetrics)
	if err != nil {
		log.WithError(err).Error("could not POST service metrics")
	}
}
