package application

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/deviceplane/deviceplane/cmd/deviceplane/cliutils"
	"github.com/deviceplane/deviceplane/pkg/interpolation"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func addApplicationArg(cmd *kingpin.CmdClause) *kingpin.ArgClause {
	arg := cmd.Arg("application", "Application name.").Required()
	arg.StringVar(applicationArg)
	arg.HintAction(func() []string {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		applications, err := config.APIClient.ListApplications(ctx, *config.Flags.Project)
		if err != nil {
			return []string{}
		}

		names := make([]string, len(applications))
		for _, app := range applications {
			names = append(names, app.Name)
		}
		return names
	})
	return arg
}

func applicationListAction(c *kingpin.ParseContext) error {
	applications, err := config.APIClient.ListApplications(context.TODO(), *config.Flags.Project)
	if err != nil {
		return err
	}

	if *applicationOutputFlag == cliutils.FormatTable {
		table := cliutils.DefaultTable()
		table.SetHeader([]string{"Name", "Description", "Created"})
		for _, app := range applications {
			table.Append([]string{app.Name, app.Description, cliutils.DurafmtSince(app.CreatedAt).String() + " ago"})
		}
		table.Render()
		return nil
	}

	return cliutils.PrintWithFormat(applications, *applicationOutputFlag)
}

func applicationCreateAction(c *kingpin.ParseContext) error {
	application, err := config.APIClient.CreateApplication(context.TODO(), *config.Flags.Project, *applicationArg)
	if err != nil {
		return err
	}

	fmt.Printf("Application %s successfully created at %s!\n", application.Name, application.CreatedAt.Format("Mon Jan _2 15:04:05 2006"))

	return nil
}

func applicationDeployAction(c *kingpin.ParseContext) error {
	yamlConfigBytes, err := ioutil.ReadFile(*applicationDeployFileArg)
	if err != nil {
		return err
	}

	finalYamlConfig, err := interpolation.Interpolate(string(yamlConfigBytes), os.Getenv)
	if err != nil {
		return err
	}

	release, err := config.APIClient.CreateRelease(context.TODO(), *config.Flags.Project, *applicationArg, finalYamlConfig)
	if err != nil {
		return err
	}

	fmt.Printf("Latest release %s for application %s successfully released at %s!\n", release.ID, *applicationArg, release.CreatedAt.Format("Mon Jan _2 15:04:05 2006"))

	return nil
}

func applicationInspectAction(c *kingpin.ParseContext) error {
	if applicationConfigOnlyFlag == nil || !*applicationConfigOnlyFlag {
		application, err := config.APIClient.GetApplication(context.TODO(), *config.Flags.Project, *applicationArg)
		if err != nil {
			return err
		}

		return cliutils.PrintWithFormat(application, *applicationOutputFlag)
	}

	release, err := config.APIClient.GetLatestRelease(context.TODO(), *config.Flags.Project, *applicationArg)
	if err != nil {
		return err
	}

	if *applicationOutputFlag == cliutils.FormatYAML {
		fmt.Println(release.RawConfig)
		return nil
	}

	return cliutils.PrintWithFormat(release.Config, *applicationOutputFlag)
}

func applicationEditAction(c *kingpin.ParseContext) error {
	release, err := config.APIClient.GetLatestRelease(context.TODO(), *config.Flags.Project, *applicationArg)
	if err != nil {
		return err
	}

	var yamlConfig string
	if release != nil {
		yamlConfig = release.RawConfig
	}

	tmpfile, err := ioutil.TempFile("", "")
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(yamlConfig)); err != nil {
		return err
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	cmd := exec.Command(editor, tmpfile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		fmt.Println("Edit cancelled, no changes made.")
		return nil
	}

	if err := tmpfile.Close(); err != nil {
		return err
	}

	yamlConfigFile, err := os.Open(tmpfile.Name())
	if err != nil {
		return err
	}

	yamlConfigBytes, err := ioutil.ReadAll(yamlConfigFile)
	if err != nil {
		return err
	}

	release, err = config.APIClient.CreateRelease(context.TODO(), *config.Flags.Project, *applicationArg, string(yamlConfigBytes))
	if err != nil {
		return err
	}

	fmt.Printf("Latest release %s for application %s successfully released at %s!\n", release.ID, *applicationArg, release.CreatedAt.Format("Mon Jan _2 15:04:05 2006"))

	return nil
}
