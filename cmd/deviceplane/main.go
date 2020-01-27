package main

import (
	"os"

	"github.com/deviceplane/deviceplane/cmd/deviceplane/application"
	"github.com/deviceplane/deviceplane/cmd/deviceplane/cliutils"
	"github.com/deviceplane/deviceplane/cmd/deviceplane/configure"
	"github.com/deviceplane/deviceplane/cmd/deviceplane/device"
	"github.com/deviceplane/deviceplane/cmd/deviceplane/global"
	"github.com/deviceplane/deviceplane/cmd/deviceplane/metrics"
	"github.com/deviceplane/deviceplane/cmd/deviceplane/project"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("deviceplane", "The Deviceplane CLI.").UsageTemplate(cliutils.CustomTemplate).Version("dev")

	config = global.Config{
		App:             app,
		ParsedCorrectly: app.Flag("internal-parsing-validator", "").Hidden().Default("true").Bool(),

		Flags: global.ConfigFlags{
			APIEndpoint: app.Flag("url", "API Endpoint.").Hidden().Default("https://cloud.deviceplane.com:443/api").URL(),
			AccessKey:   app.Flag("access-key", "Access key used for authentication. (env: DEVICEPLANE_ACCESS_KEY)").Envar("DEVICEPLANE_ACCESS_KEY").String(),
			Project:     app.Flag("project", "Project name. (env: DEVICEPLANE_PROJECT)").Envar("DEVICEPLANE_PROJECT").String(),
			ConfigFile:  app.Flag("config", "Config file to use.").Default("~/.deviceplane/config").String(),
		},

		APIClient: nil,
	}
)

func main() {
	configure.Initialize(&config)
	project.Initialize(&config)
	application.Initialize(&config)
	device.Initialize(&config)
	metrics.Initialize(&config)

	app.PreAction(cliutils.InitializeAPIClient(&config))
	kingpin.MustParse(app.Parse(os.Args[1:]))
}
