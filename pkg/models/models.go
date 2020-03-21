package models

import (
	"time"
)

type User struct {
	ID             string    `json:"id" yaml:"id"`
	CreatedAt      time.Time `json:"createdAt" yaml:"createdAt"`
	Name           string    `json:"name" yaml:"name"`
	InternalUserID *string   `json:"-" yaml:"-"`
	ExternalUserID *string   `json:"-" yaml:"-"`
	SuperAdmin     bool      `json:"superAdmin" yaml:"superAdmin"`
}

type UserFull struct {
	User
	ProviderName string `json:"providerName,omitempty"`
	ProviderID   string `json:"providerId,omitempty"`
	Email        string `json:"email"`
}

type ExternalUser struct {
	ID           string
	ProviderName string
	ProviderID   string
	Email        string
	Info         map[string]interface{}
}

type InternalUser struct {
	ID    string
	Email string
}

type RegistrationToken struct {
	ID             string
	CreatedAt      time.Time
	InternalUserID string
}

type PasswordRecoveryToken struct {
	ID             string
	CreatedAt      time.Time
	ExpiresAt      time.Time
	InternalUserID string
}

type Session struct {
	ID        string    `json:"id" yaml:"id"`
	CreatedAt time.Time `json:"createdAt" yaml:"createdAt"`
	UserID    string    `json:"userId" yaml:"userId"`
}

type UserAccessKey struct {
	ID          string    `json:"id" yaml:"id"`
	CreatedAt   time.Time `json:"createdAt" yaml:"createdAt"`
	UserID      string    `json:"userId" yaml:"userId"`
	Description string    `json:"description" yaml:"description"`
}

type UserAccessKeyWithValue struct {
	UserAccessKey
	Value string `json:"value" yaml:"value"`
}

type Project struct {
	ID            string    `json:"id" yaml:"id"`
	CreatedAt     time.Time `json:"createdAt" yaml:"createdAt"`
	Name          string    `json:"name" yaml:"name"`
	DatadogAPIKey *string   `json:"datadogApiKey" yaml:"datadogApiKey"`
}

type ProjectDeviceCounts struct {
	AllCount int `json:"allCount" yaml:"allCount"`
}

type ProjectApplicationCounts struct {
	AllCount int `json:"allCount" yaml:"allCount"`
}

type Role struct {
	ID          string    `json:"id" yaml:"id"`
	CreatedAt   time.Time `json:"createdAt" yaml:"createdAt"`
	ProjectID   string    `json:"projectId" yaml:"projectId"`
	Name        string    `json:"name" yaml:"name"`
	Description string    `json:"description" yaml:"description"`
	Config      string    `json:"config" yaml:"config"`
}

type Membership struct {
	UserID    string    `json:"userId" yaml:"userId"`
	ProjectID string    `json:"projectId" yaml:"projectId"`
	CreatedAt time.Time `json:"createdAt" yaml:"createdAt"`
}

type MembershipRoleBinding struct {
	UserID    string    `json:"userId" yaml:"userId"`
	RoleID    string    `json:"roleId" yaml:"roleId"`
	CreatedAt time.Time `json:"createdAt" yaml:"createdAt"`
	ProjectID string    `json:"projectId" yaml:"projectId"`
}

type ServiceAccount struct {
	ID          string    `json:"id" yaml:"id"`
	CreatedAt   time.Time `json:"createdAt" yaml:"createdAt"`
	ProjectID   string    `json:"projectId" yaml:"projectId"`
	Name        string    `json:"name" yaml:"name"`
	Description string    `json:"description" yaml:"description"`
}

type ServiceAccountAccessKey struct {
	ID               string    `json:"id" yaml:"id"`
	CreatedAt        time.Time `json:"createdAt" yaml:"createdAt"`
	ProjectID        string    `json:"projectId" yaml:"projectId"`
	ServiceAccountID string    `json:"serviceAccountId" yaml:"serviceAccountId"`
	Description      string    `json:"description" yaml:"description"`
}

type ServiceAccountAccessKeyWithValue struct {
	ServiceAccountAccessKey
	Value string `json:"value" yaml:"value"`
}

