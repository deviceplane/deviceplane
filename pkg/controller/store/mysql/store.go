package mysql

import (
	"context"
	"database/sql"

	"github.com/deviceplane/deviceplane/pkg/controller/store"
	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
)

type scanner interface {
	Scan(...interface{}) error
}

var (
	_ store.Users        = &Store{}
	_ store.Projects     = &Store{}
	_ store.Devices      = &Store{}
	_ store.Applications = &Store{}
	_ store.Releases     = &Store{}
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) CreateUser(ctx context.Context) (*models.User, error) {
	id := ksuid.New().String()

	if _, err := s.db.ExecContext(
		ctx,
		createUser,
		id,
	); err != nil {
		return nil, err
	}

	return &models.User{
		ID: id,
	}, nil
}

func (s *Store) GetUser(ctx context.Context, id string) (*models.User, error) {
	userRow := s.db.QueryRowContext(ctx, getUser, id)

	user, err := s.scanUser(userRow)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Store) scanUser(scanner scanner) (*models.User, error) {
	var user models.User
	if err := scanner.Scan(
		&user.ID,
		&user.CreatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "scan user")
	}
	return &user, nil
}

func (s *Store) CreateProject(ctx context.Context) (*models.Project, error) {
	id := ksuid.New().String()

	if _, err := s.db.ExecContext(
		ctx,
		createProject,
		id,
	); err != nil {
		return nil, err
	}

	return &models.Project{
		ID: id,
	}, nil
}

func (s *Store) GetProject(ctx context.Context, id string) (*models.Project, error) {
	projectRow := s.db.QueryRowContext(ctx, getProject, id)

	project, err := s.scanProject(projectRow)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (s *Store) scanProject(scanner scanner) (*models.Project, error) {
	var project models.Project
	if err := scanner.Scan(
		&project.ID,
		&project.CreatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "scan project")
	}
	return &project, nil
}

func (s *Store) CreateDevice(ctx context.Context, projectID string) (*models.Device, error) {
	id := ksuid.New().String()

	if _, err := s.db.ExecContext(
		ctx,
		createDevice,
		id,
		projectID,
	); err != nil {
		return nil, err
	}

	return &models.Device{
		ID:        id,
		ProjectID: projectID,
	}, nil
}

func (s *Store) GetDevice(ctx context.Context, id string) (*models.Device, error) {
	deviceRow := s.db.QueryRowContext(ctx, getDevice, id)

	device, err := s.scanDevice(deviceRow)
	if err != nil {
		return nil, err
	}

	return device, nil
}

func (s *Store) scanDevice(scanner scanner) (*models.Device, error) {
	var device models.Device
	if err := scanner.Scan(
		&device.ID,
		&device.ProjectID,
	); err != nil {
		return nil, errors.Wrap(err, "scan device")
	}
	return &device, nil
}

func (s *Store) CreateApplication(ctx context.Context, projectID string) (*models.Application, error) {
	id := ksuid.New().String()

	if _, err := s.db.ExecContext(
		ctx,
		createApplication,
		id,
		projectID,
	); err != nil {
		return nil, err
	}

	return &models.Application{
		ID:        id,
		ProjectID: projectID,
	}, nil
}

func (s *Store) ListApplications(ctx context.Context, projectID string) ([]models.Application, error) {
	applicationRows, err := s.db.QueryContext(ctx, listApplications, projectID)
	if err != nil {
		return nil, errors.Wrap(err, "query applications")
	}
	defer applicationRows.Close()

	var applications []models.Application
	for applicationRows.Next() {
		application, err := s.scanApplication(applicationRows)
		if err != nil {
			return nil, err
		}
		applications = append(applications, *application)
	}

	if err := applicationRows.Err(); err != nil {
		return nil, err
	}

	return applications, nil
}

func (s *Store) scanApplication(scanner scanner) (*models.Application, error) {
	var application models.Application
	if err := scanner.Scan(
		&application.ID,
		&application.ProjectID,
	); err != nil {
		return nil, errors.Wrap(err, "scan application")
	}
	return &application, nil
}

func (s *Store) CreateRelease(ctx context.Context, applicationID, config string) (*models.Release, error) {
	id := ksuid.New().String()

	if _, err := s.db.ExecContext(
		ctx,
		createRelease,
		id,
		applicationID,
		config,
	); err != nil {
		return nil, err
	}

	return &models.Release{
		ID:            id,
		ApplicationID: applicationID,
		Config:        config,
	}, nil
}

func (s *Store) GetRelease(ctx context.Context, id string) (*models.Release, error) {
	applicationRow := s.db.QueryRowContext(ctx, getRelease, id)

	release, err := s.scanRelease(applicationRow)
	if err != nil {
		return nil, err
	}

	return release, nil
}

func (s *Store) GetLatestRelease(ctx context.Context, applicationID string) (*models.Release, error) {
	applicationRow := s.db.QueryRowContext(ctx, getLatestRelease, applicationID)

	release, err := s.scanRelease(applicationRow)
	if err != nil {
		return nil, err
	}

	return release, nil
}

func (s *Store) scanRelease(scanner scanner) (*models.Release, error) {
	var release models.Release
	if err := scanner.Scan(
		&release.ID,
		&release.CreatedAt,
		&release.ApplicationID,
		&release.Config,
	); err != nil {
		return nil, errors.Wrap(err, "scan application")
	}
	return &release, nil
}
