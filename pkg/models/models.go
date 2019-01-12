package models

import "time"

type User struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

type Project struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

type Device struct {
	ID        string `json:"id"`
	ProjectID string `json:"projectId"`
}

type Application struct {
	ID        string `json:"id"`
	ProjectID string `json:"projectId"`
}

type Release struct {
	ID            string    `json:"id"`
	CreatedAt     time.Time `json:"createdAt"`
	ApplicationID string    `json:"applicationId"`
	Config        string    `json:"config"`
}

type CreateRelease struct {
	Config string `json:"config"`
}

type Bundle struct {
	ID           string                        `json:"id"`
	Applications []ApplicationAndLatestRelease `json:"applications"`
}

type ApplicationAndLatestRelease struct {
	Application   Application `json:"application"`
	LatestRelease Release     `json:"latestRelease"`
}
