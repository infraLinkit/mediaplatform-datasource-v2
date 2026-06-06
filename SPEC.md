# Spec-Driven Development — `mediaplatform-datasource-v2`

> Living document. Update saat penambahan fitur **sebelum** implementasi.

| Meta | Value |
|------|-------|
| Project | `infraLinkit/mediaplatform-datasource-v2` |
| Module path | `github.com/infraLinkit/mediaplatform-datasource-v2` |
| Stack | Go 1.24.1, Fiber v2, GORM, PostgreSQL, Redis (go-redis + rueidis), RabbitMQ (`wiliehidayat87/rmqp` + `RabbitManager`), Google Sheets API, Cobra |
| Architecture | Canonical DDD (`src/{application,domain,infrastructure,interfaces,cmd}`) |
| Last updated | 2026-06-05 |
| UC Docs | [`docs/`](./docs/README.md) |

---

## 1. Business Context

### 1.1 Problem Domain

Mediaplatform-datasource adalah **central data API & ingestion layer** untuk ekosistem media campaign Linkit (DCB / direct carrier billing). Sistem mengelola:

- **Campaign tracking** end-to-end: traffic → landing → click → MO (mobile-originated) → billed
- **Postback handling** dari adnet (advertising network) callback — async (RTO) dan sync/RPC (RTD)
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
│  ┌─────────────────┴──────────────────────────────┐
│  │ infrastructure/ (persistence, messaging, ext)  │
│  └──────┬─────────────────────┬───────────────────┘
└─────────┼─────────────────────┼────────────────────────────────────┘
          │                     │
      ┌───▼───┐  ┌─────────┐  ┌─▼────────────────┐   ┌─────────┐
      │ PgSQL │  │ Redis   │  │ RabbitMQ          │   │ Sheets  │
      └───────┘  └─────────┘  │ (rmqp + RM pool) │   └─────────┘
                               └──────────────────┘
