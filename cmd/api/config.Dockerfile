FROM postgres:alpine as builder

RUN apk update && \
    apk add curl \
            git \
            bash \
            make \
    rm -rf /var/cache/apk/*

# install migrate which will be used by entrypoint.sh to perform DB migration
ARG MIGRATE_VERSION=4.7.1
ADD https://github.com/golang-migrate/migrate/releases/download/v${MIGRATE_VERSION}/migrate.linux-amd64.tar.gz /tmp
RUN tar -xzf /tmp/migrate.linux-amd64.tar.gz -C /usr/local/bin && mv /usr/local/bin/migrate.linux-amd64 /usr/local/bin/migrate

# TODO: set up golang-migrate
WORKDIR /app

# COPY source files to container
COPY . .

# Distribution
FROM postgres:alpine

RUN apk --no-cache add bash make

WORKDIR /app/
COPY --from=builder /app/bin .
COPY --from=builder /app/Makefile .
COPY --from=builder /app/.envrc .
COPY --from=builder /usr/local/bin/migrate /usr/local/bin
COPY --from=builder /app/migrations ./migrations/
COPY --from=builder /app/testdata ./testdata/
COPY --from=builder /app/cmd/api/config.sh .

CMD ["./config.sh"]


