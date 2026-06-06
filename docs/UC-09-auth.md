# UC-09 — Authentication & Authorization

> **Status**: Active  
> **Owner**: Wilie  
> **Last Updated**: 2026-06-05  
> **Route Group**: All protected groups (`/dashboard`, `/v1/report`, `/v1/int`, `/v1/management`)

---

## Overview

JWT-based authentication dengan HS256. Auth middleware diterapkan di group level. Aktif/non-aktif dikontrol oleh env var `AUTH_ENFORCE_DEFAULT`.

---

## Endpoints Auth

| Method | Path | Handler | Deskripsi |
|--------|------|---------|-----------|
| POST | `/v1/management/user/login` | (via user management) | Login → return JWT |
| POST | `/v1/management/user/logout` | (via user management) | Logout → blacklist JWT |
| POST | `/v1/management/user/refresh` | (via user management) | Refresh access token |

> **Catatan**: Endpoint login/logout/refresh ada di dalam user management handler. Lihat UC-04.

---

## JWT Structure

```json
{
  "sub": "user_id",
  "jti": "unique_token_id",
  "type": "access",
  "aud": "mediaplatform",
  "iss": "mediaplatform-datasource",
  "exp": 1234567890,
  "nbf": 1234567890,
  "companies": ["company_a", "company_b"],
  "adnets": ["adnet_x", "adnet_y"]
}
```

---

## Auth Middleware Flow

```
Setiap request ke protected group
  │
  ├─ Ambil header: Authorization: Bearer <token>
  ├─ Verify JWT signature (HS256, secret dari JWT_SECRET env)
  ├─ Check claims:
  │   ├─ type == "access"
  │   ├─ exp > now
  │   ├─ nbf <= now
  │   └─ aud, iss sesuai
  ├─ Check Redis blacklist: GET "jwt:blacklist:{jti}"
  │   └─ Jika ada → 401 Unauthorized
  ├─ Inject ke Locals:
  │   ├─ c.Locals("sub", claims["sub"])
  │   ├─ c.Locals("companies", []string{...})
  │   └─ c.Locals("adnets", []string{...})
  └─ c.Next()
```

---

## Redis Blacklist

Saat logout, token di-blacklist dengan TTL = sisa waktu exp:

```
Redis key: "jwt:blacklist:{jti}"
DB index:  REDISDBINDEX (default 0)
TTL:       token.exp - now
Value:     "1"
```

---

## Scope Filtering

Handler yang scope-aware menggunakan:

```go
allowedCompanies, _ := c.Locals("companies").([]string)
allowedAdnets, _ := c.Locals("adnets").([]string)
```

Filtering dilakukan di query level (DB) — user hanya melihat data sesuai scope-nya.

---

## Enable/Disable Auth

```bash
# Di .env:
AUTH_ENFORCE_DEFAULT=true   # Aktifkan auth untuk semua protected groups
AUTH_ENFORCE_DEFAULT=false  # Disable auth (default — untuk migrasi FE)
```

Saat `false`, middleware adalah no-op (`c.Next()` langsung).

---

## Konfigurasi

| Env Var | Required | Deskripsi |
|---------|----------|-----------|
| `JWT_SECRET` | ✅ | Secret key untuk sign/verify JWT (min 32 bytes recommended) |
| `AUTH_ENFORCE_DEFAULT` | Optional | `true`/`false`, default `false` |

---

## File
- [`incoming_auth_handler.go`](../src/interfaces/http/handler/incoming_auth_handler.go)
- [`url_mapping.go`](../src/interfaces/http/url_mapping.go) (middleware mounting)

---

## Menambah Protected Route Baru

1. Register route di group yang sudah protected (`/dashboard`, `/v1/report`, `/v1/int`, `/v1/management`)
2. **Tidak perlu tambahan middleware** — group level sudah handle
3. Jika butuh scope filtering → gunakan `c.Locals("companies"/"adnets")`
4. Jika butuh public endpoint → register di `v1` langsung (lihat `routes/postback.go`)
