DROP TYPE IF EXISTS status;
CREATE TYPE status AS ENUM ('STATUS_UNSPECIFIED', 'ACTIVE', 'DISABLED');

CREATE TABLE orgs (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name varchar(40) UNIQUE NOT NULL CHECK (name = lower(name)),
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL
);

CREATE INDEX orgs_paginate_idx ON orgs (created_at, id);

CREATE TABLE devices (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id uuid NOT NULL REFERENCES orgs (id),
  uniq_id varchar(40) UNIQUE NOT NULL CHECK (uniq_id = lower(uniq_id)),
  status status NOT NULL,
  token uuid UNIQUE NOT NULL DEFAULT gen_random_uuid(),
  decoder varchar(40) NOT NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL
);

CREATE INDEX devices_read_and_paginate_idx ON devices (org_id, created_at, id);
