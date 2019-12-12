#!/bin/bash
set -e

MYSQL_USER=deviceplane \
  MYSQL_PASSWORD=deviceplane \
  MYSQL_RANDOM_ROOT_PASSWORD=yes \
  MYSQL_DATABASE=deviceplane \
  docker-entrypoint.sh mysqld &

exec controller
