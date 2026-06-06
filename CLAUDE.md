# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Service overview

`mediaplatform-datasource-v2` — Go HTTP API backing the Linkit360 mediaplatform. Cobra CLI with two subcommands: `server` (Fiber v2 listener on `APPAPIPORT`) and `migrate` (run-once GORM AutoMigrate then exit). Module path `github.com/infraLinkit/mediaplatform-datasource-v2`. Go 1.24.1. Canonical DDD layout (`src/{application,domain,infrastructure,interfaces,cmd}`).

Single binary serves campaign management, postback intake, reporting (CPA, mainstream, traffic, redirection-time, IO/budget, revenue, dashboard), user/role/menu management, country-service catalog, IP-range uploads.

Source lives under `src/`. `main.go` at repo root is a 3-line wrapper that imports `src/cmd` and calls `cmd.Execute()`. Authoritative spec: [`SPEC.md`](./SPEC.md). Audit log: [`AUDIT_FIX.md`](./AUDIT_FIX.md). **Use case docs**: [`docs/`](./docs/README.md) — per-UC reference.

## Common commands

`Makefile` is the canonical entry point. `make help` lists all targets.

```bash
# Run locally — env loaded from .env (sourced by shell, not auto-loaded by app)
set -a && source .env && set +a && make run            # ./datasource server
set -a && source .env && set +a && make run-migrate    # ./datasource migrate

# Build
make build           # ./bin/datasource (host arch)
make build-linux     # ./bin/datasource-linux (linux/amd64, CGO off)

# Test / lint
make test            # go test -race -count=1 ./...
make test-pkg PKG=./src/infrastructure/external
make cover           # coverage.out + summary
make vet
make staticcheck     # requires `go install honnef.co/go/tools/cmd/staticcheck@latest`
make check           # vet + staticcheck

# Format / deps
make fmt
make tidy

# Docker (two images: server + migrate)
make docker-build           # both
make docker-build-server
make docker-build-migrate
make docker-push-server VERSION=v1.x.y
```

Raw Go works too (`go run . server`, `go test ./...`, etc.) — `Makefile` just wraps with consistent flags (`-trimpath`, `-race`, ldflags).

CLI is `cobra-cli`. Subcommands registered in `src/cmd/root.go`: `server` (`src/cmd/server.go`), `migrate` (`src/cmd/migrate.go`).

## Architecture

Canonical DDD, single-process Fiber server. Layout under `src/`.

```
main.go (root) → src/cmd/{root,server,migrate}
                       │
   server: config.InitCfg → config.Initiate("api") → http.MapUrls(App3rdParty) → fiber.App
                                                              ↓
                                       src/interfaces/http/routes/* (route registration)
                                                              ↓
                                       src/interfaces/http/handler/incoming_*.go (HTTP layer)
                                                              ↓
                                       src/infrastructure/persistence/*.go (DB / Redis / Sheets)
                                                              ↓
                                       src/domain/entity/*.go (GORM tables)

   migrate: config.InitCfg → config.Initiate("migrate") → DB.AutoMigrate(migrateEntities...)
```

Reserved empty layers (DDD skeleton, not yet populated):
- `src/application/service/` — application services / use case orchestrators
- `src/domain/service/` — pure domain logic
- `src/infrastructure/messaging/` — **[POPULATED]** `RabbitManager` publisher (channel pool, auto-reconnect, 3x retry with backoff). File: `rabbitmq.go`. **Gunakan `h.RM.PublishWithRetry` / `h.RM.DirectReplyToWithRetry` — jangan pakai `h.Rmqp.PublishMsg` lagi.**
- `src/interfaces/http/middleware/` — extracted middleware (auth currently inline in `url_mapping.go`)

