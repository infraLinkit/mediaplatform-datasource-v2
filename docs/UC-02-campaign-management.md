# UC-02 — Campaign Management

> **Status**: Active  
> **Owner**: Wilie / Fadhil  
> **Last Updated**: 2026-06-05  
> **Route Group**: `/v1/management/campaign` (Auth-gated jika `AUTH_ENFORCE_DEFAULT=true`)

---

## Overview

Campaign Management mengelola lifecycle kampanye iklan: display list, update setting (ratio, MO capping, PO), update status aktif/non-aktif, delete, dan kirim data kampanye ke worker via RabbitMQ.

---

## Endpoints

### Campaign CRUD

| Method | Path | Handler | Deskripsi |
|--------|------|---------|-----------|
| GET | `/v1/management/campaign/` | `DisplayCampaignManagement` | List semua campaign |
| GET | `/v1/management/campaign/:v` | `DisplayCampaignManagement` | Detail campaign (v=`detail`) |
| GET | `/v1/management/campaign/campaigncounts` | `GetCampaignCounts` | Jumlah total/active/inactive |
| POST | `/v1/management/campaign/send` | `SendCampaignHandler` | Kirim data ke worker via RabbitMQ |
| POST | `/v1/management/campaign/updatestatus` | `UpdateStatusCampaign` | Toggle aktif/non-aktif |
| POST | `/v1/management/campaign/editcampaign` | `EditCampaign` | Edit PO, ratio, MO capping |
| POST | `/v1/management/campaign/editmocapping` | `EditCampaignMOCapping` | Edit MO capping saja |
| POST | `/v1/management/campaign/editratio` | `EditCampaignRatio` | Edit ratio send/receive saja |
| POST | `/v1/management/campaign/editpo` | `EditCampaignPO` | Edit PO saja + recalculate CPA |
| POST | `/v1/management/campaign/delcampaign` | `DelCampaign` | Hapus campaign + cleanup Redis |
| POST | `/v1/management/campaign/updatekeymainstream` | `UpdateKeyMainstream` | Update key mainstream |
| POST | `/v1/management/campaign/updategooglesheet` | `UpdateGoogleSheet` | Update Google Sheet ID |
| POST | `/v1/management/campaign/updategooglesheetbillable` | `UpdateGoogleSheetBillable` | Update Google Sheet Billable ID |
| POST | `/v1/management/campaign/editmocappingservices2s` | `EditMOCappingServiceS2S` | Bulk edit MO capping by service (S2S) |
| POST | `/v1/management/campaign/editpoaf` | `EditPOAF` | Edit PO After Fee |
| POST | `/v1/management/campaign/editcampaignmanagementdetail` | `EditCampaignManagementDetail` | Edit detail campaign |
| POST | `/v1/management/campaign/updatecampaignmanagement` | `UpdateCampaign` | Update campaign dari form |

### Campaign Setting

| Method | Path | Handler | Deskripsi |
|--------|------|---------|-----------|
| POST | `/v1/management/campaign-setting/editratio` | `EditCampaignSettingRatio` | Edit ratio via setting page |
| POST | `/v1/management/campaign-setting/editpo` | `EditCampaignSettingPO` | Edit PO via setting page |
| POST | `/v1/management/campaign-setting/editmocapping` | `EditCampaignSettingMOCapping` | Edit MO capping via setting page |

---

## Flow: SendCampaignHandler

```
POST /v1/management/campaign/send
  │
  ├─ BodyParser → campaignData (map[string]interface{})
  ├─ JSON encode (SetEscapeHTML=false untuk preserve & dll)
  └─ Publish ke RabbitMQ CampaignManagement exchange
      └─ h.RM.PublishWithRetry(ctx, E_CAMPAIGN_MGT, Q_CAMPAIGN_MGT, body, "")
```

---

## Flow: EditCampaignPO (contoh flow kompleks)

```
POST /v1/management/campaign/editpo
  │
  ├─ BodyParser → CampaignDetail + SummaryCampaign
  ├─ Load config dari Redis (cfgRediskey)
  ├─ Update PO di Cfg + simpan kembali ke Redis
  ├─ DB: UpdateCampaignPO (jika today or zero date)
  ├─ Load SummaryCampaign dari DB
  ├─ FormulaCPA(sum) → recalculate CPA metrics
  └─ ReCalculateSummaryCampaign(calculated)
```

---

## Flow: DeleteCampaign

```
POST /v1/management/campaign/delcampaign
  │
  ├─ Load campaign dari Redis
  ├─ DelData(cfgRediskey) + DropIndex FT
  ├─ DelData(counterRedisKey) + DropIndex FT
  ├─ DelData(summaryRedisKey) + DropIndex FT
  ├─ DB: DelCampaignDetail
  ├─ DB: DelSummaryCampaign
  ├─ Count remaining details by campaign_id
  └─ Jika 0 → DB: DelCampaign (hapus campaign induk)
```

---

## RabbitMQ (SendCampaignHandler)

| Exchange | Queue | Env Var |
|----------|-------|---------|
| `RABBITMQCAMPAIGNMANAGEMENTEXCHANGENAME` | `RABBITMQCAMPAIGNMANAGEMENTQUEUENAME` | Config Cfg |

---

## Redis Keys

| Key Pattern | DB | Konten |
|-------------|-----|--------|
| `{urlservicekey}-configIdx` | `REDISDBINDEX` | Campaign config (JSON) |
| `{urlservicekey}-counterIdx` | `REDISDBINDEX` | Counter MO/traffic |
| `{urlservicekey}-summary` | `REDISDBINDEX` | Summary kampanye |

---

## Handler Files
- [`incoming_campaign_management_handler.go`](../src/interfaces/http/handler/incoming_campaign_management_handler.go)
- [`routes/management.go`](../src/interfaces/http/routes/management.go)

---

## Perubahan Terkini
| Tanggal | Perubahan |
|---------|-----------|
| 2026-06-05 | Ganti `h.Rmqp.PublishMsg` → `h.RM.PublishWithRetry` di `SendCampaignHandler` |

---

## Menambah/Mengubah Fitur

1. Edit handler di `incoming_campaign_management_handler.go`
2. Daftarkan route baru di `routes/management.go` → fungsi `registerCampaign`
3. Jika butuh table baru: tambah entity → append ke `migrateEntities` di `src/cmd/migrate.go`
4. Update tabel Endpoints di file ini
5. `go build ./...` → pastikan clean
