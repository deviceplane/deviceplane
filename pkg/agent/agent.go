package agent

import (
	"context"
	"encoding/json"
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
	"github.com/deviceplane/deviceplane/pkg/agent/status"
	"github.com/deviceplane/deviceplane/pkg/agent/supervisor"
	"github.com/deviceplane/deviceplane/pkg/agent/updater"
	"github.com/deviceplane/deviceplane/pkg/agent/validator"
	"github.com/deviceplane/deviceplane/pkg/agent/validator/customcommands"
	"github.com/deviceplane/deviceplane/pkg/agent/validator/image"
	"github.com/deviceplane/deviceplane/pkg/agent/variables"
	"github.com/deviceplane/deviceplane/pkg/agent/variables/fsnotify"
	"github.com/deviceplane/deviceplane/pkg/engine"
	"github.com/deviceplane/deviceplane/pkg/file"
	"github.com/deviceplane/deviceplane/pkg/models"
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
	client                 *client.Client // TODO: interface
	variables              variables.Interface
	projectID              string
	registrationToken      string
	confDir                string
	stateDir               string
	serverPort             int
	supervisor             *supervisor.Supervisor
	statusGarbageCollector *status.GarbageCollector
	infoReporter           *info.Reporter
	localServer            *local.Server
	remoteServer           *remote.Server
	updater                *updater.Updater
}

func NewAgent(
	client *client.Client, engine engine.Engine,
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

	supervisor := supervisor.NewSupervisor(
		engine,
		func(ctx context.Context, applicationID, currentReleaseID string) error {
			return client.SetDeviceApplicationStatus(ctx, applicationID, models.SetDeviceApplicationStatusRequest{
				CurrentReleaseID: currentReleaseID,
			})
		},
		func(ctx context.Context, applicationID, service, currentReleaseID string) error {
			return client.SetDeviceServiceStatus(ctx, applicationID, service, models.SetDeviceServiceStatusRequest{
				CurrentReleaseID: currentReleaseID,
			})
		},
		[]validator.Validator{
			image.NewValidator(variables),
			customcommands.NewValidator(variables),
		},
	)

	service := service.NewService(variables, supervisor, engine, confDir)

	return &Agent{
		client:                 client,
		variables:              variables,
		projectID:              projectID,
		registrationToken:      registrationToken,
		confDir:                confDir,
		stateDir:               stateDir,
		serverPort:             serverPort,
		supervisor:             supervisor,
		statusGarbageCollector: status.NewGarbageCollector(client.DeleteDeviceApplicationStatus, client.DeleteDeviceServiceStatus),
		infoReporter:           info.NewReporter(client, version),
		localServer:            local.NewServer(service),
		remoteServer:           remote.NewServer(client, service),
		updater:                updater.NewUpdater(projectID, version, binaryPath),
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
	registerDeviceResponse, err := a.client.RegisterDevice(context.Background(), a.registrationToken)
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
	go a.runBundleApplier()
	go a.runInfoReporter()
	go a.runRemoteServer()
	go a.runLocalServer()
	select {}
}

func (a *Agent) runBundleApplier() {
	if bundle := a.loadSavedBundle(); bundle != nil {
		a.supervisor.SetApplications(bundle.Applications)
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		if bundle := a.downloadLatestBundle(); bundle != nil {
			a.supervisor.SetApplications(bundle.Applications)
			a.statusGarbageCollector.SetBundle(*bundle)
			a.updater.SetDesiredVersion(bundle.DesiredAgentVersion)
		}

		select {
		case <-ticker.C:
			continue
		}
	}
}

func (a *Agent) loadSavedBundle() *models.Bundle {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		if _, err := os.Stat(a.fileLocation(bundleFilename)); err == nil {
			savedBundleBytes, err := ioutil.ReadFile(a.fileLocation(bundleFilename))
			if err != nil {
				log.WithError(err).Error("read saved bundle")
				goto cont
			}

			var savedBundle models.Bundle
			if err = json.Unmarshal(savedBundleBytes, &savedBundle); err != nil {
				log.WithError(err).Error("discarding invalid saved bundle")
				return nil
			}

			return &savedBundle
		} else if os.IsNotExist(err) {
			return nil
		} else {
			log.WithError(err).Error("check if saved bundle exists")
			goto cont
		}

	cont:
		select {
		case <-ticker.C:
			continue
		}
	}
}

func (a *Agent) downloadLatestBundle() *models.Bundle {
	bundle, err := a.client.GetBundle(context.TODO())
	if err != nil {
		log.WithError(err).Error("get bundle")
		return nil
	}

	bundleBytes, err := json.Marshal(bundle)
	if err != nil {
		log.WithError(err).Error("marshal bundle")
		return nil
	}

	if err = a.writeFile(bundleBytes, bundleFilename); err != nil {
		log.WithError(err).Error("save bundle")
		return nil
	}

	return bundle
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
