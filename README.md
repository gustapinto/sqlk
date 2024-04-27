# SQLK

A minimal terminal-based sql querier that uses sqlc-compatible files.

This project does not aims to be a full DBA tool, like DBeaver or Jetbraind Datagrip, instead we provide only a simple command line based tool to execute queries and return its results in plain-text open formats, like JSON.

The SQLK outputs are intended to be as plain and simple as possible, in order to be parsed by external tools, like `jq` and `miller`.

## Usage

```bash
sqlk -conn connection_name -file file.sql -query namedquery
```

## The query file format

SQLK uses the same .sql file format as [sqlc](https://docs.sqlc.dev/en/latest/).

And example query file, named `foo.sql`:

```sql
-- name: SelectNumber :one
SELECT 1 AS "NUMBER";
```

To execute the SelectNumber query we can execute SQLK as:

```bash
sqlk -file foo.sql -query selectnumber
```

This call would return the following:

```json
{"NUMBER":1}
```

## Configuration

SQLK uses a json based configuration for its connection settings, the default configuration file can be found at `$HOME/.sqlk.json`, but other files can be passed using the `-config file.path.json` flag.

An example configuration file:

```json
{
    "connections": [
        {
            "name": "conn",
            "driver": "postgres",
            "dsn": "postgres://conn:conn@localhost:5432/database_name?sslmode=disable",
            "default": true
        }
    ]
}
```