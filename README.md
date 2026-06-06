# mediaplatform-datasource-v2

Go HTTP API backing the **Linkit360 MediaPlatform**. Single binary, Fiber v2, DDD layout.

## Quick Start

```bash
set -a && source .env && set +a && make run        # server
set -a && source .env && set +a && make run-migrate # migrate (run once, before server)
```

## Redis DB Indexes

| DB Index | Purpose |
|----------|---------|
| `REDISDBINDEX` (default 0) | Campaign cluster data, JSON values, JWT blacklist |
| `REDISCACHEPIXEL` (default 1) | Temporary pixel storage, postback dedup keys |

## Dokumentasi

| File | Isi |
|------|-----|
| [`CLAUDE.md`](./CLAUDE.md) | Architecture, conventions, how-to-change guide |
| [`SPEC.md`](./SPEC.md) | Authoritative spec (use cases, API contracts, entities) |
| [`AUDIT_FIX.md`](./AUDIT_FIX.md) | Audit log dan bug fixes |
| [`NOTES.md`](./NOTES.md) | Changelog per release |
| [`TESTING.md`](./TESTING.md) | Testing guide |
| [`docs/`](./docs/README.md) | **Per-UC documentation** (UC-01 s/d UC-09) |

## RabbitMQ Publish Convention

```go
// Async (RTO) — fire and forget
corId := "RTO" + external.GetUniqId(h.Config.TZ)
h.RM.PublishWithRetry(ctx, exchange, queue, body, corId)

// Sync/RPC (RTD) — tunggu reply dari worker
corId := "RTD" + external.GetUniqId(h.Config.TZ)
reply, err := h.RM.DirectReplyToWithRetry(ctx, exchange, queue, body, corId)
```

Lihat [UC-08](./docs/UC-08-rabbitmq-messaging.md) untuk detail lengkap.
