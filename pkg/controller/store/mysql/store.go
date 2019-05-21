package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/deviceplane/deviceplane/pkg/controller/authz"
	"github.com/deviceplane/deviceplane/pkg/controller/store"
	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	"gopkg.in/yaml.v2"
)

type scanner interface {
	Scan(...interface{}) error
}

const (
	userPrefix                    = "usr"
	registrationTokenPrefix       = "reg"
	sessionPrefix                 = "ses"
	accessKeyPrefix               = "key"
	projectPrefix                 = "prj"
	rolePrefix                    = "rol"
	membershipPrefix              = "mem"
	membershipRoleBindingPrefix   = "mrb"
	devicePrefix                  = "dev"
	deviceRegistrationTokenPrefix = "drt"
	deviceAccessKeyPrefix         = "dak"
	applicationPrefix             = "app"
	releasePrefix                 = "rel"
)

func newUserID() string {
	return fmt.Sprintf("%s_%s", userPrefix, ksuid.New().String())
}

func newRegistrationTokenID() string {
	return fmt.Sprintf("%s_%s", registrationTokenPrefix, ksuid.New().String())
}

func newSessionID() string {
	return fmt.Sprintf("%s_%s", sessionPrefix, ksuid.New().String())
}

func newAccessKeyID() string {
	return fmt.Sprintf("%s_%s", accessKeyPrefix, ksuid.New().String())
}

func newProjectID() string {
	return fmt.Sprintf("%s_%s", projectPrefix, ksuid.New().String())
}

func newRoleID() string {
	return fmt.Sprintf("%s_%s", rolePrefix, ksuid.New().String())
}

func newMembershipID() string {
	return fmt.Sprintf("%s_%s", membershipPrefix, ksuid.New().String())
}

func newMembershipRoleBindingID() string {
	return fmt.Sprintf("%s_%s", membershipRoleBindingPrefix, ksuid.New().String())
}

func newDeviceID() string {
	return fmt.Sprintf("%s_%s", devicePrefix, ksuid.New().String())
}

func newDeviceRegistrationTokenID() string {
	return fmt.Sprintf("%s_%s", deviceRegistrationTokenPrefix, ksuid.New().String())
}

func newDeviceAccessKeyID() string {
	return fmt.Sprintf("%s_%s", deviceAccessKeyPrefix, ksuid.New().String())
}

func newApplicationID() string {
	return fmt.Sprintf("%s_%s", applicationPrefix, ksuid.New().String())
}

func newReleaseID() string {
	return fmt.Sprintf("%s_%s", releasePrefix, ksuid.New().String())
}

var (
	_ store.Users                     = &Store{}
	_ store.RegistrationTokens        = &Store{}
	_ store.Sessions                  = &Store{}
	_ store.AccessKeys                = &Store{}
	_ store.Projects                  = &Store{}
	_ store.ProjectDeviceCounts       = &Store{}
	_ store.Roles                     = &Store{}
	_ store.Memberships               = &Store{}
	_ store.MembershipRoleBindings    = &Store{}
	_ store.Devices                   = &Store{}
	_ store.DeviceLabels              = &Store{}
	_ store.DeviceAccessKeys          = &Store{}
	_ store.DeviceRegistrationTokens  = &Store{}
	_ store.Applications              = &Store{}
	_ store.Releases                  = &Store{}
	_ store.ReleaseDeviceCounts       = &Store{}
	_ store.DeviceApplicationStatuses = &Store{}
	_ store.DeviceServiceStatuses     = &Store{}
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) CreateUser(ctx context.Context, email, passwordHash, firstName, lastName string) (*models.User, error) {
	id := newUserID()

	if _, err := s.db.ExecContext(
		ctx,
		createUser,
		id,
		email,
		passwordHash,
		firstName,
		lastName,
	); err != nil {
		return nil, err
	}

	return s.GetUser(ctx, id)
}

