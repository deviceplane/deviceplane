package cliutils

import (
	"os"

	"github.com/deviceplane/deviceplane/cmd/deviceplane/global"
	"github.com/deviceplane/deviceplane/pkg/client"

	"github.com/olekukonko/tablewriter"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func InitializeAPIClient(config *global.Config) func(c *kingpin.ParseContext) error {
	return func(c *kingpin.ParseContext) error {
		config.APIClient = client.NewClient(*config.Flags.APIEndpoint, *config.Flags.AccessKey, nil)
		return nil
	}
}

func DefaultTable() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("-")
	table.SetColumnSeparator("|")
	table.SetRowSeparator("-")
	table.SetHeaderLine(true)
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetTablePadding(" ")
	table.SetNoWhiteSpace(false)
	return table
}
