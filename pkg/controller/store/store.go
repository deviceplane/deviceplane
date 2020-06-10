package store

import (
	"context"
	"errors"

	"github.com/deviceplane/deviceplane/pkg/models"
)

type Users interface {
	InitializeUser(ctx context.Context, internalUserID, externalUserID *string) (*models.User, error)
	GetUser(ctx context.Context, id string) (*models.User, error)
	GetUserByInternalID(ctx context.Context, internalUserID string) (*models.User, error)
	GetUserByExternalID(ctx context.Context, externalUserID string) (*models.User, error)
	UpdateUserName(ctx context.Context, id, name string) (*models.User, error)
}

type InternalUsers interface {
	CreateInternalUser(ctx context.Context, email, passwordHash string) (*models.InternalUser, error)
	GetInternalUser(ctx context.Context, id string) (*models.InternalUser, error)
	LookupInternalUser(ctx context.Context, email string) (*models.InternalUser, error)
	ValidateInternalUser(ctx context.Context, id, passwordHash string) (*models.InternalUser, error)
	ValidateInternalUserWithEmail(ctx context.Context, email, passwordHash string) (*models.InternalUser, error)
	UpdateInternalUserPasswordHash(ctx context.Context, id, passwordHash string) (*models.InternalUser, error)
}

type ExternalUsers interface {
	CreateExternalUser(ctx context.Context, providerName, providerID, email string, info map[string]interface{}) (*models.ExternalUser, error)
	GetExternalUser(ctx context.Context, id string) (*models.ExternalUser, error)
	GetExternalUserByProviderID(ctx context.Context, providerName, providerID string) (*models.ExternalUser, error)
}

var ErrUserNotFound = errors.New("user not found")

type RegistrationTokens interface {
	CreateRegistrationToken(ctx context.Context, internalUserID, hash string) (*models.RegistrationToken, error)
	GetRegistrationToken(ctx context.Context, id string) (*models.RegistrationToken, error)
	ValidateRegistrationToken(ctx context.Context, hash string) (*models.RegistrationToken, error)
}

var ErrRegistrationTokenNotFound = errors.New("registration token not found")

type PasswordRecoveryTokens interface {
	CreatePasswordRecoveryToken(ctx context.Context, userID, hash string) (*models.PasswordRecoveryToken, error)
	GetPasswordRecoveryToken(ctx context.Context, id string) (*models.PasswordRecoveryToken, error)
	ValidatePasswordRecoveryToken(ctx context.Context, hash string) (*models.PasswordRecoveryToken, error)
}

var ErrPasswordRecoveryTokenNotFound = errors.New("password recovery token not found")

type Sessions interface {
	CreateSession(ctx context.Context, userID string, hash string) (*models.Session, error)
	GetSession(ctx context.Context, id string) (*models.Session, error)
	ValidateSession(ctx context.Context, hash string) (*models.Session, error)
	DeleteSession(ctx context.Context, id string) error
}

var ErrSessionNotFound = errors.New("session not found")

type UserAccessKeys interface {
	CreateUserAccessKey(ctx context.Context, userID string, hash, description string) (*models.UserAccessKey, error)
	GetUserAccessKey(ctx context.Context, id string) (*models.UserAccessKey, error)
	ValidateUserAccessKey(ctx context.Context, hash string) (*models.UserAccessKey, error)
	ListUserAccessKeys(ctx context.Context, userID string) ([]models.UserAccessKey, error)
	DeleteUserAccessKey(ctx context.Context, id string) error
}

var ErrUserAccessKeyNotFound = errors.New("user access key not found")

type Projects interface {
	CreateProject(ctx context.Context, name string) (*models.Project, error)
	GetProject(ctx context.Context, id string) (*models.Project, error)
	LookupProject(ctx context.Context, name string) (*models.Project, error)
	ListProjects(ctx context.Context) ([]models.Project, error)
	UpdateProject(ctx context.Context, id, name, datadogApiKey string) (*models.Project, error)
	DeleteProject(ctx context.Context, id string) error
}

