ALTER TABLE users DROP COLUMN app_key;

ALTER TABLE alarms DROP COLUMN type;
DROP TYPE IF EXISTS alarm_type;

ALTER TABLE alerts DROP COLUMN IF EXISTS error;
ALTER TABLE alerts DROP COLUMN IF EXISTS status;
DROP TYPE IF EXISTS alert_status;
