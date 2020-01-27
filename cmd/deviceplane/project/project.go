package project

import (
	"context"
	"fmt"

	"github.com/deviceplane/deviceplane/cmd/deviceplane/cliutils"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func projectListAction(c *kingpin.ParseContext) error {
	projects, err := config.APIClient.ListProjects(context.TODO(), *config.Flags.Project)
	if err != nil {
		return err
	}

	if *projectOutputFlag == cliutils.FormatTable {
		table := cliutils.DefaultTable()
		table.SetHeader([]string{"Name", "Devices", "Applications", "Created"})
		for _, p := range projects {
			table.Append([]string{
				p.Name,
				fmt.Sprintf("%d", p.DeviceCounts.AllCount),
				fmt.Sprintf("%d", p.ApplicationCounts.AllCount),
				cliutils.DurafmtSince(p.CreatedAt).String() + " ago",
			})
		}
		table.Render()
		return nil
	}

	return cliutils.PrintWithFormat(projects, *projectOutputFlag)
}

func projectCreateAction(c *kingpin.ParseContext) error {
	project, err := config.APIClient.CreateProject(context.TODO(), *config.Flags.Project)
	if err != nil {
		return err
	}

	fmt.Printf("Project %s successfully created at %s!\n", project.Name, project.CreatedAt.Format("Mon Jan _2 15:04:05 2006"))

	return nil
}
