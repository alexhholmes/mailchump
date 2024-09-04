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
	rm -rf bin coverage.out

gen: pkg/api
	oapi-codegen --config=config/oapi-codegen.yaml config/api.yaml
	mockery --config config/mockery.yaml

# Testing targets
.PHONY: test-cov-out
test: gen
	go test ./...

.PHONY: test-cov-out
test-verbose: gen
	go test -v ./...

.PHONY: test-cov-out
test-cov: gen
	go test -cover ./...

.PHONY: test-cov-out
test-cov-out: gen clean
	go test -coverprofile=coverage.out ./...
	cat coverage.out | grep -v '\.gen\.go' > temp.out && mv temp.out coverage.out
	go tool cover -html=coverage.out

