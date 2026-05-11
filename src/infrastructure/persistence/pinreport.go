package persistence

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
	"gorm.io/gorm/clause"
)

func (r *BaseModel) getLastPayout(country, operator, adnet, service string, dateSend time.Time) (float64, float64) {

	var prev entity.ApiPinReport

	err := r.DB.
		Select("payout_adn, payout_af").
		Where(
			`country = ?
			 AND operator = ?
			 AND adnet = ?
			 AND service = ?
			 AND date_send < ?`,
			country, operator, adnet, service, dateSend,
		).
		Order("date_send DESC").
		Limit(1).
		First(&prev).Error

	if err != nil {
		return 0, 0
	}

	return prev.PayoutAdn, prev.PayoutAF
}

func (r *BaseModel) getTodayPayout(country, operator, adnet, service string, dateSend time.Time) (float64, float64) {

	var res struct {
		PayoutAdn float64
		PayoutAF  float64
	}

	r.DB.Model(&entity.ApiPinReport{}).
		Select("payout_adn, payout_af").
		Where(`
			date_send = ? AND
			country = ? AND
			operator = ? AND
			adnet = ? AND
			service = ? 
		`,
			dateSend,
			country,
			operator,
			adnet,
			service,
		).
		Take(&res)

	return res.PayoutAdn, res.PayoutAF
}

func (r *BaseModel) PinReport(o entity.ApiPinReport) int {

	if o.PayoutAdn == 0 && o.PayoutAF == 0 {

		todayAdn, todayAF := r.getTodayPayout(
			o.Country,
			o.Operator,
			o.Adnet,
			o.Service,
			o.DateSend,
		)

		if todayAdn == 0 && todayAF == 0 {
			prevAdn, prevAF := r.getLastPayout(
				o.Country,
				o.Operator,
				o.Adnet,
				o.Service,
				o.DateSend,
			)

			o.PayoutAdn = prevAdn
			o.PayoutAF = prevAF
		} else {
			o.PayoutAdn = todayAdn
			o.PayoutAF = todayAF
		}
	}

	entity.BuildPinReportCalculation(&o)

	result := r.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "date_send"},
			{Name: "country"},
			{Name: "adnet"},
			{Name: "operator"},
			{Name: "service"},
			{Name: "campaign_id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"payout_adn",
			"payout_af",
			"total_mo",
			"total_postback",
			"sbaf",
			"saaf",
			"price_per_mo",
			"waki_revenue",
			"updated_at",
		}),
	}).Create(&o)

	r.Logs.Debugf(
		"[PIN_REPORT] affected=%d error=%v | %s %s %s %s %s",
		result.RowsAffected,
		result.Error,
		o.DateSend.Format("2006-01-02"),
		o.Country,
		o.Operator,
		o.Adnet,
		o.Service,
	)

	return int(o.ID)
}

