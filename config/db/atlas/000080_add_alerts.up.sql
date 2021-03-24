-- alerts is lightly linked to non-org tables for retention purposes
CREATE TABLE alerts (
  org_id uuid NOT NULL REFERENCES orgs (id),
  uniq_id varchar(40) NOT NULL CHECK (uniq_id = lower(uniq_id)),
  alarm_id uuid NOT NULL,
  user_id uuid NOT NULL,
  created_at timestamptz NOT NULL,
  trace_id uuid NOT NULL,
  PRIMARY KEY (org_id, uniq_id, alarm_id, user_id, created_at)
);

CREATE INDEX alerts_list_idx ON alerts (org_id, uniq_id, alarm_id, user_id, created_at DESC);
