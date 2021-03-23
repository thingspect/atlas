CREATE TABLE alarms (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id uuid NOT NULL REFERENCES orgs (id),
  rule_id uuid NOT NULL REFERENCES rules (id),
  name varchar(80) NOT NULL,
  status status NOT NULL,
  user_tags varchar(255)[] NOT NULL,
  subject_template varchar(1024) NOT NULL,
  body_template varchar(4096) NOT NULL,
  repeat_interval integer NOT NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL
);

CREATE INDEX alarms_read_and_paginate_idx ON alarms (org_id, created_at, id);
CREATE INDEX alarms_read_and_paginate_filter_rule_id_idx ON alarms (org_id, rule_id, created_at, id);
