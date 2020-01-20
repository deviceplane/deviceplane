package datadog

import (
	"context"
	"sync"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/apex/log"

	"github.com/deviceplane/deviceplane/pkg/controller/connman"
	"github.com/deviceplane/deviceplane/pkg/controller/store"
	"github.com/deviceplane/deviceplane/pkg/metrics/datadog"
	"github.com/deviceplane/deviceplane/pkg/metrics/datadog/translation"
	"github.com/deviceplane/deviceplane/pkg/models"
)

type Runner struct {
	projects      store.Projects
	applications  store.Applications
	devices       store.Devices
	releases      store.Releases
	metricConfigs store.MetricConfigs
	st            *statsd.Client
	connman       *connman.ConnectionManager
	statsCache    *translation.StatsCache
}

func NewRunner(projects store.Projects, applications store.Applications, releases store.Releases, devices store.Devices, metricConfigs store.MetricConfigs, st *statsd.Client, connman *connman.ConnectionManager) *Runner {
	return &Runner{
		projects:      projects,
		applications:  applications,
		devices:       devices,
		releases:      releases,
		metricConfigs: metricConfigs,
		st:            st,
		connman:       connman,
		statsCache:    translation.NewStatsCache(),
	}
}

func (r *Runner) Do(ctx context.Context) {
	projects, err := r.projects.ListProjects(ctx)
	if err != nil {
		log.WithError(err).Error("list projects")
		return
	}

	for _, project := range projects {
		if project.DatadogAPIKey == nil {
			continue
		}

		devices, err := r.devices.ListDevices(ctx, project.ID, "")
		if err != nil {
			log.WithError(err).Error("list devices")
			continue
		}

		// Get metric configs
		projectMetricsConfig, err := r.metricConfigs.GetProjectMetricsConfig(ctx, project.ID)
		if err != nil {
			log.WithField("project_id", project.ID).
				WithError(err).Error("getting project metrics config")
			continue
		}

		deviceMetricsConfig, err := r.metricConfigs.GetDeviceMetricsConfig(ctx, project.ID)
		if err != nil {
			log.WithField("project_id", project.ID).
				WithError(err).Error("getting device metrics config")
			continue
		}

		serviceMetricsConfigs, err := r.metricConfigs.GetServiceMetricsConfigs(ctx, project.ID)
		if err != nil {
			log.WithField("project_id", project.ID).
				WithError(err).Error("getting service metrics configs")
			continue
		}

		// Add apps to map by ID
		// Add services to map by name
		apps, err := r.applications.ListApplications(ctx, project.ID)
		if err != nil {
			log.WithField("project", project.ID).WithError(err).Error("listing applications")
			continue
		}
		appsByID := make(map[string]*models.Application, len(apps))
		latestAppReleaseByAppID := make(map[string]*models.Release, len(apps))
		if len(serviceMetricsConfigs) != 0 {
			for i, app := range apps {
				appsByID[app.ID] = &apps[i]

				release, err := r.releases.GetLatestRelease(ctx, project.ID, app.ID)
				if err == store.ErrReleaseNotFound {
					continue
				} else if err != nil {
					log.WithField("application", app.ID).
						WithError(err).Error("get latest release")
					continue
				}

				latestAppReleaseByAppID[app.ID] = release
			}
		}

		var lock sync.Mutex
		var req datadog.PostMetricsRequest
		var wg sync.WaitGroup
		for i := range devices {
			wg.Add(1)
			go func(device models.Device) {
				defer wg.Done()

				if len(projectMetricsConfig.ExposedMetrics) != 0 {
					projectMetrics := r.getProjectMetrics(ctx, &project, &device)
					filteredProjectMetrics := FilterMetrics(projectMetrics, &project, &device, models.ProjectMetricsConfigKey, projectMetricsConfig.ExposedMetrics, nil, nil)
					if len(filteredProjectMetrics) != 0 {
						lock.Lock()
						req.Series = append(req.Series, filteredProjectMetrics...)
						lock.Unlock()
					}
				}

				// If the device is offline can't get non-state metrics
				// from it
				if device.Status != models.DeviceStatusOnline {
					return
				}

				deviceConn, err := r.connman.Dial(ctx, project.ID+device.ID)
				if err != nil {
					return
				}
				defer deviceConn.Close()

				if len(deviceMetricsConfig.ExposedMetrics) != 0 {
					deviceMetrics := r.getDeviceMetrics(ctx, deviceConn, &project, &device)
					filteredDeviceMetrics := FilterMetrics(deviceMetrics, &project, &device, models.DeviceMetricsConfigKey, deviceMetricsConfig.ExposedMetrics, nil, nil)
					if len(filteredDeviceMetrics) != 0 {
						lock.Lock()
						req.Series = append(req.Series, filteredDeviceMetrics...)
						lock.Unlock()
					}
				}

				if len(serviceMetricsConfigs) != 0 {
					serviceMetrics := r.getServiceMetrics(ctx, deviceConn, &project, &device, apps, appsByID, latestAppReleaseByAppID, serviceMetricsConfigs)
					if len(serviceMetrics) != 0 {
						lock.Lock()
						req.Series = append(req.Series, serviceMetrics...)
						lock.Unlock()
					}
				}
			}(devices[i])
		}
		wg.Wait()

		if len(req.Series) == 0 {
			continue
		}

		client := datadog.NewClient(*project.DatadogAPIKey)
		if err := client.PostMetrics(ctx, req); err != nil {
			log.WithError(err).Error("post metrics")
			continue
		}
	}
}
