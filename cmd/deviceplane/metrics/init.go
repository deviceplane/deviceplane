package metrics

import (
	"context"
	"fmt"
	"time"

	"github.com/deviceplane/deviceplane/cmd/deviceplane/global"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	deviceArgVar *string = &[]string{""}[0]

	metricsServiceArgVar     *string = &[]string{""}[0]
	metricsApplicationArgVar *string = &[]string{""}[0]

	config *global.Config
)

func Initialize(c *global.Config) {
	config = c

	metricsCmd := c.App.Command("metrics", "Get device and service metrics.")

	metricsHostCmd := metricsCmd.Command("host", "Get metrics that tell you about the device itself.")
	addDeviceArg(metricsHostCmd)
	metricsHostCmd.Action(hostMetricsAction)

	metricsUserDefinedCmd := metricsCmd.Command("user-defined", "Get user-defined metrics from a device.")
	metricsApplicationArg := metricsUserDefinedCmd.Arg("application", "Application name.").Required()
	metricsApplicationArg.StringVar(metricsApplicationArgVar)
	metricsApplicationArg.HintAction(func() []string {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		applications, err := config.APIClient.ListApplications(ctx, *config.Flags.Project)
		if err != nil {
			return nil
		}

		var appnames []string
		for _, a := range applications {
			appnames = append(appnames, a.Name)
		}
		fmt.Println("-") // TODO: find out kingpin won't autocomplete without this
		return appnames
	})
	metricsServiceArg := metricsUserDefinedCmd.Arg("service", "The name of the service exposing the metrics endpoint.").Required()
	metricsServiceArg.StringVar(metricsServiceArgVar)
	metricsServiceArg.HintAction(func() []string {
		if metricsApplicationArgVar == nil || *metricsApplicationArgVar == "" {
			return nil
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		latestRelease, err := config.APIClient.GetLatestRelease(ctx, *config.Flags.Project, *metricsApplicationArgVar)
		if err != nil {
			return nil
		}

		var services []string
		for k, _ := range latestRelease.Config {
			services = append(services, k)
		}
		fmt.Println("-") // TODO: find out kingpin won't autocomplete without this
		return services
	})
	addDeviceArg(metricsUserDefinedCmd)
	metricsUserDefinedCmd.Action(userDefinedMetricsAction)
}

func addDeviceArg(cmd *kingpin.CmdClause) *kingpin.ArgClause {
	arg := cmd.Arg("device", "Device name.").Required()
	arg.StringVar(deviceArgVar)
	arg.HintAction(func() []string {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		devices, err := config.APIClient.ListDevices(ctx, nil, *config.Flags.Project)
		if err != nil {
			return nil
		}

		names := make([]string, len(devices))
		for _, d := range devices {
			names = append(names, d.Name)
		}
		return names
	})
	return arg
}