func (r *BaseModel) GetDisplayPinReport(o entity.DisplayPinReport) ([]entity.ApiPinReport, int64, error) {
	var totalRows int64
	var ss []entity.ApiPinReport

	query := r.DB.Model(&entity.ApiPinReport{}).Select(`
		api_pin_reports.*,
		(payout_af * total_postback) AS saaf,
		(payout_adn * total_postback) AS sbaf,
		(CASE WHEN total_mo > 0 THEN (payout_af * total_postback) / total_mo ELSE 0 END) AS price_per_mo,
		((payout_af * total_postback) - (payout_adn * total_postback)) AS waki_revenue
	`)

	if o.Action == "Search" {
		if o.CampaignId != "" {
			query = query.Where("campaign_id = ?", o.CampaignId)
		}
		if o.Country != "" {
			query = query.Where("country = ?", o.Country)
		}
		if o.Company != "" {
			query = query.Where("company = ?", o.Company)
		}
		if o.Operator != "" {
			query = query.Where("operator = ?", o.Operator)
		}
		if len(o.Adnets) > 0 {
			query = query.Where("adnet IN ?", o.Adnets)
		}
		if o.Service != "" {
			query = query.Where("service = ?", o.Service)
		}

		if o.DateRange != "" {
			switch strings.ToUpper(o.DateRange) {
			case "TODAY":
				query = query.Where("date_send = CURRENT_DATE")
			case "YESTERDAY":
				query = query.Where("date_send = CURRENT_DATE - INTERVAL '1 DAY'")
			case "LAST7DAY":
				query = query.Where("date_send BETWEEN CURRENT_DATE - INTERVAL '7 DAY' AND CURRENT_DATE")
			case "LAST30DAY":
				query = query.Where("date_send BETWEEN CURRENT_DATE - INTERVAL '30 DAY' AND CURRENT_DATE")
			case "THISMONTH":
				query = query.Where("date_send >= DATE_TRUNC('month', CURRENT_DATE)")
			case "LASTMONTH":
				query = query.Where("date_send BETWEEN DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 MONTH') AND DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '1 DAY'")
			case "CUSTOMRANGE":
				query = query.Where("date_send BETWEEN ? AND ?", o.DateBefore, o.DateAfter)
			default:
				query = query.Where("date_send = ?", o.DateRange)
			}
		} else {
			query = query.Where("date_send = CURRENT_DATE")
		}
	} else {
		query = query.Where("date_send = CURRENT_DATE")
	}

	if err := query.Count(&totalRows).Error; err != nil {
		return []entity.ApiPinReport{}, 0, err
	}

	if o.OrderColumn != "" {
		dir := "ASC"
		if strings.ToUpper(o.OrderDir) == "DESC" {
			dir = "DESC"
		}
	
		switch o.OrderColumn {
		case "saaf":
			query = query.Order(fmt.Sprintf("(payout_af * total_postback) %s", dir))
		case "sbaf":
			query = query.Order(fmt.Sprintf("(payout_adn * total_postback) %s", dir))
		case "price_per_mo":
			query = query.Order(fmt.Sprintf("(CASE WHEN total_mo > 0 THEN (payout_af * total_postback) / total_mo ELSE 0 END) %s", dir))
		case "waki_revenue":
			query = query.Order(fmt.Sprintf("((payout_af * total_postback) - (payout_adn * total_postback)) %s", dir))
		default:
			query = query.Order(fmt.Sprintf("%s %s", o.OrderColumn, dir))
		}
	} else {
		query = query.Order("date_send DESC").Order("id DESC")
	}	

	if err := query.
		Limit(o.PageSize).
		Offset((o.Page - 1) * o.PageSize).
		Find(&ss).Error; err != nil {
		return []entity.ApiPinReport{}, 0, err
	}

	return ss, totalRows, nil
}

func (r *BaseModel) EditPayoutAPIReport(o entity.ApiPinReport) error {

	dateOnly := o.DateSend.Format("2006-01-02")

	var existing entity.ApiPinReport

	err := r.DB.
		Where(`
			date_send = ? AND
			country = ? AND
			operator = ? AND
			service = ? AND
			adnet = ?
		`,
			dateOnly,
			o.Country,
			o.Operator,
			o.Service,
			o.Adnet,
		).
		First(&existing).Error

	if err != nil {
		return err
	}

	if o.PayoutAF == 0 {
		o.PayoutAF = existing.PayoutAF
	}
	if o.PayoutAdn == 0 {
		o.PayoutAdn = existing.PayoutAdn
	}

	o.TotalPostback = existing.TotalPostback
	o.TotalMO = existing.TotalMO

	entity.BuildPinReportCalculation(&o)

	result := r.DB.
		Model(&entity.ApiPinReport{}).
		Where(`
			date_send = ? AND
			country = ? AND
			operator = ? AND
			service = ? AND
			adnet = ?
		`,
			dateOnly,
			o.Country,
			o.Operator,
			o.Service,
			o.Adnet,
		).
		Updates(map[string]interface{}{
			"payout_af":     o.PayoutAF,
			"payout_adn":    o.PayoutAdn,
			"sbaf":         o.SBAF,
			"saaf":         o.SAAF,
			"price_per_mo":  o.PricePerMO,
			"waki_revenue":  o.WakiRevenue,
		})

	r.Logs.Debug(fmt.Sprintf(
		"affected: %d, error: %#v",
		result.RowsAffected,
		result.Error,
	))

	return result.Error
}

