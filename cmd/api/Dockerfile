FROM golang:alpine AS builder

RUN apk update && \
    apk add curl \
            git \
            base \
            make \
    rm -rf /var/cache/apk/*

# install migrate which will be used by entrypoint.sh to perform DB migration
ARG MIGRATE_VERSION=4.7.1
ADD https://github.com/golang-migrate/migrate/releases/download/v${MIGRATE_VERSION}/migrate.linux-amd64.tar.gz /tmp
RUN tar -xzf /tmp/migrate.linux-amd64.tar.gz -C /usr/local/bin && mv /usr/local/bin/migrate.linux-amd64 /usr/local/bin/migrate

# TODO: set up golan-migrate
WORKDIR /app

# copy module files first so that they don't need to be downloaded again if no change
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go mod verify

# COPY source files to container
COPY . .
RUN make build/api


# Distribution
FROM alpine:latest

RUN apk --no-cache add bash make

# create directory for logfiles.
RUN mkdir -p /var/log/app

WORKDIR /app/
COPY --from=builder /app/bin .
COPY --from=builder /app/Makefile .
COPY --from=builder /app/.envrc .
COPY --from=builder /usr/local/bin/migrate /usr/local/bin
COPY --from=builder /app/migrations ./migrations/
COPY --from=builder /app/cmd/api/entrypoint.sh .

CMD ["./entrypoint.sh"]
