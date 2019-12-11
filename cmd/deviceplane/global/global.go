package global

import (
	"net/url"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/deviceplane/deviceplane/pkg/client"
)

type Config struct {
	App             *kingpin.Application
	ParsedCorrectly *bool
	Flags           ConfigFlags
	APIClient       *client.Client
}

type ConfigFlags struct {
	APIEndpoint **url.URL
	AccessKey   *string
	Project     *string
	ConfigFile  *string
}
