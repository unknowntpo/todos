#!/bin/bash

ls -l
# set up logfile location to /var/log/app.

# ./api with args
./api -db-dsn=$TODOS_DB_DSN
