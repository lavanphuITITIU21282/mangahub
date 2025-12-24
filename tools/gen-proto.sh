#!/usr/bin/env bash
set -euo pipefail

# Generate Go code from proto definitions.
# Requirements:
#   - protoc
#   - protoc-gen-go
#   - protoc-gen-go-grpc

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

protoc --go_out=. --go-grpc_out=.   --go_opt=paths=source_relative   --go-grpc_opt=paths=source_relative   proto/mangahubv1/mangahub.proto

echo "✅ Generated: proto/mangahubv1/*.pb.go"
