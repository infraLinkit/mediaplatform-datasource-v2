# UC-05 — Country, Service & Catalog Management

> **Status**: Active  
> **Owner**: Wilie / Fadhil  
> **Last Updated**: 2026-06-05  
> **Route Group**: `/v1/management/country-service` (Auth-gated)

---

## Overview

Mengelola master data catalog: country, operator, partner, service, domain, company, agency, channel, adnet, mainstream group. Data ini menjadi referensi saat membuat/edit campaign.

---

## Endpoints

### Country

| Method | Path | Handler | Deskripsi |
|--------|------|---------|-----------|
| GET | `/v1/management/country-service/country` | `DisplayCountry` | List semua country |
| GET | `/v1/management/country-service/country/:code` | `DisplayCountryInfo` | Info country by code |
| POST | `/v1/management/country-service/country` | `CreateCountry` | Tambah country |
| PUT | `/v1/management/country-service/country/:id` | `UpdateCountry` | Update country |
| DELETE | `/v1/management/country-service/country/:id` | `DeleteCountry` | Hapus country |
| GET | `/v1/management/country-service/continent` | `DisplayCountry` | List continent (same handler) |

### Company

| Method | Path | Handler | Deskripsi |
|--------|------|---------|-----------|
| GET | `/v1/management/country-service/company` | `DisplayCompany` | List company |
| POST | `/v1/management/country-service/company` | `CreateCompany` | Tambah company |
| PUT | `/v1/management/country-service/company/:id` | `UpdateCompany` | Update company |
| DELETE | `/v1/management/country-service/company/:id` | `DeleteCompany` | Hapus company |
| GET | `/v1/management/country-service/company-group` | `DisplayCompanyGroup` | List company group |
| POST | `/v1/management/country-service/company-group` | `CreateCompanyGroup` | Tambah company group |
| PUT | `/v1/management/country-service/company-group/:id` | `UpdateCompanyGroup` | Update company group |
| DELETE | `/v1/management/country-service/company-group/:id` | `DeleteCompanyGroup` | Hapus company group |

### Domain

| Method | Path | Handler | Deskripsi |
|--------|------|---------|-----------|
| GET | `/v1/management/country-service/domain` | `DisplayDomain` | List domain |
| POST | `/v1/management/country-service/domain` | `CreateDomain` | Tambah domain |
| PUT | `/v1/management/country-service/domain/:id` | `UpdateDomain` | Update domain |
| DELETE | `/v1/management/country-service/domain/:id` | `DeleteDomain` | Hapus domain |
| GET | `/v1/management/country-service/domain-service` | `DisplayDomainService` | List domain-service mapping |

### Operator

| Method | Path | Handler | Deskripsi |
|--------|------|---------|-----------|
| GET | `/v1/management/country-service/operator` | `DisplayOperator` | List operator |
| POST | `/v1/management/country-service/operator` | `CreateOperator` | Tambah operator |
| PUT | `/v1/management/country-service/operator/:id` | `UpdateOperator` | Update operator |
| DELETE | `/v1/management/country-service/operator/:id` | `DeleteOperator` | Hapus operator |
| GET | `/v1/management/country-service/api-operator-list` | `DisplayAPIOperatorList` | Operator list untuk API dropdown |

### Partner

| Method | Path | Handler | Deskripsi |
|--------|------|---------|-----------|
| GET | `/v1/management/country-service/partner` | `DisplayPartner` | List partner |
| POST | `/v1/management/country-service/partner` | `CreatePartner` | Tambah partner |
| PUT | `/v1/management/country-service/partner/:id` | `UpdatePartner` | Update partner |
| DELETE | `/v1/management/country-service/partner/:id` | `DeletePartner` | Hapus partner |

### Service

