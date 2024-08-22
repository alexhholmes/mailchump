# Build/Run targets
build: gen
	go build -mod=mod -o bin/mailchump cmd/main.go

build-docker: gen  # For use in Dockerfile
	CGO_ENABLED=0 GOOS=linux go build -mod=mod -o bin/mailchump cmd/main.go

package:
	docker build --target prod --tag mailchump:latest .

package-dev:
	docker build --target dev --tag mailchump:dev .

run: package-dev stop
	INIT_DB=true docker compose up --force-recreate

run-db: package-dev stop
	INIT_ONLY=true docker compose up -d --force-recreate

stop:
	docker compose down

# Util targets
clean:
	rm -rf bin

gen: gen/api.yaml
	mkdir -p gen
	oapi-codegen --config=gen/config.yaml gen/api.yaml

# Testing targets
test: gen
	go test ./...

test-verbose: gen
	go test -v ./...

test-coverage: gen
	go test -coverprofile=coverage.out ./...

test-coverage-html: test-coverage
	go tool cover -html=coverage.out

