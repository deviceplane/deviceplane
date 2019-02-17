package mysql

const createUser = `
  insert into users (
    id,
    email,
    password_hash
  )
  values (?, ?, ?)
`

const getUser = `
  select id, email, password_hash from users
  where id = ?
`

const createAccessKey = `
  insert into access_keys (
    id,
    user_id,
    hash
  )
  values (?, ?, ?)
`

const getAccessKey = `
  select id, user_id, hash from access_keys
  where id = ?
`

const validateAccessKey = `
  select id, user_id, hash from access_keys
  where hash = ?
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

const createMembership = `
  insert into memberships (
    user_id,
    project_id,
    level
  )
  values (?, ?, ?)
`

const getMembership = `
  select user_id, project_id, level from memberships
  where user_id = ? and project_id = ?
`

const listMembershipsByUser = `
  select user_id, project_id, level from memberships
  where user_id = ?
`

const listMembershipsByProject = `
  select user_id, project_id, level from memberships
  where project_id = ?
`

const createDevice = `
  insert into devices (
    id,
    project_id
  )
  values (?, ?)
`

const getDevice = `
  select id, project_id from devices
  where id = ?
`

const listDevices = `
  select id, project_id from devices
  where project_id = ?
`

const createApplication = `
  insert into applications (
    id,
    project_id,
    name
  )
  values (?, ?, ?)
`

const getApplication = `
  select id, project_id, name from applications
  where id = ?
`

const listApplications = `
  select id, project_id, name from applications
  where project_id = ?
`

const createRelease = `
  insert into releases (
    id,
    application_id,
    config
  )
  values (?, ?, ?)
`

const getRelease = `
  select id, application_id, config from releases
  where id = ?
`

const getLatestRelease = `
  select id, application_id, config from releases
  where application_id = ?
  order by created_at desc
  limit 1
`

const listReleases = `
  select id, application_id, config from releases
  where application_id = ?
`
