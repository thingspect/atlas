CREATE TABLE rules (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id uuid NOT NULL REFERENCES orgs (id),
  name varchar(80) NOT NULL,
  status status NOT NULL,
  tag varchar(255) NOT NULL,
  attr varchar(40) NOT NULL,
  expr varchar(1024) NOT NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL
);

CREATE INDEX rules_read_and_paginate_idx ON rules (org_id, created_at, id);
