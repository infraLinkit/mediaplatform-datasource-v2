# UC-01 — Postback Intake

> **Status**: Active  
> **Owner**: Wilie  
> **Last Updated**: 2026-06-05  
> **Route Group**: `/v1/` (Public — no auth middleware)

---

## Overview

Postback endpoints menerima callback dari adnet/operator setelah user berhasil subscribe (pixel/MO). Semua endpoint ini **public** karena callback datang dari eksternal. Security via signed params / IP allowlist, bukan JWT.

---

## Endpoints

| Method | Path | Handler | Deskripsi |
|--------|------|---------|-----------|
| GET | `/v1/postback/:urlservicekey/` | `Postback` | Postback v1 — identifikasi via URL path param |
| GET | `/v1/postback` | `PostbackV3` | Postback v3 — identifikasi via query param `aff_sub` |
| GET | `/v1/postback_billed` | `PostbackBilled` | Postback khusus billed/charge event |
| GET | `/v1/postback_sync` | `PostbackDirectReply` | Postback synchronous — Direct Reply-To RabbitMQ (`RTD` correlation) |
| GET | `/v1/inquire/campid` | `InquiryCampID` | Inquiry campaign ID berdasarkan URL service key |
| GET | `/v1/inquire/api-campid` | `InquiryAPICampID` | Inquiry campaign ID untuk API channel |

---

## Flow Umum Postback

```
HTTP GET /v1/postback/:urlservicekey/
  │
  ├─ Parse params (urlservicekey, aff_sub/pixel, method, etc.)
  ├─ Validate params → 400 jika missing/invalid
  ├─ Check cookie dedup → 403 jika duplicate dalam window
  ├─ Set cookie (3 detik dedup window)
  ├─ Lookup campaign config dari Redis (cfgRediskey)
  │   └─ 404 jika campaign not found
  ├─ Lookup pixel dari Redis (pixelKey)
  │   └─ 404 jika pixel not found
  ├─ Check px.IsUsed
  │   └─ 200 NOK jika pixel sudah dipakai
  └─ Publish ke RabbitMQ Ratio queue
      ├─ RTO prefix → PublishWithRetry (async, fire-and-forget)
      └─ RTD prefix → DirectReplyToWithRetry (sync, tunggu reply SHAVED/NOTSHAVED)
```

---

## Correlation ID Convention

| Prefix | Jenis | Publisher | Worker Action |
|--------|-------|-----------|---------------|
| `RTOxxxx` | Async (fire-and-forget) | `h.RM.PublishWithRetry` | Worker proses, **tidak reply** |
| `RTDxxxx` | Sync (RPC Direct Reply-To) | `h.RM.DirectReplyToWithRetry` | Worker reply `SHAVED` atau `NOTSHAVED` |

---

## RabbitMQ

| Exchange | Queue | Env Var |
|----------|-------|---------|
| `RABBITMQRATIOEXCHANGENAME` | `RABBITMQRATIOQUEUENAME` | Config Cfg |

Publisher menggunakan `RabbitManager` (channel pool, auto-reconnect, 3x retry).  
**Jangan** gunakan `h.Rmqp.PublishMsg` lagi — sudah diganti seluruhnya.

---

## Request Parameters

### Postback v1 / v3 (GET)
| Param | Source | Required | Deskripsi |
|-------|--------|----------|-----------|
| `urlservicekey` | path / query | ✅ | Identifikasi campaign |
| `aff_sub` / `pixel` | query | ✅ | Pixel/token dari landing page |
| `method` | query | Conditional | Metode subscribe (SPC, SPC-MVLS, dll) |
| `adnet` | query | ✅ | Nama adnet pengirim |
| `pub_id` | query | Optional | Publisher ID |
| `trxid` | query | Optional | Transaction ID (wildcard support) |

### PostbackBilled (GET)
Sama seperti v3, namun khusus untuk billing event (charge).

### PostbackDirectReply (GET)
Sama seperti v3, menggunakan `RTD` correlation ID dan menunggu reply sinkron dari worker.

---

## Response Envelope

```json
// Success
{ "code": 200, "message": "OK", "data": { "url_service_key": "...", "pixel": "...", ... } }

// Pixel already used
{ "code": 200/409, "message": "NOK - Pixel already used", "data": {...} }

// Not found
{ "code": 404, "message": "Campaign ID not found" }

// Bad request
{ "code": 400, "message": "..." }

// Forbidden (duplicate cookie)
{ "code": 403, "message": "forbidden access" }
```

---

## Handler File
- [`incoming_postback_handler.go`](../src/interfaces/http/handler/incoming_postback_handler.go)
- [`routes/postback.go`](../src/interfaces/http/routes/postback.go)

---

## Perubahan Terkini
| Tanggal | Perubahan |
|---------|-----------|
| 2026-06-05 | Ganti `h.Rmqp.PublishMsg` → `h.RM.PublishWithRetry` / `DirectReplyToWithRetry` di semua postback handler |
| 2026-06-05 | Tambah `PostbackDirectReply` (RTD sync path) |
| 2026-06-05 | CorrelationID dibedakan: `RTO` (async) vs `RTD` (sync Direct Reply-To) |

---

## Menambah/Mengubah Fitur

1. Edit handler di `incoming_postback_handler.go`
2. Daftarkan route baru di `routes/postback.go`
3. Update tabel Endpoints di file ini
4. Update `CLAUDE.md` jika ada route prefix baru
5. `go build ./...` → pastikan clean
