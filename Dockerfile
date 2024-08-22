FROM golang:latest AS builder
RUN go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
WORKDIR /app
COPY ./ ./
RUN make build-docker

FROM ubuntu AS dev
RUN apt-get update && apt-get install -y postgresql-client
WORKDIR /app
COPY --from=builder /app/migrations/tables.sql tables.sql
COPY --from=builder /app/scripts/wait-for-postgres.sh wait-for-postgres.sh
COPY --from=builder /app/bin/mailchump mailchump
EXPOSE 8080 6060
ENTRYPOINT ["./mailchump"]

FROM alpine:latest AS prod
WORKDIR /app
COPY --from=builder /app/bin/mailchump mailchump
EXPOSE 8080
ENTRYPOINT ["./mailchump"]
