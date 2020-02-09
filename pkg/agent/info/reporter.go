package info

import (
	"context"
	"time"

	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/agent/client"
	dpcontext "github.com/deviceplane/deviceplane/pkg/context"
	"github.com/deviceplane/deviceplane/pkg/models"
)

type Reporter struct {
	client       *client.Client // TODO: interface
	agentVersion string

	info models.DeviceInfo
}

func NewReporter(client *client.Client, agentVersion string) *Reporter {
	return &Reporter{
		client:       client,
		agentVersion: agentVersion,
	}
}

func (r *Reporter) Report() error {
	newInfo := r.readInfo()

	if newInfo != r.info {
		ctx, cancel := dpcontext.New(context.Background(), time.Minute)
		defer cancel()

		if err := r.client.SetDeviceInfo(ctx, models.SetDeviceInfoRequest{
			DeviceInfo: newInfo,
		}); err != nil {
			return err
		}

		r.info = newInfo
	}

	return nil
}

func (r *Reporter) readInfo() models.DeviceInfo {
	info := models.DeviceInfo{
		AgentVersion: r.agentVersion,
	}

	ipAddress, err := getIPAddress()
	if err == nil {
		info.IPAddress = ipAddress
	} else {
		log.WithError(err).Error("failed to get IP address")
	}

	osRelease, err := getOSRelease()
	if err == nil {
		info.OSRelease = *osRelease
	} else {
		log.WithError(err).Error("failed to get OS release")
	}

	return info
}