### Boot sequence — `server` (`src/cmd/server.go`)
1. `config.InitCfg()` reads env vars into `Cfg`.
2. `cfg.Initiate("api")` opens external connections eagerly: Postgres (GORM, pgx-simple-protocol), Redis classic client (`go-redis`, db index `REDISCACHEPIXEL`), Redis rueidis client (db index `REDISDBINDEX`), RabbitMQ (`wiliehidayat87/rmqp`), Google Sheets API. Failure behavior depends on `REDIS_REQUIRED` (default `true` → fatal; `false` → degraded mode with `nil` Redis clients, `Setup.RedisAvailable=false`).
3. **Server does NOT run AutoMigrate.** Schema migration is the `migrate` sub-cmd's job. See NFC-06 / OC-01 in SPEC.
4. RabbitMQ exchange/queue declarations: `PixelStorage`, `ClickStorage`, `Ratio`, `CampaignManagement`, plus hard-coded `E_RESENDCAMPAIGNDATA`/`Q_RESENDCAMPAIGNDATA`. Also initializes `RabbitManager` (channel pool, 5 channels, QoS 10) for `Ratio`, `CampaignManagement`, `ResendData` queues. Other queues in `Cfg` (RedisCounter, PostbackAdnet, Alert) declared by their publishers.
5. `rmqpReconnectWatcher` goroutine: blocks on `NotifyClose`, reconnects + re-sets up channels. Loops forever.
6. `http.MapUrls` (`src/interfaces/http/url_mapping.go`) mounts routes; `router.Listen(":" + AppApiPort)` blocks.

### Boot sequence — `migrate` (`src/cmd/migrate.go`)
Run-once. `config.Initiate("migrate")` → `DB.WithContext(ctx).AutoMigrate(migrateEntities...)` → exit. Context timeout from `AUTO_MIGRATE_TIMEOUT_MIN` env (default 5 min). Run as init container / CI step before server starts.

### Redis split (important)
Two Redis clients on the same host, different DB indexes:
- `R` (`*rueidis.Storage`, DB `REDISDBINDEX`, default 0): campaign cluster data, JSON values, hot path lookups, JWT blacklist (`jwt:blacklist:<jti>`).
- `RCP` (`*redis.Client`, DB `REDISCACHEPIXEL`, default 1): temporary pixel storage, postback dedup keys.

Per `README.md`: db 0 = campaign cluster, db 1 = pixel storage. Don't write to the wrong index. Use degraded-mode-aware helpers in `src/infrastructure/persistence/redis_safe.go`.

### Layout (under `src/`)
- `cmd/` — Cobra entrypoints. `root.go` registers, `server.go` boots HTTP, `migrate.go` runs AutoMigrate (`migrateEntities` slice).
- `infrastructure/config/c.go` — `Cfg` struct (all env vars), constructors for every dep (`InitGormPgx`, `InitRedis`, `InitRedisJSON`, `InitMessageBroker`, `InitGoogleSheet`). `Cfg.Redacted()` masks secrets for logging.
- `interfaces/http/url_mapping.go` — `MapUrls` builds `*fiber.App`, applies middleware, delegates route registration to `interfaces/http/routes/*`.
- `interfaces/http/routes/` — **wired up.** One file per group: `postback.go` (public), `dashboard.go`, `report.go`, `internal.go`, `management.go`. Each exports `Register<Name>(grp fiber.Router, h *handler.IncomingHandler)`.
- `interfaces/http/handler/incoming_*.go` — HTTP handlers, one file per resource. Each method hangs off `*IncomingHandler`.
- `interfaces/http/middleware/` — reserved (empty). Auth middleware currently inline in `url_mapping.go`.
- `domain/entity/*.go` — GORM models (table schema). Adding fields to existing tables: picked up by next `migrate` run. Drops/renames: manual SQL.
- `domain/service/` — reserved (empty). Pure domain logic goes here when extracted.
- `infrastructure/persistence/*.go` — DB + Redis + Sheets access. Methods hang off `*BaseModel` (or per-domain repo). Handlers reach repos via `h.DS`. Includes `redis_safe.go` for degraded-mode helpers.
- `infrastructure/external/` — outbound integrations + low-level helpers: `http_client.go` (TLS verify default-on), `logrus.go` (logger factory), `sig.go`, `utils.go` (incl. AES helpers — key must be 16/24/32 bytes).
- `infrastructure/messaging/` — **[ACTIVE]** `RabbitManager` — channel pool publisher dengan auto-reconnect. File: `rabbitmq.go`. Struct: `RabbitManager`, `RabbitDeclare`, `RabbitConfig`. Fungsi utama: `PublishWithRetry`, `DirectReplyToWithRetry`, `InitMessageBroker`.
- `application/service/` — reserved (empty). Application services / use case orchestrators.
- Root: `main.go`, `Makefile`, `Dockerfile.datasource.server`, `Dockerfile.datasource.migrate`, `build.sh`, `.env`, `SPEC.md`, `AUDIT_FIX.md`, `TESTING.md`, `UNUSED_FUNCTIONS.md`, `NOTES.md`, **`docs/`** (per-UC documentation).
- `Dockerfile.datasource.server` — multi-stage; final image is `scratch`. Optional `test` target stage runs `go test -v ./...`. `CMD ["/datasource", "server"]`.
- `Dockerfile.datasource.migrate` — same build, `CMD ["/datasource", "migrate"]`. Run-once container.

