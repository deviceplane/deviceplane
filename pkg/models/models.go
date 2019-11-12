package models

import (
	"time"

	"github.com/deviceplane/deviceplane/pkg/yamltypes"
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
	RawConfig                 string             `json:"RawConfig"`
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

type Service struct {
	CapAdd         []string                  `yaml:"cap_add,omitempty"`
	CapDrop        []string                  `yaml:"cap_drop,omitempty"`
	Command        yamltypes.Command         `yaml:"command,flow,omitempty"`
	CPUSet         string                    `yaml:"cpuset,omitempty"`
	CPUShares      yamltypes.StringorInt     `yaml:"cpu_shares,omitempty"`
	CPUQuota       yamltypes.StringorInt     `yaml:"cpu_quota,omitempty"`
	DNS            yamltypes.Stringorslice   `yaml:"dns,omitempty"`
	DNSOpts        []string                  `yaml:"dns_opt,omitempty"`
	DNSSearch      yamltypes.Stringorslice   `yaml:"dns_search,omitempty"`
	DomainName     string                    `yaml:"domainname,omitempty"`
	Entrypoint     yamltypes.Command         `yaml:"entrypoint,flow,omitempty"`
	Environment    yamltypes.MaporEqualSlice `yaml:"environment,omitempty"`
	ExtraHosts     []string                  `yaml:"extra_hosts,omitempty"`
	GroupAdd       []string                  `yaml:"group_add,omitempty"`
	Image          string                    `yaml:"image,omitempty"`
	Hostname       string                    `yaml:"hostname,omitempty"`
	Ipc            string                    `yaml:"ipc,omitempty"`
	Labels         yamltypes.SliceorMap      `yaml:"labels,omitempty"`
	MemLimit       yamltypes.MemStringorInt  `yaml:"mem_limit,omitempty"`
	MemReservation yamltypes.MemStringorInt  `yaml:"mem_reservation,omitempty"`
	MemSwapLimit   yamltypes.MemStringorInt  `yaml:"memswap_limit,omitempty"`
	NetworkMode    string                    `yaml:"network_mode,omitempty"`
	OomKillDisable bool                      `yaml:"oom_kill_disable,omitempty"`
	OomScoreAdj    yamltypes.StringorInt     `yaml:"oom_score_adj,omitempty"`
	Pid            string                    `yaml:"pid,omitempty"`
	Ports          []string                  `yaml:"ports,omitempty"`
	Privileged     bool                      `yaml:"privileged,omitempty"`
	ReadOnly       bool                      `yaml:"read_only,omitempty"`
	Restart        string                    `yaml:"restart,omitempty"`
	SecurityOpt    []string                  `yaml:"security_opt,omitempty"`
	ShmSize        yamltypes.MemStringorInt  `yaml:"shm_size,omitempty"`
	StopSignal     string                    `yaml:"stop_signal,omitempty"`
	User           string                    `yaml:"user,omitempty"`
	Uts            string                    `yaml:"uts,omitempty"`
	Volumes        *yamltypes.Volumes        `yaml:"volumes,omitempty"`
	WorkingDir     string                    `yaml:"working_dir,omitempty"`
}
