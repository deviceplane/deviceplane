package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/deviceplane/deviceplane/pkg/controller/store"
	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
)

type scanner interface {
	Scan(...interface{}) error
}

const (
	userPrefix        = "usr"
	accessKeyPrefix   = "key"
	projectPrefix     = "prj"
	devicePrefix      = "dev"
	applicationPrefix = "app"
	releasePrefix     = "rel"
)

func newUserID() string {
	return fmt.Sprintf("%s_%s", userPrefix, ksuid.New().String())
}

func newAccessKeyID() string {
	return fmt.Sprintf("%s_%s", accessKeyPrefix, ksuid.New().String())
}

func newProjectID() string {
	return fmt.Sprintf("%s_%s", projectPrefix, ksuid.New().String())
}

func newDeviceID() string {
	return fmt.Sprintf("%s_%s", devicePrefix, ksuid.New().String())
}

func newApplicationID() string {
	return fmt.Sprintf("%s_%s", applicationPrefix, ksuid.New().String())
}

func newReleaseID() string {
	return fmt.Sprintf("%s_%s", releasePrefix, ksuid.New().String())
}

var (
	_ store.Users        = &Store{}
	_ store.AccessKeys   = &Store{}
	_ store.Projects     = &Store{}
	_ store.Memberships  = &Store{}
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

func (s *Store) CreateUser(ctx context.Context, email, passwordHash string) (*models.User, error) {
	id := newUserID()

	if _, err := s.db.ExecContext(
		ctx,
		createUser,
		id,
		email,
		passwordHash,
	); err != nil {
		return nil, err
	}

	return s.GetUser(ctx, id)
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
		&user.Email,
		&user.PasswordHash,
	); err != nil {
		return nil, errors.Wrap(err, "scan user")
	}
	return &user, nil
}

func (s *Store) CreateAccessKey(ctx context.Context, userID, hash string) (*models.AccessKey, error) {
	id := newAccessKeyID()

	if _, err := s.db.ExecContext(
		ctx,
		createAccessKey,
		id,
		userID,
		hash,
	); err != nil {
		return nil, err
	}

	return s.GetAccessKey(ctx, id)
}

func (s *Store) GetAccessKey(ctx context.Context, id string) (*models.AccessKey, error) {
	accessKeyRow := s.db.QueryRowContext(ctx, getAccessKey, id)

	accessKey, err := s.scanAccessKey(accessKeyRow)
	if err != nil {
		return nil, err
	}

	return accessKey, nil
}

func (s *Store) ValidateAccessKey(ctx context.Context, hash string) (*models.AccessKey, error) {
	accessKeyRow := s.db.QueryRowContext(ctx, validateAccessKey, hash)

	accessKey, err := s.scanAccessKey(accessKeyRow)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return accessKey, nil
}

func (s *Store) scanAccessKey(scanner scanner) (*models.AccessKey, error) {
	var accessKey models.AccessKey
	if err := scanner.Scan(
		&accessKey.ID,
		&accessKey.UserID,
		&accessKey.Hash,
	); err != nil {
		return nil, errors.Wrap(err, "scan access key")
	}
	return &accessKey, nil
}

func (s *Store) CreateProject(ctx context.Context, name string) (*models.Project, error) {
	id := newProjectID()

	if _, err := s.db.ExecContext(
		ctx,
		createProject,
		id,
		name,
	); err != nil {
		return nil, err
	}

	return s.GetProject(ctx, id)
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
		&project.Name,
	); err != nil {
		return nil, errors.Wrap(err, "scan project")
	}
	return &project, nil
}

func (s *Store) CreateMembership(ctx context.Context, userID, projectID, level string) (*models.Membership, error) {
	if _, err := s.db.ExecContext(
		ctx,
		createMembership,
		userID,
		projectID,
		level,
	); err != nil {
		return nil, err
	}

	return s.GetMembership(ctx, userID, projectID)
}

