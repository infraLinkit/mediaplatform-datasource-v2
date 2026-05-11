# Spec-Driven Development — `mediaplatform-datasource`

> Living document. Update saat penambahan fitur sebelum implementasi.

| Meta | Value |
|------|-------|
| Project | `infraLinkit/mediaplatform-datasource-v2` |
| Module path | `github.com/infraLinkit/mediaplatform-datasource-v2` |
| Stack | Go 1.24, Fiber v2, GORM, PostgreSQL, Redis (go-redis + rueidis), RabbitMQ (rmqp), Google Sheets API, Cobra |
| Architecture | Canonical DDD (`src/{application,domain,infrastructure,interfaces,cmd}`) |
| Last updated | 2026-05-11 |

---

## 1. Business Context

### 1.1 Problem Domain

Mediaplatform-datasource adalah **central data API & ingestion layer** untuk ekosistem media campaign Linkit (DCB / direct carrier billing). Sistem mengelola:

- **Campaign tracking** end-to-end: traffic → landing → click → MO (mobile-originated) → billed
- **Postback handling** dari adnet (advertising network) callback
- **Performance reporting**: CPA, ARPU, traffic, revenue, conversion, alert
- **Budget management**: target budget, spending, IO (insertion order) approval
- **Master data management**: country, operator, partner, service, adnet, agency, channel
- **User & access control**: user, role, menu, permission, company hierarchy

### 1.2 Stakeholders & Value

| Stakeholder | Value yg didapat |
|-------------|------------------|
| Operations team | Real-time campaign monitoring, alert otomatis untuk capping/ratio breach |
| Finance | Budget tracking per IO, agency fee/cost calculation, revenue reconciliation |
| Adnet partners | Postback API untuk report konversi billed/unbilled |
| Account managers | Dashboard summary, top-campaign, performance per agency/channel |
| Admins | Master data management, user/role approval flow |

### 1.3 System Boundaries

```
┌──────────────────────────────────────────────────────────────────┐
│                     mediaplatform-datasource-v2                  │
│  ┌───────────────┐  ┌───────────────┐  ┌───────────┐  ┌──────────┐
│  │ interfaces/   │  │ application/  │  │ domain/   │  │ cmd/     │
│  │ http/handlers │  │ services      │  │ entities  │  │ CLI      │
│  └──────┬────────┘  └──────┬────────┘  └─────┬─────┘  └────┬─────┘
│         │                  │                 │             │
│         └──────────┬───────┴─────────────────┴─────────────┘
│                    │
│  ┌─────────────────┴─────────────────┐
│  │ infrastructure/ (persistence, ext)│
│  └──────┬──────────────────┬─────────┘
└─────────┼──────────────────┼─────────────────────────────────────┘
          │                  │
      ┌───▼───┐          ┌───▼─────┐      ┌──────────┐   ┌─────────┐
      │ PgSQL │          │ Redis   │      │ RabbitMQ │   │ Sheets  │
      └───────┘          └─────────┘      └──────────┘   └─────────┘
```

**External dependencies**: PostgreSQL (state of truth), Redis (cache + JWT blacklist + counters), RabbitMQ (async pipeline ke worker process), Google Sheets API (export billable data), ARPU API (eksternal report aggregator).

### 1.4 Out of Scope

- Authentikasi user issuance (token issued oleh service auth terpisah, datasource cuma validate)
- Frontend (separate repo `cms`)
- Batch worker logic (separate `cores/worker`)
- Landing page rendering (separate `cores/lp`)

---

## 2. Actors

### 2.1 Human Actors

| Actor | Role | Akses |
|-------|------|-------|
| **Admin** | Super admin platform | Full CRUD master data, approve user, manage role/menu |
| **Account Manager (AM)** | Manage campaign per agency | Read campaign monitoring, edit budget/IO/postback, view reports per company scope |
| **Finance** | Reconcile budget & payout | View cost report, edit payout, approve budget IO, export Excel |
| **Operations** | Monitor traffic & alert | View dashboard, alert report, traffic/redirection time, trigger resend data |
| **Read-only User** | View-only access | Subset dashboard + report based on company/adnet/agency assignment |

### 2.2 System Actors

