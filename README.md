# Mailchump

A shitty newsletter service using vanilla golang with minimal dependencies
and PostgreSQL for storage. Additionally, it is packaged with docker and tested/deployed
to Google GCP using Github actions.

_Why use big import when small import work good?_

[Just Use Postgres](https://mccue.dev/pages/8-16-24-just-use-postgres)

A list of packages used (from `go.mod`):
```text
github.com/cucumber/godog
github.com/lib/pq
github.com/google/uuid
github.com/oapi-codegen/runtime
```

This repository aims to reduce ~~carbon emissions~~ tech debt by 100% through sustainable practices.

## Project Setup
Codegen from OpenAPI spec (needed for IDE indexing support):
```
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
```

## Project Structure
```
.
├── api # API route handlers and application entrypoint
├── cmd
│   └── main.go # Main entrypoint
│   └── local.go # Local entrypoint with database initialization and migration
├── gen # Generated code and codegen configuration files
├── model # Database models
├── postgres # Database connection initialization
├── scripts # Development environment scripts
```

## References
