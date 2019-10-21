begin;

ALTER TABLE devices
DROP FOREIGN KEY devices_registration_token_id;

ALTER TABLE devices
ADD CONSTRAINT devices_registration_token_id FOREIGN KEY (registration_token_id)
references device_registration_tokens(id) on delete set null on update cascade;

UPDATE device_registration_tokens
set id = CONCAT("drt_", SUBSTRING(id, LENGTH('reg_') + 1, LENGTH(id) - 5)) -- note, an old migration made the "default" tokens' id too long. We fix that here
where LEFT(id, 4) = "reg_";

ALTER TABLE devices
DROP FOREIGN KEY devices_registration_token_id;

ALTER TABLE devices
ADD CONSTRAINT devices_registration_token_id FOREIGN KEY (registration_token_id)
references device_registration_tokens(id) on delete set null;

commit;