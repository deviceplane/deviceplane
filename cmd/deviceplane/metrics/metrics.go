package metrics

import (
	"context"
	"fmt"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func hostMetricsAction(c *kingpin.ParseContext) error {
	metrics, err := config.APIClient.GetDeviceHostMetrics(context.TODO(), *config.Flags.Project, *deviceArgVar)
	if err != nil {
		return err
	}

	fmt.Println(*metrics)
	return nil
}

func userDefinedMetricsAction(c *kingpin.ParseContext) error {
	metrics, err := config.APIClient.GetDeviceServiceMetrics(
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
