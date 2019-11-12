begin;

ALTER TABLE releases
CHANGE COLUMN config raw_config longtext not null;

ALTER TABLE releases
ADD config longtext not null before raw_config;

commit;