### Auth middleware (default-secure gate)
`authEnforceDefault()` reads `AUTH_ENFORCE_DEFAULT` env (default `false` while FE migrates). When true, group-level `AuthMiddleware` applies to `/dashboard`, `/v1/report`, `/v1/int`, `/v1/management`. JWT HS256, claims `sub` / `jti` / `type=access` / `exp` / `nbf` / `aud` / `iss`. See SPEC §4.2 + §7.1.

### Route prefixes (see `src/interfaces/http/url_mapping.go` + `src/interfaces/http/routes/`)
- `/dashboard/*` — dashboard widgets (auth-gated when `AUTH_ENFORCE_DEFAULT=true`).
- `/v1/postback`, `/v1/postback/:urlservicekey/`, `/v1/postback_billed`, `/v1/inquire/campid` — adnet postback intake (always public).
- `/v1/report/*` — reports (auth-gated).
- `/v1/int/*` — internal-only mutation/export endpoints (auth-gated).
- `/v1/management/{campaign,campaign-setting,menu,role,user,userlog,budget-io,country-service,ipranges}` — admin/config CRUD (auth-gated).
- `/v1/ext/*` — placeholder, currently empty.

### Conventions
- Handlers thin: parse → call repo (`h.DS.<Method>`) → respond. Business logic in `src/infrastructure/persistence/`.
- Handler files: `incoming_<resource>_handler.go` or `incoming_<resource>.go`. Match existing pattern.
- Test files coexist next to source (`*_test.go`). `make test` runs with `-race -count=1`.
- DB pool config — see `src/infrastructure/config/c.go`. NFC-01..NFC-02: `DB_MAX_OPEN_CONNS` × replicas ≤ PG `max_connections`; `DB_CONN_MAX_LIFETIME_MIN` default 30. Don't add `WithContext` cancellation expecting pooled timeouts without revisiting pool config.
- GORM logger `logger.Info` — every query logged.
- Body limit 100 MiB (`fiber.Config.BodyLimit`).
- Access log → `LOGPATH/access_log` via `fiberlogrus`; app log → `LOGPATH/api`.
- Response envelope (mostly): `{"code":N, "desc":"...", "data":{...}}` or `{"code":N, "desc":"...", "error":"..."}`.

## Environment

`.env` at repo root drives local config (sample committed — contains dev creds). `src/cmd/server.go` does **not** call `godotenv.Load`; the shell must export the vars (`set -a; source .env; set +a`) before `make run`. In Docker, env comes from compose / k8s.

