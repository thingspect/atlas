# ChirpStack schema migrations

All of the following assume a working `.pgpass` or passed credentials. `psql` is
available in Homebrew as part of the `libpq` package.

## Init

```
psql -h localhost -p 2439 postgres postgres -f init_db.sql
psql -h localhost -p 2439 chirpstack postgres -f init_ext.sql
```
