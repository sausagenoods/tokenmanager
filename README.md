# tokenmanager
This is a token management library in Go.

## PostgreSQL Store
### How to use

Create this table first:
```sql
CREATE TABLE IF NOT EXISTS tokens (
	token BYTEA NOT NULL,
	data BYTEA NOT NULL,
	expiry TIMESTAMPTZ NOT NULL,
	owner TEXT NOT NULL,
);
```
