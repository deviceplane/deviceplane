package mysql

const createUser = `
  insert into users (
    id,
    email,
    password_hash,
    first_name,
    last_name,
    company
  )
  values (?, ?, ?, ?, ?, ?)
`

const getUser = `
  select id, created_at, email, first_name, last_name, company, registration_completed, super_admin from users
  where id = ?
`

const lookupUser = `
  select id, created_at, email, first_name, last_name, company, registration_completed, super_admin from users
  where email = ?
`

const validateUser = `
  select id, created_at, email, first_name, last_name, company, registration_completed, super_admin from users
  where id = ? and password_hash = ?
`

const validateUserWithEmail = `
  select id, created_at, email, first_name, last_name, company, registration_completed, super_admin from users
  where email = ? and password_hash = ?
`

const markRegistrationComplete = `
  update users
  set registration_completed = true
  where id = ?
`

const updatePasswordHash = `
  update users
  set password_hash = ?
  where id = ?
`

const updateFirstName = `
  update users
  set first_name = ?
  where id = ?
`

const updateLastName = `
  update users
  set last_name = ?
  where id = ?
`

const updateCompany = `
  update users
  set company = ?
  where id = ?
`

const createRegistrationToken = `
  insert into registration_tokens (
    id,
    user_id,
    hash
  )
  values (?, ?, ?)
`

const getRegistrationToken = `
  select id, created_at, user_id from registration_tokens
  where id = ?
`

const validateRegistrationToken = `
  select id, created_at, user_id from registration_tokens
  where hash = ?
`

const createPasswordRecoveryToken = `
  insert into password_recovery_tokens (
    id,
    expires_at,
    user_id,
    hash
  )
  values (?, now() + interval 1 hour, ?, ?)
`

const getPasswordRecoveryToken = `
  select id, created_at, expires_at, user_id from password_recovery_tokens
  where id = ?
`

const validatePasswordRecoveryToken = `
  select id, created_at, expires_at, user_id from password_recovery_tokens
  where hash = ?
`

const createSession = `
  insert into sessions (
    id,
    user_id,
    hash
  )
  values (?, ?, ?)
`

const getSession = `
  select id, created_at, user_id from sessions
  where id = ?
`

const validateSession = `
  select id, created_at, user_id from sessions
  where hash = ?
`

const deleteSession = `
  delete from sessions
  where id = ?
  limit 1
`

const createUserAccessKey = `
  insert into user_access_keys (
    id,
    user_id,
    hash,
    description
  )
  values (?, ?, ?, ?)
`

const getUserAccessKey = `
  select id, created_at, user_id, description from user_access_keys
  where id = ?
`

const validateUserAccessKey = `
  select id, created_at, user_id, description from user_access_keys
  where hash = ?
`

const listUserAccessKeys = `
  select id, created_at, user_id, description from user_access_keys
  where user_id = ?
`

const deleteUserAccessKey = `
  delete from user_access_keys
  where id = ?
  limit 1
`

const createProject = `
  insert into projects (
    id,
    name
  )
  values (?, ?)
`

const getProject = `
  select id, created_at, name, datadog_api_key from projects
  where id = ?
`

const lookupProject = `
  select id, created_at, name, datadog_api_key from projects
  where name = ?
`

const listProjects = `
  select id, created_at, name, datadog_api_key from projects
`

const updateProject = `
  update projects
  set name = ?, datadog_api_key = ?
  where id = ?
`

const deleteProject = `
  delete from projects
  where id = ?
  limit 1
`

const getProjectDeviceCounts = `
  select count(*) from devices
  where project_id = ?
`

const getProjectApplicationCounts = `
  select count(*) from applications
  where project_id = ?
`

const createRole = `
  insert into roles (
    id,
    project_id,
    name,
    description,
    config
  )
  values (?, ?, ?, ?, ?)
`

const getRole = `
  select id, created_at, project_id, name, description, config from roles
  where id = ? and project_id = ?
`

const lookupRole = `
  select id, created_at, project_id, name, description, config from roles
  where name = ? and project_id = ?
`

const listRoles = `
  select id, created_at, project_id, name, description, config from roles
  where project_id = ?
`

const updateRole = `
  update roles
  set name = ?, description = ?, config = ?
  where id = ? and project_id = ?
`

const deleteRole = `
  delete from roles
  where id = ? and project_id = ?
  limit 1
`

const createMembership = `
  insert into memberships (
    user_id,
    project_id
  )
  values (?, ?)
`

const getMembership = `
  select user_id, project_id, created_at from memberships
  where user_id = ? and project_id = ?
`

const listMembershipsByUser = `
  select user_id, project_id, created_at from memberships
  where user_id = ?
`

const listMembershipsByProject = `
  select user_id, project_id, created_at from memberships
  where project_id = ?
`

const deleteMembership = `
  delete from memberships
  where user_id = ? and project_id = ?
  limit 1
`

const createMembershipRoleBinding = `
  insert into membership_role_bindings (
    user_id,
    role_id,
    project_id
  )
  values (?, ?, ?)
`

