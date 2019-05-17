package models

type SetApplicationSettingsRequest struct {
	ApplicationSettings ApplicationSettings `json:"applicationSettings"`
}

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

type SetDeviceApplicationStatusRequest struct {
	CurrentReleaseID string `json:"currentReleaseId"`
}

type SetDeviceServiceStatusRequest struct {
	CurrentReleaseID string `json:"currentReleaseId"`
}
