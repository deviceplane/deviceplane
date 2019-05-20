
use deviceplane;
begin;

--
-- Users
--

create table if not exists users (
  id varchar(32) not null,
  created_at timestamp not null default current_timestamp,

  email varchar(255) not null,
  password_hash varchar(255) not null,
  first_name varchar(100) not null,
  last_name varchar(100) not null,
  registration_completed boolean not null default false,

  primary key (id),
  unique(email)
);

--
-- RegistrationTokens
--

create table if not exists registration_tokens (
  id varchar(32) not null,
  created_at timestamp not null default current_timestamp,
  user_id varchar(32) not null,

  hash varchar(255) not null,

  primary key (id),
  unique(user_id),
  foreign key registration_tokens_user_id(user_id)
  references users(id)
  on delete cascade
);


--
-- Sessions
--

create table if not exists sessions (
  id varchar(32) not null,
  created_at timestamp not null default current_timestamp,
  user_id varchar(32) not null,

  hash varchar(255) not null,

  primary key (id),
  foreign key sessions_user_id(user_id)
  references users(id)
  on delete cascade
);

--
-- AccessKeys
--

create table if not exists access_keys (
  id varchar(32) not null,
  created_at timestamp not null default current_timestamp,
  user_id varchar(32) not null,

  hash varchar(255) not null,

  primary key (id),
  foreign key access_keys_user_id(user_id)
  references users(id)
  on delete cascade
);

--
-- Projects
--

create table if not exists projects (
  id varchar(32) not null,
  created_at timestamp not null default current_timestamp,

  name varchar(100) not null,

  primary key (id)
);

--
-- Roles
--

create table if not exists roles (
  id varchar(32) not null,
  created_at timestamp not null default current_timestamp,
  project_id varchar(32) not null,

  config longtext not null,

  primary key (id)
);

--
-- Memberships
--

create table if not exists memberships (
  id varchar(32) not null,
  user_id varchar(32) not null,
  project_id varchar(32) not null,
  created_at timestamp not null default current_timestamp,

  primary key (id),
  unique(user_id, project_id),
  foreign key memberships_user_id(user_id)
  references users(id)
  on delete cascade,
  foreign key memberships_project_id(project_id)
  references projects(id)
  on delete cascade
);

--
-- MembershipRoleBindings
--

create table if not exists membership_role_bindings (
  membership_id varchar(32) not null,
  role_id varchar(32) not null,
  created_at timestamp not null default current_timestamp,
  project_id varchar(32) not null,

  primary key (membership_id, role_id),
  foreign key membership_role_bindings_membership_id(membership_id)
  references memberships(id)
  on delete cascade,
  foreign key membership_role_bindings_role_id(role_id)
  references roles(id)
  on delete cascade,
  foreign key membership_role_bindings_project_id(project_id)
  references projects(id)
  on delete cascade
);

--
-- Devices
--

create table if not exists devices (
  id varchar(32) not null,
  created_at timestamp not null default current_timestamp,
  project_id varchar(32) not null,

  name varchar(100) not null,
  info longtext not null,

  primary key (id),
  foreign key devices_project_id(project_id)
  references projects(id)
  on delete cascade
);

create table if not exists device_labels (
  `key` varchar(100) not null,
  device_id varchar(32) not null,
  created_at timestamp not null default current_timestamp,
  project_id varchar(32) not null,

  value varchar(100) not null,

  primary key (`key`, device_id),
  foreign key device_labels_device_id(device_id)
  references devices(id)
  on delete cascade
);

--
-- DeviceRegistrationTokens
--

create table if not exists device_registration_tokens (
  id varchar(32) not null,
  created_at timestamp not null default current_timestamp,
  project_id varchar(32) not null,
  device_access_key_id varchar(32) default null,

  primary key (id)
);

--
-- DeviceAccessKeys
--

create table if not exists device_access_keys (
  id varchar(32) not null,
  created_at timestamp not null default current_timestamp,
  project_id varchar(32) not null,
  device_id varchar(32) not null,

  hash varchar(255) not null,

  primary key (id),
  foreign key device_access_keys_device_id(device_id)
  references devices(id)
  on delete cascade
);

--
-- Applications
--

create table if not exists applications (
  id varchar(32) not null,
  created_at timestamp not null default current_timestamp, 
  project_id varchar(32) not null,

  name varchar(100) not null,
  settings longtext not null,

  primary key (id),
  foreign key applications_project_id(project_id)
  references projects(id)
  on delete cascade
);

--
-- Releases
--

create table if not exists releases (
  id varchar(32) not null,
  created_at timestamp not null default current_timestamp, 
  project_id varchar(32) not null,
  application_id varchar(32) not null,
  config longtext not null,

  primary key (id),
  foreign key releases_application_id(application_id)
  references applications(id)
  on delete cascade
);

--
-- DeviceApplicationStatuses
--

create table if not exists device_application_statuses (
  project_id varchar(32) not null,
  device_id varchar(32) not null,
  application_id varchar(32) not null,

  current_release_id varchar(32) not null,

  primary key (project_id, device_id, application_id),
  foreign key device_application_statuses_project_id(project_id)
  references projects(id)
  on delete cascade,
  foreign key device_application_statuses_device_id(device_id)
  references devices(id)
  on delete cascade,
  foreign key device_application_statuses_application_id(application_id)
  references applications(id)
  on delete cascade
);

--
-- DeviceServiceStatuses
--

create table if not exists device_service_statuses (
  project_id varchar(32) not null,
  device_id varchar(32) not null,
  application_id varchar(32) not null,
  service varchar(100) not null,

  current_release_id varchar(32) not null,

  primary key (project_id, device_id, application_id, service),
  foreign key device_service_statuses_project_id(project_id)
  references projects(id)
  on delete cascade,
  foreign key device_service_statuses_device_id(device_id)
  references devices(id)
  on delete cascade,
  foreign key device_service_statuses_application_id(application_id)
  references applications(id)
  on delete cascade
);

--
-- Commit
--

commit;

