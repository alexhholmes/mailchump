FROM golang:latest AS builder
RUN --mount=type=cache,target=/var/cache/oapi-codegen \
    go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
RUN --mount=type=cache,target=/var/cache/mockery \
    go install github.com/vektra/mockery/v2@latest
WORKDIR /app
COPY ./ ./
RUN make build-docker

FROM ubuntu AS dev
RUN apt-get update && apt-get install -y postgresql-client
WORKDIR /app
COPY --from=builder /app/migrations/ migrations/
COPY --from=builder /app/scripts/wait-for-postgres.sh /app/config/app.yaml /app/bin/mailchump ./
EXPOSE 8080 6060
ENTRYPOINT ["./mailchump"]

FROM alpine:latest AS prod
WORKDIR /app
COPY --from=builder /app/config/app.yaml config/app.yaml
COPY --from=builder /app/bin/mailchump mailchump
EXPOSE 8080
ENTRYPOINT ["./mailchump"]
