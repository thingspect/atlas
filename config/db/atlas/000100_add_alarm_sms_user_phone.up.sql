ALTER TYPE alarm_type RENAME TO alarm_type_old;
CREATE TYPE alarm_type AS ENUM ('ALARM_TYPE_UNSPECIFIED', 'APP', 'SMS');
ALTER TABLE alarms ALTER COLUMN type TYPE alarm_type USING type::text::alarm_type;
DROP TYPE alarm_type_old;

ALTER TABLE users ADD COLUMN phone varchar(16) NOT NULL DEFAULT '';
ALTER TABLE users ALTER COLUMN phone DROP DEFAULT;
