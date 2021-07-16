#!/bin/bash

ls -l

# make sure psql and migrate is installed
migrate --version
psql --version

# psql run init.sql
PGPASSWORD="postgres" psql -h db -U postgres -f ./testdata/init.sql

# do sql migration
# TODO: How to do migration to db container ?
migrate -path ./migrations -database $TODOS_DB_DSN up

# dump testdata
PGPASSWORD="pa55word" psql -h db -U todos -f ./testdata/testdata.sql

# make config container keep running
ping google.com



