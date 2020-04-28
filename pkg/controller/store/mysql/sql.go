package mysql

const initializeUser = `
  insert into users (
    id,
    internal_user_id,
    external_user_id
  )
  values (?, ?, ?)
`

// Index: primary key
const getUser = `
  select id, created_at, internal_user_id, external_user_id, name, super_admin from users
  where id = ?
`

// Index: internal_user_id_unique
const getUserByInternalID = `
  select id, created_at, internal_user_id, external_user_id, name, super_admin from users
  where internal_user_id = ?
`

// Index: external_user_id_unique
const getUserByExternalID = `
  select id, created_at, internal_user_id, external_user_id, name, super_admin from users
  where external_user_id = ?
`

// Index: primary key
const updateUserName = `
  update users
  set name = ?
  where id = ?
`

const createExternalUser = `
  insert into external_users (
    id,
    provider_name,
    provider_id,
    email,
    info
  )
  values (?, ?, ?, ?, ?)
`

// Index: primary key
const getExternalUser = `
  select id, provider_name, provider_id, email, info from external_users
  where id = ?
`

// Index: provider_name_provider_id_unique
const getExternalUserByProvider = `
  select id, provider_name, provider_id, email, info from external_users
  where provider_name = ? and provider_id = ?
`

const createInternalUser = `
  insert into internal_users (
    id,
    email,
    password_hash
  )
  values (?, ?, ?)
`

// Index: primary key
const getInternalUser = `
  select id, email from internal_users
  where id = ?
`

// Index: email
const lookupInternalUser = `
  select id, email from internal_users
  where email = ?
`

// Index: id_password_hash
const validateInternalUser = `
  select id, email from internal_users
  where id = ? and password_hash = ?
`

// Index: email_password_hash
const validateInternalUserWithEmail = `
  select id, email from internal_users
  where email = ? and password_hash = ?
`

// Index: primary key
const updateInternalUserPasswordHash = `
  update internal_users
  set password_hash = ?
  where id = ?
`

const createRegistrationToken = `
  insert into registration_tokens (
    id,
    internal_user_id,
    hash
  )
  values (?, ?, ?)
`

// Index: primary key
const getRegistrationToken = `
  select id, created_at, internal_user_id from registration_tokens
  where id = ?
`

// Index: hash
const validateRegistrationToken = `
  select id, created_at, internal_user_id from registration_tokens
  where hash = ?
`

const createPasswordRecoveryToken = `
  insert into password_recovery_tokens (
    id,
    expires_at,
    internal_user_id,
    hash
  )
  values (?, now() + interval 1 hour, ?, ?)
`

// Index: primary key
const getPasswordRecoveryToken = `
  select id, created_at, expires_at, internal_user_id from password_recovery_tokens
  where id = ?
`

// Index: hash
const validatePasswordRecoveryToken = `
  select id, created_at, expires_at, internal_user_id from password_recovery_tokens
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

// Index: primary key
const getSession = `
  select id, created_at, user_id from sessions
  where id = ?
`

// Index: hash
const validateSession = `
  select id, created_at, user_id from sessions
  where hash = ?
`

// Index: primary key
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

// Index: primary key
const getUserAccessKey = `
  select id, created_at, user_id, description from user_access_keys
  where id = ?
`

// Index: hash
const validateUserAccessKey = `
  select id, created_at, user_id, description from user_access_keys
  where hash = ?
`

// Index: user_id
const listUserAccessKeys = `
  select id, created_at, user_id, description from user_access_keys
  where user_id = ?
`

// Index: primary key
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

// Index: primary key
const getProject = `
  select id, created_at, name, datadog_api_key from projects
  where id = ?
`

// Index: name
const lookupProject = `
  select id, created_at, name, datadog_api_key from projects
  where name = ?
`

const listProjects = `
  select id, created_at, name, datadog_api_key from projects
`

// Index: primary key
const updateProject = `
  update projects
  set name = ?, datadog_api_key = ?
  where id = ?
`

// Index: primary key
const deleteProject = `
  delete from projects
  where id = ?
  limit 1
`

// Index: project_id
const getProjectDeviceCounts = `
  select count(*) from devices
  where project_id = ?
`

// Index: project_id
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

// Index: project_id_id
const getRole = `
  select id, created_at, project_id, name, description, config from roles
  where id = ? and project_id = ?
`

// Index: project_id_name
const lookupRole = `
  select id, created_at, project_id, name, description, config from roles
  where name = ? and project_id = ?
