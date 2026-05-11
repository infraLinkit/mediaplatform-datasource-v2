package persistence

import (
	"fmt"
	"time"

	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
	"gorm.io/gorm"
)

func (r *BaseModel) GetDisplayGoogleTrafficReport(
	o entity.DisplayGoogleTrafficReport,
) ([]entity.GoogleTrafficReportRow, int64, entity.GoogleTrafficTotalSummary, error) {

	var (
		rows         []entity.GoogleTrafficReportRow
		totalSummary entity.GoogleTrafficTotalSummary
		totalData    int64
	)

	// ============================================================
	// 1. BASE QUERY (FILTER ONLY - NO GROUP)
	// ============================================================
	base := r.DB.Table("summary_campaign_billings")
	base = applyFilters(base, o)
	base, periodSelectExpr, orderExpr, groupExtra, isWeekly, err := applyPeriod(base, o)
	if err != nil {
		return nil, 0, totalSummary, err
	}

	// ============================================================
	// 2. GROUPED SUBQUERY
	// ============================================================
	grouped := base.Select(fmt.Sprintf(`
		url_service_key,
		campaign_id,
		campaign_name,
		country,
		operator,
		partner,
		adnet,
		service,
		company,
		adgroup_id,
		placement,
		%s,
		SUM(status_success) AS status_success,
		SUM(status_failed)  AS status_failed,
		SUM(total_bill)     AS total_bill,
		CASE WHEN SUM(total_bill) > 0
			THEN ROUND((SUM(status_success)::numeric / SUM(total_bill)::numeric) * 100, 2)
			ELSE 0
		END AS bill_rate
	`, periodSelectExpr)).
		Group(`
			url_service_key, campaign_id, campaign_name, country,
			operator, partner, adnet, service, company,
			adgroup_id, placement, period_label` + groupExtra)

	// ============================================================
	// 3. COUNT (SUBQUERY)
	// ============================================================
	if err := r.DB.Table("(?) as sub", grouped).Count(&totalData).Error; err != nil {
		return nil, 0, totalSummary, err
	}

	// ============================================================
	// 4. PAGINATION
	// ============================================================
	pageSize := o.PageSize
	if pageSize == 0 {
		pageSize = 10
	}

	offset := 0
	if o.Page > 1 {
		offset = (o.Page - 1) * pageSize
	}

	// ============================================================
	// 5. RAW STRUCT
	// ============================================================
	type rawRow struct {
		URLServiceKey string  `gorm:"column:url_service_key"`
		CampaignId    string  `gorm:"column:campaign_id"`
		CampaignName  string  `gorm:"column:campaign_name"`
		Country       string  `gorm:"column:country"`
		Operator      string  `gorm:"column:operator"`
		Partner       string  `gorm:"column:partner"`
		Adnet         string  `gorm:"column:adnet"`
		Service       string  `gorm:"column:service"`
		Company       string  `gorm:"column:company"`
		AdgroupID     string  `gorm:"column:adgroup_id"`
		Placement     string  `gorm:"column:placement"`
		PeriodLabel   string  `gorm:"column:period_label"`
		WeekNum       *int    `gorm:"column:week_num"`
		StatusSuccess int     `gorm:"column:status_success"`
		StatusFailed  int     `gorm:"column:status_failed"`
		TotalBill     int     `gorm:"column:total_bill"`
		BillRate      float64 `gorm:"column:bill_rate"`
	}

	var rawRows []rawRow

	err = r.DB.Table("(?) as sub", grouped).
		Order(orderExpr).
		Offset(offset).
		Limit(pageSize).
		Scan(&rawRows).Error

	if err != nil {
		return nil, 0, totalSummary, err
	}

	// ============================================================
	// 6. SUMMARY (NO GROUP BY - grand total)
	// ============================================================
	type summaryRaw struct {
		TotalSuccess   int     `gorm:"column:total_success"`
		TotalFailed    int     `gorm:"column:total_failed"`
		TotalBill      int     `gorm:"column:total_bill"`
		AvgBillRate    float64 `gorm:"column:avg_bill_rate"`
		TotalCampaigns int     `gorm:"column:total_campaigns"`
		TotalAdgroups  int     `gorm:"column:total_adgroups"`
	}

	var sr summaryRaw

	summaryBase := r.DB.Table("summary_campaign_billings")
	summaryBase = applyFilters(summaryBase, o)
	summaryBase, _, _, _, _, err = applyPeriod(summaryBase, o)
	if err != nil {
		return nil, 0, totalSummary, err
	}

	err = summaryBase.Select(`
		SUM(status_success) AS total_success,
		SUM(status_failed)  AS total_failed,
		SUM(total_bill)     AS total_bill,
		CASE WHEN SUM(total_bill) > 0
			THEN ROUND((SUM(status_success)::numeric / SUM(total_bill)::numeric) * 100, 2)
			ELSE 0
		END AS avg_bill_rate,
		COUNT(DISTINCT campaign_id) AS total_campaigns,
		COUNT(DISTINCT adgroup_id)  AS total_adgroups
	`).Scan(&sr).Error

	if err != nil {
		return nil, 0, totalSummary, err
	}

	totalSummary = entity.GoogleTrafficTotalSummary{
		TotalSuccess:   sr.TotalSuccess,
		TotalFailed:    sr.TotalFailed,
		TotalBill:      sr.TotalBill,
		AvgBillRate:    sr.AvgBillRate,
		TotalCampaigns: sr.TotalCampaigns,
		TotalAdgroups:  sr.TotalAdgroups,
	}

	// ============================================================
	// 7. MAPPING
	// ============================================================
	for _, rr := range rawRows {
		label := rr.PeriodLabel
		if isWeekly && rr.WeekNum != nil {
			label = fmt.Sprintf("Week %d (%s)", *rr.WeekNum, rr.PeriodLabel)
		}

		rows = append(rows, entity.GoogleTrafficReportRow{
			URLServiceKey: rr.URLServiceKey,
			CampaignId:    rr.CampaignId,
			CampaignName:  rr.CampaignName,
			Country:       rr.Country,
			Operator:      rr.Operator,
			Partner:       rr.Partner,
			Adnet:         rr.Adnet,
			Service:       rr.Service,
			Company:       rr.Company,
			AdgroupID:     rr.AdgroupID,
			Placement:     rr.Placement,
			PeriodLabel:   label,
			SummaryDate:   rr.PeriodLabel,
			StatusSuccess: rr.StatusSuccess,
			StatusFailed:  rr.StatusFailed,
			TotalBill:     rr.TotalBill,
			BillRate:      rr.BillRate,
		})
	}

	return rows, totalData, totalSummary, nil
}

