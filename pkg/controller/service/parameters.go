package service

import (
	"github.com/deviceplane/deviceplane/pkg/models"
)

type deviceLabelParameters []models.DeviceLabel

func (p deviceLabelParameters) Get(name string) (interface{}, error) {
	for _, deviceLabel := range []models.DeviceLabel(p) {
		if deviceLabel.Key == name {
			return deviceLabel.Value, nil
		}
	}
	return "", nil
}
