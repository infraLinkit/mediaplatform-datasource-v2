# UC-06 — Internal Endpoints & Admin Tools

> **Status**: Active  
> **Owner**: Wilie / Fadhil  
> **Last Updated**: 2026-06-05  
> **Route Group**: `/v1/int` (Auth-gated)

---

## Overview

Endpoint internal untuk operasi admin: update data transaksional, export, upload Excel, pin report, performance report, ARPU data retrieval, dan summary landing URL. Endpoint ini **tidak** diakses oleh public/adnet.

---

## Endpoints

| Method | Path | Handler | Deskripsi |
|--------|------|---------|-----------|
| PUT | `/v1/int/setdata/:v/` | `SetData` | Set target daily budget |
| PUT | `/v1/int/updatedata/:v/` | `UpdateAgencyFeeAndCostConversion` | Update agency fee & cost conversion |
| PUT | `/v1/int/updateratio/:v/` | `UpdateRatio` | Update ratio transaksional |
| PUT | `/v1/int/updatepostback/:v/` | `UpdatePostback` | Update postback transaksional |
| PUT | `/v1/int/updateagencycost/:v` | `UpdateAgencyCost` | Update agency fee & cost per conversion di DB |
| PUT | `/v1/int/updatestatusalert/:v` | `UpdateStatusAlert` | Update status alert |
| GET | `/v1/int/datasentapipinreport/` | `TrxPinReport` | Transaksi Pin Report |
| POST | `/v1/int/pinreport/editpayout` | `EditPayoutAPIReport` | Edit payout API report |
| GET | `/v1/int/datasentapiperformance/` | `TrxPerformancePinReport` | Transaksi Performance Pin Report |
| POST | `/v1/int/pinperformance/editcpa` | `EditCpaAPIPerformanceReport` | Edit CPA performance report |
| POST | `/v1/int/pinperformance/editarpu` | `EditArpuAPIPerformanceReport` | Edit ARPU performance report |
| GET | `/v1/int/exportcpa/` | `ExportCpaButton` | Export CPA report |
| GET | `/v1/int/exportcost/` | `ExportCostButton` | Export cost report |
| GET | `/v1/int/exportcostdetail/` | `ExportCostDetailButton` | Export cost detail report |
| GET | `/v1/int/pinperformance` | `PinPerformance` | Pin performance summary |
| POST | `/v1/int/uploadexcel` | `UploadExcel` | Upload Excel SMS campaign |
| PUT | `/v1/int/updateexcel/:id` | `UpdateExcel` | Update Excel SMS campaign |
| PUT | `/v1/int/upsertexcel/` | `UpsertExcel` | Upsert Excel SMS campaign |
| GET | `/v1/int/getdataarpu/` | `GetDataArpu` | Get ARPU data dari API |
| GET | `/v1/int/get_urlservice_in_summarylanding` | `GetURLServiceInSummaryLanding` | Get URL service untuk total load time |
| PUT | `/v1/int/update_response_url_service_in_summarylanding` | `UpdateResponseURLServiceInSummaryLanding` | Update response URL service |

---

## Flow: SetData (Target Daily Budget)

```
PUT /v1/int/setdata/:v/
  │
  ├─ v = "targetdailybudget" (diimplementasikan)
  ├─ Parse body: URLServiceKey, TargetDailyBudget
  ├─ Load config dari Redis
  ├─ Update field di config
  ├─ Simpan kembali ke Redis
  └─ DB: update campaign_details
```

---

## Flow: UploadExcel (SMS Campaign)

```
POST /v1/int/uploadexcel
  │
  ├─ Parse multipart: file Excel
  ├─ Parse sheet → rows SMS campaign data
  ├─ Validasi data
  └─ DB: batch insert ke SMS campaign table
```

---

## Flow: GetDataArpu

```
GET /v1/int/getdataarpu/
  │
  ├─ Cek GET_DATA_ARPU env (bool)
  ├─ Call APIARPU dengan basic auth (ARPUUsername/ARPUPassword)
  └─ Parse + simpan ke DB ARPU table
```

---

## Handler Files
- [`incoming_api_int_handler.go`](../src/interfaces/http/handler/incoming_api_int_handler.go)
- [`routes/internal.go`](../src/interfaces/http/routes/internal.go)

---

## Menambah/Mengubah Fitur

1. Edit handler di `incoming_api_int_handler.go`
2. Daftarkan route baru di `routes/internal.go` → fungsi `RegisterInternal`
3. Update tabel Endpoints di file ini
4. **Ingat**: semua endpoint di sini adalah auth-protected — jangan expose ke publik
5. `go build ./...` → pastikan clean
