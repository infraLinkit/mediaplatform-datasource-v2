# UC-03 — Reporting

> **Status**: Active  
> **Owner**: Wilie / Fadhil  
> **Last Updated**: 2026-06-05  
> **Route Group**: `/v1/report` (Auth-gated jika `AUTH_ENFORCE_DEFAULT=true`)

---

## Overview

Module reporting menyajikan semua data performa campaign, CPA, traffic, revenue, mainstream, conversion log, IO report, dan redirection time ke frontend.

---

## Endpoints

| Method | Path | Handler | Deskripsi |
|--------|------|---------|-----------|
| GET | `/v1/report/pinreport` | `DisplayPinReport` | Pin report summary |
| GET | `/v1/report/datasentapiperformance` | `DisplayPinPerformanceReport` | Pin performance API report |
| GET | `/v1/report/cpareportlist` | `DisplayCPAReport` | CPA report (filterable) |
| GET | `/v1/report/costreport/:v` | `DisplayCostReport` | Cost report (v=`list`/`detail`) |
| GET | `/v1/report/conversionlog` | `DisplayConversionLogReport` | Conversion log |
| GET | `/v1/report/campaign-monitoring-summary` | `DisplayCampaignSummary` | Campaign monitoring summary |
| GET | `/v1/report/campaign-monitoring-summary/chart` | `DisplayCampaignSummaryChart` | Chart data untuk monitoring |
| GET | `/v1/report/alertreport/:v` | `DisplayAlertReportAll` | Alert report list |
| GET | `/v1/report/trafficreport` | `DisplayTrafficReport` | Traffic report list |
| GET | `/v1/report/trafficreport/chart` | `GetTrafficReportChart` | Traffic report chart |
| GET | `/v1/report/mainstreamreport` | `DisplayMainstreamReport` | Mainstream report |
| GET | `/v1/report/google-traffic-report` | `DisplayGoogleTrafficReport` | Google traffic report |
| GET | `/v1/report/budgetmonitoring` | `DisplayBudgetMonitoring` | Budget monitoring |
| GET | `/v1/report/performance-report` | `DisplayPerformanceReport` | Performance report |
| GET | `/v1/report/revenuemonitoring` | `DisplayRevenueMonitoring` | Revenue monitoring list |
| GET | `/v1/report/revenuemonitoring/chart` | `DisplayRevenueMonitoringChart` | Revenue monitoring chart |
| GET | `/v1/report/defaultinput/` | `DisplayDefaultInput` | Default input values (cost, agency fee) |
| GET | `/v1/report/redirectiontime` | `DisplayRedirectionTime` | Redirection time report |
| POST | `/v1/report/resend-data` | `ResendData` | Resend summary data ke Linkit Dashboard |
| POST | `/v1/report/resend-data-apireport` | `ResendDataAPIReport` | Resend API report data ke Linkit Dashboard |
| GET | `/v1/report/ioreport` | `DisplaySummaryBudgetIO` | IO Report |
| PUT | `/v1/report/ioreport/update` | `UpdateSummaryBudgetIO` | Update IO Report |
| POST | `/v1/report/campaign-monitoring-summary/edit-target-budget` | `EditTargetBudgetLevel` | Edit target budget per campaign |
| POST | `/v1/report/campaign-monitoring-summary/edit-target-budget-batch` | `EditTargetBudgetBatch` | Bulk edit target budget |
| GET | `/v1/report/campaign-spending-channel` | `DisplayCampaignSpendingChannel` | Campaign spending by channel |
| GET | `/v1/report/campaign-spending-channel/country-children` | `DisplayCampaignSpendingChannelCountryChildren` | Channel spending per country |

---

## Flow: ResendData

```
POST /v1/report/resend-data
  │
  ├─ Parse form: total, id[0..N]
  ├─ GetSummaryReportById(ids) → []SummaryCampaign
  ├─ Untuk setiap report:
  │   ├─ Build query string (date, campaign_id, mo_received, etc.)
  │   ├─ fullURL = APILINKITDashboard + "?" + queryString
  │   └─ h.RM.PublishWithRetry(ctx, "E_RESENDCAMPAIGNDATA", "Q_RESENDCAMPAIGNDATA", message)
  └─ Return { status: "OK"/"NOK", error: "" }
```

---

## Flow: DisplayCPAReport

```
GET /v1/report/cpareportlist
  │
  ├─ Parse query: campaign_id, country, operator, adnets[], date_range, page, page_size
  ├─ Scope filter: allowedCompanies, allowedAdnets dari JWT Locals
  └─ DS.GetDisplayCPAReport(fe, allowedCompanies, allowedAdnets)
      └─ Response: DataTable envelope (draw, recordsTotal, data[], totalSummary)
```

---

## Flow: DisplayCostReport

```
GET /v1/report/costreport/:v
  │
  ├─ v = "list" → GetDisplayCostReport atau GetDisplayCostReportByCountry
  │   (tergantung group_by=country)
  ├─ v = "detail" → GetDisplayCostReportDetail
  └─ Redis cache (60 detik TTL) untuk non-action requests
```

---

## RabbitMQ (ResendData)

| Exchange | Queue | Deskripsi |
|----------|-------|-----------|
| `E_RESENDCAMPAIGNDATA` | `Q_RESENDCAMPAIGNDATA` | Hardcoded, kirim ulang data ke Linkit Dashboard |

---

## Response Envelope (Standard)

```json
// DataTable response
{
  "draw": 1,
  "code": 200,
  "desc": "OK",
  "data": [...],
  "recordsTotal": 100,
  "recordsFiltered": 100,
  "totalSummary": {...}
}

// Resend response
{ "status": "OK", "error": "" }
```

---

## Handler Files
- [`incoming_reports_handler.go`](../src/interfaces/http/handler/incoming_reports_handler.go)
- [`incoming_campaign_summary_handler.go`](../src/interfaces/http/handler/incoming_campaign_summary_handler.go)
- [`incoming_traffic_report_handler.go`](../src/interfaces/http/handler/incoming_traffic_report_handler.go)
- [`incoming_revenue_monitoring_handler.go`](../src/interfaces/http/handler/incoming_revenue_monitoring_handler.go)
- [`incoming_budget_monitoring.go`](../src/interfaces/http/handler/incoming_budget_monitoring.go)
- [`incoming_redirection_time_handler.go`](../src/interfaces/http/handler/incoming_redirection_time_handler.go)
- [`routes/report.go`](../src/interfaces/http/routes/report.go)

---

## Perubahan Terkini
| Tanggal | Perubahan |
|---------|-----------|
| 2026-06-05 | Ganti `h.Rmqp.PublishMsg` → `h.RM.PublishWithRetry` di `ResendData` + `ResendDataAPIReport` |
| 2026-06-05 | Tambah `DisplayCampaignSpendingChannel` dan `DisplayCampaignSpendingChannelCountryChildren` |

---

## Menambah/Mengubah Fitur

1. Edit handler di file yang relevan (lihat tabel handler di atas)
2. Daftarkan route baru di `routes/report.go`
3. Update tabel Endpoints di file ini
4. Jika ada field baru di query → update entity di `src/domain/entity/`
5. `go build ./...` → pastikan clean
