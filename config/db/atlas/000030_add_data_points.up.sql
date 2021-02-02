CREATE TABLE data_points (
  org_id uuid NOT NULL,
  uniq_id varchar(40) NOT NULL CHECK (uniq_id = lower(uniq_id)),
  attr varchar(40) NOT NULL,
  int_val integer,
  fl64_val double precision,
  str_val varchar(255),
  bool_val boolean,
  bytes_val bytea CHECK (octet_length(bytes_val) <= 255),
  created_at timestamptz NOT NULL,
  trace_id uuid NOT NULL,
  PRIMARY KEY (org_id, uniq_id, attr, created_at),
  CHECK (num_nonnulls(int_val, fl64_val, str_val, bool_val, bytes_val) = 1)
);

CREATE INDEX data_points_list_and_latest_idx ON data_points (org_id, uniq_id, attr, created_at DESC);
