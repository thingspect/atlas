DROP EXTENSION IF EXISTS pgcrypto;
ALTER TABLE orgs ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE devices ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE devices ALTER COLUMN token SET DEFAULT gen_random_uuid();
ALTER TABLE users ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE rules ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE alarms ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE keys ALTER COLUMN id SET DEFAULT gen_random_uuid();
