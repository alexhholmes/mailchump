# Mailchump

A shitty newsletter service

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
