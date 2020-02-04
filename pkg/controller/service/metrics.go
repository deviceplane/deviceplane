package service

import (
	"fmt"
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
			w.WriteHeader(http.StatusInternalServerError)
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
					true,
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
			w.WriteHeader(http.StatusInternalServerError)
			return false
		}
		return true
	}()

	var status string
	if pass {
		status = "status:success"
	} else {
		status = "status:failure"
	}
	s.st.Incr("service_metrics_push", append([]string{status}, utils.InternalTags(project.Name)...), 1)

	fmt.Println("SERVICE METRICS PASS", pass)
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
			w.WriteHeader(http.StatusInternalServerError)
			return false
		}

		var forwardedMetricsRequest models.DatadogPostMetricsRequest
		filteredDeviceMetrics := processing.ProcessDeviceMetrics(
			true,
			metricsRequest.Series,
			deviceMetricsConfig.ExposedMetrics,
			&project,
			&device,
		)
		if len(filteredDeviceMetrics) == 0 {
			forwardedMetricsRequest.Series = filteredDeviceMetrics
		}

		client := datadog.NewClient(*project.DatadogAPIKey)
		fmt.Println(metricsRequest)
		fmt.Println(forwardedMetricsRequest)
		if err := client.PostMetrics(r.Context(), forwardedMetricsRequest); err != nil {
			log.WithError(err).Error("post device metrics")
			w.WriteHeader(http.StatusInternalServerError)
			return false
		}
		return true
	}()

	var status string
	if pass {
		status = "status:success"
	} else {
		status = "status:failure"
	}
	s.st.Incr("device_metrics_push", append([]string{status}, utils.InternalTags(project.Name)...), 1)

	fmt.Println("DEVICE METRICS PASS", pass)
}
