package cliutils

import (
	"fmt"
	"log"

	"github.com/deviceplane/deviceplane/cmd/deviceplane/global"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func RequireAccessKey(config *global.Config, c interface{}) interface{} {
	return RequireVariableForPreAction(
		config,
		config.Flags.AccessKey,
		fmt.Errorf("Access key not found as a flag, environment variable, or in config (%s)", *config.Flags.ConfigFile),
		c,
	)
}

func RequireProject(config *global.Config, c interface{}) interface{} {
	return RequireVariableForPreAction(
		config,
		config.Flags.Project,
		fmt.Errorf("Project not found as a flag, environment variable, or in config (%s)", *config.Flags.ConfigFile),
		c,
	)
}

func RequireVariableForPreAction(config *global.Config, variable *string, err error, c interface{}) interface{} {
	requirePreAction := func(c *kingpin.ParseContext) error {
		if c.Error() || !*config.ParsedCorrectly {
			return nil // Let kingpin's errors precede
		}
		if variable == nil || *variable == "" {
			return err
		}
		return nil
	}

	switch v := c.(type) {
	case *kingpin.CmdClause:
		return v.PreAction(requirePreAction)
	case *kingpin.ArgClause:
		return v.PreAction(requirePreAction)
	case *kingpin.FlagClause:
		return v.PreAction(requirePreAction)
	default:
		log.Fatal("Cannot require access key on this type")
		return nil
	}
}
