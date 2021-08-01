FROM golang:alpine

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

# COPY source files to container
COPY docker/integration-test/it-test.sh .
COPY . .

#CMD ["go", "test", "-v", "./...", "-run=Integration", "-cover", "-race"]
CMD ["./it-test.sh"]
