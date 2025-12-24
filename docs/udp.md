# UDP Notifications

This project uses a simple UDP pub/sub server to broadcast notifications to subscribed clients.

## Run UDP server

```bash
go run ./cmd/udp-server
```

By default it listens on `:9092`. You can override with:

```bash
export MANGAHUB_UDP_ADDR=":9092"
```

## Subscribe and listen

```bash
go run ./cmd/udp-client -mode subscribe
```

## Chapter release notification (spec feature)

Broadcast a **chapter release** event to all subscribers:

```bash
go run ./cmd/udp-client -mode release -manga one-piece -title "One Piece" -chapter 1101
```

This sends a JSON payload like:

```json
{
  "type": "chapter_release",
  "manga_id": "one-piece",
  "manga_title": "One Piece",
  "chapter": 1101,
  "message": "New chapter released: One Piece #1101",
  "released_at": "2025-12-25T05:07:57+07:00"
}
```

## Raw broadcast

```bash
go run ./cmd/udp-client -mode broadcast -msg "hello"
```
