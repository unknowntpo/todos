#!/bin/bash

ls -l
# do sql migration
#make db/migrations/up
# dump test data

# set up logfile location to /var/log/app.

# ./api with args
./api -db-dsn=$TODOS_DB_DSN

#ping localhost:5432
