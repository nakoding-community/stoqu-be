#!/bin/bash
DIR=$PWD/internal/driver/migrations
DATABASE_URL="postgres://stoqu-be-user:stoqu-be-pass@127.0.0.1:5432/stoqu-be-db?sslmode=disable"
export DATABASE_URL
dbmate -d $DIR up