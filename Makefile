gen: api.yaml
	oapi-codegen --config=config.yaml api.yaml

build: gen
	go build -mod=mod -o bin/server cmd/main.go

run: gen
	go run cmd/main.go

test: gen
	go test ./...

test-verbose: gen
	go test -v ./...

test-coverage: gen
	go test -coverprofile=coverage.out ./...

test-coverage-html: test-coverage
	go tool cover -html=coverage.out
