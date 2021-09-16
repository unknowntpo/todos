#!/bin/bash

ls -l

# make sure psql and migrate is installed
migrate --version
psql --version

# psql run init.sql
echo "Executing init.sql..."
PGPASSWORD=$POSTGRES_PASSWORD psql -h db -d todos -U postgres -f ./testdata/init.sql

# do sql migration
# TODO: How to do migration to db container ?
echo "Executing db migration..."
make db/migrations/up
#migrate -path ./migrations -database $TODOS_DB_DSN up

psql $TODOS_APP_DB_DSN -c '\d'

# dump testdata

echo "Dumping testdata..."
psql $TODOS_APP_DB_DSN -f ./testdata/dummytask.sql
psql $TODOS_APP_DB_DSN -f ./testdata/dummyuser.sql

# make config container keep running
ping google.com
