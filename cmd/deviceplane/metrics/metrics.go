package metrics

import (
	"context"
	"fmt"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func deviceMetricsAction(c *kingpin.ParseContext) error {
	metrics, err := config.APIClient.GetDeviceMetrics(context.TODO(), *config.Flags.Project, *deviceArgVar)
	if err != nil {
		return err
	}

	fmt.Println(*metrics)
	return nil
}

func serviceMetricsAction(c *kingpin.ParseContext) error {
	metrics, err := config.APIClient.GetServiceMetrics(
		context.TODO(),
		*config.Flags.Project,
		*deviceArgVar, *metricsApplicationArgVar, *metricsServiceArgVar,
	)
	if err != nil {
		return err
	}

	fmt.Println(*metrics)
	return nil
}