`

// Index: project_id_id
const listRoles = `
  select id, created_at, project_id, name, description, config from roles
  where project_id = ?
`

// Index: project_id_id
const updateRole = `
  update roles
  set name = ?, description = ?, config = ?
  where id = ? and project_id = ?
`

// Index: project_id_id
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

// Index: primary key
const getMembership = `
  select user_id, project_id, created_at from memberships
  where user_id = ? and project_id = ?
`

// Index: user_id
const listMembershipsByUser = `
  select user_id, project_id, created_at from memberships
  where user_id = ?
`

// Index: project_id
const listMembershipsByProject = `
  select user_id, project_id, created_at from memberships
  where project_id = ?
`

// Index: primary key
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

// Index: project_id_user_id_role_id
const getMembershipRoleBinding = `
  select user_id, role_id, created_at, project_id from membership_role_bindings
  where user_id = ? and role_id = ? and project_id = ?
`

// Index: project_id_user_id_role_id
const listMembershipRoleBindings = `
  select user_id, role_id, created_at, project_id from membership_role_bindings
  where user_id = ? and project_id = ?
`

// Index: project_id_user_id_role_id
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

// Index: project_id_id
const getServiceAccount = `
  select id, created_at, project_id, name, description from service_accounts
  where id = ? and project_id = ?
`

// Index: project_id_name
const lookupServiceAccount = `
  select id, created_at, project_id, name, description from service_accounts
  where name = ? and project_id = ?
`

// Index: project_id_id
const listServiceAccounts = `
  select id, created_at, project_id, name, description from service_accounts
  where project_id = ?
`

// Index: project_id_id
const updateServiceAccount = `
  update service_accounts
  set name = ?, description = ?
  where id = ? and project_id = ?
`

// Index: project_id_id
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

// Index: project_id_id
const getServiceAccountAccessKey = `
  select id, created_at, project_id, service_account_id, description from service_account_access_keys
  where id = ? and project_id = ?
`

// Index: hash
const validateServiceAccountAccessKey = `
  select id, created_at, project_id, service_account_id, description from service_account_access_keys
  where hash = ?
`

// Index: project_id_service_account_id_id
const listServiceAccountAccessKeys = `
  select id, created_at, project_id, service_account_id, description from service_account_access_keys
  where project_id = ? and service_account_id = ?
`

// Index: project_id_id
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

// Index: project_id_service_account_id_role_id
const getServiceAccountRoleBinding = `
  select service_account_id, role_id, created_at, project_id from service_account_role_bindings
  where service_account_id = ? and role_id = ? and project_id = ?
`

// Index: project_id_service_account_id_role_id
const listServiceAccountRoleBindings = `
  select service_account_id, role_id, created_at, project_id from service_account_role_bindings
  where service_account_id = ? and project_id = ?
`

// Index: project_id_service_account_id_role_id
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
    labels,
    environment_variables
  )
  values (?, ?, ?, ?, ?, ?)
`

// Index: project_id_id
const getDevice = `
  select id, created_at, project_id, name, registration_token_id, desired_agent_version, info, labels, environment_variables, last_seen_at from devices
  where id = ? and project_id = ?
`

// Index: project_id_name
const lookupDevice = `
  select id, created_at, project_id, name, registration_token_id, desired_agent_version, info, labels, environment_variables, last_seen_at from devices
  where name = ? and project_id = ?
`

// Index: project_id_id
const listDevices = `
  select id, created_at, project_id, name, registration_token_id, desired_agent_version, info, labels, environment_variables, last_seen_at from devices
  where project_id = ?
`

// Index: project_id_id,fulltext
const searchDevices = `
  select id, created_at, project_id, name, registration_token_id, desired_agent_version, info, labels, environment_variables, last_seen_at from devices
  where project_id = ?
  and match (name, labels) against (concat('*', ?, '*') in boolean mode)
`

// Index: project_id_id
const updateDeviceName = `
  update devices
  set name = ?
  where id = ? and project_id = ?
`

// Index: project_id_id
const updateDeviceLabels = `
  update devices
  set labels = ?
  where id = ? and project_id = ?
`

// Index: project_id_id
const updateDeviceEnvironmentVariables = `
  update devices
  set environment_variables = ?
  where id = ? and project_id = ?
`

const listAllDeviceLabels = `
  select labels from devices
  where project_id = ?
`

// Index: project_id_id
const setDeviceInfo = `
  update devices
  set info = ?
  where id = ? and project_id = ?
`

// Index: project_id_id
const updateDeviceLastSeenAt = `
  update devices
  set last_seen_at = current_timestamp
  where id = ? and project_id = ?
`

