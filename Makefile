clean:
	rm -rf bin

gen: api.yaml
	mkdir -p gen
	oapi-codegen --config=gen/config.yaml api.yaml

build: gen
	go build -mod=mod -o bin/mailchump cmd/main.go

build-docker: gen  # For use in Dockerfile
	CGO_ENABLED=0 GOOS=linux go build -mod=mod -o bin/mailchump cmd/main.go

build-docker-dev: gen  # For use in Dockerfile
	CGO_ENABLED=0 GOOS=linux go build -mod=mod -o bin/mailchump cmd/local.go

package:
	docker build --tag mailchump:latest .

package-dev:
	docker build --target dev --tag mailchump:dev .

run: gen
	ENVIRONMENT=DEV go run cmd/local.go

run-container: package-dev compose-down
	docker compose up --force-recreate

run-db: package-dev compose-down
	INIT_ONLY=true docker compose up -d --force-recreate

compose-down:
	docker compose down

test: gen
	go test ./...

test-verbose: gen
	go test -v ./...

test-coverage: gen
	go test -coverprofile=coverage.out ./...

test-coverage-html: test-coverage
	go tool cover -html=coverage.out