var ErrProjectNotFound = errors.New("project not found")
var ErrProjectNameAlreadyInUse = errors.New("project name already in use")

type ProjectDeviceCounts interface {
	GetProjectDeviceCounts(ctx context.Context, projectID string) (*models.ProjectDeviceCounts, error)
}

type ProjectApplicationCounts interface {
	GetProjectApplicationCounts(ctx context.Context, projectID string) (*models.ProjectApplicationCounts, error)
}

type Roles interface {
	CreateRole(ctx context.Context, projectID, name, description, config string) (*models.Role, error)
	GetRole(ctx context.Context, id, projectID string) (*models.Role, error)
	LookupRole(ctx context.Context, name, projectID string) (*models.Role, error)
	ListRoles(ctx context.Context, projectID string) ([]models.Role, error)
	UpdateRole(ctx context.Context, id, projectID, name, description, config string) (*models.Role, error)
	DeleteRole(ctx context.Context, id, projectID string) error
}

var ErrRoleNotFound = errors.New("role not found")
var ErrRoleNameAlreadyInUse = errors.New("role name already in use")

type Memberships interface {
	CreateMembership(ctx context.Context, userID, projectID string) (*models.Membership, error)
	GetMembership(ctx context.Context, userID, projectID string) (*models.Membership, error)
	ListMembershipsByUser(ctx context.Context, userID string) ([]models.Membership, error)
	ListMembershipsByProject(ctx context.Context, projectID string) ([]models.Membership, error)
	DeleteMembership(ctx context.Context, userID, projectID string) error
}

var ErrMembershipNotFound = errors.New("membership not found")

type MembershipRoleBindings interface {
	CreateMembershipRoleBinding(ctx context.Context, userID, roleID, projectID string) (*models.MembershipRoleBinding, error)
	GetMembershipRoleBinding(ctx context.Context, userID, roleID, projectID string) (*models.MembershipRoleBinding, error)
	ListMembershipRoleBindings(ctx context.Context, userID, projectID string) ([]models.MembershipRoleBinding, error)
	DeleteMembershipRoleBinding(ctx context.Context, userID, roleID, projectID string) error
}

var ErrMembershipRoleBindingNotFound = errors.New("membership role binding not found")

type ServiceAccounts interface {
	CreateServiceAccount(ctx context.Context, projectID, name, description string) (*models.ServiceAccount, error)
	GetServiceAccount(ctx context.Context, id, projectID string) (*models.ServiceAccount, error)
	LookupServiceAccount(ctx context.Context, name, projectID string) (*models.ServiceAccount, error)
	ListServiceAccounts(ctx context.Context, projectID string) ([]models.ServiceAccount, error)
	UpdateServiceAccount(ctx context.Context, id, projectID, name, description string) (*models.ServiceAccount, error)
	DeleteServiceAccount(ctx context.Context, id, projectID string) error
}

var ErrServiceAccountNotFound = errors.New("service account not found")
var ErrServiceAccountNameAlreadyInUse = errors.New("service account name already in use")

type ServiceAccountAccessKeys interface {
	CreateServiceAccountAccessKey(ctx context.Context, projectID, serviceAccountID string, hash, description string) (*models.ServiceAccountAccessKey, error)
	GetServiceAccountAccessKey(ctx context.Context, id, projectID string) (*models.ServiceAccountAccessKey, error)
	ValidateServiceAccountAccessKey(ctx context.Context, hash string) (*models.ServiceAccountAccessKey, error)
	ListServiceAccountAccessKeys(ctx context.Context, projectID, serviceAccountID string) ([]models.ServiceAccountAccessKey, error)
	DeleteServiceAccountAccessKey(ctx context.Context, id, projectID string) error
}

var ErrServiceAccountAccessKeyNotFound = errors.New("service account access key not found")

type ServiceAccountRoleBindings interface {
	CreateServiceAccountRoleBinding(ctx context.Context, serviceAccountID, roleID, projectID string) (*models.ServiceAccountRoleBinding, error)
	GetServiceAccountRoleBinding(ctx context.Context, serviceAccountID, roleID, projectID string) (*models.ServiceAccountRoleBinding, error)
	ListServiceAccountRoleBindings(ctx context.Context, serviceAccountID, projectID string) ([]models.ServiceAccountRoleBinding, error)
	DeleteServiceAccountRoleBinding(ctx context.Context, serviceAccountID, roleID, projectID string) error
}

