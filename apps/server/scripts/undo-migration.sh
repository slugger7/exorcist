#!/usr/bin/env bash
echo "Undoing migration"

. ./apps/server/scripts/set-env.sh

echo "Running migrations"
migrate -source file://./apps/server/migrations  -database "${DATABASE_CONNECTION_STRING}" -verbose down 1
