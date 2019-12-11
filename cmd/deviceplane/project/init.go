package project

import (
	"github.com/deviceplane/deviceplane/cmd/deviceplane/cliutils"
	"github.com/deviceplane/deviceplane/cmd/deviceplane/global"
)

var (
	projectOutputFlag *string = &[]string{""}[0]

	config *global.Config
)

func Initialize(c *global.Config) {
	config = c

	projectCmd := config.App.Command("project", "Manage projects.")

	projectListCmd := projectCmd.Command("list", "List projects.")
	cliutils.AddFormatFlag(projectOutputFlag, projectListCmd,
		cliutils.FormatTable,
		cliutils.FormatYAML,
		cliutils.FormatJSON,
		cliutils.FormatJSONStream,
	)
	projectListCmd.Action(projectListAction)

	projectCreateCmd := projectCmd.Command("create", "Create a new project.")
	projectCreateCmd.Action(projectCreateAction)
}
