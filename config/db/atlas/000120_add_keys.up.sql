CREATE TABLE keys (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id uuid NOT NULL REFERENCES orgs (id),
  name varchar(80) NOT NULL,
  role role NOT NULL,
  created_at timestamptz NOT NULL
);

CREATE INDEX keys_read_and_paginate_idx ON keys (org_id, created_at, id);