| Actor | Interaksi |
|-------|-----------|
| **Adnet system** (eksternal) | Hit `/v1/postback*` callback saat user billed/unbilled |
| **CMS frontend** | Konsumsi semua endpoint authenticated |
| **Auth service** (eksternal) | Issue JWT (HS256). Datasource validate signature + claims |
| **Worker processes** | Konsume RabbitMQ queue (pixel, click, ratio, campaign management, alert) |
| **Cron scheduler** | Trigger `CronResetCapping` daily |
| **Google Sheets** | Sink untuk billable campaign data export |
| **ARPU API** | Source eksternal untuk data ARPU per service/operator |

### 2.3 Permission Scope

| Resource | Cara filter scope |
|----------|-------------------|
| Company | `user_companies` table — multi-tenant per AM |
| Adnet | `user_adnets` table |
| Agency | `user_agencies` table |
| Service | Via assignment di `user_management` |

JWT claims yg dipakai: `sub` (user_id), `jti` (revocation key), `exp`, `nbf`, `aud`, `iss`, `type` (= "access").

---

## 3. Use Cases

### 3.1 Use Case Map

| ID | Use Case | Actor Utama | Prerequisite |
|----|----------|-------------|--------------|
| UC-01 | Receive postback from adnet | Adnet | Public endpoint, signed URL params |
| UC-02 | Display campaign monitoring summary | AM, Operations | Authenticated, company scope |
| UC-03 | Display traffic/redirection report | Operations | Authenticated |
| UC-04 | Display CPA / cost report | Finance, AM | Authenticated, agency scope |
| UC-05 | Edit campaign capping/ratio/postback | AM, Admin | Authenticated, edit permission |
| UC-06 | Manage master data (country/operator/...) | Admin | Authenticated, admin role |
| UC-07 | Approve user registration | Admin | Authenticated, admin role |
| UC-08 | Create budget IO | AM | Authenticated |
| UC-09 | Approve budget IO | Finance | Authenticated, finance role |
| UC-10 | Export report to Excel | Finance, AM | Authenticated |
| UC-11 | Upload IP range CSV | Admin | Authenticated, admin role |
| UC-12 | Upload campaign Excel batch | AM | Authenticated |
| UC-13 | Resend failed data to RabbitMQ | Operations | Authenticated |
| UC-14 | View dashboard summary | All authenticated | Authenticated |
| UC-15 | Logout (revoke JWT) | All authenticated | Authenticated |
| UC-16 | DB schema migration | DevOps | Run `./datasource migrate` (out-of-band) |

### 3.2 Use Case Detail (sample template)

#### UC-01: Receive Postback from Adnet

- **Goal**: Adnet melaporkan event billed/unbilled untuk record campaign tracking + trigger downstream pixel storage.
- **Trigger**: HTTP GET dari adnet ke `/v1/postback?...` atau `/v1/postback/:urlservicekey/`
- **Main flow**:
  1. Adnet hit endpoint dgn URL params (campaign_id, partner, msisdn, status, dll)
  2. Sistem validasi params (tidak ada auth header — public endpoint)
  3. Sistem cek IP allowlist (kalau diaktifkan)
  4. Insert ke `postbacks` table
  5. Publish message ke RabbitMQ `E_PIXELSTORAGE` exchange
  6. Return HTTP 200 + tracking response body
- **Postcondition**: Record postback tersimpan, message di queue siap dikonsumsi worker
- **Alternate flow**:
  - Validasi gagal → 400 + log warn
  - DB error → 500 + log error, message tetap di-publish kalau memungkinkan
- **Related entities**: `Postback`, `PixelStorage`

#### UC-05: Edit Campaign Capping/Ratio/Postback

- **Goal**: AM update parameter campaign yg sedang berjalan tanpa restart.
- **Trigger**: `POST /v1/management/campaign/editratio` (atau editmocapping, editpo, editcampaign)
- **Main flow**:
  1. AM submit form dengan campaign_id + nilai baru
  2. JWT divalidasi → user_id + scope
  3. Sistem cek user authorize untuk campaign tersebut (via company/agency scope)
  4. Update `campaign_details` table
  5. Publish ke RabbitMQ `E_CAMPAIGNMANAGEMENT` (worker propagasi ke landing page service)
  6. Return updated record
- **Postcondition**: DB updated, downstream service sync
- **Edge cases**: lihat section 7.

---

## 4. API Contract

### 4.1 Base & Versioning

- Base path: `/v1`
- Versi mayor: path-based (`/v1`, `/v2` future)
- Fiber default 100 MB body limit (`BodyLimit: 100*1024*1024`)
- Content-Type: `application/json` (default)
- File upload: `multipart/form-data` (Excel/CSV upload)

