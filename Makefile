include .envrc

GIT_HOOKS := .git/hooks/applied

# check installation of githooks and display help message when typing make
all: help $(GIT_HOOKS)

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

$(GIT_HOOKS):
	@scripts/install-git-hooks
	@echo

# ==================================================================================== #
# TEST
# ==================================================================================== #

## test/unit: execute unit test
.PHONY: test/unit
test/unit:
	@go test -v -short ./... -count=1 -race -coverprofile=coverage.txt -covermode=atomic

## test/integration: execute integration test
.PHONY: test/integration
test/integration:
	@echo "Running integration tests using testcontainers..."
	@go test -v ./internal/.../repository/... -run=Integration -count=1 -race -coverprofile=coverage.txt -covermode=atomic

# ==================================================================================== #
# BUILD
# ==================================================================================== #


current_time = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
git_description = $(shell git describe --always --dirty --tags --long)
linker_flags = '-s -X main.buildTime=${current_time} -X main.version=${TODOS_BIN_VERSION}' 

## image/build/server: build the docker image for server
.PHONY: image/build/server
image/build/server:
	@DOCKER_BUILDKIT=1 docker build \
	    --file docker/Dockerfile \
	    --target production-server \
	    --build-arg VERSION=${git_description} \
	    --build-arg BUILDKIT_INLINE_CACHE=1 \
	    --network host \
	    --cache-from unknowntpo/todos-production-server:latest \
	    --tag unknowntpo/todos-production-server:latest .

## image/build/config: build the config image for database.
.PHONY: image/build/config
image/build/config:
	@DOCKER_BUILDKIT=1 docker build \
	    --file docker/Dockerfile \
	    --target config \
	    --build-arg BUILDKIT_INLINE_CACHE=1 \
	    --network host \
	    --cache-from unknowntpo/todos-config:latest \
	    --tag unknowntpo/todos-config:latest .

## image/push/server: push the server image to dockerhub
.PHONY: image/push/server
image/push/server:
	@docker push unknowntpo/todos-production-server:latest

## image/push/config: push the config image to dockerhub
.PHONY: image/push/config
image/push/config:
	@docker push unknowntpo/todos-config:latest

## build/server: build the cmd/server application
.PHONY: build/server
build/server:
	@echo 'Building cmd/server...'
	CGO_ENABLED=0 go build -ldflags=${linker_flags} -o=./bin/server ./cmd/server


# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## mockery: generate mocks of interfaces using mockery
.PHONY: mockery
mockery:
	@echo "Generating mocks of interfaces..."
	@mockery --all  --dir ./internal/domain --output ./internal/domain/mocks

## run/compose/up: run the services
.PHONY: run/compose/up
run/compose/up:
	@DOCKER_BUILDKIT=1 docker-compose \
	    -f docker-compose.yml \
	    -f docker-compose-dev.yml \
	    --project-name todos-prod \
	    build \
	    --parallel \
	    --build-arg VERSION=${git_description}
	@docker-compose \
	    -f docker-compose.yml \
    	    -f docker-compose-dev.yml \
    	    --env-file .envrc \
	    up \
	    -d \
	    --remove-orphans \
	    --force-recreate

## run/compose/down: shutdown the services
.PHONY: run/compose/down
run/compose/down:
	@docker-compose -f docker-compose.yml -f docker-compose-dev.yml down --remove-orphans

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

## db/connect: connect to the database in postgres container
.PHONY: db/connect
db/connect:
	@docker exec -it todos_dev-db psql $(TODOS_APP_DB_DSN_LOCAL)

## docs/gen: use swagger codegen to generate API documentation
.PHONY: docs/gen
docs/gen:
	@echo "Generating swagger API documentation..."
	@swag init --dir ./cmd/server,./internal \
	    --parseDepth 5 \
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
# BENCHMARK
# ==================================================================================== #

## bench/naivepool: Run the benchmark of naivepool and output the plot.
bench/naivepool:
	@go test -v ./pkg/naivepool -bench=. \
	    -cpuprofile cpu.pprof \
	    -memprofile mem.pprof | \
	    sed -n 's/^BenchmarkExecute1000Tasks\///p' | \
	    awk '/(\w)+-[0-9]+/{print $$1, $$3}' > ./gnuplot/perf.dat
	@gnuplot \
		-e "output_path='./gnuplot/perf.png'" \
		-e "input_path='./gnuplot/perf.dat'" \
		./gnuplot/naivepool_perf.gp

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

# ==================================================================================== #
# PRODUCTION
# ==================================================================================== #

production_host_ip = todos.unknowntpo.net

## production/connect: connect to the production server
.PHONY: production/connect
production/connect:
	ssh todos@${production_host_ip}

## production/configure/caddyfile: configure the production Caddyfile
.PHONY: production/configure/caddyfile
production/configure/caddyfile:
	rsync -P ./remote/production/Caddyfile todos@${production_host_ip}:~
	ssh -t todos@${production_host_ip} '\
		sudo mv ~/Caddyfile /etc/caddy/ \
		&& sudo systemctl reload caddy'

## production/deploy: deploy the services
.PHONY: production/deploy
production/deploy:
	@docker-compose \
	    --context production \
	    -f docker-compose.yml -f docker-compose-prod.yml \
	    --env-file .envrc \
	    up \
	    -d \
	    --remove-orphans \
	    --force-recreate
