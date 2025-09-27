# Project PattyWagon

TeamSolid Project 3 Implementation for ProjectSprint

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

## MakeFile

Run build make command with tests

```bash
make all
```

Build the application

```bash
make build
```

Run the application

```bash
make run
```

Create DB container

```bash
make docker-run
```

Shutdown DB Container

```bash
make docker-down
```

DB Integrations Test:

```bash
make itest
```

Live reload the application:

```bash
make watch
```

Run the test suite:

```bash
make test
```

Clean up binary from the last build:

```bash
make clean
```

Create db migrations script

```bash
make db-migrate-create file={file_name}
```

DB migrations up

```bash
make db-migrate-up
```

DB migrations down

```bash
make db-migrate-down
```

DB generate sql code

```bash
make db-generate-sql
```

## Setup Dependencies

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.4.0
```

## Create .env file

Create your .env file and use the .env.sample file as reference

## Setup Goose Env

```bash
 export GOOSE_DRIVER=postgres
 export GOOSE_MIGRATION_DIR=db/sql/migrations
 export GOOSE_DBSTRING="user=postgres password=postgres dbname=patty-wagon-dev host=localhost port=5432 sslmode=disable"
```

## Create SQL query

Write your query in db/sql/query directory and run `make db-generate-sql`

## Run PattyWagon API locally with docker compose 
```
make deploy-local
```

## Setup Database Monitoring uisng Docker Compose
- Monitor database performance using these following tools:
  - Prometheus
  - Prometheus Exporter
  - Grafana

  Prometheus exproter will hit database endpoint to get the resource metrics and these data will be scraped by Prometheus. Grafana then serves the visualization to us!
- Command:
  ```bash
  docker-compose --profile monitoring up -d
  ```
  * It will also run other services like PostrgeSQL database and MinIO

- Open Grafana locally -> `http://localhost:3000`
  | setup Your own credentials (default-> username: admin | password: admin)
