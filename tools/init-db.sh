#!/bin/sh

host=$cockroachdb

./cockroach sql --insecure --host=cockroachdb --execute="CREATE DATABASE IF NOT EXISTS $1;"