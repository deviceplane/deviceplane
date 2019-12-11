package device

import (
	"context"
	"time"

	"github.com/deviceplane/deviceplane/cmd/deviceplane/cliutils"
	"github.com/deviceplane/deviceplane/cmd/deviceplane/global"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	sshTimeoutFlag *int = &[]int{0}[0]

	deviceArg *string = &[]string{""}[0]

	deviceMetricsServiceArg     *string = &[]string{""}[0]
	deviceMetricsApplicationArg *string = &[]string{""}[0]

	deviceOutputFlag *string = &[]string{""}[0]

	config *global.Config
)

func Initialize(c *global.Config) {
	config = c

	deviceCmd := c.App.Command("device", "Manage devices.")

	deviceListCmd := deviceCmd.Command("list", "List devices.")
	cliutils.AddFormatFlag(deviceOutputFlag, deviceListCmd,
		cliutils.FormatTable,
		cliutils.FormatYAML,
		cliutils.FormatJSON,
		cliutils.FormatJSONStream,
	)
	deviceListCmd.Action(deviceListAction)

	cliutils.GlobalAndCategorizedCmd(config.App, deviceCmd, func(attachmentPoint cliutils.HasCommand) {
		deviceSSHCmd := attachmentPoint.Command("ssh", "SSH into a device.")
		addDeviceArg(deviceSSHCmd)
		deviceSSHCmd.Flag("timeout", "Maximum length to attempt establishing a connection.").Default("60").IntVar(sshTimeoutFlag)
		deviceSSHCmd.Action(deviceSSHAction)
	})

	deviceInspectCmd := deviceCmd.Command("inspect", "Inspect a device's properties and labels.")
	addDeviceArg(deviceInspectCmd)
	cliutils.AddFormatFlag(deviceOutputFlag, deviceInspectCmd,
		cliutils.FormatYAML,
		cliutils.FormatJSON,
	)
	deviceInspectCmd.Action(deviceInspectAction)

	deviceMetricsCmd := deviceCmd.Command("metrics", "Get device metrics.")

	deviceMetricsHostCmd := deviceMetricsCmd.Command("host", "Get metrics on the device itself.")
	addDeviceArg(deviceMetricsHostCmd)
	deviceMetricsHostCmd.Action(deviceHostMetricsAction)

	deviceMetricsServiceCmd := deviceMetricsCmd.Command("service", "Get the metrics from a service running on the device.")
	addDeviceArg(deviceMetricsServiceCmd)
	deviceMetricsServiceCmd.Arg("application", "The application under which the service is running.").Required().StringVar(deviceMetricsApplicationArg)
	deviceMetricsServiceCmd.Arg("service", "The name of the service which is exposing a metrics endpoint.").Required().StringVar(deviceMetricsServiceArg)
	deviceMetricsServiceCmd.Action(deviceServiceMetricsAction)
}

func addDeviceArg(cmd *kingpin.CmdClause) *kingpin.ArgClause {
	arg := cmd.Arg("device", "Device name.").Required()
	arg.StringVar(deviceArg)
	arg.HintAction(func() []string {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		devices, err := config.APIClient.ListDevices(ctx, *config.Flags.Project)
		if err != nil {
			return []string{}
		}

		names := make([]string, len(devices))
		for _, d := range devices {
			names = append(names, d.Name)
		}
		return names
	})
	return arg
}
