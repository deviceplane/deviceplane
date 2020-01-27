package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/deviceplane/deviceplane/pkg/controller/store"
	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/pkg/errors"

	"github.com/segmentio/ksuid"
	"gopkg.in/yaml.v2"
)

type scanner interface {
	Scan(...interface{}) error
}

const (
	userPrefix                      = "usr"
	registrationTokenPrefix         = "reg"
	passwordRecoveryTokenPrefix     = "pwr"
	sessionPrefix                   = "ses"
	userAccessKeyPrefix             = "uky"
	projectPrefix                   = "prj"
	rolePrefix                      = "rol"
	serviceAccountPrefix            = "sac"
	serviceAccountAccessKeyPrefix   = "sak"
	devicePrefix                    = "dev"
	deviceRegistrationTokenPrefix   = "drt"
	deviceAccessKeyPrefix           = "dak"
	applicationPrefix               = "app"
	releasePrefix                   = "rel"
	ExposedMetricConfigHolderPrefix = "mtc"
)

func newUserID() string {
	return fmt.Sprintf("%s_%s", userPrefix, ksuid.New().String())
}

func newRegistrationTokenID() string {
	return fmt.Sprintf("%s_%s", registrationTokenPrefix, ksuid.New().String())
}

func newPasswordRecoveryTokenID() string {
	return fmt.Sprintf("%s_%s", passwordRecoveryTokenPrefix, ksuid.New().String())
}

func newSessionID() string {
	return fmt.Sprintf("%s_%s", sessionPrefix, ksuid.New().String())
}

func newUserAccessKeyID() string {
	return fmt.Sprintf("%s_%s", userAccessKeyPrefix, ksuid.New().String())
}

func newProjectID() string {
	return fmt.Sprintf("%s_%s", projectPrefix, ksuid.New().String())
}

func newRoleID() string {
	return fmt.Sprintf("%s_%s", rolePrefix, ksuid.New().String())
}

func newServiceAccountID() string {
	return fmt.Sprintf("%s_%s", serviceAccountPrefix, ksuid.New().String())
}

