#!/usr/bin/env sh
set -eu

go test ./...
go build -o /tmp/swe-workshop-server ./cmd/server

echo "Local Go checks completed successfully."
