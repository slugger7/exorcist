#!/usr/bin/env bash
echo "Running migrations"

. ./apps/server/scripts/set-env.sh

echo $1
echo "Running migrations"
migrate -source file://./apps/server/migrations  -database "${DATABASE_CONNECTION_STRING}" -verbose force $1