type ServiceAccountRoleBinding struct {
	ServiceAccountID string    `json:"serviceAccountId" yaml:"serviceAccountId"`
	RoleID           string    `json:"roleId" yaml:"roleId"`
	CreatedAt        time.Time `json:"createdAt" yaml:"createdAt"`
	ProjectID        string    `json:"projectId" yaml:"projectId"`
}

type Device struct {
	ID                   string            `json:"id" yaml:"id"`
	CreatedAt            time.Time         `json:"createdAt" yaml:"createdAt"`
	ProjectID            string            `json:"projectId" yaml:"projectId"`
	Name                 string            `json:"name" yaml:"name"`
	RegistrationTokenID  *string           `json:"registrationTokenId" yaml:"registrationTokenId"`
	DesiredAgentVersion  string            `json:"desiredAgentVersion" yaml:"desiredAgentVersion"`
	Info                 DeviceInfo        `json:"info" yaml:"info"`
	LastSeenAt           time.Time         `json:"lastSeenAt" yaml:"lastSeenAt"`
	Status               DeviceStatus      `json:"status" yaml:"status"`
	Labels               map[string]string `json:"labels" yaml:"labels"`
	EnvironmentVariables map[string]string `json:"environmentVariables" yaml:"environmentVariables"`
}

type DeviceStatus string

const (
	DeviceStatusOnline  = DeviceStatus("online")
	DeviceStatusOffline = DeviceStatus("offline")
)

type DeviceRegistrationToken struct {
	ID                   string            `json:"id" yaml:"id"`
	CreatedAt            time.Time         `json:"createdAt" yaml:"createdAt"`
	ProjectID            string            `json:"projectId" yaml:"projectId"`
	MaxRegistrations     *int              `json:"maxRegistrations" yaml:"maxRegistrations"`
	Name                 string            `json:"name" yaml:"name"`
	Description          string            `json:"description" yaml:"description"`
	Labels               map[string]string `json:"labels" yaml:"labels"`
	EnvironmentVariables map[string]string `json:"environmentVariables" yaml:"environmentVariables"`
}

type DevicesRegisteredWithTokenCount struct {
	AllCount int `json:"allCount" yaml:"allCount"`
}

type DeviceAccessKey struct {
	ID        string    `json:"id" yaml:"id"`
	CreatedAt time.Time `json:"createdAt" yaml:"createdAt"`
	ProjectID string    `json:"projectId" yaml:"projectId"`
	DeviceID  string    `json:"deviceId" yaml:"deviceId"`
}

type Connection struct {
	ID        string    `json:"id" yaml:"id"`
	CreatedAt time.Time `json:"createdAt" yaml:"createdAt"`
	ProjectID string    `json:"projectId" yaml:"projectId"`
	Name      string    `json:"name" yaml:"name"`
	Protocol  Protocol  `json:"protocol" yaml:"protocol"`
	Port      uint      `json:"port" yaml:"port"`
}

type Protocol string

const (
	ProtocolTCP  = Protocol("tcp")
	ProtocolHTTP = Protocol("http")
)

type Application struct {
	ID                    string                          `json:"id" yaml:"id"`
	CreatedAt             time.Time                       `json:"createdAt" yaml:"createdAt"`
	ProjectID             string                          `json:"projectId" yaml:"projectId"`
	Name                  string                          `json:"name" yaml:"name"`
	Description           string                          `json:"description" yaml:"description"`
	SchedulingRule        SchedulingRule                  `json:"schedulingRule" yaml:"schedulingRule"`
	MetricEndpointConfigs map[string]MetricEndpointConfig `json:"metricEndpointConfigs" yaml:"metricEndpointConfigs"`
}

type ApplicationDeviceCounts struct {
	AllCount int `json:"allCount" yaml:"allCount"`
}

type Release struct {
	ID                        string             `json:"id" yaml:"id"`
	Number                    uint32             `json:"number" yaml:"number"`
	CreatedAt                 time.Time          `json:"createdAt" yaml:"createdAt"`
	ProjectID                 string             `json:"projectId" yaml:"projectId"`
	ApplicationID             string             `json:"applicationId" yaml:"applicationId"`
	Config                    map[string]Service `json:"config" yaml:"config"`
	RawConfig                 string             `json:"rawConfig" yaml:"rawConfig"`
	CreatedByUserID           *string            `json:"createdByUserId" yaml:"createdByUserId"`
	CreatedByServiceAccountID *string            `json:"createdByServiceAccountId" yaml:"createdByServiceAccountId"`
}