func newServiceAccountAccessKeyID() string {
	return fmt.Sprintf("%s_%s", serviceAccountAccessKeyPrefix, ksuid.New().String())
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

func newExposedMetricConfigHolderID() string {
	return fmt.Sprintf("%s_%s", ExposedMetricConfigHolderPrefix, ksuid.New().String())
}

var (
	_ store.Users                      = &Store{}
	_ store.RegistrationTokens         = &Store{}
	_ store.PasswordRecoveryTokens     = &Store{}
	_ store.Sessions                   = &Store{}
	_ store.UserAccessKeys             = &Store{}
	_ store.Projects                   = &Store{}
	_ store.ProjectDeviceCounts        = &Store{}
	_ store.Roles                      = &Store{}
	_ store.Memberships                = &Store{}
	_ store.MembershipRoleBindings     = &Store{}
	_ store.ServiceAccounts            = &Store{}
	_ store.ServiceAccountAccessKeys   = &Store{}
	_ store.ServiceAccountRoleBindings = &Store{}
	_ store.Devices                    = &Store{}
	_ store.DeviceAccessKeys           = &Store{}
	_ store.DeviceRegistrationTokens   = &Store{}
	_ store.Applications               = &Store{}
	_ store.Releases                   = &Store{}
	_ store.ReleaseDeviceCounts        = &Store{}
	_ store.DeviceApplicationStatuses  = &Store{}
	_ store.DeviceServiceStatuses      = &Store{}
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) CreateUser(ctx context.Context, email, passwordHash, firstName, lastName, company string) (*models.User, error) {
	id := newUserID()

	if _, err := s.db.ExecContext(
		ctx,
		createUser,
		id,
		email,
		passwordHash,
		firstName,
		lastName,
		company,
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

func (s *Store) LookupUser(ctx context.Context, email string) (*models.User, error) {
	userRow := s.db.QueryRowContext(ctx, lookupUser, email)

	user, err := s.scanUser(userRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Store) ValidateUser(ctx context.Context, id, passwordHash string) (*models.User, error) {
	userRow := s.db.QueryRowContext(ctx, validateUser, id, passwordHash)

	user, err := s.scanUser(userRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Store) ValidateUserWithEmail(ctx context.Context, email, passwordHash string) (*models.User, error) {
	userRow := s.db.QueryRowContext(ctx, validateUserWithEmail, email, passwordHash)

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

func (s *Store) UpdatePasswordHash(ctx context.Context, id, passwordHash string) (*models.User, error) {
	if _, err := s.db.ExecContext(
		ctx,
		updatePasswordHash,
		passwordHash,
		id,
	); err != nil {
		return nil, err
	}

	return s.GetUser(ctx, id)
}

func (s *Store) UpdateFirstName(ctx context.Context, id, firstName string) (*models.User, error) {
	if _, err := s.db.ExecContext(
		ctx,
		updateFirstName,
		firstName,
		id,
	); err != nil {
		return nil, err
	}

	return s.GetUser(ctx, id)
}

func (s *Store) UpdateLastName(ctx context.Context, id, lastName string) (*models.User, error) {
	if _, err := s.db.ExecContext(
		ctx,
		updateLastName,
		lastName,
		id,
	); err != nil {
		return nil, err
	}

	return s.GetUser(ctx, id)
}

func (s *Store) UpdateCompany(ctx context.Context, id, company string) (*models.User, error) {
	if _, err := s.db.ExecContext(
		ctx,
		updateCompany,
		company,
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
		&user.CreatedAt,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Company,
		&user.RegistrationCompleted,
		&user.SuperAdmin,
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
		&registrationToken.CreatedAt,
		&registrationToken.UserID,
	); err != nil {
		return nil, err
	}
	return &registrationToken, nil
}

func (s *Store) CreatePasswordRecoveryToken(ctx context.Context, userID, hash string) (*models.PasswordRecoveryToken, error) {
	id := newPasswordRecoveryTokenID()

	if _, err := s.db.ExecContext(
		ctx,
		createPasswordRecoveryToken,
		id,
		userID,
		hash,
	); err != nil {
		return nil, err
	}

	return s.GetPasswordRecoveryToken(ctx, id)
}

func (s *Store) GetPasswordRecoveryToken(ctx context.Context, id string) (*models.PasswordRecoveryToken, error) {
	passwordRecoveryTokenRow := s.db.QueryRowContext(ctx, getPasswordRecoveryToken, id)

	passwordRecoveryToken, err := s.scanPasswordRecoveryToken(passwordRecoveryTokenRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrPasswordRecoveryTokenNotFound
	} else if err != nil {
		return nil, err
	}

	return passwordRecoveryToken, nil
}

func (s *Store) ValidatePasswordRecoveryToken(ctx context.Context, hash string) (*models.PasswordRecoveryToken, error) {
	passwordRecoveryTokenRow := s.db.QueryRowContext(ctx, validatePasswordRecoveryToken, hash)

	passwordRecoveryToken, err := s.scanPasswordRecoveryToken(passwordRecoveryTokenRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrPasswordRecoveryTokenNotFound
	} else if err != nil {
		return nil, err
	}

	return passwordRecoveryToken, nil
}

func (s *Store) scanPasswordRecoveryToken(scanner scanner) (*models.PasswordRecoveryToken, error) {
	var passwordRecoveryToken models.PasswordRecoveryToken
	if err := scanner.Scan(
		&passwordRecoveryToken.ID,
		&passwordRecoveryToken.CreatedAt,
		&passwordRecoveryToken.ExpiresAt,
		&passwordRecoveryToken.UserID,
	); err != nil {
		return nil, err
	}
	return &passwordRecoveryToken, nil
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
		&session.CreatedAt,
		&session.UserID,
	); err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *Store) CreateUserAccessKey(ctx context.Context, userID, hash, description string) (*models.UserAccessKey, error) {
	id := newUserAccessKeyID()

	if _, err := s.db.ExecContext(
		ctx,
		createUserAccessKey,
		id,
		userID,
		hash,
		description,
	); err != nil {
		return nil, err
	}

	return s.GetUserAccessKey(ctx, id)
}

func (s *Store) GetUserAccessKey(ctx context.Context, id string) (*models.UserAccessKey, error) {
	userAccessKeyRow := s.db.QueryRowContext(ctx, getUserAccessKey, id)

	userAccessKey, err := s.scanUserAccessKey(userAccessKeyRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrUserAccessKeyNotFound
	} else if err != nil {
		return nil, err
	}

	return userAccessKey, nil
}

func (s *Store) ValidateUserAccessKey(ctx context.Context, hash string) (*models.UserAccessKey, error) {
	userAccessKeyRow := s.db.QueryRowContext(ctx, validateUserAccessKey, hash)

	userAccessKey, err := s.scanUserAccessKey(userAccessKeyRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrUserAccessKeyNotFound
	} else if err != nil {
		return nil, err
	}

	return userAccessKey, nil
}

func (s *Store) ListUserAccessKeys(ctx context.Context, projectID string) ([]models.UserAccessKey, error) {
	userAccessKeyRows, err := s.db.QueryContext(ctx, listUserAccessKeys, projectID)
	if err != nil {
		return nil, errors.Wrap(err, "query user access keys")
	}
	defer userAccessKeyRows.Close()

	userAccessKeys := make([]models.UserAccessKey, 0)
	for userAccessKeyRows.Next() {
		userAccessKey, err := s.scanUserAccessKey(userAccessKeyRows)
		if err != nil {
			return nil, err
		}
		userAccessKeys = append(userAccessKeys, *userAccessKey)
	}

	if err := userAccessKeyRows.Err(); err != nil {
		return nil, err
	}

	return userAccessKeys, nil
}

func (s *Store) DeleteUserAccessKey(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(
		ctx,
		deleteUserAccessKey,
		id,
	)
	return err
}

func (s *Store) scanUserAccessKey(scanner scanner) (*models.UserAccessKey, error) {
	var userAccessKey models.UserAccessKey
	if err := scanner.Scan(
		&userAccessKey.ID,
		&userAccessKey.CreatedAt,
		&userAccessKey.UserID,
		&userAccessKey.Description,
	); err != nil {
		return nil, err
	}
	return &userAccessKey, nil
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

func (s *Store) ListProjects(ctx context.Context) ([]models.Project, error) {
	projectRows, err := s.db.QueryContext(ctx, listProjects)
	if err != nil {
		return nil, errors.Wrap(err, "query projects")
	}
	defer projectRows.Close()

	projects := make([]models.Project, 0)
	for projectRows.Next() {
		project, err := s.scanProject(projectRows)
		if err != nil {
			return nil, err
		}
		projects = append(projects, *project)
	}

	if err := projectRows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func (s *Store) UpdateProject(ctx context.Context, id, name, datadogApiKey string) (*models.Project, error) {
	if _, err := s.db.ExecContext(
		ctx,
		updateProject,
		name,
		datadogApiKey,
		id,
	); err != nil {
		return nil, err
	}

	return s.GetProject(ctx, id)
}

func (s *Store) DeleteProject(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(
		ctx,
		deleteProject,
		id,
	)
	return err
}

func (s *Store) scanProject(scanner scanner) (*models.Project, error) {
	var project models.Project
	if err := scanner.Scan(
		&project.ID,
		&project.CreatedAt,
		&project.Name,
		&project.DatadogAPIKey,
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

func (s *Store) CreateRole(ctx context.Context, projectID, name, description, config string) (*models.Role, error) {
	id := newRoleID()

	if _, err := s.db.ExecContext(
		ctx,
		createRole,
		id,
		projectID,
		name,
		description,
		config,
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

func (s *Store) UpdateRole(ctx context.Context, id, projectID, name, description, config string) (*models.Role, error) {
	if _, err := s.db.ExecContext(
		ctx,
		updateRole,
		name,
		description,
		config,
		id,
		projectID,
	); err != nil {
		return nil, err
	}

	return s.GetRole(ctx, id, projectID)
}

func (s *Store) DeleteRole(ctx context.Context, id, projectID string) error {
	_, err := s.db.ExecContext(
		ctx,
		deleteRole,
		id,
		projectID,
	)
	return err
}

func (s *Store) scanRole(scanner scanner) (*models.Role, error) {
	var role models.Role
	if err := scanner.Scan(
		&role.ID,
		&role.CreatedAt,
		&role.ProjectID,
		&role.Name,
		&role.Description,
		&role.Config,
	); err != nil {
		return nil, err
	}
	return &role, nil
}

func (s *Store) CreateMembership(ctx context.Context, userID, projectID string) (*models.Membership, error) {
	if _, err := s.db.ExecContext(
		ctx,
		createMembership,
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

func (s *Store) DeleteMembership(ctx context.Context, userID, projectID string) error {
	_, err := s.db.ExecContext(
		ctx,
		deleteMembership,
		userID,
		projectID,
	)
	return err
}

func (s *Store) scanMembership(scanner scanner) (*models.Membership, error) {
	var membership models.Membership
	if err := scanner.Scan(
		&membership.UserID,
		&membership.ProjectID,
		&membership.CreatedAt,
	); err != nil {
		return nil, err
	}
	return &membership, nil
}

func (s *Store) CreateMembershipRoleBinding(ctx context.Context, userID, roleID, projectID string) (*models.MembershipRoleBinding, error) {
	if _, err := s.db.ExecContext(
		ctx,
		createMembershipRoleBinding,
		userID,
		roleID,
		projectID,
	); err != nil {
		return nil, err
	}

	return s.GetMembershipRoleBinding(ctx, userID, roleID, projectID)
}

func (s *Store) GetMembershipRoleBinding(ctx context.Context, userID, roleID, projectID string) (*models.MembershipRoleBinding, error) {
	membershipRoleBindingRow := s.db.QueryRowContext(ctx, getMembershipRoleBinding, userID, roleID, projectID)

	membershipRoleBinding, err := s.scanMembershipRoleBinding(membershipRoleBindingRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrMembershipRoleBindingNotFound
	} else if err != nil {
		return nil, err
	}

	return membershipRoleBinding, nil
}

func (s *Store) ListMembershipRoleBindings(ctx context.Context, userID, projectID string) ([]models.MembershipRoleBinding, error) {
	membershipRoleBindingRows, err := s.db.QueryContext(ctx, listMembershipRoleBindings, userID, projectID)
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

func (s *Store) DeleteMembershipRoleBinding(ctx context.Context, userID, roleID, projectID string) error {
	_, err := s.db.ExecContext(
		ctx,
		deleteMembershipRoleBinding,
		userID,
		roleID,
		projectID,
	)
	return err
}

func (s *Store) scanMembershipRoleBinding(scanner scanner) (*models.MembershipRoleBinding, error) {
	var membershipRoleBinding models.MembershipRoleBinding
	if err := scanner.Scan(
		&membershipRoleBinding.UserID,
		&membershipRoleBinding.RoleID,
		&membershipRoleBinding.CreatedAt,
		&membershipRoleBinding.ProjectID,
	); err != nil {
		return nil, err
	}
	return &membershipRoleBinding, nil
}

func (s *Store) CreateServiceAccount(ctx context.Context, projectID, name, description string) (*models.ServiceAccount, error) {
	id := newServiceAccountID()

	if _, err := s.db.ExecContext(
		ctx,
		createServiceAccount,
		id,
		projectID,
		name,
		description,
	); err != nil {
		return nil, err
	}

	return s.GetServiceAccount(ctx, id, projectID)
}

func (s *Store) GetServiceAccount(ctx context.Context, id, projectID string) (*models.ServiceAccount, error) {
	serviceAccountRow := s.db.QueryRowContext(ctx, getServiceAccount, id, projectID)

	serviceAccount, err := s.scanServiceAccount(serviceAccountRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrServiceAccountNotFound
	} else if err != nil {
		return nil, err
	}

	return serviceAccount, nil
}

func (s *Store) LookupServiceAccount(ctx context.Context, name, projectID string) (*models.ServiceAccount, error) {
	serviceAccountRow := s.db.QueryRowContext(ctx, lookupServiceAccount, name, projectID)

	serviceAccount, err := s.scanServiceAccount(serviceAccountRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrServiceAccountNotFound
	} else if err != nil {
		return nil, err
	}

	return serviceAccount, nil
}

func (s *Store) ListServiceAccounts(ctx context.Context, projectID string) ([]models.ServiceAccount, error) {
	serviceAccountRows, err := s.db.QueryContext(ctx, listServiceAccounts, projectID)
	if err != nil {
		return nil, errors.Wrap(err, "query service accounts")
	}
	defer serviceAccountRows.Close()

	serviceAccounts := make([]models.ServiceAccount, 0)
	for serviceAccountRows.Next() {
		serviceAccount, err := s.scanServiceAccount(serviceAccountRows)
		if err != nil {
			return nil, err
		}
		serviceAccounts = append(serviceAccounts, *serviceAccount)
	}

	if err := serviceAccountRows.Err(); err != nil {
		return nil, err
	}

	return serviceAccounts, nil
}

func (s *Store) UpdateServiceAccount(ctx context.Context, id, projectID, name, description string) (*models.ServiceAccount, error) {
	if _, err := s.db.ExecContext(
		ctx,
		updateServiceAccount,
		name,
		description,
		id,
		projectID,
	); err != nil {
		return nil, err
	}

	return s.GetServiceAccount(ctx, id, projectID)
}

func (s *Store) DeleteServiceAccount(ctx context.Context, id, projectID string) error {
	_, err := s.db.ExecContext(
		ctx,
		deleteServiceAccount,
		id,
		projectID,
	)
	return err
}

func (s *Store) scanServiceAccount(scanner scanner) (*models.ServiceAccount, error) {
	var serviceAccount models.ServiceAccount
	if err := scanner.Scan(
		&serviceAccount.ID,
		&serviceAccount.CreatedAt,
		&serviceAccount.ProjectID,
		&serviceAccount.Name,
		&serviceAccount.Description,
	); err != nil {
		return nil, err
	}
	return &serviceAccount, nil
}

func (s *Store) CreateServiceAccountAccessKey(ctx context.Context, projectID, serviceAccountID, hash, description string) (*models.ServiceAccountAccessKey, error) {
	id := newServiceAccountAccessKeyID()

	if _, err := s.db.ExecContext(
		ctx,
		createServiceAccountAccessKey,
		id,
		projectID,
		serviceAccountID,
		hash,
		description,
	); err != nil {
		return nil, err
	}

	return s.GetServiceAccountAccessKey(ctx, id, projectID)
}

func (s *Store) GetServiceAccountAccessKey(ctx context.Context, id, projectID string) (*models.ServiceAccountAccessKey, error) {
	serviceAccountAccessKeyRow := s.db.QueryRowContext(ctx, getServiceAccountAccessKey, id, projectID)

	serviceAccountAccessKey, err := s.scanServiceAccountAccessKey(serviceAccountAccessKeyRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrServiceAccountAccessKeyNotFound
	} else if err != nil {
		return nil, err
	}

	return serviceAccountAccessKey, nil
}

func (s *Store) ValidateServiceAccountAccessKey(ctx context.Context, hash string) (*models.ServiceAccountAccessKey, error) {
	serviceAccountAccessKeyRow := s.db.QueryRowContext(ctx, validateServiceAccountAccessKey, hash)

	serviceAccountAccessKey, err := s.scanServiceAccountAccessKey(serviceAccountAccessKeyRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrServiceAccountAccessKeyNotFound
	} else if err != nil {
		return nil, err
	}

	return serviceAccountAccessKey, nil
}

func (s *Store) ListServiceAccountAccessKeys(ctx context.Context, projectID, serviceAccountID string) ([]models.ServiceAccountAccessKey, error) {
	serviceAccountAccessKeyRows, err := s.db.QueryContext(ctx, listServiceAccountAccessKeys, projectID, serviceAccountID)
	if err != nil {
		return nil, errors.Wrap(err, "query service account access keys")
	}
	defer serviceAccountAccessKeyRows.Close()

	serviceAccountAccessKeys := make([]models.ServiceAccountAccessKey, 0)
	for serviceAccountAccessKeyRows.Next() {
		serviceAccountAccessKey, err := s.scanServiceAccountAccessKey(serviceAccountAccessKeyRows)
		if err != nil {
			return nil, err
		}
		serviceAccountAccessKeys = append(serviceAccountAccessKeys, *serviceAccountAccessKey)
	}

	if err := serviceAccountAccessKeyRows.Err(); err != nil {
		return nil, err
	}

	return serviceAccountAccessKeys, nil
}

func (s *Store) DeleteServiceAccountAccessKey(ctx context.Context, id, projectID string) error {
	_, err := s.db.ExecContext(
		ctx,
		deleteServiceAccountAccessKey,
		id,
		projectID,
	)
	return err
}

func (s *Store) scanServiceAccountAccessKey(scanner scanner) (*models.ServiceAccountAccessKey, error) {
	var serviceAccountAccessKey models.ServiceAccountAccessKey
	if err := scanner.Scan(
		&serviceAccountAccessKey.ID,
		&serviceAccountAccessKey.CreatedAt,
		&serviceAccountAccessKey.ProjectID,
		&serviceAccountAccessKey.ServiceAccountID,
		&serviceAccountAccessKey.Description,
	); err != nil {
		return nil, err
	}
	return &serviceAccountAccessKey, nil
}

func (s *Store) CreateServiceAccountRoleBinding(ctx context.Context, serviceAccountID, roleID, projectID string) (*models.ServiceAccountRoleBinding, error) {
	if _, err := s.db.ExecContext(
		ctx,
		createServiceAccountRoleBinding,
		serviceAccountID,
		roleID,
		projectID,
	); err != nil {
		return nil, err
	}

	return s.GetServiceAccountRoleBinding(ctx, serviceAccountID, roleID, projectID)
}

func (s *Store) GetServiceAccountRoleBinding(ctx context.Context, serviceAccountID, roleID, projectID string) (*models.ServiceAccountRoleBinding, error) {
	serviceAccountRoleBindingRow := s.db.QueryRowContext(ctx, getServiceAccountRoleBinding, serviceAccountID, roleID, projectID)

	serviceAccountRoleBinding, err := s.scanServiceAccountRoleBinding(serviceAccountRoleBindingRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrServiceAccountRoleBindingNotFound
	} else if err != nil {
		return nil, err
	}

	return serviceAccountRoleBinding, nil
}

func (s *Store) ListServiceAccountRoleBindings(ctx context.Context, serviceAccountID, projectID string) ([]models.ServiceAccountRoleBinding, error) {
	serviceAccountRoleBindingRows, err := s.db.QueryContext(ctx, listServiceAccountRoleBindings, serviceAccountID, projectID)
	if err != nil {
		return nil, errors.Wrap(err, "query service account role bindings")
	}
	defer serviceAccountRoleBindingRows.Close()

	serviceAccountRoleBindings := make([]models.ServiceAccountRoleBinding, 0)
	for serviceAccountRoleBindingRows.Next() {
		serviceAccountRoleBinding, err := s.scanServiceAccountRoleBinding(serviceAccountRoleBindingRows)
		if err != nil {
			return nil, err
		}
		serviceAccountRoleBindings = append(serviceAccountRoleBindings, *serviceAccountRoleBinding)
	}

	if err := serviceAccountRoleBindingRows.Err(); err != nil {
		return nil, err
	}

	return serviceAccountRoleBindings, nil
}

func (s *Store) DeleteServiceAccountRoleBinding(ctx context.Context, serviceAccountID, roleID, projectID string) error {
	_, err := s.db.ExecContext(
		ctx,
		deleteServiceAccountRoleBinding,
		serviceAccountID,
		roleID,
		projectID,
	)
	return err
}

func (s *Store) scanServiceAccountRoleBinding(scanner scanner) (*models.ServiceAccountRoleBinding, error) {
	var serviceAccountRoleBinding models.ServiceAccountRoleBinding
	if err := scanner.Scan(
		&serviceAccountRoleBinding.ServiceAccountID,
		&serviceAccountRoleBinding.RoleID,
		&serviceAccountRoleBinding.CreatedAt,
		&serviceAccountRoleBinding.ProjectID,
	); err != nil {
		return nil, err
	}
	return &serviceAccountRoleBinding, nil
}

func (s *Store) CreateDevice(ctx context.Context, projectID, name, deviceRegistrationTokenID string, deviceLabels map[string]string) (*models.Device, error) {
	deviceID := newDeviceID()

	serializedDeviceLabels, err := json.Marshal(deviceLabels)
	if err != nil {
		return nil, err
	}

	if _, err := s.db.ExecContext(
		ctx,
		createDevice,
		deviceID,
		projectID,
		name,
		deviceRegistrationTokenID,
		string(serializedDeviceLabels),
	); err != nil {
		return nil, err
	}

	return s.GetDevice(ctx, deviceID, projectID)
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

func (s *Store) LookupDevice(ctx context.Context, name, projectID string) (*models.Device, error) {
	deviceRow := s.db.QueryRowContext(ctx, lookupDevice, name, projectID)

	device, err := s.scanDevice(deviceRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrDeviceNotFound
	} else if err != nil {
		return nil, err
	}

	return device, nil
}

func (s *Store) ListDevices(ctx context.Context, projectID, searchQuery string) ([]models.Device, error) {
	var deviceRows *sql.Rows
	var err error
	if searchQuery == "" {
		deviceRows, err = s.db.QueryContext(ctx, listDevices, projectID)
	} else {
		deviceRows, err = s.db.QueryContext(ctx, searchDevices, projectID, searchQuery)
	}
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

func (s *Store) UpdateDeviceName(ctx context.Context, id, projectID, name string) (*models.Device, error) {
	if _, err := s.db.ExecContext(
		ctx,
		updateDeviceName,
		name,
		id,
		projectID,
	); err != nil {
		return nil, err
	}

	return s.GetDevice(ctx, id, projectID)
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

func (s *Store) UpdateDeviceLastSeenAt(ctx context.Context, projectID, deviceID string) error {
	if _, err := s.db.ExecContext(
		ctx,
		updateDeviceLastSeenAt,
		projectID,
		deviceID,
	); err != nil {
		return err
	}

	return nil
}

func (s *Store) DeleteDevice(ctx context.Context, id, projectID string) error {
	_, err := s.db.ExecContext(
		ctx,
		deleteDevice,
		id,
		projectID,
	)
	return err
}

func (s *Store) scanDevice(scanner scanner) (*models.Device, error) {
	var device models.Device
	var infoString string
	var labelsString string
	if err := scanner.Scan(
		&device.ID,
		&device.CreatedAt,
		&device.ProjectID,
		&device.Name,
		&device.RegistrationTokenID,
		&device.DesiredAgentSpec,
		&device.DesiredAgentVersion,
		&infoString,
		&labelsString,
		&device.LastSeenAt,
	); err != nil {
		return nil, err
	}

	if infoString != "" {
		if err := json.Unmarshal([]byte(infoString), &device.Info); err != nil {
			return nil, err
		}
	}

	if labelsString == "" {
		device.Labels = map[string]string{}
	} else {
		if err := json.Unmarshal([]byte(labelsString), &device.Labels); err != nil {
			return nil, err
		}
	}

	if time.Now().After(device.LastSeenAt.Add(2 * time.Minute)) {
		device.Status = models.DeviceStatusOffline
	} else {
		device.Status = models.DeviceStatusOnline
	}

	return &device, nil
}

func (s *Store) scanDeviceLabels(scanner scanner) (map[string]string, error) {
	var labelsString string
	if err := scanner.Scan(
		&labelsString,
	); err != nil {
		return nil, err
	}

	var labels map[string]string
	if labelsString == "" {
		labels = map[string]string{}
	} else {
		if err := json.Unmarshal([]byte(labelsString), &labels); err != nil {
			return nil, err
		}
	}

	return labels, nil
}

func (s *Store) ListAllDeviceLabelKeys(ctx context.Context, projectID string) ([]string, error) {
	rows, err := s.db.QueryContext(
		ctx,
		listAllDeviceLabels,
		projectID,
	)
	if err != nil {
		return nil, errors.Wrap(err, "query device labels")
	}
	defer rows.Close()

	allDeviceLabels := make(map[string]bool)
	for rows.Next() {
		deviceLabels, err := s.scanDeviceLabels(rows)
		if err != nil {
			return nil, err
		}
		for k, _ := range deviceLabels {
			allDeviceLabels[k] = true
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	allDeviceLabelKeys := make([]string, len(allDeviceLabels))
	i := 0
	for k, _ := range allDeviceLabels {
		allDeviceLabelKeys[i] = k
		i++
	}

	return allDeviceLabelKeys, nil

}

func (s *Store) SetDeviceLabel(ctx context.Context, deviceID, projectID, key, value string) (*string, error) {
	device, err := s.GetDevice(ctx, deviceID, projectID)
	if err != nil {
		return nil, err
	}

	device.Labels[key] = value

	labelsString, err := json.Marshal(device.Labels)
	if err != nil {
		return nil, err
	}

	if _, err := s.db.ExecContext(
		ctx,
		updateDeviceLabels,
		labelsString,
		deviceID,
		projectID,
	); err != nil {
		return nil, err
	}

	device, err = s.GetDevice(ctx, deviceID, projectID)
	if err != nil {
		return nil, err
	}
	v := device.Labels[key]
	return &v, nil
}

func (s *Store) DeleteDeviceLabel(ctx context.Context, deviceID, projectID, key string) error {
	device, err := s.GetDevice(ctx, deviceID, projectID)
	if err != nil {
		return err
	}

	delete(device.Labels, key)

	labelsString, err := json.Marshal(device.Labels)
	if err != nil {
		return err
	}

	if _, err := s.db.ExecContext(
		ctx,
		updateDeviceLabels,
		labelsString,
		deviceID,
		projectID,
	); err != nil {
		return err
	}
	return nil
}

func (s *Store) CreateDeviceRegistrationToken(ctx context.Context, projectID, name, description string, maxRegistrations *int) (*models.DeviceRegistrationToken, error) {
	id := newDeviceRegistrationTokenID()

	if _, err := s.db.ExecContext(
		ctx,
		createDeviceRegistrationToken,
		id,
		projectID,
		name,
		description,
		maxRegistrations,
	); err != nil {
		return nil, err
	}

	return s.GetDeviceRegistrationToken(ctx, id, projectID)
}

func (s *Store) LookupDeviceRegistrationToken(ctx context.Context, name, projectID string) (*models.DeviceRegistrationToken, error) {
	deviceRegistrationTokenRow := s.db.QueryRowContext(ctx, lookupDeviceRegistrationToken, name, projectID)

	deviceRegistrationToken, err := s.scanDeviceRegistrationToken(deviceRegistrationTokenRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrDeviceRegistrationTokenNotFound
	} else if err != nil {
		return nil, err
	}

	return deviceRegistrationToken, nil
}

func (s *Store) ListDeviceRegistrationTokens(ctx context.Context, projectID string) ([]models.DeviceRegistrationToken, error) {
	tokenRows, err := s.db.QueryContext(ctx, listDeviceRegistrationTokens, projectID)
	if err != nil {
		return nil, errors.Wrap(err, "query device registration tokens")
	}
	defer tokenRows.Close()

	deviceRegistrationTokens := make([]models.DeviceRegistrationToken, 0)
	for tokenRows.Next() {
		deviceRegistrationToken, err := s.scanDeviceRegistrationToken(tokenRows)
		if err != nil {
			return nil, err
		}
		deviceRegistrationTokens = append(deviceRegistrationTokens, *deviceRegistrationToken)
	}

	if err := tokenRows.Err(); err != nil {
		return nil, err
	}

	return deviceRegistrationTokens, nil
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

func (s *Store) UpdateDeviceRegistrationToken(ctx context.Context, id, projectID, name, description string, maxRegistrations *int) (*models.DeviceRegistrationToken, error) {
	if _, err := s.db.ExecContext(
		ctx,
		updateDeviceRegistrationToken,
		name,
		description,
		maxRegistrations,
		id,
		projectID,
	); err != nil {
		return nil, err
	}

	return s.GetDeviceRegistrationToken(ctx, id, projectID)
}

func (s *Store) DeleteDeviceRegistrationToken(ctx context.Context, id, projectID string) error {
	_, err := s.db.ExecContext(
		ctx,
		deleteDeviceRegistrationToken,
		id,
		projectID,
	)
	return err
}

func (s *Store) SetDeviceRegistrationTokenLabel(ctx context.Context, deviceRegistrationTokenID, projectID, key, value string) (*string, error) {
	deviceRegistrationToken, err := s.GetDeviceRegistrationToken(ctx, deviceRegistrationTokenID, projectID)
	if err != nil {
		return nil, err
	}

	deviceRegistrationToken.Labels[key] = value

	labelsString, err := json.Marshal(deviceRegistrationToken.Labels)
	if err != nil {
		return nil, err
	}

	if _, err := s.db.ExecContext(
		ctx,
		updateDeviceRegistrationTokenLabels,
		labelsString,
		deviceRegistrationTokenID,
		projectID,
	); err != nil {
		return nil, err
	}

	deviceRegistrationToken, err = s.GetDeviceRegistrationToken(ctx, deviceRegistrationTokenID, projectID)
	if err != nil {
		return nil, err
	}
	v := deviceRegistrationToken.Labels[key]
	return &v, nil
}

func (s *Store) DeleteDeviceRegistrationTokenLabel(ctx context.Context, deviceRegistrationTokenID, projectID, key string) error {
	deviceRegistrationToken, err := s.GetDeviceRegistrationToken(ctx, deviceRegistrationTokenID, projectID)
	if err != nil {
		return err
	}

	delete(deviceRegistrationToken.Labels, key)

	labelsString, err := json.Marshal(deviceRegistrationToken.Labels)
	if err != nil {
		return err
	}

	if _, err := s.db.ExecContext(
		ctx,
		updateDeviceRegistrationTokenLabels,
		labelsString,
		deviceRegistrationTokenID,
		projectID,
	); err != nil {
		return err
	}
	return nil
}

func (s *Store) scanDeviceRegistrationToken(scanner scanner) (*models.DeviceRegistrationToken, error) {
	var deviceRegistrationToken models.DeviceRegistrationToken
	var labelsString string
	if err := scanner.Scan(
		&deviceRegistrationToken.ID,
		&deviceRegistrationToken.CreatedAt,
		&deviceRegistrationToken.ProjectID,
		&deviceRegistrationToken.MaxRegistrations,
		&deviceRegistrationToken.Name,
		&deviceRegistrationToken.Description,
		&labelsString,
	); err != nil {
		return nil, err
	}

	if labelsString == "" {
		deviceRegistrationToken.Labels = map[string]string{}
	} else {
		if err := json.Unmarshal([]byte(labelsString), &deviceRegistrationToken.Labels); err != nil {
			return nil, err
		}
	}

	return &deviceRegistrationToken, nil
}

func (s *Store) GetDevicesRegisteredWithTokenCount(ctx context.Context, tokenID, projectID string) (*models.DevicesRegisteredWithTokenCount, error) {
	countRow := s.db.QueryRowContext(ctx, getDevicesRegisteredWithTokenCount, tokenID, projectID)

	count, err := s.scanDevicesRegisteredCountRow(countRow)
	if err == sql.ErrNoRows {
		return nil, store.ErrDeviceRegistrationTokenNotFound
	} else if err != nil {
		return nil, err
	}

	return &models.DevicesRegisteredWithTokenCount{
		AllCount: count,
	}, nil
}

func (s *Store) scanDevicesRegisteredCountRow(scanner scanner) (int, error) {
	var count int
	if err := scanner.Scan(
		&count,
	); err != nil {
		return 0, err
	}
	return count, nil
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
		&deviceAccessKey.CreatedAt,
		&deviceAccessKey.ProjectID,
		&deviceAccessKey.DeviceID,
	); err != nil {
		return nil, err
	}
	return &deviceAccessKey, nil
}

func (s *Store) CreateApplication(ctx context.Context, projectID, name, description string) (*models.Application, error) {
	id := newApplicationID()

	if _, err := s.db.ExecContext(
		ctx,
		createApplication,
		id,
		projectID,
		name,
		description,
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

func (s *Store) UpdateApplicationName(ctx context.Context, id, projectID, name string) (*models.Application, error) {
	if _, err := s.db.ExecContext(
		ctx,
		updateApplicationName,
		name,
		id,
		projectID,
	); err != nil {
		return nil, err
	}

	return s.GetApplication(ctx, id, projectID)
}

func (s *Store) UpdateApplicationDescription(ctx context.Context, id, projectID, description string) (*models.Application, error) {
	if _, err := s.db.ExecContext(
		ctx,
		updateApplicationDescription,
		description,
		id,
		projectID,
	); err != nil {
		return nil, err
	}

	return s.GetApplication(ctx, id, projectID)
}

func (s *Store) UpdateApplicationSchedulingRule(ctx context.Context, id, projectID string, schedulingRule models.SchedulingRule) (*models.Application, error) {
	schedulingRuleBytes, err := json.Marshal(schedulingRule)
	if err != nil {
		return nil, err
	}

	if _, err := s.db.ExecContext(
		ctx,
		updateApplicationSchedulingRule,
		string(schedulingRuleBytes),
		id,
		projectID,
	); err != nil {
		return nil, err
	}

	return s.GetApplication(ctx, id, projectID)
}

func (s *Store) UpdateApplicationMetricEndpointConfigs(ctx context.Context, id, projectID string, metricEndpointConfigs map[string]models.MetricEndpointConfig) (*models.Application, error) {
	metricEndpointConfigsBytes, err := json.Marshal(metricEndpointConfigs)
	if err != nil {
		return nil, err
	}

	if _, err := s.db.ExecContext(
		ctx,
		updateApplicationMetricEndpointConfigs,
		string(metricEndpointConfigsBytes),
		id,
		projectID,
	); err != nil {
		return nil, err
	}

	return s.GetApplication(ctx, id, projectID)
}

func (s *Store) DeleteApplication(ctx context.Context, id, projectID string) error {
	_, err := s.db.ExecContext(
		ctx,
		deleteApplication,
		id,
		projectID,
	)
	return err
}

func (s *Store) scanApplication(scanner scanner) (*models.Application, error) {
	var schedulingRuleStr string
	var metricEndpointConfigsStr string

	var application models.Application
	if err := scanner.Scan(
		&application.ID,
		&application.CreatedAt,
		&application.ProjectID,
		&application.Name,
		&application.Description,
		&schedulingRuleStr,
		&metricEndpointConfigsStr,
	); err != nil {
		return nil, err
	}

	if schedulingRuleStr == "" {
		application.SchedulingRule = models.SchedulingRule{
			ScheduleType:     models.ScheduleTypeNoDevices,
			DefaultReleaseID: models.LatestRelease,
			ConditionalQuery: nil,
			ReleaseSelectors: []models.ReleaseSelector{},
		}
	} else {
		err := json.Unmarshal([]byte(schedulingRuleStr), &application.SchedulingRule)
		if err != nil {
			return nil, err
		}
	}
	if metricEndpointConfigsStr == "" {
		application.MetricEndpointConfigs = make(map[string]models.MetricEndpointConfig)
	} else {
		err := json.Unmarshal([]byte(metricEndpointConfigsStr), &application.MetricEndpointConfigs)
		if err != nil {
			return nil, err
		}
	}

	return &application, nil
}

func (s *Store) GetApplicationDeviceCounts(ctx context.Context, projectID, applicationID string) (*models.ApplicationDeviceCounts, error) {
	countRow := s.db.QueryRowContext(ctx, getApplicationDeviceCounts, projectID, applicationID)

	count, err := s.scanApplicationDeviceCountRow(countRow)
	if err != nil {
		return nil, err
	}

	return &models.ApplicationDeviceCounts{
		AllCount: count,
	}, nil
}

func (s *Store) scanApplicationDeviceCountRow(scanner scanner) (int, error) {
	var count int
	if err := scanner.Scan(
		&count,
	); err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Store) CreateRelease(ctx context.Context, projectID, applicationID, yamlConfig, jsonConfig, createdByUserID, createdByServiceAccountID string) (*models.Release, error) {
	id := newReleaseID()

	var createdByUserIDNullable *string
	if createdByUserID != "" {
		createdByUserIDNullable = &createdByUserID
	}
	var createdByServiceAccountIDNullable *string
	if createdByServiceAccountID != "" {
		createdByServiceAccountIDNullable = &createdByServiceAccountID
	}

	lastRelease, err := s.GetLatestRelease(ctx, projectID, applicationID)
	var number uint32
	if err != nil && err != store.ErrReleaseNotFound {
		return nil, err
	}
	if err == store.ErrReleaseNotFound {
		number = 1
	} else {
		number = lastRelease.Number + 1
	}

	if _, err := s.db.ExecContext(
		ctx,
		createRelease,
		id,
		number,
		projectID,
		applicationID,
		jsonConfig,
		yamlConfig,
		createdByUserIDNullable,
		createdByServiceAccountIDNullable,
	); err != nil {
		return nil, err
	}

	return s.GetRelease(ctx, id, projectID, applicationID)
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

func (s *Store) GetReleaseByNumber(ctx context.Context, number uint32, projectID, applicationID string) (*models.Release, error) {
	applicationRow := s.db.QueryRowContext(ctx, getReleaseByNumber, number, projectID, applicationID)

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
	var jsonConfig string
	var release models.Release
	if err := scanner.Scan(
		&release.ID,
		&release.Number,
		&release.CreatedAt,
		&release.ProjectID,
		&release.ApplicationID,
		&jsonConfig,
		&release.RawConfig,
		&release.CreatedByUserID,
		&release.CreatedByServiceAccountID,
	); err != nil {
		return nil, err
	}

	var parsedConfig map[string]models.Service
	if len(jsonConfig) == 0 {
		if err := yaml.UnmarshalStrict([]byte(release.RawConfig), &parsedConfig); err != nil {
			return nil, err
		}
	} else {
		if err := json.Unmarshal([]byte(jsonConfig), &parsedConfig); err != nil {
			return nil, err
		}
	}
	release.Config = parsedConfig

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

func (s *Store) ListDeviceApplicationStatuses(ctx context.Context, projectID, deviceID string) ([]models.DeviceApplicationStatus, error) {
	deviceApplicationStatusRows, err := s.db.QueryContext(ctx, listDeviceApplicationStatuses, projectID, deviceID)
	if err != nil {
		return nil, errors.Wrap(err, "query device application statuses")
	}
	defer deviceApplicationStatusRows.Close()

	deviceApplicationStatuses := make([]models.DeviceApplicationStatus, 0)
	for deviceApplicationStatusRows.Next() {
		deviceApplicationStatus, err := s.scanDeviceApplicationStatus(deviceApplicationStatusRows)
		if err != nil {
			return nil, err
		}
		deviceApplicationStatuses = append(deviceApplicationStatuses, *deviceApplicationStatus)
	}

	if err := deviceApplicationStatusRows.Err(); err != nil {
		return nil, err
	}

	return deviceApplicationStatuses, nil
}

func (s *Store) DeleteDeviceApplicationStatus(ctx context.Context, projectID, deviceID, applicationID string) error {
	_, err := s.db.ExecContext(
		ctx,
		deleteDeviceApplicationStatus,
		projectID,
		deviceID,
		applicationID,
	)
	return err
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

func (s *Store) ListDeviceServiceStatuses(ctx context.Context, projectID, deviceID string) ([]models.DeviceServiceStatus, error) {
	deviceServiceStatusRows, err := s.db.QueryContext(ctx, listDeviceServiceStatuses, projectID, deviceID)
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

func (s *Store) DeleteDeviceServiceStatus(ctx context.Context, projectID, deviceID, applicationID, service string) error {
	_, err := s.db.ExecContext(
		ctx,
		deleteDeviceServiceStatus,
		projectID,
		deviceID,
		applicationID,
		service,
	)
	return err
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

func (s *Store) scanProjectConfig(scanner scanner) (*models.ProjectConfig, error) {
	var projectConfig models.ProjectConfig
	if err := scanner.Scan(
		&projectConfig.ProjectID,
		&projectConfig.Key,
		&projectConfig.Value,
	); err != nil {
		return nil, err
	}

	return &projectConfig, nil
}

func (s *Store) scanDeviceMetricsConfig(scanner scanner) (*models.DeviceMetricsConfig, error) {
	pConfig, err := s.scanProjectConfig(scanner)
	if err != nil {
		return nil, err
	}

	var dmc models.DeviceMetricsConfig
	err = json.Unmarshal([]byte(pConfig.Value), &dmc)
	if err != nil {
		return nil, err
	}

	return &dmc, nil
}

func (s *Store) scanServiceMetricsConfig(scanner scanner) ([]models.ServiceMetricsConfig, error) {
	pConfig, err := s.scanProjectConfig(scanner)
	if err != nil {
		return nil, err
	}

	var smc []models.ServiceMetricsConfig
	err = json.Unmarshal([]byte(pConfig.Value), &smc)
	if err != nil {
		return nil, err
	}

	return smc, nil
}

func (s *Store) scanProjectMetricsConfig(scanner scanner) (*models.ProjectMetricsConfig, error) {
	pConfig, err := s.scanProjectConfig(scanner)
	if err != nil {
		return nil, err
	}

	var pmc models.ProjectMetricsConfig
	err = json.Unmarshal([]byte(pConfig.Value), &pmc)
	if err != nil {
		return nil, err
	}

	return &pmc, nil
}

func (s *Store) SetServiceMetricsConfigs(ctx context.Context, projectID string, value []models.ServiceMetricsConfig) error {
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(
		ctx,
		setProjectConfig,
		projectID,
		models.ServiceMetricsConfigKey,
		valueBytes,
	)
	return err
}

func (s *Store) GetServiceMetricsConfigs(ctx context.Context, projectID string) ([]models.ServiceMetricsConfig, error) {
	smcRow := s.db.QueryRowContext(
		ctx,
		getProjectConfig,
		projectID,
		models.ServiceMetricsConfigKey,
	)

	smc, err := s.scanServiceMetricsConfig(smcRow)
	if err == sql.ErrNoRows {
		return make([]models.ServiceMetricsConfig, 0), nil
	} else if err != nil {
		return nil, err
	}

	return smc, nil
}

func (s *Store) SetProjectMetricsConfig(ctx context.Context, projectID string, value models.ProjectMetricsConfig) error {
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(
		ctx,
		setProjectConfig,
		projectID,
		models.ProjectMetricsConfigKey,
		valueBytes,
	)
	return err
}

func (s *Store) GetProjectMetricsConfig(ctx context.Context, projectID string) (*models.ProjectMetricsConfig, error) {
	pmcRow := s.db.QueryRowContext(
		ctx,
		getProjectConfig,
		projectID,
		models.ProjectMetricsConfigKey,
	)

	pmc, err := s.scanProjectMetricsConfig(pmcRow)
	if err == sql.ErrNoRows {
		return &models.ProjectMetricsConfig{
			ExposedMetrics: make([]models.ExposedMetric, 0),
		}, nil
	} else if err != nil {
		return nil, err
	}

	return pmc, nil
}

func (s *Store) SetDeviceMetricsConfig(ctx context.Context, projectID string, value models.DeviceMetricsConfig) error {
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(
		ctx,
		setProjectConfig,
		projectID,
		models.DeviceMetricsConfigKey,
		valueBytes,
	)
	return err
}

func (s *Store) GetDeviceMetricsConfig(ctx context.Context, projectID string) (*models.DeviceMetricsConfig, error) {
	dmcRow := s.db.QueryRowContext(
		ctx,
		getProjectConfig,
		projectID,
		models.DeviceMetricsConfigKey,
	)

	dmc, err := s.scanDeviceMetricsConfig(dmcRow)
	if err == sql.ErrNoRows {
		return &models.DeviceMetricsConfig{
			ExposedMetrics: make([]models.ExposedMetric, 0),
		}, nil
	} else if err != nil {
		return nil, err
	}

	return dmc, nil
}
