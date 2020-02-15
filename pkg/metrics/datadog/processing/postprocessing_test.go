package processing

import (
	"testing"

	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/stretchr/testify/assert"
)

var project = models.Project{
	ID:   "project-id",
	Name: "project-name",
}
var device = models.Device{
	ID:   "device-id",
	Name: "device-name",
	Labels: map[string]string{
		"company":  "nasa",
		"location": "ohio",
	},
}

func TestFilteringMetrics(t *testing.T) {
	metrics := []models.DatadogMetric{
		models.DatadogMetric{
			Metric: "Test_Metric",
			Tags:   []string{"TEST:TRUE"},
		},
	}
	exposedMetrics := []models.ExposedMetric{}

	filteredMetrics := ProcessDeviceMetrics(metrics, exposedMetrics, &project, &device)
	assert.Len(t, filteredMetrics, 0)
}

func TestAllowingMetrics(t *testing.T) {
	metrics := []models.DatadogMetric{
		models.DatadogMetric{
			Metric: "Test_Metric",
			Tags:   []string{"TEST:TRUE"},
		},
		models.DatadogMetric{
			Metric: "Test_Metric_Two",
			Tags:   []string{"TEST:TRUE"},
		},
	}
	exposedMetrics := []models.ExposedMetric{
		models.ExposedMetric{
			Name: "Test_Metric",
		},
	}

	filteredMetrics := ProcessDeviceMetrics(metrics, exposedMetrics, &project, &device)
	assert.Equal(t, []models.DatadogMetric{
		models.DatadogMetric{
			Metric: "deviceplane.device.Test_Metric",
			Tags:   []string{"TEST:TRUE", "deviceplane.project:project-name"},
		},
	}, filteredMetrics)
}

func TestAllowingAllMetrics(t *testing.T) {
	metrics := []models.DatadogMetric{
		models.DatadogMetric{
			Metric: "Test_Metric",
			Tags:   []string{"TEST:TRUE"},
		},
		models.DatadogMetric{
			Metric: "Test_Metric_Two",
			Tags:   []string{"TEST:TRUE"},
		},
	}
	exposedMetrics := []models.ExposedMetric{
		models.ExposedMetric{
			Name: "*",
		},
	}

	filteredMetrics := ProcessDeviceMetrics(metrics, exposedMetrics, &project, &device)
	assert.Equal(t, []models.DatadogMetric{
		models.DatadogMetric{
			Metric: "deviceplane.device.Test_Metric",
			Tags:   []string{"TEST:TRUE", "deviceplane.project:project-name"},
		},
		models.DatadogMetric{
			Metric: "deviceplane.device.Test_Metric_Two",
			Tags:   []string{"TEST:TRUE", "deviceplane.project:project-name"},
		},
	}, filteredMetrics)
}

func TestMetricLabels(t *testing.T) {
	metrics := []models.DatadogMetric{
		models.DatadogMetric{
			Metric: "Test_Metric",
			Tags:   []string{"TEST:TRUE"},
		},
	}
	exposedMetrics := []models.ExposedMetric{
		models.ExposedMetric{
			Name:   "Test_Metric",
			Labels: []string{"company"},
		},
	}

	filteredMetrics := ProcessDeviceMetrics(metrics, exposedMetrics, &project, &device)
	assert.Equal(t, []models.DatadogMetric{
		models.DatadogMetric{
			Metric: "deviceplane.device.Test_Metric",
			Tags:   []string{"TEST:TRUE", "deviceplane.labels.company:nasa", "deviceplane.project:project-name"},
		},
	}, filteredMetrics)
}

func TestMetricProperties(t *testing.T) {
	metrics := []models.DatadogMetric{
		models.DatadogMetric{
			Metric: "Test_Metric",
			Tags:   []string{"TEST:TRUE"},
		},
	}
	exposedMetrics := []models.ExposedMetric{
		models.ExposedMetric{
			Name:       "Test_Metric",
			Properties: []string{"device"},
		},
	}

	filteredMetrics := ProcessDeviceMetrics(metrics, exposedMetrics, &project, &device)
	assert.Equal(t, []models.DatadogMetric{
		models.DatadogMetric{
			Metric: "deviceplane.device.Test_Metric",
			Tags:   []string{"TEST:TRUE", "deviceplane.device:device-name", "deviceplane.project:project-name"},
		},
	}, filteredMetrics)
}

func TestMetricSeparation(t *testing.T) {
	metrics := []models.DatadogMetric{
		models.DatadogMetric{
			Metric: "Test_Metric",
			Tags:   []string{"TEST:TRUE"},
		},
		models.DatadogMetric{
			Metric: "Test_Metric_Two",
			Tags:   []string{"TEST:TRUE"},
		},
		models.DatadogMetric{
			Metric: "Test_Metric_Three",
			Tags:   []string{"TEST:TRUE"},
		},
	}
	exposedMetrics := []models.ExposedMetric{
		models.ExposedMetric{
			Name:   "Test_Metric",
			Labels: []string{"company"},
		},
		models.ExposedMetric{
			Name: "*",
		},
		models.ExposedMetric{
			Name:   "Test_Metric_Three",
			Labels: []string{"location"},
		},
	}

	filteredMetrics := ProcessDeviceMetrics(metrics, exposedMetrics, &project, &device)
	assert.Equal(t, []models.DatadogMetric{
		models.DatadogMetric{
			Metric: "deviceplane.device.Test_Metric",
			Tags:   []string{"TEST:TRUE", "deviceplane.labels.company:nasa", "deviceplane.project:project-name"},
		},
		models.DatadogMetric{
			Metric: "deviceplane.device.Test_Metric_Two",
			Tags:   []string{"TEST:TRUE", "deviceplane.project:project-name"},
		},
		models.DatadogMetric{
			Metric: "deviceplane.device.Test_Metric_Three",
			Tags:   []string{"TEST:TRUE", "deviceplane.labels.location:ohio", "deviceplane.project:project-name"},
		},
	}, filteredMetrics)
}

func TestMetricTagWhitelist(t *testing.T) {
	metrics := []models.DatadogMetric{
		models.DatadogMetric{
			Metric: "Test_Metric",
			Tags:   []string{"TEST:TRUE"},
		},
		models.DatadogMetric{
			Metric: "Test_Metric_Two",
			Tags:   []string{"TEST:TRUE"},
		},
		models.DatadogMetric{
			Metric: "Test_Metric_Three",
			Tags:   []string{"TEST:TRUE"},
		},
	}
	exposedMetrics := []models.ExposedMetric{
		models.ExposedMetric{
			Name:            "Test_Metric",
			WhitelistedTags: []string{"TEST"},
		},
		models.ExposedMetric{
			Name:            "Test_Metric_Two",
			WhitelistedTags: []string{"WEST"},
		},
		models.ExposedMetric{
			Name:            "Test_Metric_Three",
			WhitelistedTags: []string{},
		},
	}

	filteredMetrics := ProcessDeviceMetrics(metrics, exposedMetrics, &project, &device)
	assert.Equal(t, []models.DatadogMetric{
		models.DatadogMetric{
			Metric: "deviceplane.device.Test_Metric",
			Tags:   []string{"TEST:TRUE", "deviceplane.project:project-name"},
		},
		models.DatadogMetric{
			Metric: "deviceplane.device.Test_Metric_Two",
			Tags:   []string{"deviceplane.project:project-name"},
		},
		models.DatadogMetric{
			Metric: "deviceplane.device.Test_Metric_Three",
			Tags:   []string{"TEST:TRUE", "deviceplane.project:project-name"},
		},
	}, filteredMetrics)
}
