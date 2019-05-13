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
	DeviceInfo DeviceInfo `json:"deviceInfo"`
}

type SetDeviceApplicationReleaseRequest struct {
	ApplicationID string `json:"applicationId"`
	ReleaseID     string `json:"releaseId"`
}

type SetDeviceApplicationServiceReleaseRequest struct {
	ApplicationID string `json:"applicationId"`
	Service       string `json:"service"`
	ReleaseID     string `json:"releaseId"`
}