type ReleaseDeviceCounts struct {
	AllCount int `json:"allCount" yaml:"allCount"`
}

type DeviceApplicationStatus struct {
	ProjectID        string `json:"projectId" yaml:"projectId"`
	DeviceID         string `json:"deviceId" yaml:"deviceId"`
	ApplicationID    string `json:"applicationId" yaml:"applicationId"`
	CurrentReleaseID string `json:"currentReleaseId" yaml:"currentReleaseId"`
}

type DeviceServiceStatus struct {
	ProjectID        string `json:"projectId" yaml:"projectId"`
	DeviceID         string `json:"deviceId" yaml:"deviceId"`
	ApplicationID    string `json:"applicationId" yaml:"applicationId"`
	Service          string `json:"service" yaml:"service"`
	CurrentReleaseID string `json:"currentReleaseId" yaml:"currentReleaseId"`
}

type DeviceServiceState struct {
	ProjectID     string       `json:"projectId" yaml:"projectId"`
	DeviceID      string       `json:"deviceId" yaml:"deviceId"`
	ApplicationID string       `json:"applicationId" yaml:"applicationId"`
	Service       string       `json:"service" yaml:"service"`
	State         ServiceState `json:"state" yaml:"state"`
	ErrorMessage  string       `json:"errorMessage" yaml:"errorMessage"`
}

type ServiceState string

const (
	ServiceStateUnknown                   ServiceState = "unknown"
	ServiceStatePullingImage              ServiceState = "pulling image"
	ServiceStateCreatingContainer         ServiceState = "creating container"
	ServiceStateStoppingPreviousContainer ServiceState = "stopping previous container"
	ServiceStateRemovingPreviousContainer ServiceState = "removing previous container"
	ServiceStateStartingContainer         ServiceState = "starting container"
	ServiceStateRunning                   ServiceState = "running"
	ServiceStateExited                    ServiceState = "exited"
)

var AllServiceStates = map[ServiceState]bool{
	ServiceStateUnknown:                   true,
	ServiceStatePullingImage:              true,
	ServiceStateCreatingContainer:         true,
	ServiceStateStoppingPreviousContainer: true,
	ServiceStateRemovingPreviousContainer: true,
	ServiceStateStartingContainer:         true,
	ServiceStateRunning:                   true,
	ServiceStateExited:                    true,
}

type ServiceStateCount struct {
	Count         int          `json:"count" yaml:"count"`
	CountErroring int          `json:"countErroring" yaml:"countErroring"`
	State         ServiceState `json:"state" yaml:"state"`
	Service       string       `json:"service" yaml:"service"`
	ApplicationID string       `json:"applicationId" yaml:"applicationId"`
}

type MembershipFull1 struct {
	Membership
	User    User        `json:"user" yaml:"user"`
	Project ProjectFull `json:"project" yaml:"project"`
}

type ProjectFull struct {
	Project
	DeviceCounts      ProjectDeviceCounts      `json:"deviceCounts" yaml:"deviceCounts"`
	ApplicationCounts ProjectApplicationCounts `json:"applicationCounts" yaml:"applicationCounts"`
}

type MembershipFull2 struct {
	Membership
	User  UserFull `json:"user" yaml:"user"`
	Roles []Role   `json:"roles" yaml:"roles"`
}

type ServiceAccountFull struct {
	ServiceAccount
	Roles []Role `json:"roles" yaml:"roles"`
}

type DeviceFull struct {
	Device
	ApplicationStatusInfo []DeviceApplicationStatusInfo `json:"applicationStatusInfo" yaml:"applicationStatusInfo"`
}

