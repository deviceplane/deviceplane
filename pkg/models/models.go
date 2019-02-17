package models

import "time"

type User struct {
	ID           string    `json:"id"`
	CreatedAt    time.Time `json:"createdAt"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"passwordHash"`
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

type Membership struct {
	UserID    string    `json:"userId"`
	ProjectID string    `json:"projectId"`
	CreatedAt time.Time `json:"createdAt"`
	Level     string    `json:"level"`
}

type Device struct {
	ID        string `json:"id"`
	ProjectID string `json:"projectId"`
}

type Application struct {
	ID        string `json:"id"`
	ProjectID string `json:"projectId"`
	Name      string `json:"name"`
}

type Release struct {
	ID            string    `json:"id"`
	CreatedAt     time.Time `json:"createdAt"`
	ApplicationID string    `json:"applicationId"`
	Config        string    `json:"config"`
}

type Bundle struct {
	ID           string                        `json:"id"`
	Applications []ApplicationAndLatestRelease `json:"applications"`
}

type ApplicationAndLatestRelease struct {
	Application   Application `json:"application"`
	LatestRelease Release     `json:"latestRelease"`
}
