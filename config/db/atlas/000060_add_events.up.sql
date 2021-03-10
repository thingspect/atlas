-- events is lightly linked to non-org tables for retention purposes
CREATE TABLE events (
  org_id uuid NOT NULL REFERENCES orgs (id),
  rule_id uuid NOT NULL,
  uniq_id varchar(40) NOT NULL CHECK (uniq_id = lower(uniq_id)),
  created_at timestamptz NOT NULL,
  trace_id uuid NOT NULL,
  PRIMARY KEY (org_id, rule_id, uniq_id, created_at)
);

CREATE INDEX events_list_and_latest_idx ON events (org_id, uniq_id, rule_id, created_at DESC);
