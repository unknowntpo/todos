# build base
FROM golang:1.16-alpine as build-base

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
COPY migrations ./migrations
COPY Makefile .
COPY .envrc .
COPY app_config* .

ARG VERSION
ENV TODOS_BIN_VERSION=${VERSION}

RUN make build/server

# Development stage
# Use golang:1.16-alpine image, great for debugging.
FROM golang:1.16-alpine AS dev-server

WORKDIR /app/
COPY --from=build-base /app/bin .
COPY --from=build-base /app/migrations /app/migrations
COPY --from=build-base /app/app_config-dev.yml .

EXPOSE 4000

CMD ["./server", "-c", "./app_config-dev.yml"]

# Production stage
FROM scratch AS production-server

WORKDIR /app/
COPY --from=build-base /app/migrations /app/migrations
COPY --from=build-base /app/bin .
COPY --from=build-base /app/app_config-prod.yml .

EXPOSE 4000

CMD ["./server", "-c", "./app_config-prod.yml"]