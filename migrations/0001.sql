-- DIFF:
-- create table if not exists device_registration_tokens (
--   ADD name varchar(100) not null,
--   ADD max_registrations int,

--   ADD unique(name, project_id)

--   REMOVE device_access_key_id varchar(32) default null,
-- );

-- To ease our migration, let's delete all device registration tokens (we don't
-- need these one-time tokens anyway)
DELETE FROM device_registration_tokens WHERE true;

ALTER TABLE device_registration_tokens
ADD name varchar(100) not null after project_id;

ALTER TABLE device_registration_tokens
ADD max_registrations int after name;

ALTER TABLE device_registration_tokens
ADD unique(name, project_id);

ALTER TABLE device_registration_tokens
DROP COLUMN device_access_key_id;

-- Additionally, add a "default" token for all projects
INSERT INTO device_registration_tokens
(id, project_id, name)
SELECT
CONCAT("reg_", LEFT(MD5(RAND()), 32-4)) as id,
id as project_id,
'default' as name
from projects;

-- DIFF:
-- create table if not exists devices (
--   ADD registration_token_id varchar(32),

--   ADD foreign key devices_registration_token_id(registration_token_id)
--   references device_registration_tokens(id)
--   on delete set null
-- );

ALTER TABLE devices
ADD registration_token_id varchar(32) after name;

ALTER TABLE devices
ADD CONSTRAINT devices_registration_token_id FOREIGN KEY (registration_token_id)
references device_registration_tokens(id) on delete set null;