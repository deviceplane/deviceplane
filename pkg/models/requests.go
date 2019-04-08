package models

type CreateRelease struct {
	Config string `json:"config"`
}

type RegisterDeviceRequest struct {
	DeviceRegistrationTokenID string `json:"deviceRegistrationTokenId"`
}

type RegisterDeviceResponse struct {
	DeviceID             string `json:"deviceId"`
	DeviceAccessKeyValue string `json:"deviceAccessKeyValue"`
}

type SetDeviceInfoRequest struct {
	Info map[string]interface{} `json:"info"`
}
