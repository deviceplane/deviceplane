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
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`
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
	ID               string       `json:"id"`
	CreatedAt        time.Time    `json:"createdAt"`
	ProjectID        string       `json:"projectId"`
	Name             string       `json:"name"`
	DesiredAgentSpec string       `json:"desiredAgentSpec"`
	Info             DeviceInfo   `json:"info"`
	LastSeenAt       time.Time    `json:"lastSeenAt"`
	Status           DeviceStatus `json:"status"`
}

type DeviceStatus string

const (
	DeviceStatusOnline  = DeviceStatus("online")
	DeviceStatusOffline = DeviceStatus("offline")
)

type DeviceLabel struct {
	Key       string    `json:"key"`
	DeviceID  string    `json:"deviceId"`
	CreatedAt time.Time `json:"createdAt"`
	ProjectID string    `json:"projectId"`
	Value     string    `json:"value"`
}

type DeviceRegistrationToken struct {
	ID                string    `json:"id"`
	CreatedAt         time.Time `json:"createdAt"`
	ProjectID         string    `json:"projectId"`
	DeviceAccessKeyID *string   `json:"deviceAccessKeyId"`
}

type DeviceAccessKey struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	ProjectID string    `json:"projectId"`
	DeviceID  string    `json:"deviceId"`
}

type Application struct {
	ID          string              `json:"id"`
	CreatedAt   time.Time           `json:"createdAt"`
	ProjectID   string              `json:"projectId"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Settings    ApplicationSettings `json:"settings"`
}

type ApplicationDeviceCounts struct {
	AllCount int `json:"allCount"`
}

type Release struct {
	ID                        string    `json:"id"`
	CreatedAt                 time.Time `json:"createdAt"`
	ProjectID                 string    `json:"projectId"`
	ApplicationID             string    `json:"applicationId"`
	Config                    string    `json:"config"`
	CreatedByUserID           *string   `json:"createdByUserId"`
	CreatedByServiceAccountID *string   `json:"createdByServiceAccountId"`
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

type ApplicationSettings struct {
	SchedulingRule string `json:"schedulingRule"`
}