func (s *Store) GetMembership(ctx context.Context, userID, projectID string) (*models.Membership, error) {
	membershipRow := s.db.QueryRowContext(ctx, getMembership, userID, projectID)

	membership, err := s.scanMembership(membershipRow)
	if err != nil {
		return nil, err
	}

	return membership, nil
}

func (s *Store) ListMembershipsByUser(ctx context.Context, userID string) ([]models.Membership, error) {
	return s.listMemberships(ctx, userID, listMembershipsByUser)
}

func (s *Store) ListMembershipsByProject(ctx context.Context, projectID string) ([]models.Membership, error) {
	return s.listMemberships(ctx, projectID, listMembershipsByProject)
}

func (s *Store) listMemberships(ctx context.Context, id, query string) ([]models.Membership, error) {
	membershipRows, err := s.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, errors.Wrap(err, "query memberships")
	}
	defer membershipRows.Close()

	var memberships []models.Membership
	for membershipRows.Next() {
		membership, err := s.scanMembership(membershipRows)
		if err != nil {
			return nil, err
		}
		memberships = append(memberships, *membership)
	}

	if err := membershipRows.Err(); err != nil {
		return nil, err
	}

	return memberships, nil
}

func (s *Store) scanMembership(scanner scanner) (*models.Membership, error) {
	var membership models.Membership
	if err := scanner.Scan(
		&membership.UserID,
		&membership.ProjectID,
		&membership.Level,
	); err != nil {
		return nil, errors.Wrap(err, "scan membership")
	}
	return &membership, nil
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

func (s *Store) ListDevices(ctx context.Context, projectID string) ([]models.Device, error) {
	deviceRows, err := s.db.QueryContext(ctx, listDevices, projectID)
	if err != nil {
		return nil, errors.Wrap(err, "query devices")
	}
	defer deviceRows.Close()

	var devices []models.Device
	for deviceRows.Next() {
		device, err := s.scanDevice(deviceRows)
		if err != nil {
			return nil, err
		}
		devices = append(devices, *device)
	}

	if err := deviceRows.Err(); err != nil {
		return nil, err
	}

	return devices, nil
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

func (s *Store) CreateApplication(ctx context.Context, projectID, name string) (*models.Application, error) {
	id := newApplicationID()

	if _, err := s.db.ExecContext(
		ctx,
		createApplication,
		id,
		projectID,
		name,
	); err != nil {
		return nil, err
	}

	return &models.Application{
		ID:        id,
		ProjectID: projectID,
		Name:      name,
	}, nil
}

func (s *Store) GetApplication(ctx context.Context, id string) (*models.Application, error) {
	applicationRow := s.db.QueryRowContext(ctx, getApplication, id)

	application, err := s.scanApplication(applicationRow)
	if err != nil {
		return nil, err
	}

	return application, nil
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
		&application.Name,
	); err != nil {
		return nil, errors.Wrap(err, "scan application")
	}
	return &application, nil
}

func (s *Store) CreateRelease(ctx context.Context, applicationID, config string) (*models.Release, error) {
	id := newReleaseID()

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

func (s *Store) ListReleases(ctx context.Context, applicationID string) ([]models.Release, error) {
	releaseRows, err := s.db.QueryContext(ctx, listReleases, applicationID)
	if err != nil {
		return nil, errors.Wrap(err, "query releases")
	}
	defer releaseRows.Close()

	var releases []models.Release
	for releaseRows.Next() {
		release, err := s.scanRelease(releaseRows)
		if err != nil {
			return nil, err
		}
		releases = append(releases, *release)
	}

	if err := releaseRows.Err(); err != nil {
		return nil, err
	}

	return releases, nil
}

func (s *Store) scanRelease(scanner scanner) (*models.Release, error) {
	var release models.Release
	if err := scanner.Scan(
		&release.ID,
		&release.ApplicationID,
		&release.Config,
	); err != nil {
		return nil, errors.Wrap(err, "scan application")
	}
	return &release, nil
}