const getMembershipRoleBinding = `
  select user_id, role_id, created_at, project_id from membership_role_bindings
  where user_id = ? and role_id = ? and project_id = ?
`

const listMembershipRoleBindings = `
  select user_id, role_id, created_at, project_id from membership_role_bindings
  where user_id = ? and project_id = ?
`

const deleteMembershipRoleBinding = `
  delete from membership_role_bindings
  where user_id = ? and role_id = ? and project_id = ?
  limit 1
`

const createServiceAccount = `
  insert into service_accounts (
    id,
    project_id,
    name,
    description
  )
  values (?, ?, ?, ?)
`

const getServiceAccount = `
  select id, created_at, project_id, name, description from service_accounts
  where id = ? and project_id = ?
`

const lookupServiceAccount = `
  select id, created_at, project_id, name, description from service_accounts
  where name = ? and project_id = ?
`

const listServiceAccounts = `
  select id, created_at, project_id, name, description from service_accounts
  where project_id = ?
`

const updateServiceAccount = `
  update service_accounts
  set name = ?, description = ?
  where id = ? and project_id = ?
`

const deleteServiceAccount = `
  delete from service_accounts
  where id = ? and project_id = ?
  limit 1
`

const createServiceAccountAccessKey = `
  insert into service_account_access_keys (
    id,
    project_id,
    service_account_id,
    hash,
    description
  )
  values (?, ?, ?, ?, ?)
`

const getServiceAccountAccessKey = `
  select id, created_at, project_id, service_account_id, description from service_account_access_keys
  where id = ? and project_id = ?
`

const validateServiceAccountAccessKey = `
  select id, created_at, project_id, service_account_id, description from service_account_access_keys
  where hash = ?
`

const listServiceAccountAccessKeys = `
  select id, created_at, project_id, service_account_id, description from service_account_access_keys
  where project_id = ? and service_account_id = ?
`

const deleteServiceAccountAccessKey = `
  delete from service_account_access_keys
  where id = ? and project_id = ?
  limit 1
`

const createServiceAccountRoleBinding = `
  insert into service_account_role_bindings (
    service_account_id,
    role_id,
    project_id
  )
  values (?, ?, ?)
`

const getServiceAccountRoleBinding = `
  select service_account_id, role_id, created_at, project_id from service_account_role_bindings
  where service_account_id = ? and role_id = ? and project_id = ?
`

const listServiceAccountRoleBindings = `
  select service_account_id, role_id, created_at, project_id from service_account_role_bindings
  where service_account_id = ? and project_id = ?
`

const deleteServiceAccountRoleBinding = `
  delete from service_account_role_bindings
  where service_account_id = ? and role_id = ? and project_id = ?
  limit 1
`

const createDevice = `
  insert into devices (
    id,
    project_id,
    name,
    registration_token_id,
    labels
  )
  values (?, ?, ?, ?, ?)
`

const getDevice = `
  select id, created_at, project_id, name, registration_token_id, desired_agent_spec, desired_agent_version, info, labels, last_seen_at from devices
  where id = ? and project_id = ?
`

const lookupDevice = `
  select id, created_at, project_id, name, registration_token_id, desired_agent_spec, desired_agent_version, info, labels, last_seen_at from devices
  where name = ? and project_id = ?
`

const listDevices = `
  select id, created_at, project_id, name, registration_token_id, desired_agent_spec, desired_agent_version, info, labels, last_seen_at from devices
  where project_id = ?
`

const updateDeviceName = `
  update devices
  set name = ?
  where id = ? and project_id = ?
`

const updateDeviceLabels = `
  update devices
  set labels = ?
  where id = ? and project_id = ?
`

const setDeviceInfo = `
  update devices
  set info = ?
  where id = ? and project_id = ?
`

const updateDeviceLastSeenAt = `
  update devices
  set last_seen_at = current_timestamp
  where id = ? and project_id = ?
`

const deleteDevice = `
  delete from devices
  where id = ? and project_id = ?
`

const createDeviceRegistrationToken = `
  insert into device_registration_tokens (
    id,
    project_id,
    name,
    description,
    max_registrations,
    labels
  )
  values (?, ?, ?, ?, ?, '{}')
`

const lookupDeviceRegistrationToken = `
  select id, created_at, project_id, max_registrations, name, description, labels from device_registration_tokens
  where name = ? and project_id = ?
`

const listDeviceRegistrationTokens = `
  select id, created_at, project_id, max_registrations, name, description, labels from device_registration_tokens
  where project_id = ?
`

const getDeviceRegistrationToken = `
  select id, created_at, project_id, max_registrations, name, description, labels from device_registration_tokens
  where id = ? and project_id = ?
`

const updateDeviceRegistrationToken = `
  update device_registration_tokens
  set name = ?, description = ?, max_registrations = ?
  where id = ? and project_id = ?
`

const updateDeviceRegistrationTokenLabels = `
  update device_registration_tokens
  set labels = ?
  where id = ? and project_id = ?
`