| Method | Path | Handler | Deskripsi |
|--------|------|---------|-----------|
| GET | `/v1/management/country-service/service` | `DisplayService` | List service |
| POST | `/v1/management/country-service/service` | `CreateService` | Tambah service |
| PUT | `/v1/management/country-service/service/:id` | `UpdateService` | Update service |
| DELETE | `/v1/management/country-service/service/:id` | `DeleteService` | Hapus service |
| GET | `/v1/management/country-service/api-service-list` | `DisplayAPIServiceList` | Service list untuk API dropdown |

### Adnet

| Method | Path | Handler | Deskripsi |
|--------|------|---------|-----------|
| GET | `/v1/management/country-service/adnet-list` | `DisplayAdnetList` | List adnet |
| GET | `/v1/management/country-service/api-adnet-list` | `DisplayAPIAdnetList` | Adnet list untuk API dropdown |
| POST | `/v1/management/country-service/adnet-list` | `CreateAdnetList` | Tambah adnet |
| PUT | `/v1/management/country-service/adnet-list/:id` | `UpdateAdnetList` | Update adnet |
| DELETE | `/v1/management/country-service/adnet-list/:id` | `DeleteAdnetList` | Hapus adnet |
| POST | `/v1/management/country-service/edit-adnet-dsp-status` | `UpdateDSPAdnetStatus` | Toggle DSP status adnet |

### Email

| Method | Path | Handler | Deskripsi |
|--------|------|---------|-----------|
| GET | `/v1/management/country-service/email` | `DisplayEmail` | List email |
| GET | `/v1/management/country-service/email/:id` | `DisplayEmailByID` | Email by ID |
| POST | `/v1/management/country-service/email` | `CreateEmail` | Tambah email |
| PUT | `/v1/management/country-service/email/:id` | `UpdateEmail` | Update email |
| DELETE | `/v1/management/country-service/email/:id` | `DeleteEmail` | Hapus email |

### Agency

| Method | Path | Handler | Deskripsi |
|--------|------|---------|-----------|
| GET | `/v1/management/country-service/agency` | `DisplayAgency` | List agency |
| POST | `/v1/management/country-service/agency` | `CreateAgency` | Tambah agency |
| PUT | `/v1/management/country-service/agency/:id` | `UpdateAgency` | Update agency |
| DELETE | `/v1/management/country-service/agency/:id` | `DeleteAgency` | Hapus agency |

### Channel

| Method | Path | Handler | Deskripsi |
|--------|------|---------|-----------|
| GET | `/v1/management/country-service/channel` | `DisplayChannel` | List channel |
| POST | `/v1/management/country-service/channel` | `CreateChannel` | Tambah channel |
| PUT | `/v1/management/country-service/channel/:id` | `UpdateChannel` | Update channel |
| DELETE | `/v1/management/country-service/channel/:id` | `DeleteChannel` | Hapus channel |

### Mainstream Group

| Method | Path | Handler | Deskripsi |
|--------|------|---------|-----------|
| GET | `/v1/management/country-service/mainstream-group` | `DisplayMainstreamGroup` | List mainstream group |
| POST | `/v1/management/country-service/mainstream-group` | `CreateMainstreamGroup` | Tambah mainstream group |
| PUT | `/v1/management/country-service/mainstream-group/:id` | `UpdateMainstreamGroup` | Update mainstream group |
| DELETE | `/v1/management/country-service/mainstream-group/:id` | `DeleteMainstreamGroup` | Hapus mainstream group |

---

## Handler Files
- [`incoming_country_service_management_handler.go`](../src/interfaces/http/handler/incoming_country_service_management_handler.go)
- [`routes/management.go`](../src/interfaces/http/routes/management.go)

---

## Menambah Catalog Baru

1. Buat entity di `src/domain/entity/<nama>.go`
2. Append ke `migrateEntities` di `src/cmd/migrate.go`
3. Buat methods di persistence (CRUD)
4. Tambah handlers di `incoming_country_service_management_handler.go`
5. Register routes di `routes/management.go` → fungsi `registerCountryService`
6. Update tabel Endpoints di file ini
7. `go build ./...` → pastikan clean