// ============================================================
// FILTER HELPER (DRY)
// ============================================================
func applyFilters(db *gorm.DB, o entity.DisplayGoogleTrafficReport) *gorm.DB {

	if o.UrlServiceKey != "" {
		db = db.Where("url_service_key = ?", o.UrlServiceKey)
	}
	if o.CampaignId != "" {
		db = db.Where("campaign_id = ?", o.CampaignId)
	}
	if o.CampaignName != "" {
		db = db.Where("campaign_name = ?", o.CampaignName)
	}
	if o.Country != "" {
		db = db.Where("country = ?", o.Country)
	}
	if o.Operator != "" {
		db = db.Where("operator = ?", o.Operator)
	}
	if o.Partner != "" {
		db = db.Where("partner = ?", o.Partner)
	}
	if o.Service != "" {
		db = db.Where("service = ?", o.Service)
	}
	if o.Company != "" {
		db = db.Where("company = ?", o.Company)
	}
	if o.Adnet != "" {
		db = db.Where("adnet = ?", o.Adnet)
	}
	if o.AdgroupID != "" {
		db = db.Where("adgroup_id = ?", o.AdgroupID)
	}

	return db
}

// ============================================================
// PERIOD HELPER
// FIX: weekly sekarang support filter per week number (1,2,3,4,5,...)
//      + week number dihitung dengan formula yang benar (ISO-like per-bulan)
// ============================================================
func applyPeriod(db *gorm.DB, o entity.DisplayGoogleTrafficReport) (*gorm.DB, string, string, string, bool, error) {

	var (
		periodSelectExpr string
		orderExpr        string
		groupExtra       string
		isWeekly         bool
	)

	switch o.PeriodType {

	case "weekly":
		if o.Month == "" {
			return nil, "", "", "", false, fmt.Errorf("month required for weekly period")
		}

		monthStart, err := time.Parse("2006-01", o.Month)
		if err != nil {
			return nil, "", "", "", false, fmt.Errorf("invalid month format, expected YYYY-MM: %w", err)
		}

		monthEnd := time.Date(monthStart.Year(), monthStart.Month()+1, 0, 23, 59, 59, 0, time.UTC)

		db = db.Where("summary_date BETWEEN ? AND ?", monthStart.Format("2006-01-02"), monthEnd.Format("2006-01-02"))

		// FIX: gunakan CEIL(DAY/7.0) agar week 5 tetap muncul jika bulan punya 29-31 hari.
		// Formula: day 1-7 = week 1, 8-14 = week 2, 15-21 = week 3, 22-28 = week 4, 29+ = week 5
		periodSelectExpr = `
			CEIL(EXTRACT(DAY FROM summary_date) / 7.0)::int AS week_num,
			TO_CHAR(summary_date, 'YYYY-MM') AS period_label
		`
		groupExtra = ", week_num"
		orderExpr = "week_num ASC, period_label ASC"
		isWeekly = true

		// FIX: filter by specific week number jika bukan "all"
		if o.Week != "" && o.Week != "all" {
			weekNum := 0
			_, scanErr := fmt.Sscanf(o.Week, "%d", &weekNum)
			if scanErr == nil && weekNum > 0 {
				// Hitung rentang tanggal untuk week tersebut
				startDay := (weekNum-1)*7 + 1
				endDay := weekNum * 7

				// Clamp endDay ke akhir bulan
				lastDay := monthEnd.Day()
				if endDay > lastDay {
					endDay = lastDay
				}

				startDate := fmt.Sprintf("%s-%02d", o.Month, startDay)
				endDate := fmt.Sprintf("%s-%02d", o.Month, endDay)

				// Ganti kondisi WHERE date: pakai rentang yang lebih spesifik
				db = db.Where("EXTRACT(DAY FROM summary_date) BETWEEN ? AND ?", startDay, endDay)
				_ = startDate
				_ = endDate
			}
		}

	case "monthly":
		if o.Year == "" {
			return nil, "", "", "", false, fmt.Errorf("year required for monthly period")
		}
		db = db.Where("summary_date BETWEEN ? AND ?", o.Year+"-01-01", o.Year+"-12-31")
		periodSelectExpr = `TO_CHAR(summary_date, 'YYYY-MM') AS period_label`
		orderExpr = "period_label ASC"

	case "custom":
		if o.DateFrom == "" || o.DateTo == "" {
			return nil, "", "", "", false, fmt.Errorf("date_from and date_to required for custom period")
		}
		db = db.Where("summary_date BETWEEN ? AND ?", o.DateFrom, o.DateTo)
		periodSelectExpr = `TO_CHAR(summary_date, 'YYYY-MM-DD') AS period_label`
		orderExpr = "period_label ASC"

	default:
		return nil, "", "", "", false, fmt.Errorf("invalid period_type: must be weekly, monthly, or custom")
	}

	return db, periodSelectExpr, orderExpr, groupExtra, isWeekly, nil
}