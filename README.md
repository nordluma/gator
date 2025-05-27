# Gator

A cli rss feed aggregator written in go. This is just a research project and
should not be used at it's current version.

## Pre-requisites

- Go
- Postgres

## Configuration

The tool can be installed with go:

```bash
go install github.com/nordluma/gator
```

This tool expects a `.gatorconfig.json` file to be located in the home dir. The
minimal this file should contain the following before using the tool: 

```json
{
    "db_url": "postgres://username:password@addr:port/gator?sslmode=disable"
}
```

After this is set and the database is running, run the migrations:

```bash
goose postgres <db_url> up
```