```

**External dependencies**: PostgreSQL (state of truth), Redis (cache + JWT blacklist + counters), RabbitMQ (async pipeline ke worker — via `RabbitManager` channel pool), Google Sheets API (export billable data), ARPU API (eksternal report aggregator).

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

JWT claims yg dipakai: `sub` (user_id), `jti` (revocation key), `exp`, `nbf`, `aud`, `iss`, `type` (= "access"), `companies` ([]string), `adnets` ([]string).

---

## 3. Use Cases

> **Konvensi**: UC ID di SPEC.md **identik** dengan file di `docs/`. Granular sub-use-case dilist sebagai sub-item dalam UC induk.  
> UC-01 s/d UC-09 masing-masing punya file lengkap di [`docs/`](./docs/README.md).

### 3.1 Use Case Map

| UC | Nama Fitur | Actor Utama | Route Group | Doc |
|----|-----------|-------------|-------------|-----|
| **UC-01** | **Postback Intake** | Adnet | `/v1/` (Public) | [UC-01-postback.md](./docs/UC-01-postback.md) |
| UC-01.1 | Receive postback async (RTO) | Adnet | `GET /v1/postback*` | ↑ |
| UC-01.2 | Receive postback sync / Direct Reply-To (RTD) | Adnet | `GET /v1/postback_sync` | ↑ |
| UC-01.3 | Inquiry campaign ID | Adnet / System | `GET /v1/inquire/*` | ↑ |
| **UC-02** | **Campaign Management** | AM, Admin | `/v1/management/campaign` (Auth) | [UC-02-campaign-management.md](./docs/UC-02-campaign-management.md) |
| UC-02.1 | Display / filter campaign list | AM | `GET /campaign` | ↑ |
| UC-02.2 | Edit capping / ratio / PO / MO | AM, Admin | `POST /campaign/edit*` | ↑ |
| UC-02.3 | Update status active/inactive | AM, Admin | `POST /campaign/updatestatus` | ↑ |
| UC-02.4 | Delete campaign + cleanup Redis | Admin | `POST /campaign/delcampaign` | ↑ |
| UC-02.5 | Send campaign data ke worker via RabbitMQ | AM | `POST /campaign/send` | ↑ |
| UC-02.6 | Update Google Sheet / billable | AM | `POST /campaign/updategooglesheet*` | ↑ |
| **UC-03** | **Reporting** | AM, Finance, Operations | `/v1/report` (Auth) | [UC-03-reporting.md](./docs/UC-03-reporting.md) |
| UC-03.1 | Campaign monitoring summary + chart | AM, Ops | `GET /report/campaign-monitoring-summary*` | ↑ |
| UC-03.2 | CPA report | Finance, AM | `GET /report/cpareportlist` | ↑ |
| UC-03.3 | Cost report (list + detail) | Finance | `GET /report/costreport/:v` | ↑ |
| UC-03.4 | Traffic report + chart | Ops | `GET /report/trafficreport*` | ↑ |
| UC-03.5 | Revenue monitoring + chart | Finance, AM | `GET /report/revenuemonitoring*` | ↑ |
| UC-03.6 | Mainstream report | AM, Ops | `GET /report/mainstreamreport` | ↑ |
| UC-03.7 | Alert report | Ops | `GET /report/alertreport/:v` | ↑ |
| UC-03.8 | Performance report | AM | `GET /report/performance-report` | ↑ |
| UC-03.9 | IO report + update | Finance | `GET/PUT /report/ioreport*` | ↑ |
| UC-03.10 | Budget monitoring | Finance, AM | `GET /report/budgetmonitoring` | ↑ |
| UC-03.11 | Redirection time report | Ops | `GET /report/redirectiontime` | ↑ |
| UC-03.12 | Pin report / API performance | AM, Ops | `GET /report/pinreport*` | ↑ |
| UC-03.13 | Google traffic report | Ops | `GET /report/google-traffic-report` | ↑ |
| UC-03.14 | Campaign spending per channel + drill-down | AM, Finance | `GET /report/campaign-spending-channel*` | ↑ |
| UC-03.15 | Resend summary data ke Linkit Dashboard | Ops | `POST /report/resend-data` | ↑ |
| UC-03.16 | Resend API report ke Linkit Dashboard | Ops | `POST /report/resend-data-apireport` | ↑ |
| UC-03.17 | Edit target budget (single + batch) | AM | `POST /report/campaign-monitoring-summary/edit-target-budget*` | ↑ |
| UC-03.18 | Export CPA / cost report ke Excel | Finance, AM | Internal export endpoints | ↑ |
| UC-03.19 | Conversion log report | Ops, AM | `GET /report/conversionlog` | ↑ |
| **UC-04** | **User, Role & Menu Management** | Admin | `/v1/management/{user,role,menu,userlog}` (Auth) | [UC-04-user-role-menu.md](./docs/UC-04-user-role-menu.md) |
| UC-04.1 | CRUD user (create, list, update, delete) | Admin | `* /management/user*` | ↑ |
| UC-04.2 | Approve user registration | Admin | `PUT /management/user/approveuser/:id` | ↑ |
| UC-04.3 | Assign service & adnet ke user | Admin | `PUT /management/user/assignservice/:id` | ↑ |
| UC-04.4 | CRUD role + permission | Admin | `* /management/role*` | ↑ |
| UC-04.5 | CRUD menu navigasi | Admin | `* /management/menu*` | ↑ |
| UC-04.6 | View user activity log | Admin | `GET /management/userlog*` | ↑ |
| UC-04.7 | Logout / revoke JWT | All authenticated | (Auth handler) | ↑ + [UC-09](./docs/UC-09-auth.md) |
| **UC-05** | **Country, Service & Catalog Master Data** | Admin | `/v1/management/country-service` (Auth) | [UC-05-country-service-catalog.md](./docs/UC-05-country-service-catalog.md) |
| UC-05.1 | CRUD country + continent | Admin | `* /country-service/country*` | ↑ |
| UC-05.2 | CRUD company + company group | Admin | `* /country-service/company*` | ↑ |
| UC-05.3 | CRUD domain + domain service | Admin | `* /country-service/domain*` | ↑ |
| UC-05.4 | CRUD operator | Admin | `* /country-service/operator*` | ↑ |
| UC-05.5 | CRUD partner | Admin | `* /country-service/partner*` | ↑ |
| UC-05.6 | CRUD service | Admin | `* /country-service/service*` | ↑ |
| UC-05.7 | CRUD adnet list + toggle DSP status | Admin | `* /country-service/adnet-list*` | ↑ |
| UC-05.8 | CRUD email | Admin | `* /country-service/email*` | ↑ |
| UC-05.9 | CRUD agency | Admin | `* /country-service/agency*` | ↑ |
| UC-05.10 | CRUD channel | Admin | `* /country-service/channel*` | ↑ |
| UC-05.11 | CRUD mainstream group | Admin | `* /country-service/mainstream-group*` | ↑ |
| **UC-06** | **Internal Endpoints & Admin Tools** | System, AM | `/v1/int` (Auth) | [UC-06-internal-endpoints.md](./docs/UC-06-internal-endpoints.md) |
| UC-06.1 | Set target daily budget | System | `PUT /int/setdata/:v/` | ↑ |
| UC-06.2 | Update agency fee & cost conversion | System | `PUT /int/updatedata/:v/` | ↑ |
| UC-06.3 | Update ratio / postback / agency cost | System | `PUT /int/update*` | ↑ |
| UC-06.4 | Upload / upsert Excel SMS campaign | AM | `POST /int/uploadexcel`, `PUT /int/upsertexcel/` | ↑ |
| UC-06.5 | Pin report / performance transaksional | System | `GET /int/datasentapi*` | ↑ |
| UC-06.6 | Edit payout / CPA / ARPU di API report | Finance, AM | `POST /int/pin*` | ↑ |
| UC-06.7 | Export CPA / cost / cost-detail | Finance | `GET /int/export*` | ↑ |
| UC-06.8 | Fetch ARPU data dari Linkit Dashboard API | System | `GET /int/getdataarpu/` | ↑ |
| UC-06.9 | Get / update URL service dalam summary landing | System | `GET/PUT /int/*summarylanding*` | ↑ |
| **UC-07** | **Budget IO & IP Range Management** | AM, Finance, Admin | `/v1/management/{budget-io,ipranges}` (Auth) | [UC-07-budget-io-iprange.md](./docs/UC-07-budget-io-iprange.md) |
| UC-07.1 | Buat Budget IO | AM | `POST /management/budget-io/` | ↑ |
| UC-07.2 | List + filter Budget IO | AM, Finance | `GET /management/budget-io/*` | ↑ |
| UC-07.3 | Approve Budget IO | Finance | (via list page) | ↑ |
| UC-07.4 | Upload IP range CSV | Admin | `POST /management/ipranges/upload` | ↑ |
| UC-07.5 | Implementasi / aktifkan IP range | Admin | `POST /management/ipranges/implement` | ↑ |
| UC-07.6 | Download IP range CSV | Admin | `POST /management/ipranges/download` | ↑ |
| **UC-08** | **RabbitMQ Messaging (RabbitManager)** | System | Infrastructure / boot-time | [UC-08-rabbitmq-messaging.md](./docs/UC-08-rabbitmq-messaging.md) |
| UC-08.1 | Async publish — fire-and-forget (RTO) | System | `h.RM.PublishWithRetry` | ↑ |
| UC-08.2 | Sync RPC publish — Direct Reply-To (RTD) | System | `h.RM.DirectReplyToWithRetry` | ↑ |
| UC-08.3 | Auto-reconnect + channel pool management | System | `handleReconnect` goroutine | ↑ |
| **UC-09** | **Authentication & Authorization** | All | All protected groups (Auth) | [UC-09-auth.md](./docs/UC-09-auth.md) |
| UC-09.1 | Login → issue JWT | All | (Auth service) | ↑ |
| UC-09.2 | Validate JWT per request | All | `AuthMiddleware` | ↑ |
| UC-09.3 | Logout → blacklist JWT di Redis | All authenticated | (Auth handler) | ↑ |
| UC-09.4 | Scope filtering (companies, adnets) | AM, Finance | `c.Locals("companies"/"adnets")` | ↑ |
| UC-09.5 | DB schema migration (out-of-band) | DevOps | `./datasource migrate` CLI | — |

### 3.2 Use Case Detail

#### UC-01: Postback Intake

**Ref lengkap**: [UC-01-postback.md](./docs/UC-01-postback.md)

##### UC-01.1 — Async Postback (RTO)

- **Goal**: Adnet melaporkan event billed/unbilled → trigger downstream ratio worker.
- **Trigger**: `GET /v1/postback/:urlservicekey/` atau `GET /v1/postback`
- **Main flow**:
  1. Parse URL params (urlservicekey, aff_sub/pixel, method, adnet, dll)
  2. Check cookie dedup — 403 jika duplicate dalam 3 detik window
  3. Set dedup cookie (3 detik, HttpOnly, SameSite=lax)
  4. Lookup campaign config dari Redis (`{urlservicekey}-configIdx`) — 404 jika tidak ada
  5. Lookup pixel dari Redis (`pixelKey`) — 404 jika tidak ada
  6. Cek `px.IsUsed` — 200 NOK jika sudah dipakai
  7. `corId = "RTO" + GetUniqId()` → `h.RM.PublishWithRetry(ctx, Ratio exchange, queue, body, corId)`
  8. Return 200 + `PixelStorageRsp`
- **Postcondition**: Worker proses async, tidak reply
- **Related entities**: `PixelStorage`, `CampaignDetail`

##### UC-01.2 — Sync Postback Direct Reply-To (RTD)

- **Goal**: Adnet hit postback dan sistem **menunggu** reply SHAVED/NOTSHAVED dari worker.
- **Trigger**: `GET /v1/postback_sync`
- **Main flow**:
  1. Langkah 1–6 sama dengan UC-01.1
  2. `corId = "RTD" + GetUniqId()` → `h.RM.DirectReplyToWithRetry(ctx, exchange, queue, body, corId)`
  3. Sistem block menunggu reply di `amq.rabbitmq.reply-to`
  4. Worker publish reply ke `d.ReplyTo` dengan `CorrelationId = d.CorrelationId`
  5. Reply diterima: `"SHAVED"` atau `"NOTSHAVED"` → handler return sesuai reply
- **Edge cases**: Worker tidak reply → context timeout → error. Lihat §8.3.

#### UC-02: Campaign Management

**Ref lengkap**: [UC-02-campaign-management.md](./docs/UC-02-campaign-management.md)

##### UC-02.2 — Edit Campaign Parameter

- **Goal**: AM update capping/ratio/PO campaign yang sedang berjalan tanpa restart service.
- **Trigger**: `POST /v1/management/campaign/editratio` (atau `editmocapping`, `editpo`, `editcampaign`)
- **Main flow**:
  1. JWT validated → user scope (company/agency)
  2. Scope check — 403 jika campaign bukan milik company AM
  3. Update Redis config + DB `campaign_details`
  4. Return updated record
- **Postcondition**: Redis + DB updated; perubahan berlaku real-time

##### UC-02.5 — Send Campaign Data to Worker

- **Goal**: AM kirim data kampanye ke worker via RabbitMQ.
- **Trigger**: `POST /v1/management/campaign/send`
- **Main flow**:
  1. Parse body JSON (`campaignData`)
  2. `h.RM.PublishWithRetry(ctx, E_CAMPAIGNMANAGEMENT, Q_CAMPAIGNMANAGEMENT, body, "")`
  3. Return 200 `{"message": "Campaign data sent to RabbitMQ"}`
- **Correlation ID**: kosong (`""`) — fire-and-forget, no reply

#### UC-03: Reporting

**Ref lengkap**: [UC-03-reporting.md](./docs/UC-03-reporting.md)

Semua endpoint reporting auth-gated. Scope filtering via JWT claims `companies` + `adnets`.

##### UC-03.15 — Resend Data ke Linkit Dashboard

- **Goal**: Operations resend data summary yang gagal terkirim ke Linkit Dashboard.
- **Trigger**: `POST /v1/report/resend-data`
- **Main flow**:
  1. Parse form: `total`, `id[0..N]`
  2. `GetSummaryReportById(ids)` → `[]SummaryCampaign`
  3. Untuk setiap record: build URL + `h.RM.PublishWithRetry(ctx, "E_RESENDCAMPAIGNDATA", "Q_RESENDCAMPAIGNDATA", msg, "")`
  4. Return `{status: "OK"/"NOK", error: ""}`

#### UC-04: User, Role & Menu Management

**Ref lengkap**: [UC-04-user-role-menu.md](./docs/UC-04-user-role-menu.md)

Standard CRUD via REST. Admin only. Approval flow untuk registrasi user baru.

#### UC-05: Country, Service & Catalog

**Ref lengkap**: [UC-05-country-service-catalog.md](./docs/UC-05-country-service-catalog.md)

Master data yang menjadi referensi saat create/edit campaign. ~50+ endpoints CRUD.

#### UC-06: Internal Endpoints & Admin Tools

**Ref lengkap**: [UC-06-internal-endpoints.md](./docs/UC-06-internal-endpoints.md)

Endpoint internal untuk mutasi data transaksional dan tools admin. Auth required, tidak diakses public.

#### UC-07: Budget IO & IP Range

**Ref lengkap**: [UC-07-budget-io-iprange.md](./docs/UC-07-budget-io-iprange.md)

Budget IO: buat → approval → monitoring. IP Range: upload CSV → implement → download.

#### UC-08: RabbitMQ Messaging (RabbitManager)

**Ref lengkap**: [UC-08-rabbitmq-messaging.md](./docs/UC-08-rabbitmq-messaging.md)

Infrastructure layer untuk semua publish ke RabbitMQ. Diinisialisasi saat boot di `server.go`.

```go
// Async (RTO) — fire-and-forget
corId := "RTO" + external.GetUniqId(h.Config.TZ)
h.RM.PublishWithRetry(ctx, exchange, queue, body, corId)

// Sync RPC (RTD) — tunggu reply worker
corId := "RTD" + external.GetUniqId(h.Config.TZ)
reply, err := h.RM.DirectReplyToWithRetry(ctx, exchange, queue, body, corId)
```

#### UC-09: Authentication & Authorization

**Ref lengkap**: [UC-09-auth.md](./docs/UC-09-auth.md)

JWT HS256. `AUTH_ENFORCE_DEFAULT=true` mengaktifkan group-level `AuthMiddleware`. Logout via Redis blacklist `jwt:blacklist:{jti}`.

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
- `GET /v1/postback_sync`
- `GET /v1/inquire/campid`
- `GET /v1/inquire/api-campid`
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
  "iss": "auth.linkit",
  "companies": ["company_a"],
  "adnets": ["adnet_x"]
}
```

Failure response:
```json
{ "error": "Invalid token" }
```

### 4.3 Endpoint Groups

| Group | Path | Auth | Owner |
|-------|------|------|-------|
| Public postback | `/v1/postback*`, `/v1/inquire/*` | None | Adnet integration |
| Dashboard | `/dashboard/*` | Required | All authenticated |
| Reports | `/v1/report/*` | Required | AM, Operations, Finance |
| Internal API | `/v1/int/*` | Required | Backend integration |
| Management | `/v1/management/*` | Required | Admin, AM, Finance |
| External | `/v1/ext/*` | (placeholder) | TBD |

### 4.4 Endpoint Reference (per group)

#### 4.4.1 Postback (`/v1` — Public)

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/postback/:urlservicekey/` | `Postback` | Postback v1 via path param |
| GET | `/postback` | `PostbackV3` | Postback v3 via query params |
| GET | `/postback_billed` | `PostbackBilled` | Postback billed event |
| GET | `/postback_sync` | `PostbackDirectReply` | **[NEW]** Postback sync — Direct Reply-To (RTD) |
| GET | `/inquire/campid` | `InquiryCampID` | Inquiry campaign ID |
| GET | `/inquire/api-campid` | `InquiryAPICampID` | Inquiry campaign ID (API channel) |

#### 4.4.2 Reports (`/v1/report`)

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
| GET | `/defaultinput/` | `DisplayDefaultInput` | Default input values |
| GET | `/redirectiontime` | `DisplayRedirectionTime` | Redirection time |
| GET | `/ioreport` | `DisplaySummaryBudgetIO` | IO report |
| PUT | `/ioreport/update` | `UpdateSummaryBudgetIO` | Update IO report |
| GET | `/campaign-spending-channel` | `DisplayCampaignSpendingChannel` | Spending per channel |
| GET | `/campaign-spending-channel/country-children` | `DisplayCampaignSpendingChannelCountryChildren` | Drill-down |
| POST | `/resend-data` | `ResendData` | Resend summary ke Linkit Dashboard via RabbitMQ |
| POST | `/resend-data-apireport` | `ResendDataAPIReport` | Resend API report ke Linkit Dashboard |
| POST | `/campaign-monitoring-summary/edit-target-budget` | `EditTargetBudgetLevel` | Edit target budget |
| POST | `/campaign-monitoring-summary/edit-target-budget-batch` | `EditTargetBudgetBatch` | **[NEW]** Bulk edit target budget |

#### 4.4.3 Management (`/v1/management`)

Sub-groups: `/campaign`, `/campaign-setting`, `/menu`, `/role`, `/user`, `/userlog`, `/budget-io`, `/country-service`, `/ipranges`. Total ~120 endpoints. Lihat [`src/interfaces/http/routes/management.go`](src/interfaces/http/routes/management.go) untuk lengkap.

Pattern CRUD standar:
- `GET /resource` → list
- `GET /resource/:id` → detail
- `POST /resource` → create
- `PUT /resource/:id` → update
- `DELETE /resource/:id` → soft delete

Key endpoints:

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| POST | `/campaign/send` | `SendCampaignHandler` | Kirim data ke worker via RabbitMQ (`h.RM.PublishWithRetry`) |
| POST | `/campaign/editratio` | `EditCampaignRatio` | Edit ratio send/receive |
| POST | `/campaign/editmocapping` | `EditCampaignMOCapping` | Edit MO capping |
| POST | `/campaign/editpo` | `EditCampaignPO` | Edit PO + recalculate CPA |
| POST | `/campaign/updatestatus` | `UpdateStatusCampaign` | Toggle active/inactive |

#### 4.4.4 Internal (`/v1/int`)

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| PUT | `/setdata/:v/` | `SetData` | Set target daily budget |
| PUT | `/updatedata/:v/` | `UpdateAgencyFeeAndCostConversion` | Update agency fee & cost |
| PUT | `/updateratio/:v/` | `UpdateRatio` | Update ratio |
| PUT | `/updatepostback/:v/` | `UpdatePostback` | Update postback |
| POST | `/uploadexcel` | `UploadExcel` | Upload Excel SMS campaign |
| PUT | `/upsertexcel/` | `UpsertExcel` | Upsert Excel SMS campaign |
| GET | `/getdataarpu/` | `GetDataArpu` | Fetch ARPU data |
| GET | `/get_urlservice_in_summarylanding` | `GetURLServiceInSummaryLanding` | Get URL service total load time |
| PUT | `/update_response_url_service_in_summarylanding` | `UpdateResponseURLServiceInSummaryLanding` | Update URL service response |

### 4.5 Standard Response Envelope

```json
// Success
{
  "code": 200,
  "desc": "OK",
  "data": { }
}

// Success with DataTable
{
  "draw": 1,
  "code": 200,
  "desc": "OK",
  "data": [...],
  "recordsTotal": 100,
  "recordsFiltered": 100
}

// Error
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
| **Reporting** | ApiPinReport, ApiPinPerformance, SummaryMo, SummaryCr, SummaryCapping, SummaryRatio, SummaryLanding, SummaryTraffic, SummaryDashboard, CostReport, DisplayCPAReport | `entity/apireport.go`, `entity/dashboard.go` |
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

| Key pattern | DB Index | TTL | Purpose |
|-------------|----------|-----|---------|
| `jwt:blacklist:<jti>` | `REDISDBINDEX` (0) | sisa exp | JWT revocation list |
| `{urlservicekey}-configIdx` | `REDISDBINDEX` (0) | varies | Campaign config (JSON) |
| `{urlservicekey}-counterIdx` | `REDISDBINDEX` (0) | varies | MO/traffic counter |
| `{urlservicekey}-summary` | `REDISDBINDEX` (0) | varies | Campaign summary |
| Pixel cache key | `REDISCACHEPIXEL` (1) | `RedisKeyExpiration` (default 600s) | Pixel postback dedup |
| `HistoryCappingKey` related | `REDISDBINDEX` (0) | reset oleh cron | Daily capping reset |

---

## 6. Messaging (RabbitMQ)

### 6.1 RabbitManager — New Publisher (2026-06-05)

Package: `src/infrastructure/messaging/rabbitmq.go`

| Method | Use When | Correlation ID Prefix |
|--------|----------|-----------------------|
| `h.RM.PublishWithRetry(ctx, exchange, queue, body, corId)` | Async fire-and-forget | `RTO` |
| `h.RM.DirectReplyToWithRetry(ctx, exchange, queue, body, corId)` | Sync RPC — tunggu reply worker | `RTD` |

> **Deprecated**: `h.Rmqp.PublishMsg(rmqp.PublishItems{...})` — **jangan digunakan** di handler baru.

### 6.2 Queue & Exchange Declarations (boot-time)

| Exchange | Queue | Declared By | Publisher |
|----------|-------|-------------|-----------|
| `RABBITMQRATIOEXCHANGENAME` | `RABBITMQRATIOQUEUENAME` | `server.go` (Rmqp + RM) | Postback handlers |
| `RABBITMQCAMPAIGNMANAGEMENTEXCHANGENAME` | `RABBITMQCAMPAIGNMANAGEMENTQUEUENAME` | `server.go` (Rmqp + RM) | Campaign send handler |
| `E_RESENDCAMPAIGNDATA` | `Q_RESENDCAMPAIGNDATA` | `server.go` (Rmqp + RM) | ResendData handlers |
| `RABBITMQPIXELSTORAGEEXCHANGENAME` | `RABBITMQPIXELSTORAGEQUEUENAME` | `server.go` (Rmqp) | Legacy pixel storage |
| `RABBITMQCLICKSTORAGEEXCHANGENAME` | `RABBITMQCLICKSTORAGEQUEUENAME` | `server.go` (Rmqp) | Legacy click storage |

### 6.3 Correlation ID Convention

```
RTOxxxx  →  Async (fire-and-forget). Worker proses, TIDAK reply.
RTDxxxx  →  Sync RPC. Worker WAJIB reply ke d.ReplyTo dengan CorrelationId matching.
```

Worker reply pattern (RTD):
```go
if strings.HasPrefix(d.CorrelationId, "RTD") {
    replyBody := []byte("SHAVED") // atau "NOTSHAVED"
    ch.PublishWithContext(ctx, "", d.ReplyTo, false, false, amqp.Publishing{
        CorrelationId: d.CorrelationId,
        Body:          replyBody,
    })
}
```

Lihat: [UC-08-rabbitmq-messaging.md](./docs/UC-08-rabbitmq-messaging.md)

---

## 7. Constraints

### 7.1 Functional Constraints

| ID | Constraint | Rationale | Source |
|----|------------|-----------|--------|
| FC-01 | JWT wajib HS256 | Hindari algorithm confusion (CWE-327) | `incoming_auth_handler.go` |
| FC-02 | JWT `type` claim wajib `"access"` (refresh token reject) | Pisahkan access vs refresh | `incoming_auth_handler.go` |
| FC-03 | `aud` & `iss` validation aktif kalau env diset | Multi-tenant token isolation | `incoming_auth_handler.go` |
| FC-04 | TLS verify aktif default, opt-in disable hanya non-prod | CWE-295 fix | `infrastructure/external/http_client.go` |
| FC-05 | AES key min 16/24/32 byte (AES-128/192/256) | Crypto correctness | `infrastructure/external/utils.go` |
| FC-06 | Postback endpoint public — auth via signed params/IP allowlist | Adnet callback gak bisa kirim Bearer | `routes/postback.go` |
| FC-07 | Body limit max 100 MB | Excel upload | `url_mapping.go` |
| FC-08 | Group-level `AuthMiddleware` di `/dashboard`, `/v1/{report,int,management}` | Default-secure | gated `AUTH_ENFORCE_DEFAULT` |
| FC-09 | **[NEW]** Publisher `h.RM.PublishWithRetry` wajib digunakan — bukan `h.Rmqp.PublishMsg` | Channel pool, auto-reconnect, 3x retry | `src/infrastructure/messaging/rabbitmq.go` |
| FC-10 | **[NEW]** Context untuk publish harus request-scoped (`c.UserContext()`) — bukan global context | Hindari `context canceled` panic | `incoming_*_handler.go` |

### 7.2 Non-Functional Constraints

| ID | Constraint | Target |
|----|------------|--------|
| NFC-01 | DB pool max open conn | `DB_MAX_OPEN_CONNS × replicas ≤ PG max_connections` |
| NFC-02 | DB connection lifetime | 30 min default (`DB_CONN_MAX_LIFETIME_MIN`) |
| NFC-03 | Redis fail-fast vs degraded | `REDIS_REQUIRED=true` (default) → fail. `false` → continue tanpa cache |
| NFC-04 | RabbitMQ reconnect (legacy Rmqp) | Background goroutine watch `NotifyClose`, retry exp backoff |
| NFC-05 | HTTP outbound timeout | Per-call via `infrastructure/external/http_client.go` |
| NFC-06 | Server should not run schema migration | Migration via dedicated `migrate` sub-cmd / init container |
| NFC-07 | Logs redact secrets | `Cfg.Redacted()` mask password/private key |
| NFC-08 | Build verifiable di Go 1.24.1 | Dockerfile pin |
| NFC-09 | **[NEW]** RabbitManager pool size 5 channel, QoS 10 | Throughput vs resource balance | `server.go` |
| NFC-10 | **[NEW]** Publish retry max 3x, exponential backoff 500ms/1s/1.5s | Resilience tanpa overload | `messaging/rabbitmq.go` |

### 7.3 Operational

| ID | Constraint |
|----|------------|
| OC-01 | Multi-replica deploy require: migrate run-once → server start |
| OC-02 | Tag release `v*.*.*` → CI auto-build 2 image (server + migrate) |
| OC-03 | `.env` tidak boleh di-commit (sudah di `.gitignore`) |
| OC-04 | `AES_SECRET_KEY` wajib via secret manager di prod |
| OC-05 | DB migration tool eksternal (`golang-migrate`) recommended jangka panjang vs AutoMigrate |

---

## 8. Edge Cases

### 8.1 Auth Edge Cases

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

### 8.2 Postback Edge Cases

| Case | Expected Behavior |
|------|-------------------|
| Adnet kirim duplicate postback | Dedup via cookie (3 detik window) → 403 |
| MSISDN format tidak valid | 400 + log warn |
| Campaign ID tidak ditemukan | 404 (handler return explicitly) |
| Postback masuk saat budget capping tercapai | Tetap record, mark capped |
| Pixel sudah dipakai (`px.IsUsed = true`) | 200/409 NOK — tidak diproses ulang |
| `h.RM` nil (RabbitManager tidak terinit) | Panic nil pointer — pastikan `server.go` selalu init RM |

### 8.3 RabbitMQ / RabbitManager Edge Cases

| Case | Behavior |
|------|----------|
| Initial connection gagal | `NewRabbitManager` retry 5x (2s interval) → error → panic di `InitMessageBroker` |
| Channel pool habis (semua channel dipakai) | `<-rm.ChPool` block sampai ada yang kembali |
| Publish gagal (broker nack) | Retry 3x exponential backoff → force close conn → `handleReconnect` wakes up |
| DirectReplyTo timeout (RTD) | `ctx.Done()` → return `ctx.Err()` → handler log error |
| Worker tidak reply untuk RTD | Tunggu sampai context timeout |
| `corId` prefix `RTO` di handler sync | Bug — pastikan gunakan `RTD` prefix untuk `DirectReplyToWithRetry` |
| Legacy `h.Rmqp.PublishMsg` dipanggil | Masih functional (legacy path), tapi tidak ada retry/pool — avoid |

### 8.4 HTTP Helper Edge Cases

| Case | Behavior |
|------|----------|
| `httpClient.Do(req)` return error | Return wrapped error, **tidak** akses `response.Status` |
| `response == nil` (lib quirk) | Explicit nil-check, return error |
| Non-OK status (≠200) | Return real status code + error wrapped |

### 8.5 DB / Pool Edge Cases

| Case | Behavior |
|------|----------|
| PG `max_connections` saturated | Pool throttle (block sampai available) |
| Connection lifetime exceeded | GORM auto-recreate connection |
| AutoMigrate timeout | `log.Fatalf` (exit 1) — operator harus investigate |

### 8.6 Redis Edge Cases

| Case | Behavior |
|------|----------|
| Redis down saat startup, `REDIS_REQUIRED=true` | 5x retry exp backoff → fail |
| Redis down, `REDIS_REQUIRED=false` | Log warn, app continue dgn `nil` clients, `Setup.RedisAvailable=false` |
| Redis down saat runtime | Cache miss → fallback DB (jika handler pakai `helper.SafeRedisGet`) |

### 8.7 Concurrency

| Case | Behavior |
|------|----------|
| Concurrent edit campaign by 2 AM | Last-write-wins (no optimistic lock) — TODO future |
| Cron capping reset overlap | Single-fire (cron runs in single instance) |
| Multiple goroutine publish via RabbitManager | Thread-safe via channel pool (`sync.RWMutex` + `chan *amqp.Channel`) |

---

## 9. Acceptance Criteria

### 9.1 Per Use Case

#### UC-01.1 (Postback async — RTO)
- Adnet GET `/v1/postback?...` → response 200 dalam < 500ms
- Cookie dedup aktif 3 detik — duplicate request dalam window → 403
- Pixel lookup dari Redis — 404 jika tidak ada
- `corId` berprefix `RTO` — worker tidak reply
- Message ter-publish ke Ratio queue via `h.RM.PublishWithRetry`
- Postback gagal validate → 400 + log warn dgn request payload

#### UC-01.2 (Postback sync — RTD)
- Adnet GET `/v1/postback_sync?...` → sistem menunggu reply worker
- `corId` berprefix `RTD` — worker WAJIB reply
- Reply `"SHAVED"` atau `"NOTSHAVED"` diterima dalam batas timeout
- Timeout → log error, return appropriate error response

#### UC-02.2 (Edit campaign)
- AM dgn scope agency `X` bisa edit campaign agency `X`
- AM dgn scope agency `X` **tidak bisa** edit campaign agency `Y` → 403
- Update tersimpan ke `campaign_details` DB + Redis config
- Concurrent edit oleh 2 AM → last-write-wins (sementara, TODO lock)

#### UC-02.5 (Send campaign via RabbitMQ)
- `POST /v1/management/campaign/send` → publish ke `E_CAMPAIGNMANAGEMENT`
- Menggunakan `h.RM.PublishWithRetry` (bukan Rmqp legacy)
- Context timeout dari `RABBITMQCONTEXTTIMEOUT` env (default 30s)
- Publish gagal → log debug, return 200 tetap (fire-and-forget)

#### UC-04.7 / UC-09.3 (Logout)
- Logout endpoint call `h.RevokeJWT(jti, ttl)` dgn ttl = exp - now
- Token yg sama di-pakai → 401 `Token revoked`
- Setelah ttl expire, key blacklist auto-cleanup oleh Redis


### 9.2 Cross-Cutting (Definition of Done)

Untuk setiap fitur baru:

- Endpoint terdaftar di `src/interfaces/http/routes/<group>.go`
- Handler di `src/interfaces/http/handler/incoming_<domain>_handler.go`
- Repository (kalau ada DB access) di `src/infrastructure/persistence/<domain>.go`
- Entity (kalau ada DB schema baru) di `src/domain/entity/<domain>.go` + register di `migrateEntities`
- JWT auth applied (kalau bukan public)
- Scope filter (company/agency/adnet) sesuai actor
- Error response konsisten dgn envelope `{code, desc, error}`
- Response time < 1 detik untuk read endpoint, < 3 detik untuk export
- Logs structured (logrus) — no println bocor
- Secrets tidak di-log (gunakan `Cfg.Redacted()` pattern)
- `staticcheck ./...` zero new warning
- `go build ./...` exit 0
- **UC doc di `docs/UC-NN-*.md` updated**
- **SPEC.md updated**: section 3 (use case), 4 (API contract), 5 (entity jika ada), 8 (edge case), 9 (AC)
- Publisher menggunakan `h.RM.PublishWithRetry` / `h.RM.DirectReplyToWithRetry` — **bukan** `h.Rmqp.PublishMsg`

### 9.3 Non-Functional

- DB query log enabled (GORM logger Info mode) — no N+1 baru
- Endpoint baru di-monitor di Grafana (HTTP latency, error rate)
- Postback endpoint stress test ≥ 500 req/sec sustained
- Redis fail simulation: degraded mode tidak crash (kalau `REDIS_REQUIRED=false`)
- DB pool sat test: tidak melebihi `DB_MAX_OPEN_CONNS`
- RabbitMQ disconnect simulation: `RabbitManager.handleReconnect` auto-recovery dalam < 10s

---

## 10. Adding New Feature — Quick Workflow

```
1. Identify domain (campaign / budget / master-data / messaging / ...)

2. Update docs (SEBELUM kode):
   a. Buat/update docs/UC-NN-nama.md
   b. Update docs/README.md tabel index
   c. Update SPEC.md:
      - Section 3: Use case baru
      - Section 4.4: Endpoint definition
      - Section 5: Entity baru (jika ada)
      - Section 8: Edge cases
      - Section 9.1: Acceptance criteria

3. Implement (urutan: entity → repo → handler → route):
   - entity  → src/domain/entity/<file>.go  (+ register migrateEntities)
   - repo    → src/infrastructure/persistence/<file>.go
   - handler → src/interfaces/http/handler/incoming_<file>_handler.go
   - route   → src/interfaces/http/routes/<group>.go

4. RabbitMQ publish:
   - Gunakan h.RM.PublishWithRetry (async/RTO) atau h.RM.DirectReplyToWithRetry (sync/RTD)
   - Context: context.WithTimeout(c.UserContext(), timeout)
   - Jika queue baru: tambah ke rmDeclares + channels di server.go

5. Test:
   - go build ./...
   - staticcheck ./...
   - manual hit endpoint
   - cek log untuk leak secret

6. Run migration (jika ada entity baru):
   ./datasource migrate

7. Deploy:
   - tag git v1.x.y
   - CI auto build
   - update compose version
   - run migrate compose first
   - deploy server
```

---

## Appendix A — Env Vars Reference

Lihat `AUDIT_FIX.md` section "Env Vars Summary". Env vars terkait RabbitMQ:

| Env Var | Default | Digunakan Oleh |
|---------|---------|----------------|
| `RABBITMQHOST` | — | `Cfg.RabbitMQHost` |
| `RABBITMQPORT` | — | `Cfg.RabbitMQPort` |
| `RABBITMQUSERNAME` | — | `Cfg.RabbitMQUsername` |
| `RABBITMQPASSWORD` | — | `Cfg.RabbitMQPassword` |
| `RABBITMQVHOST` | — | `Cfg.RabbitMQVHost` |
| `RABBITMQCONTEXTTIMEOUT` | `30` | Context timeout publish (detik) |
| `RABBITMQRATIOEXCHANGENAME` | — | Ratio exchange |
| `RABBITMQRATIOQUEUENAME` | — | Ratio queue |
| `RABBITMQCAMPAIGNMANAGEMENTEXCHANGENAME` | — | CampaignManagement exchange |
| `RABBITMQCAMPAIGNMANAGEMENTQUEUENAME` | — | CampaignManagement queue |

## Appendix B — Related Docs

- [`docs/README.md`](./docs/README.md) — **Per-UC documentation index** (UC-01 s/d UC-09+)
- [`CLAUDE.md`](./CLAUDE.md) — Architecture guide, conventions, how-to-change
- [`AUDIT_FIX.md`](./AUDIT_FIX.md) — Audit fix implementation log (8 priority items)
- [`UNUSED_FUNCTIONS.md`](./UNUSED_FUNCTIONS.md) — Dead code report
- [`NOTES.md`](./NOTES.md) — Changelog per release
- `Data Source.pdf` — Original audit report (2026-02-11)

## Appendix C — Changelog SPEC

| Tanggal | Perubahan |
|---------|-----------|
| 2026-06-05 | Tambah Section 6 (RabbitMQ/RabbitManager), UC-01b (RTD sync), UC-30, UC-31; update FC/NFC constraints; update edge cases section 8.3, 8.7; update DoD section 9.2; tambah Appendix C |
| 2026-05-13 | Initial spec creation |
