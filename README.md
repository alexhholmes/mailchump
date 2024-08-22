# Mailchump

A shitty newsletter service using vanilla golang with minimal dependencies;
using PostgreSQL for storage and redis for caching. Additionally, it is packaged
with docker and tested/deployed to Google GCP using Github actions.

_Why use big import when small import work good?_

[Just use Postgres](https://mccue.dev/pages/8-16-24-just-use-postgres)

A list of packages used (from `go.mod`):
```text
github.com/cucumber/godog       # Integration testing
github.com/lib/pq               # PostgresDB driver for database/sql
github.com/google/uuid          # IDs for database primary keys
github.com/oapi-codegen/runtime # Open API 3 http server codegen
github.com/redis/go-redis       # Caching
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
├── api          # API route handlers and application entrypoint
├── cmd
│   └── main.go  # Entrypoint
├── gen          # Generated code and codegen configuration files
├── model        # Database models
├── pgdb         # Database connection initialization and utility functions
├── scripts      # Development environment scripts
├── Makefile     # Everything you need to test/build/deploy this project
```

## References
[github.com/alexhholmes/mailchump](https://github.com/alexhholmes/mailchump/blob/main/README.md)