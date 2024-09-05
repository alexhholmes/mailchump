# Build/Run targets
#

.PHONY: build
build: gen
	go build -mod=mod -o bin/mailchump cmd/main.go

.PHONY: build-docker
build-docker: gen  # For use in Dockerfile
	CGO_ENABLED=0 GOOS=linux go build -mod=mod -o bin/mailchump cmd/main.go

package:
	docker build --target prod --tag mailchump:latest .

package-dev:
	docker build --target dev --tag mailchump:dev .

.PHONY: run
run: package-dev stop
	docker compose up --force-recreate

.PHONY: run-db
run-db: package-dev stop
	INIT_ONLY=true docker compose up -d --force-recreate

.PHONY: stop
stop:
	docker compose down

# Util targets
#

clean:
	rm -rf bin coverage.out

gen: pkg/api
	oapi-codegen --config=config/oapi-codegen.yaml config/api.yaml
	mockery --config=config/mockery.yaml

# Testing targets
#

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

