# Database Guide

## SQL

- Timestamp columns should use `timestamptz`, which are frequently populated
  with UTC timestamps.
- Timestamp columns should not use default values, and should instead be
  populated by application code. This allows for better application and test
  cohesion when timekeeping can vary across hosts.
- `ON CASCADE DELETE` should be avoided, in favor of intentionally deleting with
  `CASCADE` on a per-use case basis.
- Down migrations should use `IF EXISTS` wherever possible in case of partially
  failed up migrations.

## Go

- Timestamp accuracy of the database should be understood and accommodated.
  PostgreSQL has microsecond accuracy, which requires truncation of timestamps
  during INSERT. The `pgx` driver handles this transparently, but for
  consistency, the developer must decide to either truncate a timestamp
  themselves before INSERT, or return the microsecond-specific timestamp from
  the database after INSERT.
