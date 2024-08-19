#!/bin/bash
set -e

cmd="$1"
user="$2"
dbname="$3"
shift

until pg_isready -h localhost -p 5432 -U "$user" -d "$dbname"; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up - executing command"
exec $cmd
