UP DATE : 2026-06-05
===

Author : Wilie (Antigravity AI assist)
----------
description :
----------
- [INFRA] Migrasi publisher RabbitMQ dari `wiliehidayat87/rmqp` (PublishMsg) ke `RabbitManager` baru
  - Buat package `src/infrastructure/messaging/rabbitmq.go` (channel pool, auto-reconnect, 3x retry)
  - Tambah field `RM *messaging.RabbitManager` ke `IncomingHandler` struct
  - Tambah field `RM *messaging.RabbitManager` ke `App3rdParty` struct di `url_mapping.go`
  - Init `RabbitManager` di `src/cmd/server.go` (pool 5, QoS 10)
  - Ganti semua `h.Rmqp.PublishMsg` → `h.RM.PublishWithRetry` di:
    - `incoming_postback_handler.go` (3 lokasi — RTO async)
    - `incoming_campaign_management_handler.go` (1 lokasi)
    - `incoming_reports_handler.go` (2 lokasi — ResendData, ResendDataAPIReport)
- [DOCS] Buat folder `docs/` dengan 9 UC file (UC-01 sampai UC-09)
  - UC-01: Postback Intake
  - UC-02: Campaign Management
  - UC-03: Reporting
  - UC-04: User, Role & Menu Management
  - UC-05: Country, Service & Catalog
  - UC-06: Internal Endpoints
  - UC-07: Budget IO & IP Range
  - UC-08: RabbitMQ Messaging (RabbitManager)
  - UC-09: Authentication & Authorization
- [DOCS] Update `CLAUDE.md` dengan konvensi baru (RabbitManager, docs/, CorrelationID RTO/RTD)
----------
version image :
- (pending build)
----------
correlation :
- internal refactor — no Jira ticket
----------

UP DATE : 2025-10-13
===


Author : Wilie
----------
description : 
----------
- Add subkeyword
----------
version image : 
- infralinkit/mediaplatform-datasource-server:1.8.3
----------
correlation : 
- https://linkit360.atlassian.net/browse/xxxx
----------

UP DATE : 2025-10-07
===

DS
===
Author : Fadhil
----------
description : 
----------
- Fix Total Load Time redirection
- Fix data type summary_landings table
- Fix Display Mainstream Report
----------
version image : 
- infralinkit/mediaplatform-datasource-server:1.7.7
- infralinkit/mediaplatform-datasource-server:1.7.8
- infralinkit/mediaplatform-datasource-server:1.7.9
- infralinkit/mediaplatform-datasource-server:1.8.1
- infralinkit/mediaplatform-datasource-server:1.8.2
----------
correlation : 
- https://linkit360.atlassian.net/browse/WR-189
----------

UP DATE : 2025-10-02
===

DS
===
Author : Fadhil
----------
description : 
----------
- Fix Cannot change previous day's payout and ratio CPA, API Report, Mainstream Report
- Mainstream Report adjustment
- Fix total load time redirection time tools
- Fix condition in edit po, ratio
----------
version image : 
- infralinkit/mediaplatform-datasource-server:1.7.5
- infralinkit/mediaplatform-datasource-server:1.7.6
----------
correlation : 
- https://linkit360.atlassian.net/browse/WR-189
----------

UP DATE : 2025-10-02
===

DS
===
Author : Wilie
----------
description : 
----------
- Total Load time (Report redirection time tools)
- Add wildcard trxid={trxid} parameter for postback method=SPC-MVLS
----------
version image : 
- infralinkit/mediaplatform-datasource-server:1.6.9
- infralinkit/mediaplatform-datasource-server:1.7.1
- infralinkit/mediaplatform-datasource-server:1.7.2
- infralinkit/mediaplatform-datasource-server:1.7.3
- infralinkit/mediaplatform-datasource-server:1.7.4

----------
correlation : 
- https://linkit360.atlassian.net/browse/WR-xxx
----------

UP DATE : 2025-10-01
===

DS
===
Author : Fadhil
----------
description : 
----------
- Fix api report insert, get display report, & edit payout af
- Fix date send insert api report
----------
version image : 
- infralinkit/mediaplatform-datasource-server:1.6.7
- infralinkit/mediaplatform-datasource-server:1.6.8
----------
correlation : 
- https://linkit360.atlassian.net/browse/WR-245
----------

UP DATE : 2025-09-29(2)
===

Author : Wilie
----------
description : 
----------
- Fix postback adnet code & spc-mvls
----------
version image : 
- infralinkit/mediaplatform-datasource-server:1.6.3
----------
correlation : 
- https://linkit360.atlassian.net/browse/WR-240
- https://linkit360.atlassian.net/browse/WR-241
----------

UP DATE : 2025-09-29
===

DS
===
Author : Fadhil
----------
description : 
----------
- ?
----------
version image : 
- infralinkit/mediaplatform-datasource-server:1.5.7
- infralinkit/mediaplatform-datasource-server:1.5.8
- infralinkit/mediaplatform-datasource-server:1.5.9
- infralinkit/mediaplatform-datasource-server:1.6.1
- infralinkit/mediaplatform-datasource-server:1.6.2
----------
correlation : 
- https://linkit360.atlassian.net/browse/xxxxx
----------

UP DATE : 2025-09-24
===

DS
===
Author : wilie
----------
description : 
----------
- Scan semua pixel / key redis di db logical db = 1 sesuai campaign id yang di postback
- Jika ketemu 1 saja pada saat searching lalu di break
- Jika valid data di teruskan ke proses ratio
----------
version image : infralinkit/mediaplatform-datasource-server:1.5.6
----------
correlation : https://linkit360.atlassian.net/browse/xxxxx
----------
