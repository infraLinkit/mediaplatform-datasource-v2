# UC-07 — Budget IO & IP Range Management

> **Status**: Active  
> **Owner**: Wilie / Fadhil  
> **Last Updated**: 2026-06-05  
> **Route Group**: `/v1/management/{budget-io,ipranges}` (Auth-gated)

---

## Overview

**Budget IO** mengelola insertion order (IO) antara platform dan klien: buat IO, approval, monitoring penggunaan budget.  
**IP Range** mengelola whitelist IP range untuk validasi traffic.

---

## Endpoints

### Budget IO

| Method | Path | Handler | Deskripsi |
|--------|------|---------|-----------|
| POST | `/v1/management/budget-io/` | `CreateBudgetIO` | Buat Budget IO baru |
| GET | `/v1/management/budget-io/budgetiolist` | `DisplayBudgetIO` | List Budget IO (paginated) |
| GET | `/v1/management/budget-io/budgetiolistall` | `DisplayBudgetIOAll` | List semua Budget IO (no pagination) |
| GET | `/v1/management/budget-io/budgetioapproved` | `DisplayBudgetIOApproved` | List IO yang sudah approved |
| GET | `/v1/management/budget-io/budgetioapprovedall` | `DisplayBudgetIOApprovedAll` | Semua IO approved (no pagination) |

### IP Range

| Method | Path | Handler | Deskripsi |
|--------|------|---------|-----------|
| GET | `/v1/management/ipranges/` | `GetIPRangeFiles` | List file IP range |
| POST | `/v1/management/ipranges/upload` | `UploadIPRangeRows` | Upload CSV file IP range |
| POST | `/v1/management/ipranges/implement` | `ImplementIPRange` | Terapkan/aktifkan IP range |
| POST | `/v1/management/ipranges/download` | `DownloadIPRangeCSV` | Download CSV IP range |

---

## Flow: Budget IO Creation

```
POST /v1/management/budget-io/
  │
  ├─ Parse body: BudgetIO struct (company, amount, period, etc.)
  ├─ Validasi field
  ├─ DB: insert ke budget_ios table
  └─ Return 200 OK dengan data IO yang dibuat
```

---

## Flow: Upload IP Range

```
POST /v1/management/ipranges/upload
  │
  ├─ Parse multipart: CSV file
  ├─ Parse CSV: CIDR blocks per row
  ├─ Validasi format (net.ParseCIDR)
  ├─ DB: batch insert ke ip_ranges table
  └─ Return count rows inserted
```

---

## Flow: Implement IP Range

```
POST /v1/management/ipranges/implement
  │
  ├─ Load semua IP range dari DB
  ├─ Build lookup structure di Redis
  └─ Return status implementasi
```

---

## Handler Files
- [`incoming_budget_io_handler.go`](../src/interfaces/http/handler/incoming_budget_io_handler.go)
- [`incoming_budget_target_handler.go`](../src/interfaces/http/handler/incoming_budget_target_handler.go)
- [`incoming_iprange_handler.go`](../src/interfaces/http/handler/incoming_iprange_handler.go)
- [`routes/management.go`](../src/interfaces/http/routes/management.go)

---

## Menambah/Mengubah Fitur

1. Edit handler di file yang sesuai
2. Daftarkan route baru di `routes/management.go` → `registerBudgetIO` atau `registerIPRange`
3. Update tabel Endpoints di file ini
4. `go build ./...` → pastikan clean
