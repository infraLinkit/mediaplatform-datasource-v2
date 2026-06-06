# 📚 docs/ — Use Case Documentation

Dokumentasi per fitur/use case untuk `mediaplatform-datasource-v2`.

> **Convention**: Setiap UC file adalah source of truth untuk fitur tersebut.  
> Update file ini **sebelum** coding perubahan. Pastikan selalu sinkron dengan implementasi.

---

## Daftar Use Case

| UC | File | Deskripsi | Route Group |
|----|------|-----------|-------------|
| UC-01 | [UC-01-postback.md](./UC-01-postback.md) | Postback Intake (adnet callbacks) | `/v1/` (Public) |
| UC-02 | [UC-02-campaign-management.md](./UC-02-campaign-management.md) | Campaign Management CRUD | `/v1/management/campaign` |
| UC-03 | [UC-03-reporting.md](./UC-03-reporting.md) | Reporting (CPA, traffic, revenue, dll) | `/v1/report` |
| UC-04 | [UC-04-user-role-menu.md](./UC-04-user-role-menu.md) | User, Role & Menu Management | `/v1/management/{user,role,menu}` |
| UC-05 | [UC-05-country-service-catalog.md](./UC-05-country-service-catalog.md) | Country, Service & Catalog Master Data | `/v1/management/country-service` |
| UC-06 | [UC-06-internal-endpoints.md](./UC-06-internal-endpoints.md) | Internal Admin Tools | `/v1/int` |
| UC-07 | [UC-07-budget-io-iprange.md](./UC-07-budget-io-iprange.md) | Budget IO & IP Range | `/v1/management/{budget-io,ipranges}` |
| UC-08 | [UC-08-rabbitmq-messaging.md](./UC-08-rabbitmq-messaging.md) | RabbitMQ / RabbitManager Infrastructure | `src/infrastructure/messaging` |
| UC-09 | [UC-09-auth.md](./UC-09-auth.md) | Authentication & Authorization (JWT) | All protected groups |

---

## Cara Menambah UC Baru

1. Buat file `UC-NN-nama-fitur.md` di folder ini
2. Gunakan template di bawah
3. Tambahkan entry ke tabel di atas
4. Update `CLAUDE.md` section "When changing things" jika perlu

---

## Template UC

```markdown
# UC-NN — Nama Fitur

> **Status**: Active / Draft / Deprecated
> **Owner**: Nama
> **Last Updated**: YYYY-MM-DD
> **Route Group**: /v1/xxx (Auth/Public)

## Overview
Penjelasan singkat fitur.

## Endpoints
| Method | Path | Handler | Deskripsi |
|--------|------|---------|-----------|

## Flow
Diagram atau deskripsi alur.

## Handler Files
- [filename.go](../src/...)

## Perubahan Terkini
| Tanggal | Perubahan |

## Menambah/Mengubah Fitur
Langkah-langkah checklist.
```

---

## Quick Reference: Semua Endpoints

### Public (No Auth)
| Method | Path | UC |
|--------|------|----|
| GET | `/v1/postback/:urlservicekey/` | UC-01 |
| GET | `/v1/postback` | UC-01 |
| GET | `/v1/postback_billed` | UC-01 |
| GET | `/v1/postback_sync` | UC-01 |
| GET | `/v1/inquire/campid` | UC-01 |
| GET | `/v1/inquire/api-campid` | UC-01 |

### Dashboard (Auth-gated)
| Method | Path | UC |
|--------|------|----|
| GET | `/dashboard/*` | UC-09 |

### Report (Auth-gated)
| Method | Path | UC |
|--------|------|----|
| GET | `/v1/report/pinreport` | UC-03 |
| GET | `/v1/report/cpareportlist` | UC-03 |
| GET | `/v1/report/costreport/:v` | UC-03 |
| GET | `/v1/report/trafficreport` | UC-03 |
| GET | `/v1/report/revenuemonitoring` | UC-03 |
| POST | `/v1/report/resend-data` | UC-03 |
| ... | (lihat UC-03 untuk list lengkap) | UC-03 |

### Internal (Auth-gated)
| Method | Path | UC |
|--------|------|----|
| PUT | `/v1/int/setdata/:v/` | UC-06 |
| POST | `/v1/int/uploadexcel` | UC-06 |
| ... | (lihat UC-06 untuk list lengkap) | UC-06 |

### Management (Auth-gated)
| Method | Path | UC |
|--------|------|----|
| GET | `/v1/management/campaign/` | UC-02 |
| POST | `/v1/management/campaign/send` | UC-02 |
| GET | `/v1/management/user/` | UC-04 |
| GET | `/v1/management/role/` | UC-04 |
| GET | `/v1/management/country-service/country` | UC-05 |
| GET | `/v1/management/budget-io/budgetiolist` | UC-07 |
| POST | `/v1/management/ipranges/upload` | UC-07 |
| ... | (lihat UC file masing-masing) | |
