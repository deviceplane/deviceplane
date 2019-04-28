package info

import (
	"context"

	"github.com/apex/log"
	agent_client "github.com/deviceplane/deviceplane/pkg/agent/client"
	"github.com/deviceplane/deviceplane/pkg/models"
)

type Reporter struct {
	client *agent_client.Client // TODO: interface
	info   models.DeviceInfo
}

func NewReporter(client *agent_client.Client) *Reporter {
	return &Reporter{
		client: client,
	}
}

func (r *Reporter) Report() error {
	newInfo := r.readInfo()
	if newInfo != r.info {
		if err := r.client.SetDeviceInfo(context.TODO(), models.SetDeviceInfoRequest{
			DeviceInfo: newInfo,
		}); err != nil {
			return err
		}
		r.info = newInfo
	}
	return nil
}

func (r *Reporter) readInfo() models.DeviceInfo {
	var info models.DeviceInfo

	ipAddress, err := getIPAddress()
	if err == nil {
		info.IPAddress = ipAddress
	} else {
		log.WithError(err).Error("failed to get IP address")
	}

	return info
}
