# build base
FROM golang:alpine as build-base

ARG VERSION
ENV TODOS_BIN_VERSION=${VERSION}

RUN apk add --no-cache curl \
git \
bash \
make \
build-base

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod verify

COPY internal ./internal
COPY pkg ./pkg
COPY cmd ./cmd
COPY config ./config
COPY docs ./docs
COPY vendor ./vendor
COPY Makefile .
COPY .envrc .
COPY app_config* .

RUN make build/server

# Development stage
# Use golang:alpine image, great for debugging.
FROM golang:alpine AS dev-server

WORKDIR /app/
COPY --from=build-base /app/bin .
COPY --from=build-base /app/app_config-dev.yml .

EXPOSE 4000

CMD ["./server", "-c", "./app_config-dev.yml"]

# Production stage
FROM scratch AS production-server

WORKDIR /app/
COPY --from=build-base /app/bin .
COPY --from=build-base /app/app_config-prod.yml .

EXPOSE 4000

CMD ["./server", "-c", "./app_config-prod.yml"]

# config base
FROM alpine as config-base
RUN apk add --no-cache curl \
git \
bash \
make \
postgresql-client

ARG MIGRATE_VERSION=4.7.1
ADD https://github.com/golang-migrate/migrate/releases/download/v${MIGRATE_VERSION}/migrate.linux-amd64.tar.gz /tmp
RUN tar -xzf /tmp/migrate.linux-amd64.tar.gz -C /usr/local/bin && mv /usr/local/bin/migrate.linux-amd64 /usr/local/bin/migrate

WORKDIR /app

COPY Makefile .
COPY .envrc .
COPY migrations ./migrations
COPY testdata ./testdata
COPY docker ./docker

# config stage
FROM config-base AS config

WORKDIR /app/
COPY --from=config-base /app/Makefile .
COPY --from=config-base /app/.envrc .
COPY --from=config-base /usr/local/bin/migrate /usr/local/bin
COPY --from=config-base /app/migrations ./migrations/
COPY --from=config-base /app/testdata ./testdata/
COPY --from=config-base /app/docker/config.sh .

CMD ["./config.sh"]
