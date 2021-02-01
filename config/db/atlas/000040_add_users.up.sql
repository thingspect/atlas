DROP TYPE IF EXISTS role;
CREATE TYPE role AS ENUM ('ROLE_UNSPECIFIED', 'CONTACT', 'VIEWER', 'BUILDER', 'ADMIN', 'SYS_ADMIN');

CREATE TABLE users (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id uuid NOT NULL REFERENCES orgs (id),
  email varchar(80) NOT NULL,
  password_hash bytea NOT NULL DEFAULT '',
  role role NOT NULL,
  status status NOT NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL
);

CREATE UNIQUE INDEX users_read_and_email_idx ON users (org_id, email);
CREATE INDEX users_read_and_paginate_idx ON users (org_id, created_at, id);
