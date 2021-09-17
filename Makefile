include .envrc
# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

# ==================================================================================== #
# TEST
# ==================================================================================== #

## test/unit: execute unit test
.PHONY: test/unit
test/unit:
	@go test -v ./... -short -count=1 -cover -race

## test/integration: execute integration test
.PHONY: test/integration
test/integration:
	@echo "Running integration tests using testcontainers..."
	@go test -v ./... -run=Integration -cover -race

# ==================================================================================== #
# BUILD
# ==================================================================================== #


current_time = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
git_description = $(shell git describe --always --dirty --tags --long)
linker_flags = '-s -X main.buildTime=${current_time} -X main.version=${git_description}'

## build/docker/image: build the docker image for api
.PHONY: build/docker/image
build/docker/image:
	@docker build -f ./cmd/api/Dockerfile -t api .

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	CGO_ENABLED=0 go build -ldflags=${linker_flags} -o=./bin/api ./cmd/api

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/api: run the cmd/api application with trusted origin = swagger UI server
.PHONY: run/api
run/api: db/start
	@TODOS_APP_DB_DSN=$(TODOS_APP_DB_DSN_LOCAL) go run ./cmd/api -c ./config/app_config-dev.yml



## run/compose/up: run the services
.PHONY: run/compose/up
run/compose/up:
	@docker-compose -f docker-compose-prod.yml --env-file .envrc up \
	    -d \
	    --build \
	    --remove-orphans \
	    --force-recreate

## run/compose/down: shutdown the services
.PHONY: run/compose/down
run/compose/down:
	@docker-compose -f docker-compose-prod.yml down

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up:
	@echo 'Running up migrations...'
	migrate -path ./migrations -database $(TODOS_APP_DB_DSN) up

## db/start: start a postgres container with testdata
.PHONY: db/start
db/start:
	@echo "Start a new postgres db with testdata..."
	@docker-compose -f docker-compose-db.yml --env-file .envrc up \
	    -d \
	    --remove-orphans \
	    --force-recreate \
	    --build

## db/stop: stop a postgres container.
.PHONY: db/stop
db/stop:
	@echo "Stop postgres db container..."
	@docker-compose -f docker-compose-prod.yml down

## db/connect: connect to the database in postgres container
.PHONY: db/connect
db/connect:
	@docker exec -it todos_devdb psql $(TODOS_APP_DB_DSN_LOCAL)

## docs/gen: use swagger codegen to generate API documentation
.PHONY: docs/gen
docs/gen:
	@echo "Generating swagger API documentation..."
	@swag init --dir ./cmd/api \
	    --parseDependency \
	    --parseInternal \
	    -o ./docs


## docs/show: use swaggerUI container to show API documentation. 
.PHONY: docs/show
docs/show:
	@echo "Showing swagger API documentation at :8080..."
	@docker run --rm -d -p 8080:8080 \
	-e SWAGGER_JSON='/docs/swagger.yaml' \
	-e BASE_URL='/' \
	--mount type=bind,source="$(shell pwd)"/docs,target=/docs \
	swaggerapi/swagger-ui





# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit:
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	$(MAKE) test/unit	

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor
