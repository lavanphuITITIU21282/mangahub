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


## CLI usage (gRPC)

The same `mangahub` CLI binary can call gRPC endpoints via the `grpc` command group.

Set the gRPC address (optional):
- Default: `localhost:50051`
- Override with env var: `MANGAHUB_GRPC_ADDR`

Examples:

```bash
# search
MANGAHUB_GRPC_ADDR=localhost:50051 mangahub grpc search one

# get manga detail
mangahub grpc get one-piece

# update progress (requires you to login first so the JWT token is saved)
mangahub login demo demo
mangahub grpc progress one-piece 12
```

Note: you still need to generate protobuf Go code once (see `proto/README.md` or `tools/gen-proto.sh`).