// Index: project_id_id
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
    labels,
    environment_variables
  )
  values (?, ?, ?, ?, ?, '{}', '{}')
`

// Index: project_id_id
const getDeviceRegistrationToken = `
  select id, created_at, project_id, max_registrations, name, description, labels, environment_variables from device_registration_tokens
  where id = ? and project_id = ?
`

// Index: project_id_name
const lookupDeviceRegistrationToken = `
  select id, created_at, project_id, max_registrations, name, description, labels, environment_variables from device_registration_tokens
  where name = ? and project_id = ?
`

// Index: project_id_id
const listDeviceRegistrationTokens = `
  select id, created_at, project_id, max_registrations, name, description, labels, environment_variables from device_registration_tokens
  where project_id = ?
`

// Index: project_id_id
const updateDeviceRegistrationToken = `
  update device_registration_tokens
  set name = ?, description = ?, max_registrations = ?
  where id = ? and project_id = ?
`

// Index: project_id_id
const updateDeviceRegistrationTokenLabels = `
  update device_registration_tokens
  set labels = ?
  where id = ? and project_id = ?
`

// Index: project_id_id
const updateDeviceRegistrationTokenEnvironmentVariables = `
  update device_registration_tokens
  set environment_variables = ?
  where id = ? and project_id = ?
`

// Index: project_id_id
const deleteDeviceRegistrationToken = `
  delete from device_registration_tokens
  where id = ? and project_id = ?
  limit 1
`

// Index: project_id_registration_token_id
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

// Index: project_id_id
const getDeviceAccessKey = `
  select id, created_at, project_id, device_id from device_access_keys
  where id = ? and project_id = ?
`

// Index: project_id_hash
const validateDeviceAccessKey = `
  select id, created_at, project_id, device_id from device_access_keys
  where project_id = ? and hash = ?
`

const createConnection = `
  insert into connections (
    id,
    project_id,
    name,
    protocol,
    port
  )
  values (?, ?, ?, ?, ?)
`

// Index: project_id_id
const getConnection = `
  select id, created_at, project_id, name, protocol, port from connections
  where id = ? and project_id = ?
`

// Index: project_id_name
const lookupConnection = `
  select id, created_at, project_id, name, protocol, port from connections
  where name = ? and project_id = ?
`

// Index: project_id_id
const listConnections = `
  select id, created_at, project_id, name, protocol, port from connections
  where project_id = ?
`

// Index: project_id_id
const updateConnection = `
  update connections
  set name = ?, protocol = ?, port = ?
  where id = ? and project_id = ?
`

// Index: project_id_id
const deleteConnection = `
  delete from connections
  where id = ? and project_id = ?
  limit 1
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

// Index: project_id_id
const getApplication = `
  select id, created_at, project_id, name, description, scheduling_rule, metric_endpoint_configs from applications
  where id = ? and project_id = ?
`

// Index: project_id_name
const lookupApplication = `
  select id, created_at, project_id, name, description, scheduling_rule, metric_endpoint_configs from applications
  where name = ? and project_id = ?
`

// Index: project_id_id
const listApplications = `
  select id, created_at, project_id, name, description, scheduling_rule, metric_endpoint_configs from applications
  where project_id = ?
`

// Index: project_id_id
const updateApplicationName = `
  update applications
  set name = ?
  where id = ? and project_id = ?
`

// Index: project_id_id
const updateApplicationDescription = `
  update applications
  set description = ?
  where id = ? and project_id = ?
`

// Index: project_id_id
const updateApplicationSchedulingRule = `
  update applications
  set scheduling_rule = ?
  where id = ? and project_id = ?
`

// Index: project_id_id
const updateApplicationMetricEndpointConfigs = `
  update applications
  set metric_endpoint_configs = ?
  where id = ? and project_id = ?
`

// Index: project_id_id
const deleteApplication = `
  delete from applications
  where id = ? and project_id = ?
  limit 1
`

// Index: project_id_application_id_current_release_id
const getApplicationDeviceCounts = `
  select count(*) from device_application_statuses
  where project_id = ? and application_id = ?
`

const createRelease = `
  insert into releases (
    id,
    ` + "`number`" + `,
    project_id,
    application_id,
    config,
    raw_config,
    created_by_user_id,
    created_by_service_account_id
  )
  values (?, ?, ?, ?, ?, ?, ?, ?)
`

// Index: project_id_application_id_id
const getRelease = `
  select id, ` + "`number`" + `, created_at, project_id, application_id, config, raw_config, created_by_user_id, created_by_service_account_id from releases
  where id = ? and project_id = ? and application_id = ?
`

