ALTER TABLE users ADD COLUMN name varchar(80) NOT NULL DEFAULT '';
ALTER TABLE users ALTER COLUMN name DROP DEFAULT;
