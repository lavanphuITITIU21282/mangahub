# gRPC Server

## Run

1) Generate protobuf Go code:

```bash
bash tools/gen-proto.sh
```

2) Start the server:

```bash
go run ./cmd/grpc-server
```

Default address: `:50051` (config.GRPCAddr).

## Test with grpcurl

If you have `grpcurl` installed:

- List services:

```bash
grpcurl -plaintext localhost:50051 list
```

- Search manga:

```bash
grpcurl -plaintext -d '{"query":"one"}' localhost:50051 mangahub.v1.MangaHubService/SearchManga
```

- Update progress (requires JWT):

```bash
grpcurl -plaintext \
  -H "authorization: Bearer <JWT>" \
  -d '{"manga_id":"one-piece","chapter":12}' \
  localhost:50051 mangahub.v1.MangaHubService/UpdateProgress
```