### 4.2 Authentication

**Public endpoints** (no auth):
- `GET /v1/postback/:urlservicekey/`
- `GET /v1/postback`
- `GET /v1/postback_billed`
- `GET /v1/inquire/campid`
- `/v1/ext/*` (placeholder)

**Authenticated endpoints**: Header `Authorization: Bearer <JWT>`. Algorithm HS256. Required claims:

```json
{
  "sub": 123,
  "jti": "uuid-...",
  "type": "access",
  "exp": 1700000000,
  "nbf": 1699999999,
  "aud": "mediaplatform",
  "iss": "auth.linkit"
}
```

Failure response:
```json
{ "error": "Invalid token" }
```

### 4.3 Endpoint Groups

| Group | Path | Auth | Owner |
|-------|------|------|-------|
| Public postback | `/v1/postback*`, `/v1/inquire/campid` | None | Adnet integration |
| Dashboard | `/dashboard/*` | Required | All authenticated |
| Reports | `/v1/report/*` | Required | AM, Operations, Finance |
| Internal API | `/v1/int/*` | Required | Backend integration |
| Management | `/v1/management/*` | Required | Admin, AM, Finance |
| External | `/v1/ext/*` | (placeholder) | TBD |

### 4.4 Endpoint Reference (per group)

#### 4.4.1 Reports (`/v1/report`)

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/pinreport` | `DisplayPinReport` | Pin report summary |
| GET | `/datasentapiperformance` | `DisplayPinPerformanceReport` | API performance report |
| GET | `/cpareportlist` | `DisplayCPAReport` | CPA report list |
| GET | `/costreport/:v` | `DisplayCostReport` | Cost report (detail/list mode via `:v`) |
| GET | `/conversionlog` | `DisplayConversionLogReport` | Conversion log per service |
| GET | `/campaign-monitoring-summary` | `DisplayCampaignSummary` | Campaign summary |
| GET | `/campaign-monitoring-summary/chart` | `DisplayCampaignSummaryChart` | Chart data |
| GET | `/alertreport/:v` | `DisplayAlertReportAll` | Alert report |
| GET | `/trafficreport` | `DisplayTrafficReport` | Traffic |
| GET | `/trafficreport/chart` | `GetTrafficReportChart` | Traffic chart |
| GET | `/mainstreamreport` | `DisplayMainstreamReport` | Mainstream |
| GET | `/google-traffic-report` | `DisplayGoogleTrafficReport` | Google traffic |
| GET | `/budgetmonitoring` | `DisplayBudgetMonitoring` | Budget monitoring |
| GET | `/performance-report` | `DisplayPerformanceReport` | Performance |
| GET | `/revenuemonitoring` | `DisplayRevenueMonitoring` | Revenue list |
| GET | `/revenuemonitoring/chart` | `DisplayRevenueMonitoringChart` | Revenue chart |
| GET | `/redirectiontime` | `DisplayRedirectionTime` | Redirection time |
| GET | `/ioreport` | `DisplaySummaryBudgetIO` | IO report |
| GET | `/campaign-spending-channel` | `DisplayCampaignSpendingChannel` | Spending per channel |
| GET | `/campaign-spending-channel/country-children` | `DisplayCampaignSpendingChannelCountryChildren` | Drill-down |
| POST | `/resend-data` | `ResendData` | Resend to RabbitMQ |
| POST | `/resend-data-apireport` | `ResendDataAPIReport` | Resend api report |
| POST | `/campaign-monitoring-summary/edit-target-budget` | `EditTargetBudget` | Edit target budget |

#### 4.4.2 Management (`/v1/management`)

Sub-groups: `/campaign`, `/campaign-setting`, `/menu`, `/role`, `/user`, `/userlog`, `/budget-io`, `/country-service`, `/ipranges`. Total ~120 endpoints. Lihat [`src/app/routes/management.go`](src/app/routes/management.go) untuk lengkap.

Pattern CRUD standar:
- `GET /resource` → list
- `GET /resource/:id` → detail
- `POST /resource` → create
- `PUT /resource/:id` → update
- `DELETE /resource/:id` → soft delete

### 4.5 Standard Response Envelope

Convention saat ini (mostly):
```json
{
  "code": 200,
  "desc": "OK",
  "data": { }
}
```

Error:
```json
{
  "code": 400,
  "desc": "Bad Request",
  "error": "validation message"
}
```

---

## 5. Data Model

### 5.1 Domain Aggregates (high-level)

| Aggregate Root | Entities | File |
|----------------|----------|------|
| **Campaign** | Campaign, CampaignDetail, IncSummaryCampaign, IncSummaryCampaignHour, SummaryCampaign, SummaryCampaignBilling | `entity/campaign_detail.go`, `entity/inc_summary_campaigns.go`, `entity/campaignsummary.go` |
| **Tracking** | DataTraffic, DataLanding, DataClicked, DataRedirect, MO, PixelStorage, ClickStorage, Postback | `entity/traffic.go`, `entity/postback.go`, `entity/dashboard.go` |
| **Reporting** | ApiPinReport, ApiPinPerformance, SummaryMo, SummaryCr, SummaryCapping, SummaryRatio, SummaryLanding, SummaryTraffic, SummaryDashboard | `entity/apireport.go`, `entity/dashboard.go` |
| **Budget** | TargetBudget, TargetBudgetDetail, BudgetIO, SummaryBudgetIO | `entity/target_budget.go`, `entity/budgetio.go` |
| **Master Data** | Country, Continent, Company, CompanyGroup, Domain, DomainService, Operator, Partner, Service, AdnetList, Agency, Channel, MainstreamGroup, OperatorAlias, IPRange, IPRangeCsvRow, LpDesignType | `entity/table.go` |
| **User & Access** | User, DetailUser, UserCompany, UserAdnet, Role, Permission, Menu, CcEmail, Email, HistoryCappingKey | `entity/usermanagement.go`, `entity/rolemanagement.go`, `entity/userlog.go` |

### 5.2 Sample Schema (TargetBudget)

```go
type TargetBudget struct {
    ID                uint           `gorm:"primaryKey"`
    Date              time.Time
    Service           string
    Country           string
    Operator          string
    Adnet             string
    DailyBudget       float64
    SpendingToAdnets  float64
    AgencyFee         float64
    CostPerConversion float64
    CreatedAt         time.Time
    UpdatedAt         time.Time
    DeletedAt         gorm.DeletedAt `gorm:"index"`
}
```

### 5.3 Migration Strategy

- **Authoritative**: AutoMigrate via `./datasource migrate` sub-command (run-once container)
- **Entity registration**: `src/cmd/migrate.go` — `migrateEntities` slice
- **Add new entity**: 1) buat di `src/domain/entity/`, 2) tambah ke `migrateEntities`, 3) jalankan `./datasource migrate`
- **Production**: Wajib via init-container atau CI step. Jangan rely on app startup.

### 5.4 Cache Layer (Redis)

| Key pattern | TTL | Purpose |
|-------------|-----|---------|
| `jwt:blacklist:<jti>` | sisa exp | JWT revocation list |
| Pixel cache key (DB index `RedisCachePixel`) | `RedisKeyExpiration` (default 600s) | Pixel postback dedup |
| RedisJSON keys (DB index `RedisDBIndex`) | varies | IP range, capping counter |
| `HistoryCappingKey` related | reset oleh cron `CronResetCapping` | Daily capping reset |

---

## 6. Constraints

### 6.1 Functional Constraints

| ID | Constraint | Rationale | Source |
|----|------------|-----------|--------|
| FC-01 | JWT wajib HS256 | Hindari algorithm confusion (CWE-327) | `incoming_auth_handler.go` |
| FC-02 | JWT `type` claim wajib `"access"` (refresh token reject) | Pisahkan access vs refresh | `incoming_auth_handler.go` |
| FC-03 | `aud` & `iss` validation aktif kalau env diset | Multi-tenant token isolation | `incoming_auth_handler.go` |
| FC-04 | TLS verify aktif default, opt-in disable hanya non-prod | CWE-295 fix | `helper/http.go` |
| FC-05 | AES key min 16/24/32 byte (AES-128/192/256) | Crypto correctness | `helper/utils.go` |
| FC-06 | Postback endpoint public — auth via signed params/IP allowlist | Adnet callback gak bisa kirim Bearer | `routes/postback.go` |
| FC-07 | Body limit max 100 MB | Excel upload | `app/url_mapping.go` |
| FC-08 | Group-level `AuthMiddleware` di `/dashboard`, `/v1/{report,int,management}` | Default-secure | gated `AUTH_ENFORCE_DEFAULT` |

### 6.2 Non-Functional Constraints

| ID | Constraint | Target |
|----|------------|--------|
| NFC-01 | DB pool max open conn | `DB_MAX_OPEN_CONNS × replicas ≤ PG max_connections` |
| NFC-02 | DB connection lifetime | 30 min default (`DB_CONN_MAX_LIFETIME_MIN`) |
| NFC-03 | Redis fail-fast vs degraded | `REDIS_REQUIRED=true` (default) → fail. `false` → continue tanpa cache |
| NFC-04 | RabbitMQ reconnect | Background goroutine watch `NotifyClose`, retry exp backoff |
| NFC-05 | HTTP outbound timeout | Per-call via `helper.PHttp` |
| NFC-06 | Server should not run schema migration | Migration via dedicated `migrate` sub-cmd / init container |
| NFC-07 | Logs redact secrets | `Cfg.Redacted()` mask password/private key |
| NFC-08 | Build verifiable di Go 1.24.1 | Dockerfile pin |

### 6.3 Operational

| ID | Constraint |
|----|------------|
| OC-01 | Multi-replica deploy require: migrate run-once → server start |
| OC-02 | Tag release `v*.*.*` → CI auto-build 2 image (server + migrate) |
| OC-03 | `.env` tidak boleh di-commit (sudah di `.gitignore`) |
| OC-04 | `AES_SECRET_KEY` wajib via secret manager di prod |
| OC-05 | DB migration tool eksternal (`golang-migrate`) recommended jangka panjang vs AutoMigrate |

---

## 7. Edge Cases

### 7.1 Auth Edge Cases

| Case | Expected Behavior |
|------|-------------------|
| Missing `Authorization` header | 401 `{"error":"Missing or invalid token"}` |
| Header format selain `Bearer <token>` | 401 |
| `JWT_SECRET` env not set | 500 `{"error":"JWT secret not configured"}` |
| Token expired | 401 `{"error":"Token expired"}` |
| Token `nbf` di future | 401 `{"error":"Token not yet active"}` (30s leeway) |
| `aud` mismatch (kalau env diset) | 401 |
| `iss` mismatch (kalau env diset) | 401 |
| `type` ≠ "access" | 401 |
| `jti` di Redis blacklist | 401 `{"error":"Token revoked"}` |
| Redis nil saat blacklist check | Skip check (degraded mode) |
| `sub` claim invalid format | 401 |

### 7.2 Postback Edge Cases

| Case | Expected Behavior |
|------|-------------------|
| Adnet kirim duplicate postback | Idempotent: cache pixel via Redis + dedup di worker |
| MSISDN format tidak valid | 400 + log |
| Campaign ID tidak ditemukan | 200 dengan flag, **bukan** 404 (adnet jangan retry) |
| Postback masuk saat budget capping tercapai | Tetap record, mark capped |

### 7.3 HTTP Helper Edge Cases (post-fix #2)

| Case | Behavior |
|------|----------|
| `httpClient.Do(req)` return error | Return wrapped error, **tidak** akses `response.Status` |
| `response == nil` (lib quirk) | Explicit nil-check, return error |
| Non-OK status (≠200) | Return real status code + error wrapped |

### 7.4 DB / Pool Edge Cases

| Case | Behavior |
|------|----------|
| PG `max_connections` saturated | Pool throttle (block sampai available) |
| Connection lifetime exceeded | GORM auto-recreate connection |
| AutoMigrate timeout | `log.Fatalf` (exit 1) — operator harus investigate |

### 7.5 Redis Edge Cases

| Case | Behavior |
|------|----------|
| Redis down saat startup, `REDIS_REQUIRED=true` | 5x retry exp backoff → fail |
| Redis down, `REDIS_REQUIRED=false` | Log warn, app continue dgn `nil` clients, `Setup.RedisAvailable=false` |
| Redis down saat runtime | Cache miss → fallback DB (jika handler pakai `helper.SafeRedisGet`) |

### 7.6 RabbitMQ Edge Cases

| Case | Behavior |
|------|----------|
| Initial connection gagal | Lib `SetupConnectionAmqpAndReconnect` retry 5x → return error → `log.Fatalf` |
| Channel setup gagal | `log.Fatalf` |
| Connection drop saat runtime | Goroutine watcher detect via `NotifyClose`, reconnect + re-setup channel |
| Publisher error saat queue full | Bug: saat ini ignore — TODO future improvement |

### 7.7 Concurrency

| Case | Behavior |
|------|----------|
| Concurrent edit campaign by 2 AM | Last-write-wins (no optimistic lock) — TODO future |
| Cron capping reset overlap | Single-fire (cron runs in single instance) |

---

## 8. Acceptance Criteria

### 8.1 Per Use Case

#### UC-01 (Postback)
- Adnet GET `/v1/postback?campaign_id=X&...` → response 200 dlm < 500ms
- Record postback insert ke `postbacks` table dgn semua param
- Message ter-publish ke `E_PIXELSTORAGE` exchange
- Duplicate postback dlm 600s (Redis TTL) tidak ter-record dua kali
- Postback gagal validate → 400 + log warn dgn request payload

#### UC-05 (Edit campaign)
- AM dgn scope agency `X` bisa edit campaign agency `X`
- AM dgn scope agency `X` **tidak bisa** edit campaign agency `Y` → 403
- Update tersimpan ke `campaign_details` + audit log (kalau ada)
- Message ter-publish ke `E_CAMPAIGNMANAGEMENT`
- Concurrent edit oleh 2 AM → last-write-wins (sementara, TODO lock)

#### UC-15 (Logout)
- Logout endpoint call `h.RevokeJWT(jti, ttl)` dgn ttl = exp - now
- Token yg sama di-pakai → 401 `Token revoked`
- Setelah ttl expire, key blacklist auto-cleanup oleh Redis

### 8.2 Cross-Cutting (Definition of Done)

Untuk setiap fitur baru:

- Endpoint terdaftar di `src/app/routes/<domain>.go`
- Handler di `src/handler/incoming_<domain>_handler.go`
- Repository (kalau ada DB access) di `src/domain/repository/<domain>.go`
- Entity (kalau ada DB schema baru) di `src/domain/entity/<domain>.go` + register di `migrateEntities`
- JWT auth applied (kalau bukan public)
- Scope filter (company/agency/adnet) sesuai actor
- Error response konsisten dgn envelope `{code, desc, error}`
- Response time < 1 detik untuk read endpoint, < 3 detik untuk export
- Logs structured (logrus) — no println bocor
- Secrets tidak di-log (gunakan `Cfg.Redacted()` pattern)
- `staticcheck ./...` zero new warning
- `go build ./...` exit 0
- Spec doc updated: section 3 (use case), 4 (API contract), 5 (entity), 8 (AC)

### 8.3 Non-Functional

- DB query log enabled (GORM logger Info mode) — no N+1 baru
- Endpoint baru di-monitor di Grafana (HTTP latency, error rate)
- Postback endpoint stress test ≥ 500 req/sec sustained
- Redis fail simulation: degraded mode tidak crash (kalau `REDIS_REQUIRED=false`)
- DB pool sat test: tidak melebihi `DB_MAX_OPEN_CONNS`

---

## 9. Adding New Feature — Quick Workflow

```
1. Identify domain (campaign/budget/master-data/...)
2. Update SPEC.md:
   - Section 3: Use case
   - Section 4.4: Endpoint definition
   - Section 5: Entity (kalau perlu)
   - Section 7: Edge case
   - Section 8.1: Acceptance criteria

3. Implement:
   - entity → src/domain/entity/<file>.go (+ register migrateEntities)
   - repo   → src/domain/repository/<file>.go
   - handler → src/handler/incoming_<file>_handler.go
   - route → src/app/routes/<group>.go
   - DTO/params → src/domain/entity/<file>.go (request/response struct)

4. Test:
   - go build ./...
   - staticcheck ./...
   - manual hit endpoint
   - cek log untuk leak secret

5. Run migration:
   ./datasource migrate

6. Deploy:
   - tag git v1.x.y
   - CI auto build
   - update compose version
   - run migrate compose first
   - deploy server
```

---

## Appendix A — Env Vars Reference

Lihat `AUDIT_FIX.md` section "Env Vars Summary".

## Appendix B — Related Docs

- [`AUDIT_FIX.md`](./AUDIT_FIX.md) — Audit fix implementation log (8 priority items)
- [`UNUSED_FUNCTIONS.md`](./UNUSED_FUNCTIONS.md) — Dead code report
- `Data Source.pdf` — Original audit report (2026-02-11)
- `README.md` / `NOTES.md` — Repo specifics
