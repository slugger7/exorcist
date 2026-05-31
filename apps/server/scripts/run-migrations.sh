#!/usr/bin/env bash
echo "Running migrations"

. ./apps/server/scripts/set-env.sh

echo "Running migrations"
migrate -path ./apps/server/migrations  -database "${DATABASE_CONNECTION_STRING}" -verbose up
