package models

import (
	"time"
)

type User struct {
	ID                    string    `json:"id"`
	CreatedAt             time.Time `json:"createdAt"`
	Email                 string    `json:"email"`
	FirstName             string    `json:"firstName"`
	LastName              string    `json:"lastName"`
	Company               string    `json:"company"`
	RegistrationCompleted bool      `json:"registrationCompleted"`
	SuperAdmin            bool      `json:"superAdmin"`
}

type RegistrationToken struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UserID    string    `json:"userId"`
}

type PasswordRecoveryToken struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
	UserID    string    `json:"userId"`
}

type Session struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UserID    string    `json:"userId"`
}

type UserAccessKey struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	UserID      string    `json:"userId"`
	Description string    `json:"description"`
}

type UserAccessKeyWithValue struct {
	UserAccessKey
	Value string `json:"value"`
}

type Project struct {
	ID            string    `json:"id"`
	CreatedAt     time.Time `json:"createdAt"`
	Name          string    `json:"name"`
	DatadogAPIKey *string   `json:"datadogApiKey"`
}

type ProjectDeviceCounts struct {
	AllCount int `json:"allCount"`
}

type ProjectApplicationCounts struct {
	AllCount int `json:"allCount"`
}

type Role struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	ProjectID   string    `json:"projectId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Config      string    `json:"config"`
}

type Membership struct {
	UserID    string    `json:"userId"`
	ProjectID string    `json:"projectId"`
	CreatedAt time.Time `json:"createdAt"`
}

type MembershipRoleBinding struct {
	UserID    string    `json:"userId"`
	RoleID    string    `json:"roleId"`
	CreatedAt time.Time `json:"createdAt"`
	ProjectID string    `json:"projectId"`
}

type ServiceAccount struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	ProjectID   string    `json:"projectId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type ServiceAccountAccessKey struct {
	ID               string    `json:"id"`
	CreatedAt        time.Time `json:"createdAt"`
	ProjectID        string    `json:"projectId"`
	ServiceAccountID string    `json:"serviceAccountId"`
	Description      string    `json:"description"`
}

type ServiceAccountAccessKeyWithValue struct {
	ServiceAccountAccessKey
	Value string `json:"value"`
}

type ServiceAccountRoleBinding struct {
	ServiceAccountID string    `json:"serviceAccountId"`
	RoleID           string    `json:"roleId"`
	CreatedAt        time.Time `json:"createdAt"`
	ProjectID        string    `json:"projectId"`
}

type Device struct {
	ID                  string            `json:"id"`
	CreatedAt           time.Time         `json:"createdAt"`
	ProjectID           string            `json:"projectId"`
	Name                string            `json:"name"`
	RegistrationTokenID *string           `json:"registrationTokenId"`
	DesiredAgentSpec    string            `json:"desiredAgentSpec"`
	DesiredAgentVersion string            `json:"desiredAgentVersion"`
	Info                DeviceInfo        `json:"info"`
	LastSeenAt          time.Time         `json:"lastSeenAt"`
	Status              DeviceStatus      `json:"status"`
	Labels              map[string]string `json:"labels"`
}

type DeviceStatus string

const (
	DeviceStatusOnline  = DeviceStatus("online")
	DeviceStatusOffline = DeviceStatus("offline")
)

type DeviceRegistrationToken struct {
	ID               string            `json:"id"`
	CreatedAt        time.Time         `json:"createdAt"`
	ProjectID        string            `json:"projectId"`
	MaxRegistrations *int              `json:"maxRegistrations"`
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	Labels           map[string]string `json:"labels"`
}

type DevicesRegisteredWithTokenCount struct {
	AllCount int `json:"allCount"`
}

type DeviceAccessKey struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	ProjectID string    `json:"projectId"`
	DeviceID  string    `json:"deviceId"`
}

