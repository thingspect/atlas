ALTER TYPE alarm_type RENAME TO alarm_type_old;
CREATE TYPE alarm_type AS ENUM ('ALARM_TYPE_UNSPECIFIED', 'APP', 'SMS', 'EMAIL');
ALTER TABLE alarms ALTER COLUMN type TYPE alarm_type USING type::text::alarm_type;
DROP TYPE alarm_type_old;

ALTER TABLE orgs ADD COLUMN display_name varchar(80) NOT NULL DEFAULT '';
ALTER TABLE orgs ALTER COLUMN display_name DROP DEFAULT;
ALTER TABLE orgs ADD COLUMN email varchar(80) NOT NULL DEFAULT '';
ALTER TABLE orgs ALTER COLUMN email DROP DEFAULT;
