FROM golang:alpine as builder

RUN apk update && \
    apk add --no-cache curl \
            git \
            bash \
            make \
            build-base

# TODO: set up golan-migrate
WORKDIR /app

# copy module files first so that they don't need to be downloaded again if no change
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go mod verify

COPY . .

FROM golang:alpine

RUN apk --no-cache add bash make build-base

WORKDIR /app/

# COPY source files to container
COPY --from=builder /app/docker/integration-test/it-test.sh .
COPY --from=builder /app .

CMD ["./it-test.sh"]
