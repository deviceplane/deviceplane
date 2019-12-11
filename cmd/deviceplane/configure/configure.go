package configure

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/deviceplane/deviceplane/pkg/interpolation"
	"github.com/pkg/errors"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
)

type ConfigValues struct {
	AccessKey *string `yaml:"access-key,omitempty"`
	Project   *string `yaml:"project,omitempty"`
}

func populateEmptyValuesFromConfig(c *kingpin.ParseContext) (err error) {
	defer func() {
		if err != nil {
			err = errors.Wrap(err, "failed while parsing config")
		}
	}()

	// (This happens if kingpin has an error while parsing. Let it throw its
	// errors first, ours don't matter at that point)
	if c.Error() || gConfig.ParsedCorrectly == nil || !*gConfig.ParsedCorrectly {
		return nil
	}

	// This should normally be expanded by the shell,
	// but this is for the our default flag value,
	// which starts with "~/" but does not get expanded
	if strings.HasPrefix(*gConfig.Flags.ConfigFile, "~/") {
		usr, err := user.Current()
		if err != nil {
			return errors.Wrap(err, "failed to get home dir")
		}
		dir := usr.HomeDir
		expandedPath := filepath.Join(dir, (*gConfig.Flags.ConfigFile)[2:])
		gConfig.Flags.ConfigFile = &expandedPath
	}

	gcf, err := os.Open(*gConfig.Flags.ConfigFile)
	if err != nil {
		if !os.IsNotExist(err) {
			return errors.Wrap(err, "unknown file error")
		}

		// Create if not exists
		err := os.MkdirAll(filepath.Dir(*gConfig.Flags.ConfigFile), os.ModeDir)
		if err != nil && !os.IsExist(err) {
			return errors.Wrap(err, "could not create config directory")
		}
		err = os.Chmod(filepath.Dir(*gConfig.Flags.ConfigFile), 0700)
		if err != nil {
			return errors.Wrap(err, "failed to change config file perms")
		}
		gcf, err = os.Create(*gConfig.Flags.ConfigFile)
		if err != nil {
			return errors.Wrap(err, "failed to create config file")
		}
	}

	r := bufio.NewReader(gcf)
	configBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return errors.Wrap(err, "failed to read config file")
	}

	configString, err := interpolation.Interpolate(string(configBytes), os.Getenv)
	if err != nil {
		return errors.Wrap(err, "failed to interpolate config data")
	}

	var configValues ConfigValues
	err = yaml.Unmarshal([]byte(configString), &configValues)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal config file")
	}

	// Fill config in order of FLAG -> ENV -> CONFIG
	// The first two steps are handled automatically by kingpin
	if configValues.AccessKey != nil {
		if gConfig.Flags.AccessKey == nil || *gConfig.Flags.AccessKey == "" {
			*gConfig.Flags.AccessKey = *configValues.AccessKey
		}
	}
	if configValues.Project != nil {
		if gConfig.Flags.Project == nil || *gConfig.Flags.Project == "" {
			*gConfig.Flags.Project = *configValues.Project
		}
	}

	return nil
}

// Configure uses the existing value as a fallback
func configureAction(c *kingpin.ParseContext) error {
	reader := bufio.NewReader(os.Stdin)

	// Read input
	var extraAccessKeyMsg string
	if gConfig.Flags.AccessKey != nil && *gConfig.Flags.AccessKey != "" {
		extraAccessKeyMsg = fmt.Sprintf(` (or leave empty to use "%s")`, *gConfig.Flags.AccessKey)
	}
	fmt.Printf("Enter access key%s: \n>", extraAccessKeyMsg)
	rawAccessKey, _ := reader.ReadString('\n')

	var extraProjectMsg string
	if gConfig.Flags.Project != nil && *gConfig.Flags.Project != "" {
		extraProjectMsg = fmt.Sprintf(` (or leave empty to use "%s")`, *gConfig.Flags.Project)
	}
	fmt.Printf("Enter project%s: \n>", extraProjectMsg)
	rawProject, _ := reader.ReadString('\n')

	// Clean input
	accessKey := strings.TrimSpace(rawAccessKey)
	project := strings.TrimSpace(rawProject)

	// Replace input if needed
	if accessKey == "" {
		accessKey = *gConfig.Flags.AccessKey
	}
	if project == "" {
		project = *gConfig.Flags.Project
	}

	fmt.Printf("Configuring with access key (%s) and project (%s)\n", accessKey, project)

	// Actually configure
	configValues := ConfigValues{
		AccessKey: &accessKey,
		Project:   &project,
	}

	configBytes, err := yaml.Marshal(configValues)
	if err != nil {
		return errors.Wrap(err, "failed to serialize config")
	}

	err = ioutil.WriteFile(*gConfig.Flags.ConfigFile, configBytes, 0700)
	if err != nil {
		return errors.Wrap(err, "failed to write config to disk")
	}
	return nil
}