type Application struct {
	ID             string    `json:"id"`
	CreatedAt      time.Time `json:"createdAt"`
	ProjectID      string    `json:"projectId"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	SchedulingRule Query     `json:"schedulingRule"`
}

type ApplicationDeviceCounts struct {
	AllCount int `json:"allCount"`
}

type Release struct {
	ID                        string             `json:"id"`
	CreatedAt                 time.Time          `json:"createdAt"`
	ProjectID                 string             `json:"projectId"`
	ApplicationID             string             `json:"applicationId"`
	Config                    map[string]Service `json:"config"`
	RawConfig                 string             `json:"rawConfig"`
	CreatedByUserID           *string            `json:"createdByUserId"`
	CreatedByServiceAccountID *string            `json:"createdByServiceAccountId"`
}

type ReleaseDeviceCounts struct {
	AllCount int `json:"allCount"`
}

type DeviceApplicationStatus struct {
	ProjectID        string `json:"projectId"`
	DeviceID         string `json:"deviceId"`
	ApplicationID    string `json:"applicationId"`
	CurrentReleaseID string `json:"currentReleaseId"`
}

type DeviceServiceStatus struct {
	ProjectID        string `json:"projectId"`
	DeviceID         string `json:"deviceId"`
	ApplicationID    string `json:"applicationId"`
	Service          string `json:"service"`
	CurrentReleaseID string `json:"currentReleaseId"`
}

type MetricTargetConfig struct {
	ID        string           `json:"id"`
	CreatedAt time.Time        `json:"createdAt"`
	ProjectID string           `json:"projectId"`
	Type      MetricTargetType `json:"type"`
	Configs   []MetricConfig   `json:"configs"`
}

type MetricTargetType string

const (
	MetricServiceTargetType MetricTargetType = "service"
	MetricHostTargetType    MetricTargetType = "host"
	MetricStateTargetType   MetricTargetType = "state"
)

type MetricConfig struct {
	Params  *ServiceMetricParams `json:"params,omitempty"`
	Metrics []Metric             `json:"metrics"`
}

type ServiceMetricParams struct {
	ApplicationID string `json:"applicationId"`
	Service       string `json:"service"`
}

type Metric struct {
	Metric string   `json:"metric"`
	Labels []string `json:"labels"`
	Tags   []string `json:"tags"`
}

type MembershipFull1 struct {
	Membership
	User    User        `json:"user"`
	Project ProjectFull `json:"project"`
}

type ProjectFull struct {
	Project
	DeviceCounts      ProjectDeviceCounts      `json:"deviceCounts"`
	ApplicationCounts ProjectApplicationCounts `json:"applicationCounts"`
}

type MembershipFull2 struct {
	Membership
	User  User   `json:"user"`
	Roles []Role `json:"roles"`
}

type ServiceAccountFull struct {
	ServiceAccount
	Roles []Role `json:"roles"`
}

type DeviceFull struct {
	Device
	ApplicationStatusInfo []DeviceApplicationStatusInfo `json:"applicationStatusInfo"`
}

type DeviceApplicationStatusInfo struct {
	Application       Application              `json:"application"`
	ApplicationStatus *DeviceApplicationStatus `json:"applicationStatus"`
	ServiceStatuses   []DeviceServiceStatus    `json:"serviceStatuses"`
}

type ApplicationFull1 struct {
	Application
	LatestRelease *Release                `json:"latestRelease"`
	DeviceCounts  ApplicationDeviceCounts `json:"deviceCounts"`
}

type DeviceRegistrationTokenFull struct {
	DeviceRegistrationToken
	DeviceCounts DevicesRegisteredWithTokenCount `json:"deviceCounts"`
}

type ReleaseFull struct {
	Release
	CreatedByUser           *User               `json:"createdByUser"`
	CreatedByServiceAccount *ServiceAccount     `json:"createdByServiceAccount"`
	DeviceCounts            ReleaseDeviceCounts `json:"deviceCounts"`
}

type Bundle struct {
	Applications        []ApplicationFull2        `json:"applications"`
	ApplicationStatuses []DeviceApplicationStatus `json:"applicationStatuses"`
	ServiceStatuses     []DeviceServiceStatus     `json:"serviceStatuses"`
	DesiredAgentSpec    string                    `json:"desiredAgentSpec"`
	DesiredAgentVersion string                    `json:"desiredAgentVersion"`
}

type ApplicationFull2 struct {
	Application   Application `json:"application"`
	LatestRelease Release     `json:"latestRelease"`
}

type DeviceInfo struct {
	AgentVersion string    `json:"agentVersion"`
	IPAddress    string    `json:"ipAddress"`
	OSRelease    OSRelease `json:"osRelease"`
}

type OSRelease struct {
	PrettyName string `json:"prettyName"`
	Name       string `json:"name"`
	VersionID  string `json:"versionId"`
	Version    string `json:"version"`
	ID         string `json:"id"`
	IDLike     string `json:"idLike"`
}