var ErrServiceAccountRoleBindingNotFound = errors.New("service account role binding not found")

type Devices interface {
	CreateDevice(ctx context.Context, projectID, name, registrationTokenID string, deviceLabels, deviceEnvironmentVariables map[string]string) (*models.Device, error)
	GetDevice(ctx context.Context, deviceID, projectID string) (*models.Device, error)
	LookupDevice(ctx context.Context, name, projectID string) (*models.Device, error)
	ListDevices(ctx context.Context, projectID, searchQuery string) ([]models.Device, error)
	UpdateDeviceName(ctx context.Context, deviceID, projectID, name string) (*models.Device, error)
	DeleteDevice(ctx context.Context, deviceID, projectID string) error
	SetDeviceInfo(ctx context.Context, deviceID, projectID string, deviceInfo models.DeviceInfo) (*models.Device, error)
	UpdateDeviceLastSeenAt(ctx context.Context, deviceID, projectID string) error
	ListAllDeviceLabelKeys(ctx context.Context, projectID string) ([]string, error)
	SetDeviceLabel(ctx context.Context, deviceID, projectID, key, value string) (*string, error)
	DeleteDeviceLabel(ctx context.Context, deviceID, projectID, key string) error
	SetDeviceEnvironmentVariable(ctx context.Context, deviceID, projectID, key, value string) (*string, error)
	DeleteDeviceEnvironmentVariable(ctx context.Context, deviceID, projectID, key string) error
}

var ErrDeviceNotFound = errors.New("device not found")
var ErrDeviceNameAlreadyInUse = errors.New("device name already in use")

type DeviceRegistrationTokens interface {
	CreateDeviceRegistrationToken(ctx context.Context, projectID, name, description string, maxRegistrations *int) (*models.DeviceRegistrationToken, error)
	GetDeviceRegistrationToken(ctx context.Context, tokenID, projectID string) (*models.DeviceRegistrationToken, error)
	LookupDeviceRegistrationToken(ctx context.Context, name, projectID string) (*models.DeviceRegistrationToken, error)
	ListDeviceRegistrationTokens(ctx context.Context, projectID string) ([]models.DeviceRegistrationToken, error)
	UpdateDeviceRegistrationToken(ctx context.Context, tokenID, projectID, name, description string, maxRegistrations *int) (*models.DeviceRegistrationToken, error)
	DeleteDeviceRegistrationToken(ctx context.Context, tokenID, projectID string) error
	SetDeviceRegistrationTokenLabel(ctx context.Context, tokenID, projectID, key, value string) (*string, error)
	DeleteDeviceRegistrationTokenLabel(ctx context.Context, tokenID, projectID, key string) error
	SetDeviceRegistrationTokenEnvironmentVariable(ctx context.Context, tokenID, projectID, key, value string) (*string, error)
	DeleteDeviceRegistrationTokenEnvironmentVariable(ctx context.Context, tokenID, projectID, key string) error
}

type DevicesRegisteredWithToken interface {
	GetDevicesRegisteredWithTokenCount(ctx context.Context, tokenID, projectID string) (*models.DevicesRegisteredWithTokenCount, error)
}

var ErrDeviceRegistrationTokenNotFound = errors.New("device registration token not found")
var ErrDeviceRegistrationTokenNameAlreadyInUse = errors.New("device registration token name already in use")

type DeviceAccessKeys interface {
	CreateDeviceAccessKey(ctx context.Context, projectID, deviceID, hash string) (*models.DeviceAccessKey, error)
	GetDeviceAccessKey(ctx context.Context, id, projectID string) (*models.DeviceAccessKey, error)
	ValidateDeviceAccessKey(ctx context.Context, projectID, hash string) (*models.DeviceAccessKey, error)
}

var ErrDeviceAccessKeyNotFound = errors.New("device access key not found")