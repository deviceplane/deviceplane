package agent

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"time"

	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/agent/client"
	"github.com/deviceplane/deviceplane/pkg/agent/info"
	"github.com/deviceplane/deviceplane/pkg/agent/server/local"
	"github.com/deviceplane/deviceplane/pkg/agent/server/remote"
	"github.com/deviceplane/deviceplane/pkg/agent/service"
	"github.com/deviceplane/deviceplane/pkg/agent/updater"
	"github.com/deviceplane/deviceplane/pkg/agent/variables"
	"github.com/deviceplane/deviceplane/pkg/agent/variables/fsnotify"
	dpcontext "github.com/deviceplane/deviceplane/pkg/context"
	"github.com/deviceplane/deviceplane/pkg/file"
	"github.com/pkg/errors"
)

const (
	accessKeyFilename = "access-key"
	deviceIDFilename  = "device-id"
	bundleFilename    = "bundle"
)

var (
	errVersionNotSet = errors.New("version not set")
)

type Agent struct {
	client            *client.Client // TODO: interface
	variables         variables.Interface
	projectID         string
	registrationToken string
	confDir           string
	stateDir          string
	serverPort        int
	infoReporter      *info.Reporter
	localServer       *local.Server
	remoteServer      *remote.Server
	updater           *updater.Updater
}

func NewAgent(
	client *client.Client,
	projectID, registrationToken, confDir, stateDir, version, binaryPath string, serverPort int,
) (*Agent, error) {
	if version == "" {
		return nil, errVersionNotSet
	}

	if err := os.MkdirAll(confDir, 0700); err != nil {
		return nil, err
	}

	variables := fsnotify.NewVariables(confDir)
	if err := variables.Start(); err != nil {
		return nil, errors.Wrap(err, "start fsnotify variables")
	}

	service := service.NewService(variables, confDir)

	return &Agent{
		client:            client,
		variables:         variables,
		projectID:         projectID,
		registrationToken: registrationToken,
		confDir:           confDir,
		stateDir:          stateDir,
		serverPort:        serverPort,
		infoReporter:      info.NewReporter(client, version),
		localServer:       local.NewServer(service),
		remoteServer:      remote.NewServer(client, service),
		updater:           updater.NewUpdater(projectID, version, binaryPath),
	}, nil
}

func (a *Agent) fileLocation(elem ...string) string {
	return path.Join(
		append(
			[]string{a.stateDir, a.projectID},
			elem...,
		)...,
	)
}

func (a *Agent) writeFile(contents []byte, elem ...string) error {
	if err := os.MkdirAll(a.fileLocation(), 0700); err != nil {
		return err
	}
	if err := file.WriteFileAtomic(a.fileLocation(elem...), contents, 0644); err != nil {
		return err
	}
	return nil
}

func (a *Agent) Initialize() error {
	if _, err := os.Stat(a.fileLocation(accessKeyFilename)); err == nil {
		log.Info("device already registered")
	} else if os.IsNotExist(err) {
		log.Info("registering device")
		if err = a.register(); err != nil {
			return errors.Wrap(err, "failed to register device")
		}
	} else if err != nil {
		return errors.Wrap(err, "failed to check for access key")
	}

	accessKeyBytes, err := ioutil.ReadFile(a.fileLocation(accessKeyFilename))
	if err != nil {
		return errors.Wrap(err, "failed to read access key")
	}

	deviceIDBytes, err := ioutil.ReadFile(a.fileLocation(deviceIDFilename))
	if err != nil {
		return errors.Wrap(err, "failed to read device ID")
	}

	a.client.SetAccessKey(string(accessKeyBytes))
	a.client.SetDeviceID(string(deviceIDBytes))

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", a.serverPort))
		if err == nil {
			a.localServer.SetListener(listener)
			return nil
		}

		<-ticker.C
	}
}

func (a *Agent) register() error {
	ctx, cancel := dpcontext.New(context.Background(), time.Minute)
	defer cancel()

	registerDeviceResponse, err := a.client.RegisterDevice(ctx, a.registrationToken)
	if err != nil {
		return errors.Wrap(err, "failed to register device")
	}
	if err := a.writeFile([]byte(registerDeviceResponse.DeviceAccessKeyValue), accessKeyFilename); err != nil {
		return errors.Wrap(err, "failed to save access key")
	}
	if err := a.writeFile([]byte(registerDeviceResponse.DeviceID), deviceIDFilename); err != nil {
		return errors.Wrap(err, "failed to save device ID")
	}
	return nil
}

func (a *Agent) Run() {
	go a.runInfoReporter()
	go a.runRemoteServer()
	go a.runLocalServer()
	select {}
}

func (a *Agent) runInfoReporter() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		if err := a.infoReporter.Report(); err != nil {
			log.WithError(err).Error("report device info")
			goto cont
		}

	cont:
		select {
		case <-ticker.C:
			continue
		}
	}
}

func (a *Agent) runLocalServer() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		if err := a.localServer.Serve(); err != nil {
			log.WithError(err).Error("serve local device API")
			goto cont
		}

	cont:
		select {
		case <-ticker.C:
			continue
		}
	}
}

func (a *Agent) runRemoteServer() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		if err := a.remoteServer.Serve(); err != nil {
			log.WithError(err).Error("serve remote device API")
			goto cont
		}

	cont:
		select {
		case <-ticker.C:
			continue
		}
	}
}
