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

func ValidateClientAndServerProcessing(
	t *testing.T,
	mp metricProcessor,
	metrics []models.DatadogMetric,
	exposedMetrics []models.ExposedMetric,
	project *models.Project,
	device *models.Device,
	expectedFilteredMetrics []models.DatadogMetric,
) {
	t.Helper()

	filteredMetrics := mp(metrics, exposedMetrics, project, device)
	assert.Equal(t, expectedFilteredMetrics, filteredMetrics, "client-side test")

	doubleFilteredMetrics := mp(filteredMetrics, exposedMetrics, project, device)
	assert.Equal(t, expectedFilteredMetrics, doubleFilteredMetrics, "server-side test")
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

	ValidateClientAndServerProcessing(t,
		ProcessDeviceMetrics,
		metrics,
		exposedMetrics,
		&project, &device,
		[]models.DatadogMetric{
			models.DatadogMetric{
				Metric: "deviceplane.device.Test_Metric",
				Tags:   []string{"TEST:TRUE", "deviceplane.project:project-name"},
			},
		})
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

	ValidateClientAndServerProcessing(t,
		ProcessDeviceMetrics,
		metrics,
		exposedMetrics,
		&project, &device,
		[]models.DatadogMetric{
			models.DatadogMetric{
				Metric: "deviceplane.device.Test_Metric",
				Tags:   []string{"TEST:TRUE", "deviceplane.project:project-name"},
			},
			models.DatadogMetric{
				Metric: "deviceplane.device.Test_Metric_Two",
				Tags:   []string{"TEST:TRUE", "deviceplane.project:project-name"},
			},
		},
	)
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

	ValidateClientAndServerProcessing(t,
		ProcessDeviceMetrics,
		metrics,
		exposedMetrics,
		&project, &device,
		[]models.DatadogMetric{
			models.DatadogMetric{
				Metric: "deviceplane.device.Test_Metric",
				Tags:   []string{"TEST:TRUE", "deviceplane.labels.company:nasa", "deviceplane.project:project-name"},
			},
		},
	)
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

	ValidateClientAndServerProcessing(t,
		ProcessDeviceMetrics,
		metrics,
		exposedMetrics,
		&project, &device,
		[]models.DatadogMetric{
			models.DatadogMetric{
				Metric: "deviceplane.device.Test_Metric",
				Tags:   []string{"TEST:TRUE", "deviceplane.device:device-name", "deviceplane.project:project-name"},
			},
		},
	)
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

	ValidateClientAndServerProcessing(t,
		ProcessDeviceMetrics,
		metrics,
		exposedMetrics,
		&project, &device,
		[]models.DatadogMetric{
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
		},
	)
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
		models.DatadogMetric{
			Metric: "Test_Metric_Four",
			Tags:   []string{"TEST:TRUE"},
		},
		models.DatadogMetric{
			Metric: "Test_Metric_Five",
			Tags:   []string{"TEST:FALSE"},
		},
	}
	exposedMetrics := []models.ExposedMetric{
		models.ExposedMetric{
			Name: "Test_Metric",
			WhitelistedTags: []models.WhitelistedTag{
				models.WhitelistedTag{
					Key: "TEST",
				},
			},
		},
		models.ExposedMetric{
			Name: "Test_Metric_Two",
			WhitelistedTags: []models.WhitelistedTag{
				models.WhitelistedTag{
					Key: "WEST",
				},
			},
		},
		models.ExposedMetric{
			Name:            "Test_Metric_Three",
			WhitelistedTags: []models.WhitelistedTag{},
		},
		models.ExposedMetric{
			Name: "Test_Metric_Four",
			WhitelistedTags: []models.WhitelistedTag{
				models.WhitelistedTag{
					Key:    "TEST",
					Values: []string{"TRUE"},
				},
			},
		},
		models.ExposedMetric{
			Name: "Test_Metric_Five",
			WhitelistedTags: []models.WhitelistedTag{
				models.WhitelistedTag{
					Key:    "TEST",
					Values: []string{"TRUE"},
				},
			},
		},
	}

	ValidateClientAndServerProcessing(t,
		ProcessDeviceMetrics,
		metrics,
		exposedMetrics,
		&project, &device,
		[]models.DatadogMetric{
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
			models.DatadogMetric{
				Metric: "deviceplane.device.Test_Metric_Four",
				Tags:   []string{"TEST:TRUE", "deviceplane.project:project-name"},
			},
			models.DatadogMetric{
				Metric: "deviceplane.device.Test_Metric_Five",
				Tags:   []string{"deviceplane.project:project-name"},
			},
		},
	)
}