type DeviceApplicationStatusInfo struct {
	Application       Application                  `json:"application" yaml:"application"`
	ApplicationStatus *DeviceApplicationStatusFull `json:"applicationStatus" yaml:"applicationStatus"`
	ServiceStatuses   []DeviceServiceStatusFull    `json:"serviceStatuses" yaml:"serviceStatuses"`
	ServiceStates     []DeviceServiceState         `json:"serviceStates" yaml:"serviceStates"`
}

type DeviceApplicationStatusFull struct {
	DeviceApplicationStatus
	CurrentRelease Release `json:"currentRelease" yaml:"currentRelease"`
}

type DeviceServiceStatusFull struct {
	DeviceServiceStatus
	CurrentRelease Release `json:"currentRelease" yaml:"currentRelease"`
}

type ApplicationFull1 struct {
	Application
	LatestRelease      *Release                `json:"latestRelease" yaml:"latestRelease"`
	DeviceCounts       ApplicationDeviceCounts `json:"deviceCounts" yaml:"deviceCounts"`
	ServiceStateCounts []ServiceStateCount     `json:"serviceStateCounts" yaml:"serviceStateCounts"`
}

type DeviceRegistrationTokenFull struct {
	DeviceRegistrationToken
	DeviceCounts DevicesRegisteredWithTokenCount `json:"deviceCounts" yaml:"deviceCounts"`
}

type ReleaseFull struct {
	Release
	CreatedByUser           *User               `json:"createdByUser" yaml:"createdByUser"`
	CreatedByServiceAccount *ServiceAccount     `json:"createdByServiceAccount" yaml:"createdByServiceAccount"`
	DeviceCounts            ReleaseDeviceCounts `json:"deviceCounts" yaml:"deviceCounts"`
}

type Bundle struct {
	Applications        []FullBundledApplication  `json:"applications" yaml:"applications"`
	ApplicationStatuses []DeviceApplicationStatus `json:"applicationStatuses" yaml:"applicationStatuses"`
	ServiceStatuses     []DeviceServiceStatus     `json:"serviceStatuses" yaml:"serviceStatuses"`
	ServiceStates       []DeviceServiceState      `json:"serviceStates" yaml:"serviceStates"`

	DeviceID             string            `json:"deviceId" yaml:"deviceId"`
	DeviceName           string            `json:"deviceName" yaml:"deviceName"`
	EnvironmentVariables map[string]string `json:"environmentVariables" yaml:"environmentVariables"`
	DesiredAgentVersion  string            `json:"desiredAgentVersion" yaml:"desiredAgentVersion"`

	ServiceMetricsConfigs []ServiceMetricsConfig `json:"serviceMetricsConfig" yaml:"serviceMetricsConfig"`
	DeviceMetricsConfig   *DeviceMetricsConfig   `json:"deviceMetricsConfig" yaml:"deviceMetricsConfig"`
}

type BundledApplication struct {
	ID                    string                          `json:"id" yaml:"id"`
	ProjectID             string                          `json:"projectId" yaml:"projectId"`
	Name                  string                          `json:"name" yaml:"name"`
	MetricEndpointConfigs map[string]MetricEndpointConfig `json:"metricEndpointConfigs" yaml:"metricEndpointConfigs"`
}

type FullBundledApplication struct {
	Application   BundledApplication `json:"application" yaml:"application"`
	LatestRelease Release            `json:"latestRelease" yaml:"latestRelease"`
}

type DeviceInfo struct {
	AgentVersion string    `json:"agentVersion" yaml:"agentVersion"`
	IPAddress    string    `json:"ipAddress" yaml:"ipAddress"`
	OSRelease    OSRelease `json:"osRelease" yaml:"osRelease"`
}

type OSRelease struct {
	PrettyName string `json:"prettyName" yaml:"prettyName"`
	Name       string `json:"name" yaml:"name"`
	VersionID  string `json:"versionId" yaml:"versionId"`
	Version    string `json:"version" yaml:"version"`
	ID         string `json:"id" yaml:"id"`
	IDLike     string `json:"idLike" yaml:"idLike"`
}

const (
	DefaultMetricPort uint   = 2112
	DefaultMetricPath string = "/metrics"
)

type MetricEndpointConfig struct {
	Port uint   `json:"port" yaml:"port"`
	Path string `json:"path" yaml:"path"`
}
