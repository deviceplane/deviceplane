package models

type CreateReleaseRequest struct {
	RawConfig string `json:"rawConfig" validate:"config"`
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

type SetDeviceServiceStateRequest struct {
	State        ServiceState `json:"state"`
	ErrorMessage string       `json:"errorMessage"`
}

type Auth0SsoRequest struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
	IdToken     string `json:"id_token"`
	Scope       string `json:"scope"`
	State       string `json:"state"`
	TokenType   string `json:"token_type"`
}
