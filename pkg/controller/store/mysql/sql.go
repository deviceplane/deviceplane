package mysql

import "fmt"

const createUser = `
  insert into users (
    id,
    email,
    password_hash,
    first_name,
    last_name
  )
  values (?, ?, ?, ?, ?)
`

const getUser = `
  select id, email, first_name, last_name, registration_completed from users
  where id = ?
`

const validateUser = `
  select id, email, first_name, last_name, registration_completed from users
  where email = ? and password_hash = ?
`

const markRegistrationComplete = `
  update users
  set registration_completed = true
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
  select id, user_id from registration_tokens
  where id = ?
`

const validateRegistrationToken = `
  select id, user_id from registration_tokens
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
  select id, user_id from sessions
  where id = ?
`

const validateSession = `
  select id, user_id from sessions
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
  select id, name from projects
  where id = ?
`

const lookupProject = `
  select id, name from projects
  where name = ?
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
  select id, project_id, name, description, config from roles
  where id = ? and project_id = ?
`

const lookupRole = `
  select id, project_id, name, description, config from roles
  where name = ? and project_id = ?
`

const listRoles = `
  select id, project_id, name, description, config from roles
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
  select user_id, project_id from memberships
  where user_id = ? and project_id = ?
`

const listMembershipsByUser = `
  select user_id, project_id from memberships
  where user_id = ?
`

const listMembershipsByProject = `
  select user_id, project_id from memberships
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
  select user_id, role_id, project_id from membership_role_bindings
  where user_id = ? and role_id = ? and project_id = ?
`

const listMembershipRoleBindings = `
  select user_id, role_id, project_id from membership_role_bindings
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
  select id, project_id, name, description from service_accounts
  where id = ? and project_id = ?
`

const lookupServiceAccount = `
  select id, project_id, name, description from service_accounts
  where name = ? and project_id = ?
`

const listServiceAccounts = `
  select id, project_id, name, description from service_accounts
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
  select service_account_id, role_id, project_id from service_account_role_bindings
  where service_account_id = ? and role_id = ? and project_id = ?
`

const listServiceAccountRoleBindings = `
  select service_account_id, role_id, project_id from service_account_role_bindings
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
    name
  )
  values (?, ?, ?)
`

const getDevice = `
  select id, project_id, name, info from devices
  where id = ? and project_id = ?
`

const listDevices = `
  select id, project_id, name, info from devices
  where project_id = ?
`

const setDeviceInfo = `
  update devices
  set info = ?
  where id = ? and project_id = ?
`

var setDeviceLabel = fmt.Sprintf(`
  insert into device_labels (
    %skey%s,
    device_id,
    project_id,
    value
  )
  values (?, ?, ?, ?)
  on duplicate key update value = ?
`, "`", "`")

var getDeviceLabel = fmt.Sprintf(`
  select %skey%s, device_id, project_id, value from device_labels
  where %skey%s = ? and device_id = ? and project_id = ?
`, "`", "`", "`", "`")

var listDeviceLabels = fmt.Sprintf(`
  select %skey%s, device_id, project_id, value from device_labels
  where device_id = ? and project_id = ?
`, "`", "`")

var deleteDeviceLabel = fmt.Sprintf(`
  delete from device_labels
  where %skey%s = ? and device_id = ? and project_id = ?
  limit 1
`, "`", "`")

const createDeviceRegistrationToken = `
  insert into device_registration_tokens (
    id,
    project_id
  )
  values (?, ?)
`

const getDeviceRegistrationToken = `
  select id, project_id, device_access_key_id from device_registration_tokens
  where id = ? and project_id = ?
`

const bindDeviceRegistrationToken = `
  update device_registration_tokens
  set device_access_key_id = ?
  where id = ? and project_id = ?
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
  select id, project_id, device_id from device_access_keys
  where id = ? and project_id = ?
`

const validateDeviceAccessKey = `
  select id, project_id, device_id from device_access_keys
  where project_id = ? and hash = ?
`

const createApplication = `
  insert into applications (
    id,
    project_id,
    name,
    description,
    settings
  )
  values (?, ?, ?, ?, ?)
`

const getApplication = `
  select id, project_id, name, description, settings from applications
  where id = ? and project_id = ?
`

const lookupApplication = `
  select id, project_id, name, description, settings from applications
  where name = ? and project_id = ?
`

const listApplications = `
  select id, project_id, name, description, settings from applications
  where project_id = ?
`

const updateApplication = `
  update applications
  set name = ?, description = ?, settings = ?
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
    config
  )
  values (?, ?, ?, ?)
`

const getRelease = `
  select id, created_at, project_id, application_id, config from releases
  where id = ? and project_id = ? and application_id = ?
`

const getLatestRelease = `
  select id, created_at, project_id, application_id, config from releases
  where project_id = ? and application_id = ?
  order by created_at desc
  limit 1
`

// TODO: real pagination
const listReleases = `
  select id, created_at, project_id, application_id, config from releases
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
