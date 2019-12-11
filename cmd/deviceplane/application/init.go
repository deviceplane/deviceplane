package application

import (
	"github.com/deviceplane/deviceplane/cmd/deviceplane/cliutils"
	"github.com/deviceplane/deviceplane/cmd/deviceplane/global"
)

var (
	applicationArg        *string = &[]string{""}[0]
	applicationOutputFlag *string = &[]string{""}[0]

	applicationConfigOnlyFlag *bool = &[]bool{false}[0]

	applicationDeployFileArg *string = &[]string{""}[0]

	config *global.Config
)

func Initialize(c *global.Config) {
	config = c

	applicationCmd := config.App.Command("application", "Manage applications.")
	cliutils.RequireAccessKey(config, applicationCmd)
	cliutils.RequireProject(config, applicationCmd)

	applicationListCmd := applicationCmd.Command("list", "List applications.")
	cliutils.AddFormatFlag(applicationOutputFlag, applicationListCmd,
		cliutils.FormatTable,
		cliutils.FormatYAML,
		cliutils.FormatJSON,
		cliutils.FormatJSONStream,
	)
	applicationListCmd.Action(applicationListAction)

	applicationInspectCmd := applicationCmd.Command("inspect", "Inspect an application.")
	addApplicationArg(applicationInspectCmd)
	cliutils.AddFormatFlag(applicationOutputFlag, applicationInspectCmd,
		cliutils.FormatYAML,
		cliutils.FormatJSON,
	)
	applicationInspectCmd.Flag("config-only", "Only output the latest release config for an application.").Short('c').Default("false").BoolVar(applicationConfigOnlyFlag)
	applicationInspectCmd.Action(applicationInspectAction)

	applicationCreateCmd := applicationCmd.Command("create", "Create a new application.")
	addApplicationArg(applicationCreateCmd)
	applicationCreateCmd.Action(applicationCreateAction)

	cliutils.GlobalAndCategorizedCmd(config.App, applicationCmd, func(attachmentPoint cliutils.HasCommand) {
		applicationEditCmd := attachmentPoint.Command("edit", "Manually modify an application's config.")
		addApplicationArg(applicationEditCmd)
		applicationEditCmd.Action(applicationEditAction)

		applicationDeployCmd := attachmentPoint.Command("deploy", "Deploy an application's config from a file.")
		addApplicationArg(applicationDeployCmd)
		applicationDeployCmd.Arg("file", "File path of the yaml file to deploy.").Required().ExistingFileVar(applicationDeployFileArg)
		applicationDeployCmd.Action(applicationDeployAction)
	})

}
