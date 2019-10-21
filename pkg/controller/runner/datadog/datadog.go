package datadog

import (
	"context"
	"strings"
	"time"

	"github.com/apex/log"

	"github.com/deviceplane/deviceplane/pkg/controller/store"
	"github.com/deviceplane/deviceplane/pkg/models"
)

const (
	metricName    = "deviceplane.devices"
	metricType    = "count"
	projectTagKey = "project"
	nameTagKey    = "name"
)

type Runner struct {
	projects store.Projects
	devices  store.Devices
}

func NewRunner(projects store.Projects, devices store.Devices) *Runner {
	return &Runner{
		projects: projects,
		devices:  devices,
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

		devices, err := r.devices.ListDevices(ctx, project.ID)
		if err != nil {
			log.WithError(err).Error("list devices")
			continue
		}

		var req postMetricsRequest
		for _, device := range devices {
			metricValue := int64(0)
			if device.Status == models.DeviceStatusOnline {
				metricValue = 1
			}

			req.Series = append(req.Series, metric{
				Metric: metricName,
				Points: [][]int64{
					[]int64{time.Now().Unix(), metricValue},
				},
				Type: metricType,
				Tags: []string{
					strings.Join([]string{projectTagKey, project.Name}, ":"),
					strings.Join([]string{nameTagKey, device.Name}, ":"),
				},
			})
		}

		client := newClient(*project.DatadogAPIKey)

		if err := client.postMetrics(ctx, req); err != nil {
			log.WithError(err).Error("post metrics")
			continue
		}
	}
}