// Index: project_id_application_id_number
const getReleaseByNumber = `
  select id, ` + "`number`" + `, created_at, project_id, application_id, config, raw_config, created_by_user_id, created_by_service_account_id from releases
  where ` + "`number`" + ` = ? and project_id = ? and application_id = ?
`

// Index: project_id_application_id_created_at
const getLatestRelease = `
  select id, ` + "`number`" + `, created_at, project_id, application_id, config, raw_config, created_by_user_id, created_by_service_account_id from releases
  where project_id = ? and application_id = ?
  order by created_at desc
  limit 1
`

// TODO: real pagination
// Index: project_id_application_id_created_at
const listReleases = `
  select id, ` + "`number`" + `, created_at, project_id, application_id, config, raw_config, created_by_user_id, created_by_service_account_id from releases
  where project_id = ? and application_id = ?
  order by created_at desc
  limit 50
`

// Index: project_id_application_id_current_release_id
const getReleaseDeviceCounts = `
  select count(*) from device_application_statuses
  where project_id = ? and application_id = ? and current_release_id = ?
`

// Index: primary key
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

// Index: primary key
const getDeviceApplicationStatus = `
  select project_id, device_id, application_id, current_release_id from device_application_statuses
  where project_id = ? and device_id = ? and application_id = ?
`

// Index: project_id_device_id
const listDeviceApplicationStatuses = `
  select project_id, device_id, application_id, current_release_id from device_application_statuses
  where project_id = ? and device_id = ?
`

// Index: project_id_device_id
const listAllDeviceApplicationStatuses = `
  select project_id, device_id, application_id, current_release_id from device_application_statuses
  where project_id = ?
`

// Index: primary key
const deleteDeviceApplicationStatus = `
  delete from device_application_statuses
  where project_id = ? and device_id = ? and application_id = ?
`

// Index: primary key
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

// Index: primary key
const getDeviceServiceStatus = `
  select project_id, device_id, application_id, service, current_release_id from device_service_statuses
  where project_id = ? and device_id = ? and application_id = ? and service = ?
`

// Index: project_id_device_id_application_id
const getDeviceServiceStatuses = `
  select project_id, device_id, application_id, service, current_release_id from device_service_statuses
  where project_id = ? and device_id = ? and application_id = ?
`

// Index: project_id_device_id_application_id
const listDeviceServiceStatuses = `
  select project_id, device_id, application_id, service, current_release_id from device_service_statuses
  where project_id = ? and device_id = ?
`

// Index: primary key
const deleteDeviceServiceStatus = `
  delete from device_service_statuses
  where project_id = ? and device_id = ? and application_id = ? and service = ?
`

// Index: primary key
const setDeviceServiceState = `
  insert into device_service_states (
    project_id,
    device_id,
    application_id,
    service,
    state,
    error_message
  )
  values (?, ?, ?, ?, ?, ?)
  on duplicate key update
    state = ?,
    error_message = ?
`

// Index: primary key
const getDeviceServiceState = `
  select project_id, device_id, application_id, service, state, error_message from device_service_states
  where project_id = ? and device_id = ? and application_id = ? and service = ?
`

// Index: project_id_device_id_application_id
const getDeviceServiceStates = `
  select project_id, device_id, application_id, service, state, error_message from device_service_states
  where project_id = ? and device_id = ? and application_id = ?
`

// Index: project_id_application_id
const getDeviceServiceStateCountsByApplication = `
  select count(state) as count, sum(error_message != "") as count_erroring, state, service, application_id
  from device_service_states
  where project_id = ? and application_id = ?
  group by service, state
  order by service, state
`

// Index: project_id_device_id_application_id
const listDeviceServiceStates = `
  select project_id, device_id, application_id, service, state, error_message from device_service_states
  where project_id = ? and device_id = ?
`

// Index: project_id_device_id_application_id
const listAllDeviceServiceStates = `
  select project_id, device_id, application_id, service, state, error_message from device_service_states
  where project_id = ?
`

// Index: primary key
const deleteDeviceServiceState = `
  delete from device_service_states
  where project_id = ? and device_id = ? and application_id = ? and service = ?
`

// Index: primary key
const setProjectConfig = `
  replace into project_configs (
    project_id,
    k,
    v
  )
  values (?, ?, ?)
`

// Index: primary key
const getProjectConfig = `
  select project_id, k, v from project_configs
  where project_id = ? and k = ?
`
