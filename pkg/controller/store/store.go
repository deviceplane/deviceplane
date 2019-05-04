package store

import (
	"context"
	"errors"

	"github.com/deviceplane/deviceplane/pkg/models"
)

type Users interface {
	CreateUser(ctx context.Context, email, passwordHash, firstName, lastName string) (*models.User, error)
	GetUser(ctx context.Context, id string) (*models.User, error)
	ValidateUser(ctx context.Context, email, passwordHash string) (*models.User, error)
	MarkRegistrationCompleted(ctx context.Context, id string) (*models.User, error)
}

var ErrUserNotFound = errors.New("user not found")

type RegistrationTokens interface {
	CreateRegistrationToken(ctx context.Context, userID, hash string) (*models.RegistrationToken, error)
	GetRegistrationToken(ctx context.Context, id string) (*models.RegistrationToken, error)
	ValidateRegistrationToken(ctx context.Context, hash string) (*models.RegistrationToken, error)
}

var ErrRegistrationTokenNotFound = errors.New("registration token not found")

type Sessions interface {
	CreateSession(ctx context.Context, userID string, hash string) (*models.Session, error)
	GetSession(ctx context.Context, id string) (*models.Session, error)
	ValidateSession(ctx context.Context, hash string) (*models.Session, error)
	DeleteSession(ctx context.Context, id string) error
}

var ErrSessionNotFound = errors.New("session not found")

type AccessKeys interface {
	CreateAccessKey(ctx context.Context, userID string, hash string) (*models.AccessKey, error)
	GetAccessKey(ctx context.Context, id string) (*models.AccessKey, error)
	ValidateAccessKey(ctx context.Context, hash string) (*models.AccessKey, error)
}

var ErrAccessKeyNotFound = errors.New("access key not found")

type Projects interface {
	CreateProject(ctx context.Context, name string) (*models.Project, error)
	GetProject(ctx context.Context, id string) (*models.Project, error)
}

var ErrProjectNotFound = errors.New("project not found")

type ProjectDeviceCounts interface {
	GetProjectDeviceCounts(ctx context.Context, projectID string) (*models.ProjectDeviceCounts, error)
}

type ProjectApplicationCounts interface {
	GetProjectApplicationCounts(ctx context.Context, projectID string) (*models.ProjectApplicationCounts, error)
}

type Memberships interface {
	CreateMembership(ctx context.Context, userID, projectID, membershipType string) (*models.Membership, error)
	GetMembership(ctx context.Context, userID, projectID string) (*models.Membership, error)
	ListMembershipsByUser(ctx context.Context, userID string) ([]models.Membership, error)
	ListMembershipsByProject(ctx context.Context, projectID string) ([]models.Membership, error)
}

var ErrMembershipNotFound = errors.New("membership not found")

type Devices interface {
	CreateDevice(ctx context.Context, projectID string) (*models.Device, error)
	GetDevice(ctx context.Context, id, projectID string) (*models.Device, error)
	ListDevices(ctx context.Context, projectID string) ([]models.Device, error)
	SetDeviceInfo(ctx context.Context, id, projectID string, deviceInfo models.DeviceInfo) (*models.Device, error)
}

var ErrDeviceNotFound = errors.New("device not found")

type DeviceLabels interface {
	SetDeviceLabel(ctx context.Context, key, deviceID, projectID, value string) (*models.DeviceLabel, error)
	GetDeviceLabel(ctx context.Context, key, deviceID, projectID string) (*models.DeviceLabel, error)
	ListDeviceLabels(ctx context.Context, deviceID, projectID string) ([]models.DeviceLabel, error)
	DeleteDeviceLabel(ctx context.Context, key, deviceID, projectID string) error
}

var ErrDeviceLabelNotFound = errors.New("device label not found")

type DeviceRegistrationTokens interface {
	CreateDeviceRegistrationToken(ctx context.Context, projectID string) (*models.DeviceRegistrationToken, error)
	GetDeviceRegistrationToken(ctx context.Context, id, projectID string) (*models.DeviceRegistrationToken, error)
	BindDeviceRegistrationToken(ctx context.Context, id, projectID, deviceAccessKeyID string) (*models.DeviceRegistrationToken, error)
}

var ErrDeviceRegistrationTokenNotFound = errors.New("device registration token not found")

type DeviceAccessKeys interface {
	CreateDeviceAccessKey(ctx context.Context, projectID, deviceID, hash string) (*models.DeviceAccessKey, error)
	GetDeviceAccessKey(ctx context.Context, id, projectID string) (*models.DeviceAccessKey, error)
	ValidateDeviceAccessKey(ctx context.Context, projectID, hash string) (*models.DeviceAccessKey, error)
}

var ErrDeviceAccessKeyNotFound = errors.New("device access key not found")

type Applications interface {
	CreateApplication(ctx context.Context, projectID, name string) (*models.Application, error)
	GetApplication(ctx context.Context, id, projectID string) (*models.Application, error)
	ListApplications(ctx context.Context, projectID string) ([]models.Application, error)
}

var ErrApplicationNotFound = errors.New("application not found")

type Releases interface {
	CreateRelease(ctx context.Context, projectID, applicationID, config string) (*models.Release, error)
	GetRelease(ctx context.Context, id, projectID string) (*models.Release, error)
	GetLatestRelease(ctx context.Context, projectID, applicationID string) (*models.Release, error)
	ListReleases(ctx context.Context, projectID, applicationID string) ([]models.Release, error)
}

var ErrReleaseNotFound = errors.New("release not found")
