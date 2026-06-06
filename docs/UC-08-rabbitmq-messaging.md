# UC-08 — RabbitMQ Messaging (RabbitManager)

> **Status**: Active  
> **Owner**: Wilie  
> **Last Updated**: 2026-06-05  
> **Package**: `src/infrastructure/messaging`

---

## Overview

Module messaging mengenkapsulasi semua komunikasi dengan RabbitMQ menggunakan `RabbitManager` — publisher yang robust dengan channel pool, auto-reconnect, dan retry. Menggantikan `wiliehidayat87/rmqp` yang lama untuk path publish.

---

## Dua Pola Publish

### 1. `PublishWithRetry` — Async (Fire-and-Forget)

Digunakan untuk alur `RTO` (Ratio Transaction One-way). Tidak menunggu reply dari worker.

```go
pubCtx, pubCancel := context.WithTimeout(c.UserContext(), time.Duration(h.Config.RabbitMQCtxTimeout)*time.Second)
defer pubCancel()

corId := "RTO" + external.GetUniqId(h.Config.TZ)
if err := h.RM.PublishWithRetry(pubCtx, exchange, routingKey, bodyReq, corId); err != nil {
    h.Logs.Debug(fmt.Sprintf("[x] Failed published: %s", corId))
} else {
    h.Logs.Debug(fmt.Sprintf("[v] Published: %s", corId))
}
```

**Worker side**: proses pesan, **tidak perlu reply**.

---

### 2. `DirectReplyToWithRetry` — Sync (RPC)

Digunakan untuk alur `RTD` (Ratio Transaction Direct-reply). Handler **memblok** hingga worker membalas.

```go
pubCtx, pubCancel := context.WithTimeout(c.UserContext(), time.Duration(h.Config.RabbitMQCtxTimeout)*time.Second)
defer pubCancel()

corId := "RTD" + external.GetUniqId(h.Config.TZ)
reply, err := h.RM.DirectReplyToWithRetry(pubCtx, exchange, routingKey, bodyReq, corId)
if err != nil {
    // timeout/error
}
// reply = []byte("SHAVED") atau []byte("NOTSHAVED")
```

**Worker side**: setelah proses, harus publish ke `d.ReplyTo` dengan `CorrelationId = d.CorrelationId`.

---

## Correlation ID Convention

| Prefix | Makna | Publisher Method | Worker Reply? |
|--------|-------|-----------------|---------------|
| `RTO` | Ratio One-way (async) | `PublishWithRetry` | ❌ Tidak |
| `RTD` | Ratio Direct-reply (sync) | `DirectReplyToWithRetry` | ✅ Ya (SHAVED/NOTSHAVED) |

---

## Worker Reply Convention (RTD)

Worker (di `cores/worker`) harus melakukan ini saat `corId` prefix = `RTD`:

```go
if strings.HasPrefix(d.CorrelationId, "RTD") {
    replyBody := []byte("SHAVED")   // atau "NOTSHAVED"
    ch.PublishWithContext(ctx, "", d.ReplyTo, false, false, amqp.Publishing{
        ContentType:   "application/json",
        CorrelationId: d.CorrelationId,
        Body:          replyBody,
    })
}
// Jika prefix "RTO" → tidak perlu reply
```

---

## Queue Declarations (saat startup)

Dideklarasikan di `src/cmd/server.go` → `messaging.InitMessageBroker`:

| Exchange | Queue | Env Var |
|----------|-------|---------|
| `RABBITMQRATIOEXCHANGENAME` | `RABBITMQRATIOQUEUENAME` | Config Cfg |
| `RABBITMQCAMPAIGNMANAGEMENTEXCHANGENAME` | `RABBITMQCAMPAIGNMANAGEMENTQUEUENAME` | Config Cfg |
| `E_RESENDCAMPAIGNDATA` | `Q_RESENDCAMPAIGNDATA` | Hardcoded |

---

## RabbitManager Internals

```
RabbitManager
  ├─ URL          — amqp://user:pass@host:port/vhost
  ├─ ChPool       — chan *amqp.Channel (pool size dari RabbitConfig.PoolSize)
  ├─ PoolSize     — default 5
  ├─ Qos          — default 10 (prefetch per channel)
  ├─ Configs      — []RabbitDeclare (exchange+queue declarations)
  └─ handleReconnect() goroutine — auto-reconnect + refill pool
```

**PublishWithRetry** internals:
1. Ambil channel dari pool (`<-rm.ChPool`)
2. Defer kembalikan ke pool (`rm.ChPool <- ch`)
3. `PublishWithDeferredConfirmWithContext` (publisher confirm mode)
4. Wait untuk broker ACK
5. Jika gagal: retry 3x dengan exponential backoff (500ms / 1000ms / 1500ms)
6. Jika semua retry gagal: close connection → handleReconnect wakes up

---

## Konfigurasi

| Env Var | Default | Deskripsi |
|---------|---------|-----------|
| `RABBITMQHOST` | — | Host RabbitMQ |
| `RABBITMQPORT` | — | Port RabbitMQ |
| `RABBITMQUSERNAME` | — | Username |
| `RABBITMQPASSWORD` | — | Password |
| `RABBITMQVHOST` | — | Virtual host |
| `RABBITMQCONTEXTTIMEOUT` | `30` | Timeout detik untuk publish context |

---

## File
- [`src/infrastructure/messaging/rabbitmq.go`](../src/infrastructure/messaging/rabbitmq.go)
- [`src/cmd/server.go`](../src/cmd/server.go) (inisialisasi)
- [`src/interfaces/http/handler/incoming_handler.go`](../src/interfaces/http/handler/incoming_handler.go) (`RM *messaging.RabbitManager` field)

---

## Menambah Queue Baru

1. Tambah entry di `rmDeclares` slice di `src/cmd/server.go`
2. Tambah entry di `channels` slice (untuk legacy Rmqp reconnect watcher)
3. Gunakan `h.RM.PublishWithRetry` di handler
4. Tambah env var ke `config.Cfg` jika exchange/queue dari config
5. Update dokumen ini

---

## Perubahan Terkini
| Tanggal | Perubahan |
|---------|-----------|
| 2026-06-05 | **[NEW]** RabbitManager package dibuat di `src/infrastructure/messaging/` |
| 2026-06-05 | Semua `h.Rmqp.PublishMsg` di handler diganti → `h.RM.PublishWithRetry` |
| 2026-06-05 | `PostbackDirectReply` menggunakan `DirectReplyToWithRetry` (RTD path) |