const deleteDeviceRegistrationToken = `
  delete from device_registration_tokens
  where id = ? and project_id = ?
  limit 1
`

const getDevicesRegisteredWithTokenCount = `
  select count(*) from devices
  where registration_token_id = ? and project_id = ?
`

const createDeviceAccessKey = `
  insert into device_access_keys (
    id,
    project_id,
    device_id,
    hash
  )
  values (?, ?, ?, ?)
`

const getDeviceAccessKey = `
  select id, created_at, project_id, device_id from device_access_keys
  where id = ? and project_id = ?
`

const validateDeviceAccessKey = `
  select id, created_at, project_id, device_id from device_access_keys
  where project_id = ? and hash = ?
`

const createApplication = `
  insert into applications (
    id,
    project_id,
    name,
    description
  )
  values (?, ?, ?, ?)
`

const getApplication = `
  select id, created_at, project_id, name, description, scheduling_rule from applications
  where id = ? and project_id = ?
`

const lookupApplication = `
  select id, created_at, project_id, name, description, scheduling_rule from applications
  where name = ? and project_id = ?
`

const listApplications = `
  select id, created_at, project_id, name, description, scheduling_rule from applications
  where project_id = ?
`

const updateApplicationName = `
  update applications
  set name = ?
  where id = ? and project_id = ?
`

const updateApplicationDescription = `
  update applications
  set description = ?
  where id = ? and project_id = ?
`

const updateApplicationSchedulingRule = `
  update applications
  set scheduling_rule = ?
  where id = ? and project_id = ?
`

const deleteApplication = `
  delete from applications
  where id = ? and project_id = ?
  limit 1
`

const getApplicationDeviceCounts = `
  select count(*) from device_application_statuses
  where project_id = ? and application_id = ?
`

const createRelease = `
  insert into releases (
    id,
    project_id,
    application_id,
    config,
    raw_config,
    created_by_user_id,
    created_by_service_account_id
  )
  values (?, ?, ?, ?, ?, ?, ?)
`

const getRelease = `
  select id, created_at, project_id, application_id, config, raw_config, created_by_user_id, created_by_service_account_id from releases
  where id = ? and project_id = ? and application_id = ?
`

const getLatestRelease = `
  select id, created_at, project_id, application_id, config, raw_config, created_by_user_id, created_by_service_account_id from releases
  where project_id = ? and application_id = ?
  order by created_at desc
  limit 1
`

// TODO: real pagination
const listReleases = `
  select id, created_at, project_id, application_id, config, raw_config, created_by_user_id, created_by_service_account_id from releases
  where project_id = ? and application_id = ?
  order by created_at desc
  limit 10
`

const getReleaseDeviceCounts = `
  select count(*) from device_application_statuses
  where project_id = ? and application_id = ? and current_release_id = ?
`

const setDeviceApplicationStatus = `
  insert into device_application_statuses (
    project_id,
    device_id,
    application_id,
    current_release_id
  )
  values (?, ?, ?, ?)
  on duplicate key update
    current_release_id = ?
`

const getDeviceApplicationStatus = `
  select project_id, device_id, application_id, current_release_id from device_application_statuses
  where project_id = ? and device_id = ? and application_id = ?
`

const listDeviceApplicationStatuses = `
  select project_id, device_id, application_id, current_release_id from device_application_statuses
  where project_id = ? and device_id = ?
`

const deleteDeviceApplicationStatus = `
  delete from device_application_statuses
  where project_id = ? and device_id = ? and application_id = ?
`

const setDeviceServiceStatus = `
  insert into device_service_statuses (
    project_id,
    device_id,
    application_id,
    service,
    current_release_id
  )
  values (?, ?, ?, ?, ?)
  on duplicate key update
    current_release_id = ?
`

const getDeviceServiceStatus = `
  select project_id, device_id, application_id, service, current_release_id from device_service_statuses
  where project_id = ? and device_id = ? and application_id = ? and service = ?
`

const getDeviceServiceStatuses = `
  select project_id, device_id, application_id, service, current_release_id from device_service_statuses
  where project_id = ? and device_id = ? and application_id = ?
`

const listDeviceServiceStatuses = `
  select project_id, device_id, application_id, service, current_release_id from device_service_statuses
  where project_id = ? and device_id = ?
`

const deleteDeviceServiceStatus = `
  delete from device_service_statuses
  where project_id = ? and device_id = ? and application_id = ? and service = ?
`

const createMetricTargetConfig = `
  insert into metric_target_configs (
    id,
    project_id,
    type,
    configs
  )
  values (?, ?, ?, ?)
`

const updateMetricTargetConfig = `
  update metric_target_configs
  set configs = ?
  where id = ? and project_id = ?
`

const getMetricTargetConfig = `
  select id, created_at, project_id, type, configs from metric_target_configs
  where id = ? and project_id = ?
`

const lookupMetricTargetConfig = `
  select id, created_at, project_id, type, configs from metric_target_configs
  where type = ? and project_id = ?
`
