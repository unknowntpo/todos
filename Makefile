confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

run/api:
	go run ./cmd/api

db/psql:
	psql ${TODOS_DB_DSN}

db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${TODOS_DB_DSN} up
