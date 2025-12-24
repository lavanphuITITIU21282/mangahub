# Protobuf / gRPC

This project uses gRPC for a subset of features (manga browsing + progress updates).

## Generate Go code

You need:
- `protoc` (Protocol Buffers compiler)
- Go plugins:
  - `protoc-gen-go`
  - `protoc-gen-go-grpc`

Install plugins:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Generate:

```bash
# from repo root
protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative proto/mangahubv1/mangahub.proto
```

This will create the Go package: `mangahub/proto/mangahubv1`.
