#!/usr/bin/env bash

migrate create -ext=sql -dir=./apps/server/migrations $1
