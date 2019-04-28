package agent

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/apex/log"
	agent_client "github.com/deviceplane/deviceplane/pkg/agent/client"
	"github.com/deviceplane/deviceplane/pkg/agent/info"
	"github.com/deviceplane/deviceplane/pkg/agent/supervisor"
	"github.com/deviceplane/deviceplane/pkg/engine"
	"github.com/pkg/errors"
)

const (
	accessKeyFilename = "access-key"
	deviceIDFilename  = "device-id"
)

type Agent struct {
	client            *agent_client.Client // TODO: interface
	engine            engine.Engine
	projectID         string
	registrationToken string
	stateDir          string
	supervisor        *supervisor.Supervisor
	infoReporter      *info.Reporter
}

func NewAgent(client *agent_client.Client, engine engine.Engine, projectID, registrationToken, stateDir string) *Agent {
	return &Agent{
		client:            client,
		engine:            engine,
		projectID:         projectID,
		registrationToken: registrationToken,
		stateDir:          stateDir,
		supervisor:        supervisor.NewSupervisor(engine),
		infoReporter:      info.NewReporter(client),
	}
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
	if err := ioutil.WriteFile(a.fileLocation(elem...), contents, 0644); err != nil {
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

	return nil
}

func (a *Agent) register() error {
	registerDeviceResponse, err := a.client.RegisterDevice(context.Background(), a.registrationToken)
	if err != nil {
		return err
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
	go a.runSupervisor()
	go a.runInfoReporter()
	select {}
}

func (a *Agent) runSupervisor() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		bundle, err := a.client.GetBundle(context.TODO())
		if err != nil {
			log.WithError(err).Error("get bundle")
			continue
		}

		for _, application := range bundle.Applications {
			if application.LatestRelease == nil {
				continue
			}

			if err := a.supervisor.SetApplication(application); err != nil {
				log.WithError(err).Error("set application")
				continue
			}
		}

		select {
		case <-ticker.C:
			continue
		}
	}
}

func (a *Agent) runInfoReporter() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		if err := a.infoReporter.Report(); err != nil {
			log.WithError(err).Error("report device info")
		}

		select {
		case <-ticker.C:
			continue
		}
	}
}
