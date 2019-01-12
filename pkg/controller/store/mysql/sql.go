package mysql

const createUser = `
  insert into users (
    id
  )
  values (?)
`

const getUser = `
  select * from users
  where id = ?
`

const createProject = `
  insert into projects (
    id
  )
  values (?)
`

const getProject = `
  select * from projects
  where id = ?
`

const createDevice = `
  insert into devices (
    id,
    project_id
  )
  values (?, ?)
`

const getDevice = `
  select * from devices
  where id = ?
`

const createApplication = `
  insert into applications (
    id,
    project_id
  )
  values (?, ?)
`

const getApplication = `
  select * from applications
  where id = ?
`

const listApplications = `
  select * from applications
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
  select * from releases
  where id = ?
`

const getLatestRelease = `
  select * from releases
  where application_id = ?
  order by created_at desc
  limit 1
`
