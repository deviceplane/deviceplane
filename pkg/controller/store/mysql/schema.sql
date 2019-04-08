
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
  unique(user_id)
);


--
-- Sessions
--

create table if not exists sessions (
  id varchar(32) not null,
  created_at timestamp not null default current_timestamp,
  user_id varchar(32) not null,

  hash varchar(255) not null,

  primary key (id)
);

--
-- AccessKeys
--

create table if not exists access_keys (
  id varchar(32) not null,
  created_at timestamp not null default current_timestamp,
  user_id varchar(32) not null,

  hash varchar(255) not null,

  primary key (id)
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
-- Memberships
--

create table if not exists memberships (
  user_id varchar(32) not null,
  project_id varchar(32) not null,
  created_at timestamp not null default current_timestamp,

  level enum ('admin', 'write', 'read') not null,

  primary key (user_id, project_id)
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

  primary key (id)
);

create table if not exists device_labels (
  key varchar(100) not null,
  device_id varchar(32) not null,
  created_at timestamp not null default current_timestamp,
  project_id varchar(32) not null,

  value varchar(100) not null,

  primary key (key, device_id)
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

  primary key (id)
);

--
-- Applications
--

create table if not exists applications (
  id varchar(32) not null,
  created_at timestamp not null default current_timestamp, 
  project_id varchar(32) not null,
  name varchar(100) not null,

  primary key (id)
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

  primary key (id)
);

--
-- ReleaseStatuses
--

create table if not exists release_statuses (
  id varchar(32) not null,
  created_at timestamp not null default current_timestamp, 
  application_id varchar(32) not null,
  config longtext not null,

  primary key (id)
);

--
-- Commit
--

commit;

