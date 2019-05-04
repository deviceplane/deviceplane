package models

import "time"

type User struct {
	ID                    string    `json:"id"`
	CreatedAt             time.Time `json:"createdAt"`
	Email                 string    `json:"email"`
	PasswordHash          string    `json:"passwordHash"`
	FirstName             string    `json:"firstName"`
	LastName              string    `json:"lastName"`
	RegistrationCompleted bool      `json:"registrationCompleted"`
}

type RegistrationToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	Hash      string    `json:"hash"`
}

type Session struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UserID    string    `json:"userId"`
	Hash      string    `json:"hash"`
}

type AccessKey struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UserID    string    `json:"userId"`
	Hash      string    `json:"hash"`
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

type Membership struct {
	UserID    string    `json:"userId"`
	ProjectID string    `json:"projectId"`
	CreatedAt time.Time `json:"createdAt"`
	Level     string    `json:"level"`
}

type Device struct {
	ID        string     `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	ProjectID string     `json:"projectId"`
	Name      string     `json:"name"`
	Info      DeviceInfo `json:"info"`
}

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
	Hash      string    `json:"hash"`
}

type Application struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	ProjectID string    `json:"projectId"`
	Name      string    `json:"name"`
}

type Release struct {
	ID            string    `json:"id"`
	CreatedAt     time.Time `json:"createdAt"`
	ProjectID     string    `json:"projectId"`
	ApplicationID string    `json:"applicationId"`
	Config        string    `json:"config"`
}

type MembershipFull struct {
	User    User        `json:"user"`
	Project ProjectFull `json:"project"`
}

type ProjectFull struct {
	Project
	DeviceCounts      ProjectDeviceCounts      `json:"deviceCounts"`
	ApplicationCounts ProjectApplicationCounts `json:"applicationCounts"`
}

type Bundle struct {
	ID           string                        `json:"id"`
	Applications []ApplicationAndLatestRelease `json:"applications"`
}

type ApplicationAndLatestRelease struct {
	Application   Application `json:"application"`
	LatestRelease *Release    `json:"latestRelease"`
}

type DeviceInfo struct {
	IPAddress string `json:"ipAddress"`
}
