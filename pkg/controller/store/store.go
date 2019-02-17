package store

import (
	"context"

	"github.com/deviceplane/deviceplane/pkg/models"
)

type Users interface {
	CreateUser(ctx context.Context, email, passwordHash string) (*models.User, error)
	GetUser(ctx context.Context, id string) (*models.User, error)
}

type AccessKeys interface {
	CreateAccessKey(ctx context.Context, userID string, hash string) (*models.AccessKey, error)
	GetAccessKey(ctx context.Context, id string) (*models.AccessKey, error)
	ValidateAccessKey(ctx context.Context, hash string) (*models.AccessKey, error)
}

type Projects interface {
	CreateProject(ctx context.Context, name string) (*models.Project, error)
	GetProject(ctx context.Context, id string) (*models.Project, error)
}

type Memberships interface {
	CreateMembership(ctx context.Context, userID, projectID, membershipType string) (*models.Membership, error)
	GetMembership(ctx context.Context, userID, projectID string) (*models.Membership, error)
	ListMembershipsByUser(ctx context.Context, userID string) ([]models.Membership, error)
	ListMembershipsByProject(ctx context.Context, projectID string) ([]models.Membership, error)
}

type Devices interface {
	CreateDevice(ctx context.Context, projectID string) (*models.Device, error)
	GetDevice(ctx context.Context, id string) (*models.Device, error)
	ListDevices(ctx context.Context, projectID string) ([]models.Device, error)
}

type Applications interface {
	CreateApplication(ctx context.Context, projectID, name string) (*models.Application, error)
	GetApplication(ctx context.Context, applicationID string) (*models.Application, error)
	ListApplications(ctx context.Context, projectID string) ([]models.Application, error)
}

type Releases interface {
	CreateRelease(ctx context.Context, applicationID, config string) (*models.Release, error)
	GetRelease(ctx context.Context, id string) (*models.Release, error)
	GetLatestRelease(ctx context.Context, applicationID string) (*models.Release, error)
	ListReleases(ctx context.Context, applicationID string) ([]models.Release, error)
}
