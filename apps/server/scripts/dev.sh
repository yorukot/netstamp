#!/usr/bin/env sh
set -eu

if [ -z "${CONFIG_FILE:-}" ] && [ -f .env ]; then
  export CONFIG_FILE=.env
fi

go run ./cmd/api
