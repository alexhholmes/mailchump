# Mailchump

A shitty newsletter service using vanilla golang with minimal dependencies
and PostgreSQL for storage. Additionally, packaged with docker and testing/deployment
to Google GCP using Github actions.

_Why use big import when small import work good?_

A list of packages used (from `go.mod`):
```text
github.com/cucumber/godog
github.com/lib/pq
github.com/google/uuid
github.com/oapi-codegen/runtime
```

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
