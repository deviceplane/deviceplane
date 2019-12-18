
use deviceplane;
begin;

--
-- Users
--

create table if not exists users (
  id varchar(32) not null,
  created_at timestamp not null default current_timestamp,

  email varchar(255) not null,
  -- SENSITIVE FIELD
  password_hash varchar(255) not null,
  first_name varchar(100) not null,
  last_name varchar(100) not null,
  company varchar(100) not null,
  registration_completed boolean not null default false,
  super_admin boolean not null default false,

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

  -- SENSITIVE FIELD
  hash varchar(255) not null,

  primary key (id),
  unique(user_id),
  unique(hash),
  foreign key registration_tokens_user_id(user_id)
  references users(id)
  on delete cascade
);

--
-- PasswordRecoveryTokens
--

create table if not exists password_recovery_tokens (
  id varchar(32) not null,
  created_at timestamp not null default current_timestamp,
  expires_at timestamp not null,
  user_id varchar(32) not null,

  -- SENSITIVE FIELD
  hash varchar(255) not null,

  primary key (id),
  unique(hash),
  foreign key password_recovery_tokens_user_id(user_id)
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

  -- SENSITIVE FIELD
  hash varchar(255) not null,

  primary key (id),
  unique(hash),
  foreign key sessions_user_id(user_id)
  references users(id)
  on delete cascade
);

--
-- UserAccessKeys
--

create table if not exists user_access_keys (
  id varchar(32) not null,
  created_at timestamp not null default current_timestamp,
  user_id varchar(32) not null,

  -- SENSITIVE FIELD
  hash varchar(255) not null,
  description longtext not null,

  primary key (id),
  unique(hash),
  foreign key user_access_keys_user_id(user_id)
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
  datadog_api_key varchar(100),

  primary key (id),
  unique(name)
);

--
-- Roles
--

create table if not exists roles (
  id varchar(32) not null,
  created_at timestamp not null default current_timestamp,
  project_id varchar(32) not null,

  name varchar(100) not null,
  description longtext not null,
  config longtext not null,

  primary key (id),
  unique(name, project_id),
  foreign key roles_project_id(project_id)
  references projects(id)
  on delete cascade
);

--
-- Memberships
--

create table if not exists memberships (
  user_id varchar(32) not null,
  project_id varchar(32) not null,
  created_at timestamp not null default current_timestamp,

  primary key (user_id, project_id),
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
  user_id varchar(32) not null,
  role_id varchar(32) not null,
  created_at timestamp not null default current_timestamp,
  project_id varchar(32) not null,

  primary key (user_id, project_id, role_id),
  foreign key membership_role_bindings_user_id_project_id(user_id, project_id)
  references memberships(user_id, project_id)
  on delete cascade,
  foreign key membership_role_bindings_role_id(role_id)
  references roles(id)
  on delete cascade,
  foreign key membership_role_bindings_project_id(project_id)
  references projects(id)
  on delete cascade
);

--
-- ServiceAccounts
--

create table if not exists service_accounts (
  id varchar(32) not null,
  created_at timestamp not null default current_timestamp,
  project_id varchar(32) not null,

  name varchar(100) not null,
  description longtext not null,

  primary key (id),
  foreign key service_accounts_project_id(project_id)
  references projects(id)
  on delete cascade
);

--
-- ServiceAccountAccessKeys
--

create table if not exists service_account_access_keys (
  id varchar(32) not null,
  created_at timestamp not null default current_timestamp,
  project_id varchar(32) not null,
  service_account_id varchar(32) not null,

  -- SENSITIVE FIELD
  hash varchar(255) not null,
  description longtext not null,

  primary key (id),
  unique(hash),
  foreign key service_account_access_keys_project_id(project_id)
  references projects(id)
  on delete cascade,
  foreign key service_account_access_keys_service_account_id(service_account_id)
  references service_accounts(id)
  on delete cascade
);

--
-- ServiceAccountRoleBindings
--

create table if not exists service_account_role_bindings (
  service_account_id varchar(32) not null,
  role_id varchar(32) not null,
  created_at timestamp not null default current_timestamp,
  project_id varchar(32) not null,

  primary key (service_account_id, role_id),
  foreign key service_account_role_bindings_service_account_id(service_account_id)
  references service_accounts(id)
  on delete cascade,
  foreign key service_account_role_bindings_role_id(role_id)
  references roles(id)
  on delete cascade,
  foreign key service_account_role_bindings_project_id(project_id)
  references projects(id)
  on delete cascade
);

--
-- DeviceRegistrationTokens
--

create table if not exists device_registration_tokens (
  id varchar(32) not null,
  created_at timestamp not null default current_timestamp,
  project_id varchar(32) not null,

  name varchar(100) not null,
  description longtext not null,
  max_registrations int,
  labels longtext not null,

  primary key (id),
  unique(name, project_id)
);

--
-- Devices
--

create table if not exists devices (
  id varchar(32) not null,
  created_at timestamp not null default current_timestamp,
  project_id varchar(32) not null,

  name varchar(100) not null,
  registration_token_id varchar(32),
  desired_agent_spec longtext not null,
  desired_agent_version varchar(100) not null,
  info longtext not null,
  last_seen_at timestamp not null default current_timestamp,
  labels longtext not null,

  primary key (id),
  unique(name, project_id),
  foreign key devices_project_id(project_id)
  references projects(id)
  on delete cascade,
  foreign key devices_registration_token_id(registration_token_id)
  references device_registration_tokens(id)
  on delete set null
);

--
-- DeviceAccessKeys
--

create table if not exists device_access_keys (
  id varchar(32) not null,
  created_at timestamp not null default current_timestamp,
  project_id varchar(32) not null,
  device_id varchar(32) not null,

  -- SENSITIVE FIELD
  hash varchar(255) not null,

  primary key (id),
  unique(hash),
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
  description longtext not null,
  scheduling_rule longtext not null,
  service_metrics_config longtext not null,

  primary key (id),
  unique(name, project_id),
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
  raw_config longtext not null,
  created_by_user_id varchar(32),
  created_by_service_account_id varchar(32),

  primary key (id),
  foreign key releases_application_id(application_id)
  references applications(id)
  on delete cascade,
  foreign key releases_created_by_user_id(created_by_user_id)
  references users(id)
  on delete set null,
  foreign key releases_created_by_service_account_id(created_by_service_account_id)
  references service_accounts(id)
  on delete set null
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
-- MetricTargetConfigs
--

create table if not exists metric_target_configs (
  id varchar(32) not null,
  created_at timestamp not null default current_timestamp,
  project_id varchar(32) not null,

  type varchar(100) not null,
  configs longtext not null,

  primary key (id),
  unique(type, project_id),
  foreign key metric_target_configs_project_id(project_id)
  references projects(id)
  on delete cascade
);

--
-- Commit
--

commit;

