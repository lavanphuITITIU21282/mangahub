#!/usr/bin/env bash
set -euo pipefail

# Generate Go code from proto definitions.
# Requirements:
<<<<<<< HEAD
#   - protoc in PATH
=======
#   - protoc
#   - protoc-gen-go
#   - protoc-gen-go-grpc
>>>>>>> 9689ae421f821d44d0dd8cf230b14194f222ae38

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

<<<<<<< HEAD
protoc \
  --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  proto/mangahubv1/mangahub.proto
=======
protoc --go_out=. --go-grpc_out=.   --go_opt=paths=source_relative   --go-grpc_opt=paths=source_relative   proto/mangahubv1/mangahub.proto
>>>>>>> 9689ae421f821d44d0dd8cf230b14194f222ae38

echo "✅ Generated: proto/mangahubv1/*.pb.go"
