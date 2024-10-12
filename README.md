# Mailchump

A shitty newsletter service made with vanilla golang and minimal dependencies;
using PostgreSQL for storage and Redis for caching. Additionally, it is packaged
with docker and tested/deployed to Google GCP using Github Actions.

_Why use big import when small import work good?_

[Just use Postgres](https://mccue.dev/pages/8-16-24-just-use-postgres)

A list of packages used (from `go.mod`):
```text
github.com/lib/pq               # PostgresDB driver for database/sql
github.com/google/uuid          # IDs for database primary keys
github.com/oapi-codegen/runtime # Open API 3 http server codegen
github.com/a-h/templ            # HTML templating engine
github.com/vektra               # Mocking for testing
github.com/redis/go-redis       # Caching
github.com/stretchr/testify     # Testing
```

This repository aims to reduce ~~carbon emissions~~ tech debt by 100% through sustainable practices.

## Project Setup
Codegen from OpenAPI spec and HTML templating (needed for IDE indexing support):
```shell
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
go install github.com/a-h/templ/cmd/templ@latest
```

Codegen for testing mocks:
```shell
brew install mockery

# or (not recommended)

go install github.com/vektra/mockery/v2@latest
```

## Project Structure
```
.
├── api          # API route handlers and application entrypoint
│   └── gen      # Generated API code and codegen config file
├── cmd
│   └── main.go  # Entrypoint
├── model        # Database models
├── migrations   # Database SQL migrations
├── pgdb         # Database connection initialization and utility functions
├── scripts      # Development environment scripts
├── Makefile     # Everything you need to test/build/deploy this project
```

## Feature Roadmap
In no particular order:
- [ ] Complete all API endpoints in the OpenAPI spec
- [ ] Complete all API endpoint implementations
- [ ] Add user authentication
- [ ] Add email sending of newsletters
- [ ] Add Redis caching
- [ ] Add unit tests
- [ ] Add integration tests
- [ ] Add performance/load testing
- [ ] Add regression testing
- [ ] Add CI/CD pipeline with Github Actions
- [ ] Add deployment to GCP
- [ ] Add Open Telemetry tracing and metrics
- [ ] Export logs to ???
- [ ] Add autoscaling to GCP deployment
- [ ] ? Shard database for scaling
- [ ] ? Shard cache for scaling
- [ ] ? Optimize email sending for performance (e.g. batching or streaming to another service)
- [ ] ? Shard database for scaling (unintended replication pun)
- [ ] ? Production monitoring and alerting
- [ ] ? Fuzzing testing (not sure if this would be applicable to this project)

## References
[github.com/alexhholmes/mailchump](https://github.com/alexhholmes/mailchump/blob/main/README.md)
