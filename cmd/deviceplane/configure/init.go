package configure

import (
	"github.com/deviceplane/deviceplane/cmd/deviceplane/global"
)

var (
	gConfig *global.Config
)

func Initialize(c *global.Config) {
	gConfig = c

	// Global initialization
	c.App.PreAction(populateEmptyValuesFromConfig)

	// Commands
	configureCmd := c.App.Command("configure", "Configure this CLI utility.")
	configureCmd.Action(configureAction)
}
