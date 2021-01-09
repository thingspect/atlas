CREATE TABLE users (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id uuid NOT NULL REFERENCES orgs (id),
  email varchar(80) NOT NULL,
  password_hash bytea NOT NULL CHECK (octet_length(password_hash) = 60),
  status status NOT NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL
);

CREATE UNIQUE INDEX users_read_email_idx ON users (org_id, email);
