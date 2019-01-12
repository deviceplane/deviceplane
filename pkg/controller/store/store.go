package store

import (
	"context"

	"github.com/deviceplane/deviceplane/pkg/models"
)

type Users interface {
	CreateUser(ctx context.Context) (*models.User, error)
	GetUser(ctx context.Context, id string) (*models.User, error)
}

type Projects interface {
	CreateProject(ctx context.Context) (*models.Project, error)
	GetProject(ctx context.Context, id string) (*models.Project, error)
}

type Devices interface {
	CreateDevice(ctx context.Context, projectID string) (*models.Device, error)
	GetDevice(ctx context.Context, id string) (*models.Device, error)
}

type Applications interface {
	CreateApplication(ctx context.Context, projectID string) (*models.Application, error)
	ListApplications(ctx context.Context, projectID string) ([]models.Application, error)
}

type Releases interface {
	CreateRelease(ctx context.Context, projectID, applicationID string) (*models.Release, error)
	GetRelease(ctx context.Context, id string) (*models.Release, error)
	GetLatestRelease(ctx context.Context, applicationID string) (*models.Release, error)
}