func (r *BaseModel) EditCpaAPIPerformanceReport(o entity.ApiPinPerformance) error {

	dateOnly := o.DateSend.Format("2006-01-02")

	var existing entity.ApiPinPerformance
	if err := r.DB.Where(`
		date_send = ? AND
		country = ? AND
		operator = ? AND
		service = ? AND
		adnet = ?
	`,
		dateOnly,
		o.Country,
		o.Operator,
		o.Service,
		o.Adnet,
	).First(&existing).Error; err != nil {
		return err
	}

	if o.CPA == 0 {
		o.CPA = existing.CPA
	}
	if o.CPAWaki == 0 {
		o.CPAWaki = existing.CPAWaki
	}

	o.PinOkRatio = existing.PinOkRatio
	o.PinOK = existing.PinOK
	o.ChargedMO = existing.ChargedMO

	entity.BuildAPIPerformanceReportCalculation(&o)


	result := r.DB.
		Model(&entity.ApiPinPerformance{}).
		Where(`
			date_send = ? AND
			country = ? AND
			operator = ? AND
			service = ? AND
			adnet = ?
		`,
			dateOnly,
			o.Country,
			o.Operator,
			o.Service,
			o.Adnet,
		).
		Updates(map[string]interface{}{
			"cpa":     o.CPA,
			"cpa_waki":    o.CPAWaki,
			"total_spending":         o.TotalSpending,
			"total_spending_after_waki":         o.TotalSpendingAfterWaki,
			"cac":  o.CAC,
			"paid_cac":  o.PaidCAC,
		})

	r.Logs.Debug(fmt.Sprintf(
		"affected: %d, error: %#v",
		result.RowsAffected,
		result.Error,
	))

	return result.Error
}

func (r *BaseModel) EditArpuAPIPerformanceReport(o entity.ApiPinPerformance) error {

	dateOnly := o.DateSend.Format("2006-01-02")

	result := r.DB.
		Model(&entity.ApiPinPerformance{}).
		Where(`
			date_send = ? AND
			country = ? AND
			operator = ? AND
			service = ? AND
			adnet = ?
		`,
			dateOnly,
			o.Country,
			o.Operator,
			o.Service,
			o.Adnet,
		).
		Updates(map[string]interface{}{
			"estimated_arpu": o.EstimatedARPU,
			"updated_at":     time.Now(),
		})

	r.Logs.Debug(fmt.Sprintf(
		"affected: %d, error: %#v",
		result.RowsAffected,
		result.Error,
	))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("no data updated (data not found)")
	}

	return nil
}

func (r *BaseModel) UpsertApiPerformanceReport(o *entity.ApiPinPerformance) error {

	// Ambil sum saaf, total_mo & total_postback dari api_pin_reports
	var saafSum float64
	r.DB.Model(&entity.ApiPinReport{}).
		Select("saaf").
		Where("country=? AND operator=? AND adnet=? AND service=? AND date_send=?",
			o.Country, o.Operator, o.Adnet, o.Service, o.DateSend).
		Scan(&saafSum)
	o.SAAF = saafSum

	// Ambil estimated ARPU dari api_performance_reports
	var estimatedARPU float64
	r.DB.Model(&entity.ApiPinPerformance{}).
		Select("estimated_arpu").
		Where("country=? AND operator=? AND adnet=? AND service=? AND date_send=?",
			o.Country, o.Operator, o.Adnet, o.Service, o.DateSend).
		Scan(&estimatedARPU)
	o.EstimatedARPU = estimatedARPU

	// Ambil CPA & CPAWaki hari ini (kalau sudah ada)
	var today struct {
		CPA     float64
		CPAWaki float64
	}

	r.DB.Model(&entity.ApiPinPerformance{}).
		Select("cpa, cpa_waki").
		Where("country=? AND operator=? AND adnet=? AND service=? AND date_send=?",
			o.Country, o.Operator, o.Adnet, o.Service, o.DateSend).
		Take(&today)

	// Ambil prev CPA & CPA after waki
	var prev struct {
		CPA     float64
		CPAWaki float64
	}
	r.DB.Raw(`SELECT cpa, cpa_waki 
		FROM api_pin_performances
		WHERE country=? AND operator=? AND service=? AND adnet=? LIMIT 1`,
		o.Country, o.Operator, o.Service, o.Adnet).Order("date_send DESC").Scan(&prev)

	if o.CPA == 0 && o.CPAWaki == 0 {
		if today.CPA == 0 && today.CPAWaki == 0 {
			o.CPA = prev.CPA
			o.CPAWaki = prev.CPAWaki
		} else {
			o.CPA = today.CPA
			o.CPAWaki = today.CPAWaki
		}
	}

	// Hitung semua field
	o.CpaPerPO = 0
	o.TotalSpending = 0
	o.TotalSpendingAfterWaki = 0
	o.PaidCAC = 0
	o.CAC = 0
	o.SubsCR = 0

	if o.PinOK > 0 && o.SAAF > 0 {
		o.CpaPerPO = o.SAAF / float64(o.PinOK)
	}

	if o.PinOkRatio > 0 {
		o.TotalSpending = o.CPA * float64(o.PinOkRatio)
		o.TotalSpendingAfterWaki = o.CPAWaki * float64(o.PinOkRatio)
		o.CAC = o.TotalSpendingAfterWaki / float64(o.PinOK)
		if o.ChargedMO > 0 {
			o.PaidCAC = o.TotalSpendingAfterWaki / o.ChargedMO
		}
	}

	if o.PinVerifyRequest > 0 {
		o.SubsCR = float64(o.PinOK) / float64(o.PinVerifyRequest)
		o.AdnetCR = float64(o.PinOkRatio) / float64(o.PinVerifyRequest)
	}

	// Upsert ke api_performance_reports
	result := r.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "date_send"}, {Name: "country"}, {Name: "operator"}, {Name: "adnet"}, {Name: "service"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"pin_request", "unique_pin_request", "pin_success", "pin_failed", "pin_verify_request",
			"pin_verify_request_unique", "pin_ok", "pin_not_ok", "pin_ok_ratio", "saaf",
			"cac", "total_spending", "paid_cac", "subs_cr", "adnet_cr", "estimated_arpu", "cpa", "cpa_waki",
			"cpa_per_po", "total_spending_after_waki", "updated_at",
		}),
	}).Create(o)

	if result.Error != nil {
		return result.Error
	}

	fmt.Printf("[API_PERFORMANCE] inserted/updated %s %s %s %s\n",
		o.DateSend.Format("2006-01-02"), o.Country, o.Operator, o.Adnet)

	return nil
}

