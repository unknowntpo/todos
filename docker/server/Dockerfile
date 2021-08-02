FROM golang:alpine AS builder

RUN apk update && \
    apk add --no-cache curl \
            git \
            bash \
            make

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
FROM scratch AS production

WORKDIR /app/
COPY --from=builder /app/bin .
COPY --from=builder /app/.envrc .
COPY --from=builder /app/docker/server/entrypoint.sh .

CMD ["./entrypoint.sh"]
