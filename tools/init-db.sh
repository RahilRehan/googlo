#!/bin/sh

./cockroach sql --insecure --host=cockroachdb --execute="CREATE DATABASE IF NOT EXISTS $1;"