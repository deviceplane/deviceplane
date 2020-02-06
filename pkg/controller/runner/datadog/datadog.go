package datadog

import (
	"context"
	"sync"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/apex/log"

	"github.com/deviceplane/deviceplane/pkg/controller/connman"
	"github.com/deviceplane/deviceplane/pkg/controller/store"
	"github.com/deviceplane/deviceplane/pkg/metrics/datadog"
	"github.com/deviceplane/deviceplane/pkg/metrics/datadog/processing"
	"github.com/deviceplane/deviceplane/pkg/metrics/datadog/translation"
	"github.com/deviceplane/deviceplane/pkg/models"
)

const (
	perDeviceTimeout = 10 * time.Second
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

	var wg sync.WaitGroup
	for _, project := range projects {
		if project.DatadogAPIKey == nil {
			continue
		}

		wg.Add(1)
		go func(project models.Project) {
			r.doForProject(ctx, project)
			wg.Done()
		}(project)
	}

	wg.Wait()
}

func (r *Runner) doForProject(ctx context.Context, project models.Project) {
	// Get metric configs
	projectMetricsConfig, err := r.metricConfigs.GetProjectMetricsConfig(ctx, project.ID)
	if err != nil {
		log.WithField("project_id", project.ID).
			WithError(err).Error("getting project metrics config")
		return
	}
	if len(projectMetricsConfig.ExposedMetrics) != 0 {
		return
	}

	var req models.DatadogPostMetricsRequest

	devices, err := r.devices.ListDevices(ctx, project.ID, "")
	if err != nil {
		log.WithError(err).Error("list devices")
		return
	}
	for _, device := range devices {
		projectMetrics := r.getProjectMetrics(&project, &device)
		filteredProjectMetrics := processing.ProcessProjectMetrics(projectMetrics, projectMetricsConfig.ExposedMetrics, &project, &device)
		req.Series = append(req.Series, filteredProjectMetrics...)
	}

	if len(req.Series) == 0 {
		return
	}

	client := datadog.NewClient(*project.DatadogAPIKey)
	if err := client.PostMetrics(ctx, req); err != nil {
		log.WithError(err).Error("post metrics")
		return
	}
}
