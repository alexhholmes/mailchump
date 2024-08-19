# Mailchump

A shitty newsletter service

## Project setup
Codegen from OpenAPI spec
```
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
```

## Project Structure
```
.
├── README.md
├── api
│   ├── api.gen.go # Generated API code
│   ├── api.yaml # OpenAPI spec
├── cmd
│   └── main.go # Main entry point
├── model # Database models
```

## References