func (r *BaseModel) PinPerformanceReport(o entity.ApiPinPerformance) int {

	result := r.DB.Create(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return int(o.ID)
}

func (r *BaseModel) GetApiPinPerformanceReport(o entity.DisplayPinPerformanceReport) ([]entity.ApiPinPerformance, int64, error) {

	var (
		rows       *sql.Rows
		total_rows int64
	)

	// Apply filters, minus the pagination constraints
	query := r.DB.Model(&entity.ApiPinPerformance{})
	if o.Action == "Search" {
		if o.Country != "" {
			query = query.Where("country = ?", o.Country)
		}
		if o.Operator != "" {
			query = query.Where("operator = ?", o.Operator)
		}
		if o.Service != "" {
			query = query.Where("service = ?", o.Service)
		}
		if len(o.Adnets) > 0 {
			query = query.Where("adnet IN ?", o.Adnets)
		}
		if o.DateRange != "" {
			switch strings.ToUpper(o.DateRange) {
			case "TODAY":
				query = query.Where("date_send = CURRENT_DATE")
			case "YESTERDAY":
				query = query.Where("date_send BETWEEN CURRENT_DATE - INTERVAL '1 DAY' AND CURRENT_DATE")
			case "LAST7DAY":
				query = query.Where("date_send BETWEEN CURRENT_DATE - INTERVAL '7 DAY' AND CURRENT_DATE")
			case "LAST30DAY":
				query = query.Where("date_send BETWEEN CURRENT_DATE - INTERVAL '30 DAY' AND CURRENT_DATE")
			case "THISMONTH":
				query = query.Where("date_send >= DATE_TRUNC('month', CURRENT_DATE)")
			case "LASTMONTH":
				query = query.Where("date_send BETWEEN DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 MONTH') AND DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '1 DAY'")
			case "CUSTOMRANGE":
				query = query.Where("date_send BETWEEN ? AND ?", o.DateBefore, o.DateAfter)
			case "ALLDATERANGE":
			default:
				query = query.Where("date_send = ?", o.DateRange)
			}
		} else {
			query = query.Where("date_send = CURRENT_DATE")
		}
	}

	// Get the total count after applying filters
	query.Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, _ = query_limit.Order("subs_cr").Rows()
	defer rows.Close()

	var ss []entity.ApiPinPerformance
	for rows.Next() {
		var s entity.ApiPinPerformance
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	r.Logs.Debug(fmt.Sprintf("Total data : %d ...\n", len(ss)))

	return ss, total_rows, rows.Err()
}

func (r *BaseModel) GetConversionLogReport(o entity.DisplayConversionLogReport) ([]entity.PixelStorage, int64, error) {
	var (
		rows       *sql.Rows
		total_rows int64
	)

	tableName := "pixel_storages"

	if strings.ToUpper(o.DateRange) == "YESTERDAY" {
		yesterday := time.Now().AddDate(0, 0, -1).Format("20060102")
		tableName = fmt.Sprintf("pixel_storages_%s", yesterday)
	} else if strings.ToUpper(o.DateRange) == "2DAYAGO" {
		twoDaysAgo := time.Now().AddDate(0, 0, -2).Format("20060102")
		tableName = fmt.Sprintf("pixel_storages_%s", twoDaysAgo)
	}

	query := r.DB.Table(tableName)
	query = query.Where("is_used = ?", "true")

	if o.CampaignType == "mainstream" {
		query = query.Where("campaign_objective LIKE ?", "%MAINSTREAM%")
	} else {
		query = query.Where("campaign_objective IN ?", []string{"CPA", "CPC", "CPI", "CPM", "SINGLE URL S2S"})
	}

	if o.Action == "Search" {
		if o.Country != "" {
			query = query.Where("country = ?", o.Country)
		}
		if o.Operator != "" {
			query = query.Where("operator = ?", o.Operator)
		}
		if o.Pixel != "" {
			query = query.Where("pixel = ?", o.Pixel)
		}
		if o.CampaignId != "" {
			query = query.Where("campaign_id = ?", o.CampaignId)
		}
		if o.CampaignName != "" {
			query = query.Where("campaign_name = ?", o.CampaignName)
		}
		if o.CampaignType == "mainstream" {
			if o.StatusPostback != "" {
				query = query.Where("status_postback = ?", o.StatusPostback)
			}
			if o.Agency != "" {
				query = query.Where("adnet = ?", o.Agency)
			}
		}
		if o.CampaignType == "s2s" {
			if o.StatusPostback != "" {
				query = query.Where("status_postback = ?", o.StatusPostback)
			}
			if o.Adnet != "" {
				query = query.Where("adnet = ?", o.Adnet)
			}
		}

		if o.DateRange != "" && strings.ToUpper(o.DateRange) != "YESTERDAY" && strings.ToUpper(o.DateRange) != "2DAYAGO" {
			switch strings.ToUpper(o.DateRange) {
			case "TODAY":
				query = query.Where("pxdate BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '1 day' - INTERVAL '1 second'")
			case "LAST30DAY":
				query = query.Where("pxdate BETWEEN CURRENT_DATE - INTERVAL '30 DAY' AND CURRENT_DATE")
			case "THISMONTH":
				query = query.Where("pxdate >= DATE_TRUNC('month', CURRENT_DATE)")
			case "LASTMONTH":
				query = query.Where("pxdate BETWEEN DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 MONTH') AND DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '1 DAY'")
			case "CUSTOMRANGE":
				dateEnd, _ := time.Parse("2006-01-02", o.DateEnd)
				query = query.Where("pxdate BETWEEN ? AND ?", o.DateStart, dateEnd.AddDate(0, 0, 1))
			case "LAST7DAY":
				query = query.Where("pxdate BETWEEN CURRENT_DATE - INTERVAL '7 DAY' AND CURRENT_DATE")
			case "ALLDATERANGE":
			default:
				query = query.Where("pxdate = ?", o.DateRange)
			}
		}
	}

	query.Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	var startIndex int
	if o.Order == "asc" {
		query_limit = query_limit.Order("pxdate asc")
		startIndex = int(total_rows) - ((o.Page - 1) * o.PageSize)
	} else {
		query_limit = query_limit.Order("pxdate desc")
		startIndex = (o.Page - 1) * o.PageSize
	}

	rows, _ = query_limit.Rows()
	defer rows.Close()

	var ss []entity.PixelStorage
	for rows.Next() {
		var s entity.PixelStorage
		r.DB.ScanRows(rows, &s)

		if o.Order == "asc" {
			s.ID = startIndex
			startIndex--
		} else {
			s.ID = startIndex + 1
			startIndex++
		}
		ss = append(ss, s)
	}
	return ss, total_rows, rows.Err()
}

func (r *BaseModel) GetPerformanceReport(o entity.PerformaceReportParams) ([]entity.PerformanceReport, int64, error) {

	var (
		rows      *sql.Rows
		totalRows int64
		dateStart time.Time
		dateEnd   time.Time
	)

	now := time.Now()

	if o.DateRange != "" {
		switch strings.ToUpper(o.DateRange) {

		case "TODAY":
			dateStart = now
			dateEnd = now

		case "YESTERDAY":
			dateStart = now.AddDate(0, 0, -1)
			dateEnd = dateStart

		case "LAST7DAY":
			dateStart = now.AddDate(0, 0, -7)
			dateEnd = now

		case "LAST30DAY":
			dateStart = now.AddDate(0, 0, -30)
			dateEnd = now

		case "THISMONTH":
			dateStart = time.Date(
				now.Year(), now.Month(), 1,
				0, 0, 0, 0, now.Location(),
			)
			dateEnd = now

		case "LASTMONTH":
			first := time.Date(
				now.Year(), now.Month()-1, 1,
				0, 0, 0, 0, now.Location(),
			)
			dateStart = first
			dateEnd = first.AddDate(0, 1, -1)

		case "CUSTOMRANGE":
			var err error

			if o.DateBefore == "" || o.DateAfter == "" {
				return nil, 0, fmt.Errorf("date_before and date_after are required")
			}

			dateStart, err = time.Parse("2006-01-02", o.DateBefore)
			if err != nil {
				return nil, 0, fmt.Errorf("invalid date_before: %v", err)
			}

			dateEnd, err = time.Parse("2006-01-02", o.DateAfter)
			if err != nil {
				return nil, 0, fmt.Errorf("invalid date_after: %v", err)
			}

		default:
			dateStart = now
			dateEnd = now
		}
	} else {
		dateStart = now
		dateEnd = now
	}

	query := r.DB.Model(&entity.SummaryCampaign{})
	query = query.Where("mo_received > 0")

	query.Select(`
		country, company, client_type, campaign_name, partner,
		operator, service, adnet,
		SUM(mo_received) AS pixel_received,
		SUM(postback) AS pixel_send,
		SUM(cr_postback) AS cr_postback,
		SUM(cr_mo) AS cr_mo,
		SUM(landing) AS landing,
		MAX(ratio_send) AS ratio_send,
		MAX(ratio_receive) AS ratio_receive,
		SUM(po) AS price_per_postback,
		SUM(cost_per_conversion) AS cost_per_conversion,
		SUM(agency_fee) AS agency_fee,
		SUM(postback * po) AS spending_to_adnets,
		SUM(total_waki_agency_fee),
		SUM(total_waki_agency_fee + po * postback) AS total_spending,
		SUM(cpa) AS e_cpa,
		SUM(total_fp) AS total_fp,
		SUM(success_fp) AS success_fp
	`)

	if o.Action == "Search" {

		if o.Country != "" {
			query = query.Where("country = ?", o.Country)
		}
		if o.Operator != "" {
			query = query.Where("operator = ?", o.Operator)
		}
		if o.Partner != "" {
			query = query.Where("partner = ?", o.Partner)
		}
		if o.Company != "" {
			query = query.Where("company = ?", o.Company)
		}
		if o.CampaignType != "" {
			query = query.Where("campaign_objective = ?", o.CampaignType)
		}
		if o.CampaignId != "" {
			query = query.Where("campaign_id = ?", o.CampaignId)
		}
		if o.CampaignName != "" {
			query = query.Where("campaign_name = ?", o.CampaignName)
		}
		if o.ClientType != "" {
			query = query.Where("client_type = ?", o.ClientType)
		}
		if o.Publisher != "" {
			query = query.Where("adnet = ?", o.Publisher)
		}
		if o.Service != "" {
			query = query.Where("service = ?", o.Service)
		}

		query = query.Where(
			"summary_date BETWEEN ? AND ?",
			dateStart,
			dateEnd,
		)
	}

	query.Group("country, company, client_type, campaign_name, partner, operator, service, adnet")

	query.Unscoped().Count(&totalRows)

	if o.PageSize > 0 {
		query = query.Limit(o.PageSize)
	}
	if o.Page > 0 {
		query = query.Offset((o.Page - 1) * o.PageSize)
	}

	rows, err := query.Order("country").Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []entity.PerformanceReport

	for rows.Next() {
		var s entity.PerformanceReport
		r.DB.ScanRows(rows, &s)
		r.GetARPUReport(&s, dateStart, dateEnd)

		result = append(result, s)
	}

	return result, totalRows, rows.Err()
}

func (r *BaseModel) GetDistinctPerformanceReport(o entity.SummaryCampaign) ([]entity.SummaryCampaign, error) {

	var (
		rows *sql.Rows
	)

	query := r.DB.Model(&entity.SummaryCampaign{})

	rows, _ = query.Unscoped().Distinct("summary_date", "country", "operator", "service").Where("summary_date = CURRENT_DATE").Order("summary_date").Rows()

	defer rows.Close()

	var (
		ss []entity.SummaryCampaign
	)

	for rows.Next() {

		var s entity.SummaryCampaign

		// ScanRows scans a row into a struct
		r.DB.ScanRows(rows, &s)

		ss = append(ss, s)
	}

	return ss, rows.Err()
}

func (r *BaseModel) GetARPUReport(s *entity.PerformanceReport, dateStart time.Time, dateEnd time.Time) {

	var apiResponse entity.ARPUResponse
	var isempty bool
	key := fmt.Sprintf("%s_%s_%s_%s_%s", s.Country, s.Operator, s.Service, dateStart, dateEnd)

	if apiResponse, isempty = r.RGetArpuReport(key, "$"); isempty {
		fmt.Println("IS Empty: ", isempty)
		apiBase := os.Getenv("APIARPU")
		if apiBase == "" {
			fmt.Println("Missing APIARPU environment variable")
			return
		}

		// Build base URL
		baseURL, err := url.Parse(apiBase + "/api/v4/arpu/arpu90")
		if err != nil {
			fmt.Println("Failed to parse base URL:", err)
			return
		}

		// Manually encode all query params
		query := fmt.Sprintf(
			"from=%s&to=%s&country=%s&operator=%s&service=%s",
			url.QueryEscape(dateStart.Format("2006-01-02")),
			url.QueryEscape(dateEnd.Format("2006-01-02")),
			url.QueryEscape(s.Country),
			url.QueryEscape(s.Operator),
			url.QueryEscape(s.Service), // encodes spaces as %20
		)

		baseURL.RawQuery = query

		// Make the request
		req, err := http.NewRequest("GET", baseURL.String(), nil)
		if err != nil {
			fmt.Println("Failed to create request:", err)
			return
		}

		req.Header.Add("Authorization", "Basic bWlkZGxld2FyZTpsMW5rMXQzNjA=")
		req.Header.Add("Accept", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Failed to make request:", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Failed to read response:", err)
			return
		}

		err = json.Unmarshal(body, &apiResponse)
		if err != nil {
			fmt.Println("Failed to parse JSON:", err)
			return
		}

		s, _ := json.Marshal(apiResponse)

		r.SetData(key, "$", string(s))
		r.SetExpireData(key, 60)
	}

	if apiResponse.Data == nil || len(apiResponse.Data.Data) == 0 {
		return
	}
	for _, item := range apiResponse.Data.Data {

		if item.Adnet == s.Adnet {
			s.ARPUROI = s.ECPA / (item.Arpu90Net / 3)
			s.ARPU90 = item.Arpu90Net
			if s.TotalFP > 0 {
				s.BillrateFP = s.SuccessFP / s.TotalFP * 100
			}
		}
	}
}

func (r *BaseModel) GetAPIReportById(id []string) ([]entity.ApiPinReport, error) {
	query := r.DB.Model(&entity.ApiPinReport{}).Where("id IN ?", id)
	rows, err := query.Rows()

	if err != nil {
		return []entity.ApiPinReport{}, err
	}

	var results []entity.ApiPinReport
	for rows.Next() {
		var s entity.ApiPinReport
		r.DB.ScanRows(rows, &s)
		results = append(results, s)
	}

	return results, nil
}

func (r *BaseModel) GetOperatorAliases() ([]entity.OperatorAlias, error) {
	var res []entity.OperatorAlias
	err := r.DB.
		Table("operator_aliases").
		Where("type = ?", "API").
		Find(&res).Error
	return res, err
}