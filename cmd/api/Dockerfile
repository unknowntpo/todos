FROM golang:alpine AS builder

RUN apk update && \
    apk add curl \
            git \
            base \
            make \
    rm -rf /var/cache/apk/*

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
COPY --from=builder /app/cmd/api/entrypoint.sh .

CMD ["./entrypoint.sh"]
