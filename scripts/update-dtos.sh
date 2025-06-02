#!/usr/bin/env bash

cd server;
make dtos;

rm -rf ../web/src/dto
mv ts ../web/src/dto