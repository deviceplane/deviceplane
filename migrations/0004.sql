begin;

ALTER TABLE applications
ADD scheduling_rule longtext not null after description;

ALTER TABLE applications
DROP COLUMN settings;

commit;