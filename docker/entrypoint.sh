#!/bin/bash

ls -l
# set up logfile location to /var/log/app.

# ./api with args
./api -c ./app_config.prod.yml
