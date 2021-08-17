## base stage
FROM golang:alpine as base

RUN apk update && \
apk add --no-cache curl \
git \
bash \
make \
build-base \
postgresql-client

ARG MIGRATE_VERSION=4.7.1
ADD https://github.com/golang-migrate/migrate/releases/download/v${MIGRATE_VERSION}/migrate.linux-amd64.tar.gz /tmp
RUN tar -xzf /tmp/migrate.linux-amd64.tar.gz -C /usr/local/bin && mv /usr/local/bin/migrate.linux-amd64 /usr/local/bin/migrate

WORKDIR /app

# copy module files first so that they don't need to be downloaded again if no change
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go mod verify

COPY . .

RUN make build/api

## it-test stage
FROM base AS it-test

RUN apk --no-cache add bash make build-base

WORKDIR /app/
# COPY source files to container
COPY --from=base /app/docker/it-test.sh .
COPY --from=base /app .

CMD ["./it-test.sh"]

# config stage
FROM base AS config

#RUN apk --no-cache add bash make

WORKDIR /app/
#COPY --from=base /app/bin .
#COPY --from=base /app/Makefile .
COPY --from=base /app/.envrc .
COPY --from=base /usr/local/bin/migrate /usr/local/bin
COPY --from=base /app/migrations ./migrations/
COPY --from=base /app/testdata ./testdata/
COPY --from=base /app/docker/config.sh .

CMD ["./config.sh"]

# Production stage
FROM alpine AS production

WORKDIR /app/
COPY --from=base /app/bin .
COPY --from=base /app/.envrc .
COPY --from=base /app/docker/entrypoint.sh .

EXPOSE 4000

CMD ["./entrypoint.sh"]
