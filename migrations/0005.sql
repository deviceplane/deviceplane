begin;

ALTER TABLE releases
CHANGE COLUMN config raw_config longtext not null;

ALTER TABLE releases
ADD config longtext not null after raw_config;

commit;