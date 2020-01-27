package models

const LatestRelease = "latest"

type ScheduledDevice struct {
	Device
	ReleaseID string `json:"releaseId"`
}

type SchedulingRule struct {
	ScheduleType     ScheduleType      `json:"scheduleType"`
	DefaultReleaseID string            `json:"defaultReleaseId"` // TODO: validate Release ID?
	ConditionalQuery *Query            `json:"conditionalQuery,omitempty"`
	ReleaseSelectors []ReleaseSelector `json:"releaseSelectors"`
}

type ScheduleType string

const (
	ScheduleTypeNoDevices   = "NoDevices"
	ScheduleTypeAllDevices  = "AllDevices"
	ScheduleTypeConditional = "Conditional"
)

type ReleaseSelector struct {
	Query     Query  `json:"releaseQuery"`
	ReleaseID string `json:"releaseId"` // TODO: validate Release ID?
}
