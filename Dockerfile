FROM golang:latest AS builder-dev
RUN go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
WORKDIR /app
COPY ./ ./
RUN make build-docker-dev

FROM ubuntu AS dev
RUN apt-get update && apt-get install -y postgresql-client
WORKDIR /app
COPY --from=builder-dev /app/scripts/wait-for-postgres.sh wait-for-postgres.sh
COPY --from=builder-dev /app/bin/mailchump mailchump
EXPOSE 8080
EXPOSE 6060
ENTRYPOINT ["./mailchump"]

FROM golang:latest AS builder-prod
RUN go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
WORKDIR /app
COPY ./ ./
RUN make build-docker

FROM alpine:latest AS prod
WORKDIR /app
COPY --from=builder-prod /app/bin/mailchump mailchump
EXPOSE 8080
ENTRYPOINT ["./mailchump"]