Required-ish env vars: `APPPATH`, `APPAPIPORT`, `REDISHOST`/`REDISPORT`/`REDISPASSWORD`/`REDISDBINDEX`/`REDISCACHEPIXEL`, `REDIS_REQUIRED`, `DB_HOST`/`DB_PORT`/`DB_USERNAME`/`DB_PASSWORD`/`DB_DATABASE`, `DB_MAX_OPEN_CONNS`/`DB_CONN_MAX_LIFETIME_MIN`, `RABBITMQ*`, `LOGPATH`/`LOGENV`/`LOGLEVEL`, `TZ`, `JWT_SECRET`, `AUTH_ENFORCE_DEFAULT`, `AES_SECRET_KEY`, all `GS*` for Google Sheets, `APIARPU`/`ARPUUsername`/`ARPUPassword`, `AUTO_MIGRATE_TIMEOUT_MIN` (migrate cmd). Full list: SPEC §Appendix A + AUDIT_FIX.md.

## Release / deploy

Image tagging via `make docker-build VERSION=v1.x.y` (or edit `build.sh`). Tag git `v*.*.*`. `NOTES.md` is the changelog: append a dated block with author, description, and docker tag(s) shipped — keep the running record current when cutting a release. Image names:
- `infralinkit/mediaplatform-datasource-server:<ver>`
- `infralinkit/mediaplatform-datasource-migrate:<ver>`

Deploy order (OC-01): **migrate run-once container first → then start server replicas**. Multi-replica deploy must not race AutoMigrate.

## When changing things

- **New table**: add `src/domain/entity/<name>.go` + append to `migrateEntities` in `src/cmd/migrate.go`. Run `./datasource migrate` (or `make run-migrate`) before deploying server.
- **New repository**: add `src/infrastructure/persistence/<name>.go`. Methods hang off `*BaseModel` (see `persistence/base.go`) or a dedicated per-domain struct. Handlers reach repos via `h.DS`.
- **New route**: add handler in `src/interfaces/http/handler/incoming_<resource>_handler.go`, register in the appropriate `src/interfaces/http/routes/<group>.go`. Don't bypass `routes/` — `MapUrls` delegates to it.
- **New env var**: add field to `config.Cfg` (`src/infrastructure/config/c.go`) + parse in `InitCfg` + document in `.env` + SPEC §Appendix A. Mask in `Cfg.Redacted()` if sensitive.
- **New RabbitMQ queue used at boot**: (a) tambah `messaging.RabbitDeclare` entry di `rmDeclares` slice di `src/cmd/server.go`. (b) tambah entry di `channels` slice (untuk legacy Rmqp reconnect watcher). (c) Use `h.RM.PublishWithRetry(ctx, exchange, queue, body, corId)` di handler. **Jangan gunakan `h.Rmqp.PublishMsg`** — sudah deprecated di semua handler.
- **Correlation ID convention**: prefix `RTO` = async (PublishWithRetry), prefix `RTD` = sync RPC (DirectReplyToWithRetry). Worker wajib reply ke `d.ReplyTo` hanya untuk `RTD`. Lihat [UC-08](./docs/UC-08-rabbitmq-messaging.md).
- **DB calls**: pool is configurable now (no longer uncapped). If adding `WithContext` cancellation, verify it interacts sanely with `DB_MAX_OPEN_CONNS` / `DB_CONN_MAX_LIFETIME_MIN`.
- **Auth on a new protected route**: nothing extra — group-level middleware handles it. Just register under the right group (`/dashboard`, `/v1/report`, `/v1/int`, `/v1/management`). Public-only endpoint: register under `v1` directly (see `routes/postback.go`).
- **New use case**: (1) buat/update file di `docs/UC-NN-nama.md`, (2) update `docs/README.md` tabel index, (3) update SPEC.md §3.1+§3.2+§4.4+§7+§8.1 **before** code, (4) implement entity → repo → handler → route.
- **Definition of Done**: SPEC §8.2 — route + handler + repo + entity registered + auth + scope filter + envelope + structured logs (no secret leak) + `staticcheck ./...` clean + `go build ./...` exit 0 + SPEC.md updated + **UC doc di `docs/` updated**.
