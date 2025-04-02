#!/bin/bash
set -e

# Create proto output directory if it doesn't exist
mkdir -p proto/ui

# Generate Go code from proto
protoc \
  --go_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_out=. \
  --go-grpc_opt=paths=source_relative \
  proto/ui/ui.proto

echo "Proto generation complete." 