func (s *Store) GetUser(ctx context.Context, id string) (*models.User, error) {
	userRow := s.db.QueryRowContext(ctx, getUser, id)

	user, err := s.scanUser(userRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Store) ValidateUser(ctx context.Context, email, passwordHash string) (*models.User, error) {
	userRow := s.db.QueryRowContext(ctx, validateUser, email, passwordHash)

	user, err := s.scanUser(userRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Store) MarkRegistrationCompleted(ctx context.Context, id string) (*models.User, error) {
	if _, err := s.db.ExecContext(
		ctx,
		markRegistrationComplete,
		id,
	); err != nil {
		return nil, err
	}

	return s.GetUser(ctx, id)
}

func (s *Store) scanUser(scanner scanner) (*models.User, error) {
	var user models.User
	if err := scanner.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.RegistrationCompleted,
	); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Store) CreateRegistrationToken(ctx context.Context, userID, hash string) (*models.RegistrationToken, error) {
	id := newRegistrationTokenID()

	if _, err := s.db.ExecContext(
		ctx,
		createRegistrationToken,
		id,
		userID,
		hash,
	); err != nil {
		return nil, err
	}

	return s.GetRegistrationToken(ctx, id)
}

func (s *Store) GetRegistrationToken(ctx context.Context, id string) (*models.RegistrationToken, error) {
	registrationTokenRow := s.db.QueryRowContext(ctx, getRegistrationToken, id)

	registrationToken, err := s.scanRegistrationToken(registrationTokenRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrRegistrationTokenNotFound
	} else if err != nil {
		return nil, err
	}

	return registrationToken, nil
}

func (s *Store) ValidateRegistrationToken(ctx context.Context, hash string) (*models.RegistrationToken, error) {
	registrationTokenRow := s.db.QueryRowContext(ctx, validateRegistrationToken, hash)

	registrationToken, err := s.scanRegistrationToken(registrationTokenRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrRegistrationTokenNotFound
	} else if err != nil {
		return nil, err
	}

	return registrationToken, nil
}

func (s *Store) scanRegistrationToken(scanner scanner) (*models.RegistrationToken, error) {
	var registrationToken models.RegistrationToken
	if err := scanner.Scan(
		&registrationToken.ID,
		&registrationToken.UserID,
		&registrationToken.Hash,
	); err != nil {
		return nil, err
	}
	return &registrationToken, nil
}

func (s *Store) CreateSession(ctx context.Context, userID, hash string) (*models.Session, error) {
	id := newSessionID()

	if _, err := s.db.ExecContext(
		ctx,
		createSession,
		id,
		userID,
		hash,
	); err != nil {
		return nil, err
	}

	return s.GetSession(ctx, id)
}

func (s *Store) GetSession(ctx context.Context, id string) (*models.Session, error) {
	sessionRow := s.db.QueryRowContext(ctx, getSession, id)

	session, err := s.scanSession(sessionRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrSessionNotFound
	} else if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *Store) ValidateSession(ctx context.Context, hash string) (*models.Session, error) {
	sessionRow := s.db.QueryRowContext(ctx, validateSession, hash)

	session, err := s.scanSession(sessionRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrSessionNotFound
	} else if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *Store) DeleteSession(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, deleteSession, id)
	return err
}

func (s *Store) scanSession(scanner scanner) (*models.Session, error) {
	var session models.Session
	if err := scanner.Scan(
		&session.ID,
		&session.UserID,
		&session.Hash,
	); err != nil {
		return nil, err
	}
	return &session, nil
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
	if err == sql.ErrNoRows {
		return nil, store.ErrAccessKeyNotFound
	} else if err != nil {
		return nil, err
	}

	return accessKey, nil
}

