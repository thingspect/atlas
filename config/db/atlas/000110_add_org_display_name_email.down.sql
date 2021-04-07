ALTER TABLE orgs DROP COLUMN email;
ALTER TABLE orgs DROP COLUMN display_name;

ALTER TYPE alarm_type RENAME TO alarm_type_old;
CREATE TYPE alarm_type AS ENUM ('ALARM_TYPE_UNSPECIFIED', 'APP', 'SMS');
ALTER TABLE alarms ALTER COLUMN type TYPE alarm_type USING type::text::alarm_type;
DROP TYPE alarm_type_old;
