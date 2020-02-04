package service

import (
	"net/http"

	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/metrics/datadog"
	"github.com/deviceplane/deviceplane/pkg/metrics/datadog/processing"
	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/deviceplane/deviceplane/pkg/utils"
)

func (s *Service) forwardServiceMetrics(w http.ResponseWriter, r *http.Request, project models.Project, device models.Device) {
	pass := func() bool {
		var metricsRequest models.IntermediateServiceMetricsRequest

		if err := read(r, &metricsRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return false
		}

		serviceMetricsConfigs, err := s.metricConfigs.GetServiceMetricsConfigs(r.Context(), project.ID)
		if err != nil {
			log.WithField("project_id", project.ID).
				WithError(err).Error("getting Service metrics config")
			return false
		}
		configsByID := make(map[string]*models.ServiceMetricsConfig, len(serviceMetricsConfigs))
		for i, c := range serviceMetricsConfigs {
			configsByID[c.ApplicationID+c.Service] = &serviceMetricsConfigs[i]
		}

		var forwardedMetricsRequest models.DatadogPostMetricsRequest

		for appID, v := range metricsRequest {
			for service, series := range v {
				config := configsByID[appID+service]

				filteredServiceMetrics := processing.ProcessServiceMetrics(
					"",
					service,
				)(
					series,
					config.ExposedMetrics,
					&project,
					&device,
				)
				if len(filteredServiceMetrics) == 0 {
					forwardedMetricsRequest.Series = filteredServiceMetrics
				}
			}
		}

		client := datadog.NewClient(*project.DatadogAPIKey)
		if err := client.PostMetrics(r.Context(), forwardedMetricsRequest); err != nil {
			log.WithError(err).Error("post service metrics")
			return false
		}
		return true
	}()
	if pass {
		s.st.Incr("runner.datadog.service_metrics_push", append([]string{"status:success"}, utils.InternalTags(project.Name)...), 1)
	} else {
		s.st.Incr("runner.datadog.service_metrics_push", append([]string{"status:failure"}, utils.InternalTags(project.Name)...), 1)
	}
}

func (s *Service) forwardDeviceMetrics(w http.ResponseWriter, r *http.Request, project models.Project, device models.Device) {
	pass := func() bool {
		var metricsRequest models.DatadogPostMetricsRequest

		if err := read(r, &metricsRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return false
		}

		deviceMetricsConfig, err := s.metricConfigs.GetDeviceMetricsConfig(r.Context(), project.ID)
		if err != nil {
			log.WithField("project_id", project.ID).
				WithError(err).Error("getting device metrics config")
			return false
		}

		var forwardedMetricsRequest models.DatadogPostMetricsRequest
		filteredDeviceMetrics := processing.ProcessDeviceMetrics(
			metricsRequest.Series,
			deviceMetricsConfig.ExposedMetrics,
			&project,
			&device,
		)
		if len(filteredDeviceMetrics) == 0 {
			forwardedMetricsRequest.Series = filteredDeviceMetrics
		}

		client := datadog.NewClient(*project.DatadogAPIKey)
		if err := client.PostMetrics(r.Context(), metricsRequest); err != nil {
			log.WithError(err).Error("post device metrics")
			return false
		}
		return true
	}()
	if pass {
		s.st.Incr("runner.datadog.device_metrics_push", append([]string{"status:success"}, utils.InternalTags(project.Name)...), 1)
	} else {
		s.st.Incr("runner.datadog.device_metrics_push", append([]string{"status:failure"}, utils.InternalTags(project.Name)...), 1)
	}
}