func (s *Store) ValidateAccessKey(ctx context.Context, hash string) (*models.AccessKey, error) {
	accessKeyRow := s.db.QueryRowContext(ctx, validateAccessKey, hash)

	accessKey, err := s.scanAccessKey(accessKeyRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrAccessKeyNotFound
	} else if err != nil {
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
		return nil, err
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
	if err == sql.ErrNoRows {
		return nil, store.ErrProjectNotFound
	} else if err != nil {
		return nil, err
	}

	return project, nil
}

func (s *Store) LookupProject(ctx context.Context, name string) (*models.Project, error) {
	projectRow := s.db.QueryRowContext(ctx, lookupProject, name)

	project, err := s.scanProject(projectRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrProjectNotFound
	} else if err != nil {
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
		return nil, err
	}
	return &project, nil
}

func (s *Store) GetProjectDeviceCounts(ctx context.Context, projectID string) (*models.ProjectDeviceCounts, error) {
	countRow := s.db.QueryRowContext(ctx, getProjectDeviceCounts, projectID)

	count, err := s.scanProjectDeviceCountRow(countRow)
	if err != nil {
		return nil, err
	}

	return &models.ProjectDeviceCounts{
		AllCount: count,
	}, nil
}

func (s *Store) scanProjectDeviceCountRow(scanner scanner) (int, error) {
	var count int
	if err := scanner.Scan(
		&count,
	); err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Store) GetProjectApplicationCounts(ctx context.Context, projectID string) (*models.ProjectApplicationCounts, error) {
	countRow := s.db.QueryRowContext(ctx, getProjectApplicationCounts, projectID)

	count, err := s.scanProjectApplicationCountRow(countRow)
	if err != nil {
		return nil, err
	}

	return &models.ProjectApplicationCounts{
		AllCount: count,
	}, nil
}

func (s *Store) scanProjectApplicationCountRow(scanner scanner) (int, error) {
	var count int
	if err := scanner.Scan(
		&count,
	); err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Store) CreateRole(ctx context.Context, projectID, name, description string, config authz.Config) (*models.Role, error) {
	id := newRoleID()

	configBytes, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	if _, err := s.db.ExecContext(
		ctx,
		createRole,
		id,
		projectID,
		name,
		description,
		string(configBytes),
	); err != nil {
		return nil, err
	}

	return s.GetRole(ctx, id, projectID)
}

func (s *Store) GetRole(ctx context.Context, id, projectID string) (*models.Role, error) {
	roleRow := s.db.QueryRowContext(ctx, getRole, id, projectID)

	role, err := s.scanRole(roleRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrRoleNotFound
	} else if err != nil {
		return nil, err
	}

	return role, nil
}

func (s *Store) LookupRole(ctx context.Context, name, projectID string) (*models.Role, error) {
	roleRow := s.db.QueryRowContext(ctx, lookupRole, name, projectID)

	role, err := s.scanRole(roleRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrRoleNotFound
	} else if err != nil {
		return nil, err
	}

	return role, nil
}

func (s *Store) ListRoles(ctx context.Context, projectID string) ([]models.Role, error) {
	roleRows, err := s.db.QueryContext(ctx, listRoles, projectID)
	if err != nil {
		return nil, errors.Wrap(err, "query roles")
	}
	defer roleRows.Close()

	roles := make([]models.Role, 0)
	for roleRows.Next() {
		role, err := s.scanRole(roleRows)
		if err != nil {
			return nil, err
		}
		roles = append(roles, *role)
	}

	if err := roleRows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}

func (s *Store) UpdateRole(ctx context.Context, id, projectID, name, description string, config authz.Config) (*models.Role, error) {
	configBytes, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	if _, err := s.db.ExecContext(
		ctx,
		updateRole,
		name,
		description,
		string(configBytes),
		id,
		projectID,
	); err != nil {
		return nil, err
	}

	return s.GetRole(ctx, id, projectID)
}

func (s *Store) scanRole(scanner scanner) (*models.Role, error) {
	var role models.Role
	var configString string
	if err := scanner.Scan(
		&role.ID,
		&role.ProjectID,
		&role.Name,
		&role.Description,
		&configString,
	); err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal([]byte(configString), &role.Config); err != nil {
		return nil, err
	}
	return &role, nil
}

func (s *Store) CreateMembership(ctx context.Context, userID, projectID string) (*models.Membership, error) {
	id := newMembershipID()

	if _, err := s.db.ExecContext(
		ctx,
		createMembership,
		id,
		userID,
		projectID,
	); err != nil {
		return nil, err
	}

	return s.GetMembership(ctx, userID, projectID)
}

func (s *Store) GetMembership(ctx context.Context, userID, projectID string) (*models.Membership, error) {
	membershipRow := s.db.QueryRowContext(ctx, getMembership, userID, projectID)

	membership, err := s.scanMembership(membershipRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrMembershipNotFound
	} else if err != nil {
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

	memberships := make([]models.Membership, 0)
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
		&membership.ID,
		&membership.UserID,
		&membership.ProjectID,
	); err != nil {
		return nil, err
	}
	return &membership, nil
}

func (s *Store) CreateMembershipRoleBinding(ctx context.Context, membershipID, roleID, projectID string) (*models.MembershipRoleBinding, error) {
	if _, err := s.db.ExecContext(
		ctx,
		createMembershipRoleBinding,
		membershipID,
		roleID,
		projectID,
	); err != nil {
		return nil, err
	}

	return s.GetMembershipRoleBinding(ctx, membershipID, roleID, projectID)
}

func (s *Store) GetMembershipRoleBinding(ctx context.Context, membershipID, roleID, projectID string) (*models.MembershipRoleBinding, error) {
	membershipRoleBindingRow := s.db.QueryRowContext(ctx, getMembershipRoleBinding, membershipID, roleID, projectID)

	membershipRoleBinding, err := s.scanMembershipRoleBinding(membershipRoleBindingRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrMembershipRoleBindingNotFound
	} else if err != nil {
		return nil, err
	}

	return membershipRoleBinding, nil
}

func (s *Store) ListMembershipRoleBindings(ctx context.Context, membershipID, projectID string) ([]models.MembershipRoleBinding, error) {
	membershipRoleBindingRows, err := s.db.QueryContext(ctx, listMembershipRoleBindings, membershipID, projectID)
	if err != nil {
		return nil, errors.Wrap(err, "query membership role bindings")
	}
	defer membershipRoleBindingRows.Close()

	membershipRoleBindings := make([]models.MembershipRoleBinding, 0)
	for membershipRoleBindingRows.Next() {
		membershipRoleBinding, err := s.scanMembershipRoleBinding(membershipRoleBindingRows)
		if err != nil {
			return nil, err
		}
		membershipRoleBindings = append(membershipRoleBindings, *membershipRoleBinding)
	}

	if err := membershipRoleBindingRows.Err(); err != nil {
		return nil, err
	}

	return membershipRoleBindings, nil
}

func (s *Store) scanMembershipRoleBinding(scanner scanner) (*models.MembershipRoleBinding, error) {
	var membershipRoleBinding models.MembershipRoleBinding
	if err := scanner.Scan(
		&membershipRoleBinding.MembershipID,
		&membershipRoleBinding.RoleID,
		&membershipRoleBinding.ProjectID,
	); err != nil {
		return nil, err
	}
	return &membershipRoleBinding, nil
}

func (s *Store) CreateDevice(ctx context.Context, projectID string) (*models.Device, error) {
	id := newDeviceID()
	name := namesgenerator.GetRandomName(0)

	if _, err := s.db.ExecContext(
		ctx,
		createDevice,
		id,
		projectID,
		name,
	); err != nil {
		return nil, err
	}

	return s.GetDevice(ctx, id, projectID)
}

func (s *Store) GetDevice(ctx context.Context, id, projectID string) (*models.Device, error) {
	deviceRow := s.db.QueryRowContext(ctx, getDevice, id, projectID)

	device, err := s.scanDevice(deviceRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrDeviceNotFound
	} else if err != nil {
		return nil, err
	}

	return device, nil
}

func (s *Store) ListDevices(ctx context.Context, projectID string) ([]models.Device, error) {
	deviceRows, err := s.db.QueryContext(ctx, listDevices, projectID)
	if err != nil {
		return nil, err
	}
	defer deviceRows.Close()

	devices := make([]models.Device, 0)
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

func (s *Store) SetDeviceInfo(ctx context.Context, id, projectID string, deviceInfo models.DeviceInfo) (*models.Device, error) {
	infoBytes, err := json.Marshal(deviceInfo)
	if err != nil {
		return nil, err
	}

	if _, err := s.db.ExecContext(
		ctx,
		setDeviceInfo,
		string(infoBytes),
		id,
		projectID,
	); err != nil {
		return nil, err
	}

	return s.GetDevice(ctx, id, projectID)
}

func (s *Store) scanDevice(scanner scanner) (*models.Device, error) {
	var device models.Device
	var infoString string
	if err := scanner.Scan(
		&device.ID,
		&device.ProjectID,
		&device.Name,
		&infoString,
	); err != nil {
		return nil, err
	}
	if infoString != "" {
		if err := json.Unmarshal([]byte(infoString), &device.Info); err != nil {
			return nil, err
		}
	}
	return &device, nil
}

func (s *Store) SetDeviceLabel(ctx context.Context, key, deviceID, projectID, value string) (*models.DeviceLabel, error) {
	if _, err := s.db.ExecContext(
		ctx,
		setDeviceLabel,
		key,
		deviceID,
		projectID,
		value,
		value,
	); err != nil {
		return nil, err
	}

	return s.GetDeviceLabel(ctx, key, deviceID, projectID)
}

func (s *Store) GetDeviceLabel(ctx context.Context, key, deviceID, projectID string) (*models.DeviceLabel, error) {
	deviceLabelRow := s.db.QueryRowContext(ctx, getDeviceLabel, key, deviceID, projectID)

	deviceLabel, err := s.scanDeviceLabel(deviceLabelRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrDeviceLabelNotFound
	} else if err != nil {
		return nil, err
	}

	return deviceLabel, nil
}

func (s *Store) ListDeviceLabels(ctx context.Context, deviceID, projectID string) ([]models.DeviceLabel, error) {
	deviceLabelRows, err := s.db.QueryContext(ctx, listDeviceLabels, deviceID, projectID)
	if err != nil {
		return nil, err
	}
	defer deviceLabelRows.Close()

	deviceLabels := make([]models.DeviceLabel, 0)
	for deviceLabelRows.Next() {
		deviceLabel, err := s.scanDeviceLabel(deviceLabelRows)
		if err != nil {
			return nil, err
		}
		deviceLabels = append(deviceLabels, *deviceLabel)
	}

	if err := deviceLabelRows.Err(); err != nil {
		return nil, err
	}

	return deviceLabels, nil
}

func (s *Store) DeleteDeviceLabel(ctx context.Context, key, deviceID, projectID string) error {
	if _, err := s.db.ExecContext(
		ctx,
		deleteDeviceLabel,
		key,
		deviceID,
		projectID,
	); err != nil {
		return err
	}

	return nil
}

func (s *Store) scanDeviceLabel(scanner scanner) (*models.DeviceLabel, error) {
	var deviceLabel models.DeviceLabel
	if err := scanner.Scan(
		&deviceLabel.Key,
		&deviceLabel.DeviceID,
		&deviceLabel.ProjectID,
		&deviceLabel.Value,
	); err != nil {
		return nil, err
	}
	return &deviceLabel, nil
}

func (s *Store) CreateDeviceRegistrationToken(ctx context.Context, projectID string) (*models.DeviceRegistrationToken, error) {
	id := newDeviceRegistrationTokenID()

	if _, err := s.db.ExecContext(
		ctx,
		createDeviceRegistrationToken,
		id,
		projectID,
	); err != nil {
		return nil, err
	}

	return s.GetDeviceRegistrationToken(ctx, id, projectID)
}

func (s *Store) GetDeviceRegistrationToken(ctx context.Context, id, projectID string) (*models.DeviceRegistrationToken, error) {
	deviceRegistrationTokenRow := s.db.QueryRowContext(ctx, getDeviceRegistrationToken, id, projectID)

	deviceRegistrationToken, err := s.scanDeviceRegistrationToken(deviceRegistrationTokenRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrDeviceRegistrationTokenNotFound
	} else if err != nil {
		return nil, err
	}

	return deviceRegistrationToken, nil
}

func (s *Store) BindDeviceRegistrationToken(ctx context.Context, id, projectID, deviceAccessKeyID string) (*models.DeviceRegistrationToken, error) {
	if _, err := s.db.ExecContext(
		ctx,
		bindDeviceRegistrationToken,
		deviceAccessKeyID,
		id,
		projectID,
	); err != nil {
		return nil, err
	}

	return s.GetDeviceRegistrationToken(ctx, id, projectID)
}

func (s *Store) scanDeviceRegistrationToken(scanner scanner) (*models.DeviceRegistrationToken, error) {
	var deviceRegistrationToken models.DeviceRegistrationToken
	if err := scanner.Scan(
		&deviceRegistrationToken.ID,
		&deviceRegistrationToken.ProjectID,
		&deviceRegistrationToken.DeviceAccessKeyID,
	); err != nil {
		return nil, err
	}
	return &deviceRegistrationToken, nil
}

func (s *Store) CreateDeviceAccessKey(ctx context.Context, projectID, deviceID, hash string) (*models.DeviceAccessKey, error) {
	id := newDeviceAccessKeyID()

	if _, err := s.db.ExecContext(
		ctx,
		createDeviceAccessKey,
		id,
		projectID,
		deviceID,
		hash,
	); err != nil {
		return nil, err
	}

	return s.GetDeviceAccessKey(ctx, id, projectID)
}

func (s *Store) GetDeviceAccessKey(ctx context.Context, id, projectID string) (*models.DeviceAccessKey, error) {
	deviceAccessKeyRow := s.db.QueryRowContext(ctx, getDeviceAccessKey, id, projectID)

	deviceAccessKey, err := s.scanDeviceAccessKey(deviceAccessKeyRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrDeviceAccessKeyNotFound
	} else if err != nil {
		return nil, err
	}

	return deviceAccessKey, nil
}

func (s *Store) ValidateDeviceAccessKey(ctx context.Context, projectID, hash string) (*models.DeviceAccessKey, error) {
	deviceAccessKeyRow := s.db.QueryRowContext(ctx, validateDeviceAccessKey, projectID, hash)

	deviceAccessKey, err := s.scanDeviceAccessKey(deviceAccessKeyRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrDeviceAccessKeyNotFound
	} else if err != nil {
		return nil, err
	}

	return deviceAccessKey, nil
}

func (s *Store) scanDeviceAccessKey(scanner scanner) (*models.DeviceAccessKey, error) {
	var deviceAccessKey models.DeviceAccessKey
	if err := scanner.Scan(
		&deviceAccessKey.ID,
		&deviceAccessKey.ProjectID,
		&deviceAccessKey.DeviceID,
		&deviceAccessKey.Hash,
	); err != nil {
		return nil, err
	}
	return &deviceAccessKey, nil
}

func (s *Store) CreateApplication(ctx context.Context, projectID, name, description string,
	applicationSettings models.ApplicationSettings) (*models.Application, error) {
	id := newApplicationID()

	settingsBytes, err := json.Marshal(applicationSettings)
	if err != nil {
		return nil, err
	}

	if _, err := s.db.ExecContext(
		ctx,
		createApplication,
		id,
		projectID,
		name,
		description,
		string(settingsBytes),
	); err != nil {
		return nil, err
	}

	return s.GetApplication(ctx, id, projectID)
}

func (s *Store) GetApplication(ctx context.Context, id, projectID string) (*models.Application, error) {
	applicationRow := s.db.QueryRowContext(ctx, getApplication, id, projectID)

	application, err := s.scanApplication(applicationRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrApplicationNotFound
	} else if err != nil {
		return nil, err
	}

	return application, nil
}

func (s *Store) LookupApplication(ctx context.Context, name, projectID string) (*models.Application, error) {
	applicationRow := s.db.QueryRowContext(ctx, lookupApplication, name, projectID)

	application, err := s.scanApplication(applicationRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrApplicationNotFound
	} else if err != nil {
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

	applications := make([]models.Application, 0)
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

func (s *Store) UpdateApplication(ctx context.Context, id, projectID, name, description string,
	applicationSettings models.ApplicationSettings) (*models.Application, error) {
	settingsBytes, err := json.Marshal(applicationSettings)
	if err != nil {
		return nil, err
	}

	if _, err := s.db.ExecContext(
		ctx,
		updateApplication,
		name,
		description,
		string(settingsBytes),
		id,
		projectID,
	); err != nil {
		return nil, err
	}

	return s.GetApplication(ctx, id, projectID)
}

func (s *Store) scanApplication(scanner scanner) (*models.Application, error) {
	var application models.Application
	var settingsString string
	if err := scanner.Scan(
		&application.ID,
		&application.ProjectID,
		&application.Name,
		&application.Description,
		&settingsString,
	); err != nil {
		return nil, err
	}
	if settingsString != "" {
		if err := json.Unmarshal([]byte(settingsString), &application.Settings); err != nil {
			return nil, err
		}
	}
	return &application, nil
}

func (s *Store) CreateRelease(ctx context.Context, projectID, applicationID, config string) (*models.Release, error) {
	id := newReleaseID()

	if _, err := s.db.ExecContext(
		ctx,
		createRelease,
		id,
		projectID,
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

func (s *Store) GetRelease(ctx context.Context, id, projectID, applicationID string) (*models.Release, error) {
	applicationRow := s.db.QueryRowContext(ctx, getRelease, id, projectID, applicationID)

	release, err := s.scanRelease(applicationRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrReleaseNotFound
	} else if err != nil {
		return nil, err
	}

	return release, nil
}

func (s *Store) GetLatestRelease(ctx context.Context, projectID, applicationID string) (*models.Release, error) {
	applicationRow := s.db.QueryRowContext(ctx, getLatestRelease, projectID, applicationID)

	release, err := s.scanRelease(applicationRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrReleaseNotFound
	} else if err != nil {
		return nil, err
	}

	return release, nil
}

func (s *Store) ListReleases(ctx context.Context, projectID, applicationID string) ([]models.Release, error) {
	releaseRows, err := s.db.QueryContext(ctx, listReleases, projectID, applicationID)
	if err != nil {
		return nil, errors.Wrap(err, "query releases")
	}
	defer releaseRows.Close()

	releases := make([]models.Release, 0)
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
		&release.ProjectID,
		&release.ApplicationID,
		&release.Config,
	); err != nil {
		return nil, err
	}
	return &release, nil
}

func (s *Store) GetReleaseDeviceCounts(ctx context.Context, projectID, applicationID, releaseID string) (*models.ReleaseDeviceCounts, error) {
	countRow := s.db.QueryRowContext(ctx, getReleaseDeviceCounts, projectID, applicationID, releaseID)

	count, err := s.scanReleaseDeviceCountRow(countRow)
	if err != nil {
		return nil, err
	}

	return &models.ReleaseDeviceCounts{
		AllCount: count,
	}, nil
}

func (s *Store) scanReleaseDeviceCountRow(scanner scanner) (int, error) {
	var count int
	if err := scanner.Scan(
		&count,
	); err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Store) SetDeviceApplicationStatus(ctx context.Context, projectID, deviceID, applicationID, currentReleaseID string) error {
	_, err := s.db.ExecContext(
		ctx,
		setDeviceApplicationStatus,
		projectID,
		deviceID,
		applicationID,
		currentReleaseID,
		currentReleaseID,
	)
	return err
}

func (s *Store) GetDeviceApplicationStatus(ctx context.Context, projectID, deviceID, applicationID string) (*models.DeviceApplicationStatus, error) {
	deviceApplicationStatusRow := s.db.QueryRowContext(ctx, getDeviceApplicationStatus, projectID, deviceID, applicationID)

	deviceApplicationStatus, err := s.scanDeviceApplicationStatus(deviceApplicationStatusRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrDeviceApplicationStatusNotFound
	} else if err != nil {
		return nil, err
	}

	return deviceApplicationStatus, nil
}

func (s *Store) scanDeviceApplicationStatus(scanner scanner) (*models.DeviceApplicationStatus, error) {
	var deviceApplicationStatus models.DeviceApplicationStatus
	if err := scanner.Scan(
		&deviceApplicationStatus.ProjectID,
		&deviceApplicationStatus.DeviceID,
		&deviceApplicationStatus.ApplicationID,
		&deviceApplicationStatus.CurrentReleaseID,
	); err != nil {
		return nil, err
	}
	return &deviceApplicationStatus, nil
}

func (s *Store) SetDeviceServiceStatus(ctx context.Context, projectID, deviceID, applicationID, service, currentReleaseID string) error {
	_, err := s.db.ExecContext(
		ctx,
		setDeviceServiceStatus,
		projectID,
		deviceID,
		applicationID,
		service,
		currentReleaseID,
		currentReleaseID,
	)
	return err
}

func (s *Store) GetDeviceServiceStatus(ctx context.Context, projectID, deviceID, applicationID, service string) (*models.DeviceServiceStatus, error) {
	deviceServiceStatusRow := s.db.QueryRowContext(ctx, getDeviceServiceStatus, projectID, deviceID, applicationID, service)

	deviceServiceStatus, err := s.scanDeviceServiceStatus(deviceServiceStatusRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrDeviceApplicationStatusNotFound
	} else if err != nil {
		return nil, err
	}

	return deviceServiceStatus, nil
}

func (s *Store) GetDeviceServiceStatuses(ctx context.Context, projectID, deviceID, applicationID string) ([]models.DeviceServiceStatus, error) {
	deviceServiceStatusRows, err := s.db.QueryContext(ctx, getDeviceServiceStatuses, projectID, deviceID, applicationID)
	if err != nil {
		return nil, errors.Wrap(err, "query device service statuses")
	}
	defer deviceServiceStatusRows.Close()

	deviceServiceStatuses := make([]models.DeviceServiceStatus, 0)
	for deviceServiceStatusRows.Next() {
		deviceServiceStatus, err := s.scanDeviceServiceStatus(deviceServiceStatusRows)
		if err != nil {
			return nil, err
		}
		deviceServiceStatuses = append(deviceServiceStatuses, *deviceServiceStatus)
	}

	if err := deviceServiceStatusRows.Err(); err != nil {
		return nil, err
	}

	return deviceServiceStatuses, nil
}

func (s *Store) scanDeviceServiceStatus(scanner scanner) (*models.DeviceServiceStatus, error) {
	var deviceServiceStatus models.DeviceServiceStatus
	if err := scanner.Scan(
		&deviceServiceStatus.ProjectID,
		&deviceServiceStatus.DeviceID,
		&deviceServiceStatus.ApplicationID,
		&deviceServiceStatus.Service,
		&deviceServiceStatus.CurrentReleaseID,
	); err != nil {
		return nil, err
	}
	return &deviceServiceStatus, nil
}
