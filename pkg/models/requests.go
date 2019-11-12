package models

type CreateReleaseRequest struct {
	Config string `json:"config" validate:"config"`
}

type ExecuteResponse struct {
	ExitCode int `json:"exitCode"`
}

type RegisterDeviceRequest struct {
	DeviceRegistrationTokenID string `json:"deviceRegistrationTokenId" validate:"id"`
}

type RegisterDeviceResponse struct {
	DeviceID             string `json:"deviceId"`
	DeviceAccessKeyValue string `json:"deviceAccessKeyValue"`
}

type SetDeviceInfoRequest struct {
	DeviceInfo DeviceInfo `json:"deviceInfo"` // TODO: validate
}

type SetDeviceApplicationStatusRequest struct {
	CurrentReleaseID string `json:"currentReleaseId" validate:"id"`
}

type SetDeviceServiceStatusRequest struct {
	CurrentReleaseID string `json:"currentReleaseId" validate:"id"`
}
