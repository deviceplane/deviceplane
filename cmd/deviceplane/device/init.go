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

	deviceFilterListFlag *[]string = &[][]string{[]string{}}[0]

	deviceOutputFlag *string = &[]string{""}[0]

	sshCommandsArg *[]string = &[][]string{[]string{}}[0]

	config *global.Config
)

func Initialize(c *global.Config) {
	config = c

	deviceCmd := c.App.Command("device", "Manage devices.")

	deviceListCmd := deviceCmd.Command("list", "List devices.")
	deviceListCmd.Flag("filter", `Label key/values used to filter devices. e.g. "--filter status=online --filter labels.location=hq2"`).StringsVar(deviceFilterListFlag)
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
		deviceSSHCmd.Arg("ssh-commands", "If provided, runs commands, prints output, and exits after SSH completes.").StringsVar(sshCommandsArg)
	})

	deviceInspectCmd := deviceCmd.Command("inspect", "Inspect a device's properties and labels.")
	addDeviceArg(deviceInspectCmd)
	cliutils.AddFormatFlag(deviceOutputFlag, deviceInspectCmd,
		cliutils.FormatYAML,
		cliutils.FormatJSON,
	)
	deviceInspectCmd.Action(deviceInspectAction)

	cliutils.GlobalAndCategorizedCmd(config.App, deviceCmd, func(attachmentPoint cliutils.HasCommand) {
		deviceRebootCmd := attachmentPoint.Command("reboot", "Reboot a device.")
		addDeviceArg(deviceRebootCmd)
		deviceRebootCmd.Action(deviceRebootAction)
	})
}

func addDeviceArg(cmd *kingpin.CmdClause) *kingpin.ArgClause {
	arg := cmd.Arg("device", "Device name.").Required()
	arg.StringVar(deviceArg)
	arg.HintAction(func() []string {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		devices, err := config.APIClient.ListDevices(ctx, nil, *config.Flags.Project)
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
