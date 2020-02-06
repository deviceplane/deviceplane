package metrics

import (
	"context"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/agent/client"
	"github.com/deviceplane/deviceplane/pkg/metrics/datadog/processing"
	"github.com/deviceplane/deviceplane/pkg/metrics/datadog/translation"
	"github.com/deviceplane/deviceplane/pkg/models"
)

type MetricsPusher struct {
	client                *client.Client
	statsCache            *translation.StatsCache
	serviceMetricsFetcher *ServiceMetricsFetcher

	lock sync.Mutex
	once sync.Once

	bundle models.Bundle
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

func (m *MetricsPusher) SetBundle(bundle models.Bundle) {
	m.lock.Lock()
	m.bundle = bundle
	m.lock.Unlock()

	go m.once.Do(m.begin)
}

func (m *MetricsPusher) begin() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			m.PushDeviceMetrics(ctx)
			wg.Done()
		}()
		go func() {
			m.PushServiceMetrics(ctx)
			wg.Done()
		}()

		wg.Wait()
		cancel()

		<-ticker.C
	}
}

func (m *MetricsPusher) PushDeviceMetrics(ctx context.Context) {
	if m.bundle.DeviceMetricsConfig == nil {
		return
	}

	if len(m.bundle.DeviceMetricsConfig.ExposedMetrics) == 0 {
		return
	}

	deviceMetrics, err := GetFilteredHostMetrics(ctx)
	if err != nil {
		log.WithError(err).Error("could not get filtered host metrics")
		return
	}
	convertedMetrics, err := translation.ConvertOpenMetricsToDataDog(
		deviceMetrics,
		m.statsCache,
		"device-metrics",
	)
	processedMetrics := processing.ProcessDeviceMetrics(
		convertedMetrics,
		m.bundle.DeviceMetricsConfig.ExposedMetrics,
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

func (m *MetricsPusher) PushServiceMetrics(ctx context.Context) {
	if len(m.bundle.ServiceMetricsConfigs) == 0 {
		return
	}

	var datadogMetrics = make(models.IntermediateServiceMetricsRequest)

	// Faster accessing
	appsByID := make(map[string]*models.FullBundledApplication, len(m.bundle.Applications))
	for i, app := range m.bundle.Applications {
		appsByID[app.Application.ID] = &m.bundle.Applications[i]
	}

	// Faster accessing
	serviceConfigsByID := make(map[string]*models.ServiceMetricsConfig, len(m.bundle.ServiceMetricsConfigs))
	for i, config := range m.bundle.ServiceMetricsConfigs {
		serviceConfigsByID[config.ApplicationID] = &m.bundle.ServiceMetricsConfigs[i]
	}

	for _, service := range m.bundle.ServiceStatuses {
		app, exists := appsByID[service.ApplicationID]
		if !exists {
			log.WithField("application_id", service.ApplicationID).
				WithField("service", service.Service).Info("could not get application for metrics")
			continue
		}

		serviceConfig, exists := serviceConfigsByID[app.Application.ID]
		if !exists {
			log.WithField("application_id", service.ApplicationID).
				WithField("service", service.Service).Info("could not get service metrics config for metrics")
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
			continue
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
