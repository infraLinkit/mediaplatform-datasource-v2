package persistence

import (
	"database/sql"
	"errors"
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"
	"time"

	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *BaseModel) CheckSumDate(o entity.SummaryCampaign) int {

	var result int64
	r.DB.Table("summary_campaigns").Where("summary_date = ? AND url_service_key = ?", o.SummaryDate, o.URLServiceKey).Count(&result)
	return int(result)
}

func (r *BaseModel) GetSummaryCampaign(o entity.SummaryCampaign) (entity.SummaryCampaign, bool) {

	result := r.DB.Model(&o).
		Where("summary_date = ? AND url_service_key = ?", o.SummaryDate, o.URLServiceKey).
		First(&o)

	b := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if b {
		r.Logs.Warn(fmt.Sprintf("Campaign id not found %#v", o))
		return o, false
	} else {
		return o, true
	}
}

func (r *BaseModel) DelSummaryCampaign(o entity.SummaryCampaign) error {

	result := r.DB.
		Where("summary_date = ? AND url_service_key = ? AND country = ? AND operator = ? AND partner = ? AND service = ? AND adnet = ? AND campaign_id = ?", o.SummaryDate, o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Service, o.Adnet, o.CampaignId).
		Delete(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) EditSettingSummaryCampaign(o entity.SummaryCampaign) error {

	summaryDate := o.SummaryDate
	if summaryDate.IsZero() {
		summaryDate = time.Now()
	}

	result := r.DB.Model(&o).
		Where("summary_date = ? AND url_service_key = ? AND country = ? AND operator = ? AND partner = ? AND service = ? AND adnet = ? AND campaign_id = ?", summaryDate, o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Service, o.Adnet, o.CampaignId).
		Updates(entity.SummaryCampaign{PO: o.PO, MOLimit: o.MOLimit, RatioSend: o.RatioSend, RatioReceive: o.RatioReceive})

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) UpdateSummaryPO(o entity.SummaryCampaign) error {
	summaryDate := o.SummaryDate
	if summaryDate.IsZero() {
		var err error
		summaryDate, err = time.Parse(
			"2006-01-02",
			time.Now().Format("2006-01-02"),
		)
		if err != nil {
			return err
		}
	}

	result := r.DB.Exec(`
		UPDATE summary_campaigns
		SET po = ?
		WHERE summary_date = ?
		  AND url_service_key = ?
	`,
		o.PO,
		summaryDate,
		o.URLServiceKey,
	)

	r.Logs.Debug(fmt.Sprintf("UpdateSummaryPO affected: %d, error: %#v", result.RowsAffected, result.Error))
	return result.Error
}

func (r *BaseModel) UpdateSummaryRatio(o entity.SummaryCampaign) error {
	summaryDate := o.SummaryDate
	if summaryDate.IsZero() {
		var err error
		summaryDate, err = time.Parse(
			"2006-01-02",
			time.Now().Format("2006-01-02"),
		)
		if err != nil {
			return err
		}
	}

	result := r.DB.Exec(`
		UPDATE summary_campaigns
		SET ratio_send = ?, ratio_receive = ?
		WHERE summary_date = ?
		  AND url_service_key = ?
	`,
		o.RatioSend,
		o.RatioReceive,
		summaryDate,
		o.URLServiceKey,
	)

	r.Logs.Debug(fmt.Sprintf("UpdateSummaryRatio affected: %d, error: %#v", result.RowsAffected, result.Error))
	return result.Error
}

func (r *BaseModel) UpdateSummaryMOCapping(o entity.SummaryCampaign) error {
	summaryDate := o.SummaryDate
	if summaryDate.IsZero() {
		var err error
		summaryDate, err = time.Parse(
			"2006-01-02",
			time.Now().Format("2006-01-02"),
		)
		if err != nil {
			return err
		}
	}

	result := r.DB.Exec(`
		UPDATE summary_campaigns
		SET mo_limit = ?
		WHERE summary_date = ?
		  AND url_service_key = ?
	`,
		o.MOLimit,
		summaryDate,
		o.URLServiceKey,
	)

	r.Logs.Debug(fmt.Sprintf("UpdateSummaryMOCapping affected: %d, error: %#v", result.RowsAffected, result.Error))
	return result.Error
}

func (r *BaseModel) EditPOAFIncSummaryCampaign(o entity.IncSummaryCampaign) error {
	query := `
		UPDATE inc_summary_campaigns
		SET poaf = ?
		WHERE summary_date = ?
		  AND url_service_key = ?
	`

	result := r.DB.Exec(query, o.POAF, o.SummaryDate, o.URLServiceKey)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error: %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) EditPOAFSummaryCampaign(o entity.SummaryCampaign) error {

	result := r.DB.Model(&o).
		Where("summary_date = ? AND url_service_key = ?", o.SummaryDate, o.URLServiceKey).
		Updates(entity.SummaryCampaign{POAF: o.POAF})

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) UpdateSummaryCampaign(o entity.SummaryCampaign) error {

	result := r.DB.Model(&o).
		Where("summary_date = ? AND url_service_key = ? AND country = ? AND operator = ? AND partner = ? AND service = ? AND adnet = ? AND campaign_id = ?", o.SummaryDate, o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Service, o.Adnet, o.CampaignId).
		Updates(entity.SummaryCampaign{Status: o.Status})

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) SummaryCampaign(o entity.SummaryCampaign) int {

	result := r.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "summary_date"},
			{Name: "campaign_id"},
			{Name: "campaign_objective"},
			{Name: "country"},
			{Name: "partner"},
			{Name: "operator"},
			{Name: "url_service_key"},
			{Name: "service"},
			{Name: "adnet"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"traffic":               o.Traffic,
			"landing":               o.Landing,
			"mo_received":           o.MoReceived,
			"cr_mo":                 o.CrMO,
			"cr_postback":           o.CrPostback,
			"postback":              o.Postback,
			"total_fp":              o.TotalFP,
			"success_fp":            o.SuccessFP,
			"billrate":              o.Billrate,
			"po":                    o.PO,
			"sbaf":                  o.SBAF,
			"saaf":                  o.SAAF,
			"cpa":                   o.CPA,
			"revenue":               o.Revenue,
			"url_after":             o.URLAfter,
			"url_before":            o.URLBefore,
			"mo_limit":              o.MOLimit,
			"ratio_send":            o.RatioSend,
			"ratio_receive":         o.RatioReceive,
			"client_type":           o.ClientType,
			"cost_per_conversion":   o.CostPerConversion,
			"agency_fee":            o.AgencyFee,
			"total_waki_agency_fee": o.TotalWakiAgencyFee,
			"target_daily_budget":   o.TargetDailyBudget,
			"budget_usage":          o.BudgetUsage,
			"campaign_name":         o.CampaignName,
			"technical_fee":         o.TechnicalFee}),
	}).Create(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return o.ID
}

func (r *BaseModel) DataTraffic(o entity.DataTraffic) int {

	result := r.DB.Create(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return int(o.ID)
}

func (r *BaseModel) DataLanding(o entity.DataLanding) int {

	result := r.DB.Create(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return int(o.ID)
}

func (r *BaseModel) DataClicked(o entity.DataClicked) int {

	result := r.DB.Create(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return int(o.ID)
}

func (r *BaseModel) DataRedirect(o entity.DataRedirect) int {

	result := r.DB.Create(&o)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return int(o.ID)
}

func (r *BaseModel) UpdateCPAReportSummaryCampaign(o entity.SummaryCampaign) error {

	result := r.DB.Model(&o).
		Where("summary_date = ? AND url_service_key = ? AND country = ? AND operator = ? AND partner = ? AND service = ? AND adnet = ? AND campaign_id = ?", o.SummaryDate, o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Service, o.Adnet, o.CampaignId).
		Updates(entity.SummaryCampaign{CostPerConversion: o.CostPerConversion, AgencyFee: o.AgencyFee})

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error

}

func (r *BaseModel) UpdateReportSummaryCampaignMonitoringBudget(o entity.SummaryCampaign) error {
	result := r.DB.Model(&o).
		Where("EXTRACT(YEAR FROM summary_date) = ? AND EXTRACT(MONTH FROM summary_date) = ? AND country = ? AND operator = ?",
			o.SummaryDate.Year(), int(o.SummaryDate.Month()), o.Country, o.Operator).
		Updates(map[string]interface{}{
			"target_daily_budget":   o.TargetDailyBudget,
			"target_monthly_budget": o.TargetMonthlyBudget,
		})

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) GetSummaryCampaignMonitoring(params entity.ParamsCampaignSummary) ([]entity.CampaignSummaryMonitoring, time.Time, time.Time, error) {
	var (
		rows *sql.Rows
	)
	query := r.DB.Model(&entity.CampaignSummaryMonitoring{})
	query.Where("deleted_at IS NULL")
	// Apply Indicator Selection
	selectedFields := []string{"country", "url_service_key", "campaign_id", "campaign_name", "partner", "operator", "service", "adnet"}
	if params.DataType == "monthly_report" {
		selectedFields = append(selectedFields, "DATE_TRUNC('month', summary_date) as summary_date")
	} else {
		selectedFields = append(selectedFields, "summary_date")
	}

	formattedIndicators := formatQueryIndicators(params.DataIndicators, params.DataType)
	selectedFields = append(selectedFields, formattedIndicators...)

	query.Select(selectedFields)

	// Set default values
	if params.DataType == "" {
		params.DataType = "daily_report"
	}

	// Apply paramss
	if params.Country != "" {
		query.Where("country = ?", params.Country)
	}
	if params.UrlServiceKey != "" {
		query.Where("url_service_key = ?", params.UrlServiceKey)
	}
	if params.Operator != "" {
		query.Where("operator = ?", params.Operator)
	}
	if params.Adnet != "" {
		query.Where("adnet = ?", params.Adnet)
	}
	if params.PartnerName != "" {
		query.Where("partner = ?", params.PartnerName)
	}
	if params.Service != "" {
		query.Where("service = ?", params.Service)
	}
	if params.CampaignName != "" {
		query.Where("campaign_name = ?", params.CampaignName)
	}

	// Handle Date Range
	var dateStart, dateEnd time.Time
	var errStart, errEnd error
	today := time.Now()

	switch strings.ToUpper(params.DateRange) {
	case "TODAY":
		dateStart, dateEnd = today, today
	case "YESTERDAY":
		dateStart, dateEnd = today.AddDate(0, 0, -1), today.AddDate(0, 0, -1)
	case "LAST_7_DAY", "LAST_7_DAYS":
		dateStart, dateEnd = today.AddDate(0, 0, -6), today
	case "LAST_30_DAY", "LAST_30_DAYS":
		dateStart, dateEnd = today.AddDate(0, -1, 0), today
	case "THIS_MONTH":
		dateStart = time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
		dateEnd = today
	case "LAST_MONTH":
		lastMonth := today.AddDate(0, -1, 0)
		dateStart = time.Date(lastMonth.Year(), lastMonth.Month(), 1, 0, 0, 0, 0, today.Location())
		dateEnd = time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location()).AddDate(0, 0, -1)
	case "CUSTOM_RANGE":
		fmt.Println("dr", params.DateCustomRange)
		splitTime := strings.Split(params.DateCustomRange, "to")

		dateStart, errStart = time.Parse("2006-01-02", strings.TrimSpace(splitTime[0]))
		dateEnd, errEnd = time.Parse("2006-01-02", strings.TrimSpace(splitTime[1]))
		if errStart != nil {
			dateStart = today
		}
		if errEnd != nil {
			dateEnd = today
		}
	default:
		dateStart = time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
		dateEnd = today
	}

	if dateEnd.After(today) {
		dateEnd = today
	}

	query.Where("summary_date BETWEEN ? AND ?", dateStart, dateEnd)

	// Grouping for monthly reports
	if params.DataType == "monthly_report" {
		query.Group("DATE_TRUNC('month', summary_date),country, campaign_name, campaign_id, url_service_key, partner, operator, service, adnet")
	}

	rows, _ = query.Unscoped().Rows()

	defer rows.Close()

	var (
		ss []entity.CampaignSummaryMonitoring
	)

	for rows.Next() {
		var s entity.CampaignSummaryMonitoring
		// ScanRows scans a row into a struct
		r.DB.ScanRows(rows, &s)

		ss = append(ss, s)
	}
	return ss, dateStart, dateEnd, rows.Err()

}

func (r *BaseModel) GetSummaryCampaignChart(params entity.ParamsCampaignSummary) ([]entity.CampaignSummaryChart, time.Time, time.Time, error) {
	var (
		rows           *sql.Rows
		selectedFields []string
	)
	query := r.DB.Model(&entity.CampaignSummaryMonitoring{})

	if params.ChartType == "spending" {
		selectedFields = []string{"SUM(saaf) as spending,SUM(mo_received) as mo, summary_date"}
	} else if params.ChartType == "cr" {
		selectedFields = []string{"SUM(cr) as cr,SUM(mo_received) as mo, summary_date"}
	} else {
		selectedFields = []string{"SUM(cr) as cr,SUM(mo_received) as mo, summary_date"}
	}

	query.Select(selectedFields)

	// Set default values
	if params.DataType == "" {
		params.DataType = "daily_report"
	}

	// Apply paramss
	if params.Country != "" {
		query.Where("country = ?", params.Country)
	}
	if params.Operator != "" {
		query.Where("operator = ?", params.Operator)
	}
	if params.Adnet != "" {
		query.Where("adnet = ?", params.Adnet)
	}
	if params.PartnerName != "" {
		query.Where("partner = ?", params.PartnerName)
	}
	if params.Service != "" {
		query.Where("service = ?", params.Service)
	}
	if params.CampaignName != "" {
		query.Where("campaign_name = ?", params.CampaignName)
	}

	// Handle Date Range
	var dateStart, dateEnd time.Time
	var errStart, errEnd error
	today := time.Now()

	switch strings.ToUpper(params.DateRange) {
	case "TODAY":
		dateStart, dateEnd = today, today
	case "YESTERDAY":
		dateStart, dateEnd = today.AddDate(0, 0, -1), today.AddDate(0, 0, -1)
	case "LAST_7_DAY", "LAST_7_DAYS":
		dateStart, dateEnd = today.AddDate(0, 0, -6), today
	case "LAST_30_DAY", "LAST_30_DAYS":
		dateStart, dateEnd = today.AddDate(0, -1, 0), today
	case "THIS_MONTH":
		dateStart = time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
		dateEnd = today
	case "LAST_MONTH":
		lastMonth := today.AddDate(0, -1, 0)
		dateStart = time.Date(lastMonth.Year(), lastMonth.Month(), 1, 0, 0, 0, 0, today.Location())
		dateEnd = time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location()).AddDate(0, 0, -1)
	case "CUSTOM_RANGE":
		splitTime := strings.Split(params.DateCustomRange, "to")
		dateStart, errStart = time.Parse("2006-01-02", strings.TrimSpace(splitTime[0]))
		dateEnd, errEnd = time.Parse("2006-01-02", strings.TrimSpace(splitTime[1]))
		if errStart != nil {
			dateStart = today
		}
		if errEnd != nil {
			dateEnd = today
		}
	default:
		dateStart = time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
		dateEnd = today
	}

	if dateEnd.After(today) {
		dateEnd = today
	}

	query.Where("summary_date BETWEEN ? AND ?", dateStart, dateEnd)
	query.Group("summary_date")

	rows, _ = query.Unscoped().Rows()

	defer rows.Close()

	var (
		ss []entity.CampaignSummaryChart
	)

	for rows.Next() {
		var s entity.CampaignSummaryChart
		// ScanRows scans a row into a struct
		r.DB.ScanRows(rows, &s)
		sDate, err := time.Parse(time.RFC3339, s.SummaryDate)
		if err == nil {
			s.SummaryDate = sDate.Format("2006-01-02")
		}
		ss = append(ss, s)
	}
	return ss, dateStart, dateEnd, rows.Err()

}

// helper
func formatQueryIndicators(selects []string, dataType string) []string {
	var formattedSelects []string

	for _, value := range selects {
		var formattedValue string

		if dataType == "monthly_report" {
			switch value {
			case "target_daily_budget":
				formattedValue = "MAX(target_monthly_budget) as target_daily_budget"
			case "waki_revenue":
				formattedValue = "SUM(saaf - sbaf) AS waki_revenue"
			case "target_budget":
				formattedValue = "SUM(target_monthly_budget) AS target_budget"
			case "budget_usage":
				formattedValue = "SUM(CASE WHEN target_daily_budget = 0 THEN 0 ELSE (saaf / target_daily_budget * 100) END) AS budget_usage"
			case "spending_to_adnets":
				formattedValue = "SUM(sbaf) AS spending_to_adnets"
			case "total_spending":
				formattedValue = "SUM(saaf) AS total_spending"
			case "spending":
				formattedValue = "SUM(saaf) AS spending"
			case "fp":
				formattedValue = "SUM(first_push) AS fp"
			case "mo_sent":
				formattedValue = "SUM(postback) AS mo_sent"
			case "traffic":
				formattedValue = "SUM(landing) AS traffic"
			case "budget":
				formattedValue = "SUM(target_daily_budget) AS budget"
			case "revenue":
				formattedValue = "SUM(revenue) AS revenue"
			default:
				formattedValue = fmt.Sprintf("SUM(%s) AS %s", value, value)
			}
		} else { // Daily Report
			switch value {
			case "waki_revenue":
				formattedValue = "saaf - sbaf AS waki_revenue"
			case "budget_usage":
				formattedValue = "CASE WHEN target_daily_budget = 0 THEN NULL ELSE (saaf / target_daily_budget * 100) END AS budget_usage, sbaf AS sbaf_t, target_daily_budget AS target_daily_budget_t"
			case "fp":
				formattedValue = "first_push AS fp"
			case "mo_sent":
				formattedValue = "postback AS mo_sent"
			case "spending_to_adnets":
				formattedValue = "sbaf AS spending_to_adnets"
			case "total_spending":
				formattedValue = "saaf AS total_spending"
			case "spending":
				formattedValue = "saaf AS spending"
			case "budget":
				formattedValue = "target_daily_budget AS budget"
			case "target_budget":
				formattedValue = "target_daily_budget AS budget"
			case "traffic":
				formattedValue = "landing AS traffic"
			case "mo":
				formattedValue = "mo_received AS mo"
			case "revenue":
				formattedValue = "revenue AS revenue"
			case "target_daily_budget":
				formattedValue = "target_daily_budget AS target_daily_budget"
			default:
				formattedValue = fmt.Sprintf("%s AS %s", value, value)
			}
		}
		formattedSelects = append(formattedSelects, formattedValue)
	}

	return formattedSelects
}

func (r *BaseModel) GetSummaryCampaignBudgetMonitoring(params entity.ParamsCampaignSummary) ([]entity.CampaignSummaryMonitoring, time.Time, time.Time, error) {
	var (
		rows *sql.Rows
	)
	query := r.DB.Model(&entity.CampaignSummaryMonitoring{})

	// Apply Indicator Selection -
	selectedFields := []string{"summary_date", "country", "url_service_key", "campaign_id", "campaign_name", "partner", "operator", "service", "adnet"}
	formattedIndicators := formatQueryIndicatorsBudget(params.DataIndicators, params.DataType)
	selectedFields = append(selectedFields, formattedIndicators...)

	query.Select(selectedFields)

	// Set default values
	if params.DataType == "" {
		params.DataType = "daily_report"
	}
	if params.UrlServiceKey != "" {
		query.Where("url_service_key = ?", params.UrlServiceKey)
	}
	// Apply paramss
	if params.Country != "" {
		query.Where("country = ?", params.Country)
	}
	if params.Operator != "" {
		query.Where("operator = ?", params.Operator)
	}
	if params.Adnet != "" {
		query.Where("adnet = ?", params.Adnet)
	}
	if params.PartnerName != "" {
		query.Where("partner = ?", params.PartnerName)
	}
	if params.Service != "" {
		query.Where("service = ?", params.Service)
	}
	if params.CampaignName != "" {
		query.Where("campaign_name = ?", params.CampaignName)
	}

	// Handle Date Range
	var startDate, endDate time.Time
	var errStart, errEnd error
	today := time.Now()

	switch strings.ToUpper(params.DateRange) {
	case "TODAY":
		startDate, endDate = today, today
	case "YESTERDAY":
		startDate, endDate = today.AddDate(0, 0, -1), today.AddDate(0, 0, -1)
	case "LAST_7_DAY":
		startDate, endDate = today.AddDate(0, 0, -6), today
	case "LAST_30_DAY":
		startDate, endDate = today.AddDate(0, -1, 0), today
	case "THIS_MONTH":
		startDate = time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
		endDate = today
	case "LAST_MONTH":
		lastMonth := today.AddDate(0, -1, 0)
		startDate = time.Date(lastMonth.Year(), lastMonth.Month(), 1, 0, 0, 0, 0, today.Location())
		endDate = time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location()).AddDate(0, 0, -1)
	case "CUSTOM_RANGE":
		startDate, errStart = time.Parse("2006-01-02", params.DateStart)
		endDate, errEnd = time.Parse("2006-01-02", params.DateEnd)
		if errStart != nil {
			startDate = today
		}
		if errEnd != nil {
			endDate = today
		}
	default:
		startDate = time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
		endDate = today
	}

	if endDate.After(today) {
		endDate = today
	}

	query.Where("summary_date BETWEEN ? AND ?", startDate, endDate)

	// Grouping for monthly reports
	if params.DataType == "monthly_report" {
		query.Group("EXTRACT(YEAR FROM summary_date), EXTRACT(MONTH FROM summary_date), country, partner, operator, service, adnet, url_service_key")
	}

	rows, _ = query.Unscoped().Rows()

	defer rows.Close()

	var (
		ss []entity.CampaignSummaryMonitoring
	)

	for rows.Next() {
		var s entity.CampaignSummaryMonitoring

		// ScanRows scans a row into a struct
		r.DB.ScanRows(rows, &s)

		// Cek apakah data sudah ada dalam slice ss berdasarkan campaign_id, campaign_name, country, operator, service, adnet
		exists := false
		for _, existingS := range ss {
			if existingS.CampaignId == s.CampaignId && existingS.CampaignName == s.CampaignName && existingS.Country == s.Country && existingS.Operator == s.Operator && existingS.Service == s.Service && existingS.Adnet == s.Adnet {
				exists = true
				break
			}
		}
		if !exists {
			ss = append(ss, s)
		}
	}
	return ss, startDate, endDate, rows.Err()

}

func formatQueryIndicatorsBudget(selects []string, dataType string) []string {
	var formattedSelects []string

	for _, value := range selects {
		var formattedValue string
		fmt.Println("indicators:", value)
		if dataType == "monthly_report" {
			switch value {
			case "budget":
				formattedValue = "SUM(target_daily_budget) AS budget"
			case "target_budget":
				formattedValue = "SUM(target_daily_budget) AS target_budget"
			case "spending":
				formattedValue = "SUM(saaf) AS spending"
			case "mo":
				formattedValue = "SUM(mo_received) AS mo"
			case "waki_revenue":
				formattedValue = "SUM(saaf - sbaf) AS waki_revenue"
			case "budget_usage":
				formattedValue = "SUM(CASE WHEN target_daily_budget = 0 THEN 0 ELSE (saaf / target_daily_budget * 100) END) AS budget_usage"
			case "spending_to_adnets":
				formattedValue = "SUM(sbaf) AS spending_to_adnets"
			case "total_spending":
				formattedValue = "SUM(saaf) AS total_spending"
			case "fp":
				formattedValue = "SUM(first_push) AS fp"
			case "mo_sent":
				formattedValue = "SUM(postback) AS mo_sent"
			case "traffic":
				formattedValue = "SUM(landing) AS traffic"
			case "revenue":
				formattedValue = "SUM(revenue) AS revenue"
			default:
				formattedValue = fmt.Sprintf("SUM(%s) AS %s", value, value)
			}
		} else { // Daily Report
			switch value {
			case "budget":
				formattedValue = "target_daily_budget AS budget"
			case "target_budget":
				formattedValue = "target_daily_budget AS target_budget"
			case "spending":
				formattedValue = "saaf AS spending"
			case "mo":
				formattedValue = "mo_received AS mo"
			case "waki_revenue":
				formattedValue = "saaf - sbaf AS waki_revenue"
			case "budget_usage":
				formattedValue = "CASE WHEN target_daily_budget = 0 THEN NULL ELSE (saaf / target_daily_budget * 100) END AS budget_usage, sbaf AS sbaf_t, target_daily_budget AS target_daily_budget_t"
			case "fp":
				formattedValue = "first_push AS fp"
			case "mo_sent":
				formattedValue = "postback AS mo_sent"
			case "spending_to_adnets":
				formattedValue = "sbaf AS spending_to_adnets"
			case "total_spending":
				formattedValue = "saaf AS total_spending"
			case "traffic":
				formattedValue = "landing AS traffic"
			case "revenue":
				formattedValue = "revenue AS revenue"
			default:
				formattedValue = fmt.Sprintf("%s AS %s", value, value)
			}
		}
		formattedSelects = append(formattedSelects, formattedValue)
	}

	return formattedSelects
}

func (r *BaseModel) FormulaCPA(sum entity.SummaryCampaign) entity.SummaryCampaign {

	gs, _ := r.GetDataConfig("global_setting", "$")

	// CR conversion to mo received
	//var cr_mo float64
	if sum.MoReceived > 0 && sum.Landing > 0 {
		sum.CrMO, _ = strconv.ParseFloat(fmt.Sprintf("%f", float64(sum.MoReceived)/float64(sum.Landing)), 64)
	} else {
		sum.CrMO, _ = strconv.ParseFloat("0", 64)
	}

	//var cr_postback float64
	if sum.Postback > 0 && sum.Landing > 0 {
		sum.CrPostback, _ = strconv.ParseFloat(fmt.Sprintf("%f", float64(sum.Postback)/float64(sum.Landing)), 64)
	} else {
		sum.CrPostback, _ = strconv.ParseFloat("0", 64)
	}

	mo_sent := float64(sum.Postback)
	//sbaf := payout * mo_sent
	sum.SBAF = sum.PO * mo_sent //sbaf

	// GET total waki agency fee
	sum.CostPerConversion, _ = strconv.ParseFloat(strings.TrimSpace(gs.CPCR), 64)
	//cost_per_conversion := sum.CostPerConversion
	sum.AgencyFee, _ = strconv.ParseFloat(strings.TrimSpace(gs.AgencyFee), 64)
	sum.AgencyFee = sum.AgencyFee / 100
	mo_received := float64(sum.MoReceived)
	sum.TotalWakiAgencyFee = (sum.CostPerConversion * mo_received) + (sum.AgencyFee * (sum.CostPerConversion + (sum.CostPerConversion * mo_received)))

	// GET SAAF (spending after agency fee)
	//saaf := total_waki_agency_fee + sbaf
	//tech_fee, _ := strconv.Atoi(gs.TechnicalFee)
	sum.TechnicalFee, _ = strconv.ParseFloat(strings.TrimSpace(gs.TechnicalFee), 64)
	sum.TechnicalFee = sum.TechnicalFee / 100

	sum.TechnicalFee = sum.TechnicalFee * (sum.SBAF + sum.TotalWakiAgencyFee)

	sum.SAAF = sum.TotalWakiAgencyFee + sum.SBAF + sum.TechnicalFee //saaf
	if strings.ToLower(sum.ClientType) == "external" {
		sum.SAAF = mo_received * sum.POAF
	}

	// GET eCPA
	//cpa := float64(0)
	if sum.SAAF > 0 && mo_received > 0 {
		sum.CPA = sum.SAAF / mo_received
	}

	// Revenue
	//revenue := float64(0)
	if sum.SAAF > 0 && sum.SBAF > 0 {
		sum.Revenue = sum.SAAF - sum.SBAF
	}

	sum.PricePerMO = sum.SAAF / mo_received

	if sum.Landing == 0 || sum.MoReceived == 0 || sum.Postback == 0 {
		//price_per_mo_string = "0"
		sum.PricePerMO = 0
	}

	return sum
}

func (r *BaseModel) ReCalculateSummaryCampaign(o entity.SummaryCampaign) error {

	result := r.DB.Exec("UPDATE summary_campaigns SET traffic = ?, landing = ?, mo_received = ?, cr_mo = ?, cr_postback = ?, postback = ?, total_fp = ?, success_fp = ?, billrate = ?, po = ?, poaf = ?, sbaf = ?, saaf = ?, cpa = ?, revenue = ?, url_after = ?, url_before = ?, mo_limit = ?, ratio_send = ?, ratio_receive = ?, client_type = ?, cost_per_conversion = ?, agency_fee = ?, total_waki_agency_fee = ?, campaign_name = ?, technical_fee = ?, company = ? WHERE summary_date = '"+o.SummaryDate.Format("2006-01-02")+"' AND url_service_key = ? AND country = ? AND operator = ? AND partner = ? AND service = ? AND adnet = ? AND campaign_id = ? AND campaign_objective = ?", o.Traffic, o.Landing, o.MoReceived, o.CrMO, o.CrPostback, o.Postback, o.TotalFP, o.SuccessFP, o.Billrate, o.PO, o.POAF, o.SBAF, o.SAAF, o.CPA, o.Revenue, o.URLAfter, o.URLBefore, o.MOLimit, o.RatioSend, o.RatioReceive, o.ClientType, o.CostPerConversion, o.AgencyFee, o.TotalWakiAgencyFee, o.CampaignName, o.TechnicalFee, o.Company, o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Service, o.Adnet, o.CampaignId, o.CampaignObjective)

	r.Logs.Debug(fmt.Sprintf("ReCalculateSummaryCampaign : %s-%s, affected: %d, is error : %#v", o.URLServiceKey, o.SummaryDate, result.RowsAffected, result.Error))

	return result.Error
}

/* package model

import (
	"context"
	"fmt"
	"strconv"

	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"

	_ "github.com/lib/pq"
)

const (
	DELSUMMARYCAMPAIGN         = "DELETE FROM summary_campaign WHERE DATE(summary_date) = '%s' AND urlservicekey = '%s' AND country = '%s' AND operator = '%s' AND partner = '%s' AND service = '%s' AND adnet = '%s' AND campaign_id = '%s'"
	EDITSETTINGSUMMARYCAMPAIGN = "UPDATE summary_campaign SET po = '%s', mo_limit = %d, ratio_send = %d, ratio_receive = %d WHERE DATE(summary_date) = '%s' AND urlservicekey = '%s' AND country = '%s' AND operator = '%s' AND partner = '%s' AND service = '%s' AND adnet = '%s' AND campaign_id = '%s'"
	UPDATESUMMARYCAMPAIGN      = "UPDATE summary_campaign SET status = %t WHERE DATE(summary_date) = '%s' AND urlservicekey = '%s' AND country = '%s' AND operator = '%s' AND partner = '%s' AND service = '%s' AND adnet = '%s' AND campaign_id = '%s'"
	SUMMARYCAMPAIGN            = "INSERT INTO summary_campaign AS sc (id, status, summary_date, campaign_id, campaign_name, country, partner, operator, urlservicekey, aggregator, service, adnet, short_code, traffic, landing, mo_received, cr_mo, cr_postback, postback, total_fp, success_fp, billrate, po, sbaf, saaf, cpa, revenue, url_after, url_before, mo_limit, ratio_send, ratio_receive, client_type, cost_per_conversion, agency_fee, total_waki_agency_fee, target_daily_budget, budget_usage) VALUES (DEFAULT, %t, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %d, %d, %d, '%s'::double precision, '%s'::double precision, %d, %d, '%s', '%s'::double precision, '%s'::double precision, '%s'::double precision, '%s'::double precision, '%s'::double precision, '%s'::double precision, '%s', '%s', %d, %d, %d, '%s', '%s'::double precision, '%s'::double precision, '%s'::double precision, '%s'::double precision, '%s'::double precision) ON CONFLICT (summary_date, campaign_id, country, partner, operator, urlservicekey, service, adnet) DO UPDATE SET traffic = %d, landing = %d, mo_received = %d, cr_mo = '%s'::double precision, cr_postback = '%s'::double precision, postback = %d, total_fp = %d, success_fp = '%s', billrate = '%s'::double precision, po = '%s'::double precision, sbaf = '%s'::double precision, saaf = '%s'::double precision, cpa = '%s'::double precision, revenue = '%s'::double precision, url_after = '%s', url_before = '%s', mo_limit = %d, ratio_send = %d, ratio_receive = %d, client_type = '%s', cost_per_conversion = '%s'::double precision, agency_fee = '%s'::double precision, total_waki_agency_fee = '%s'::double precision, target_daily_budget = '%s'::double precision, budget_usage = '%s'::double precision, campaign_name = '%s';"
	//SUMMARYCAMPAIGN                             = "INSERT INTO summary_campaign AS sc (id, status, summary_date, campaign_id, campaign_name, country, partner, operator, urlservicekey, aggregator, service, adnet, short_code, traffic, landing, mo_received, cr_mo, cr_postback, postback, total_fp, success_fp, billrate, po, sbaf, saaf, cpa, revenue, url_after, url_before, mo_limit, ratio_send, ratio_receive, client_type, cost_per_conversion, agency_fee, total_waki_agency_fee, target_daily_budget, budget_usage) VALUES (DEFAULT, %t, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %d, %d, %d, '%s'::double precision, '%s'::double precision, %d, %d, '%s', '%s'::double precision, '%s'::double precision, '%s'::double precision, '%s'::double precision, '%s'::double precision, '%s'::double precision, '%s', '%s', %d, %d, %d, '%s', '%s'::double precision, '%s'::double precision, '%s'::double precision, '%s'::double precision, '%s'::double precision) ON CONFLICT (summary_date, campaign_id, country, partner, operator, urlservicekey, service, adnet) DO UPDATE SET traffic = sc.traffic + %d, landing = sc.landing + %d, mo_received = sc.mo_received + %d, cr_mo = '%s'::double precision, cr_postback = '%s'::double precision, postback = sc.postback + %d, total_fp = sc.total_fp + %d, success_fp = '%s', billrate = '%s'::double precision, po = '%s'::double precision, sbaf = '%s'::double precision, saaf = '%s'::double precision, cpa = '%s'::double precision, revenue = '%s'::double precision, url_after = '%s', url_before = '%s', mo_limit = %d, ratio_send = %d, ratio_receive = %d, client_type = '%s', cost_per_conversion = '%s'::double precision, agency_fee = '%s'::double precision, total_waki_agency_fee = '%s'::double precision, target_daily_budget = '%s'::double precision, budget_usage = '%s'::double precision;"
	INSERTTRAFFIC                               = "INSERT INTO data_traffic (id, traffic_time, traffic_added_time, http_status, urlservicekey, campaign_id, country, partner, operator, aggregator, service, short_code, adnet, keyword, subkeyword, is_billable, plan) VALUES (DEFAULT, '%s', %d, %d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, '%s');"
	GETTOTALDATATRAFFIC                         = "SELECT COUNT(1) FROM data_traffic WHERE DATE(traffic_time) = '%s' AND urlservicekey = '%s' AND campaign_id = '%s' AND country = '%s' AND partner = '%s' AND operator = '%s' AND service = '%s' AND short_code = '%s' AND adnet = '%s'"
	INSERTLANDING                               = "INSERT INTO data_landing (id, landing_time, landed_time, http_status, urlservicekey, campaign_id, country, partner, operator, aggregator, service, short_code, adnet, keyword, subkeyword, is_billable, plan) VALUES (DEFAULT, '%s', %d, %d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, '%s');"
	GETTOTALDATALANDING                         = "SELECT COUNT(1) FROM data_landing WHERE DATE(landing_time) = '%s' AND urlservicekey = '%s' AND campaign_id = '%s' AND country = '%s' AND partner = '%s' AND operator = '%s' AND service = '%s' AND short_code = '%s' AND adnet = '%s'"
	INSERTCLICKED                               = "INSERT INTO data_clicked (id, clicked_time, clicked_button_time, http_status, urlservicekey, campaign_id, country, partner, operator, aggregator, service, short_code, adnet, keyword, subkeyword, is_billable, plan) VALUES (DEFAULT, '%s', %d, %d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, '%s');"
	GETTOTALDATACLICKED                         = "SELECT COUNT(1) FROM data_clicked WHERE DATE(clicked_time) = '%s' AND urlservicekey = '%s' AND campaign_id = '%s' AND country = '%s' AND partner = '%s' AND operator = '%s' AND service = '%s' AND short_code = '%s' AND adnet = '%s'"
	INSERTREDIRECT                              = "INSERT INTO data_redirect (id, redirect_time, redirect_added_time, http_status, urlservicekey, campaign_id, country, partner, operator, aggregator, service, short_code, adnet, keyword, subkeyword, is_billable, plan) VALUES (DEFAULT, '%s', %d, %d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, '%s');"
	GETTOTALDATAREDIRECT                        = "SELECT COUNT(1) FROM data_redirect WHERE DATE(redirect_time) = '%s' AND urlservicekey = '%s' AND campaign_id = '%s' AND country = '%s' AND partner = '%s' AND operator = '%s' AND service = '%s' AND short_code = '%s' AND adnet = '%s'"
	UPDATECPAREPORTSUMMARYCAMPAIGN              = "UPDATE summary_campaign SET cost_per_conversion = '%s', agency_fee = '%s' WHERE DATE(summary_date) = '%s' AND urlservicekey = '%s' AND country = '%s' AND operator = '%s' AND partner = '%s' AND service = '%s' AND adnet = '%s' AND campaign_id = '%s'"
	UPDATEREPORTSUMMARYCAMPAIGNMONITORINGBUDGET = "UPDATE summary_campaign SET target_daily_budget = '%s' WHERE DATE(summary_date) = '%s AND country = '%s' AND operator = '%s'"
)

func (r *BaseModel) DelSummaryCampaign(summary_date string, o entity.DataConfig) error {

	SQL := fmt.Sprintf(DELSUMMARYCAMPAIGN, summary_date, o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Service, o.Adnet, o.CampaignId)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("(%s) Error %s when preparing SQL statement", SQL, err))

		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when update to table", SQL, err))

		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when finding rows affected", SQL, err))

		return err
	}

	r.Logs.Debug(fmt.Sprintf("SQL : %s, row affected : %d", SQL, rows))
	return nil
}

func (r *BaseModel) EditSettingSummaryCampaign(summary_date string, o entity.DataConfig) error {

	SQL := fmt.Sprintf(EDITSETTINGSUMMARYCAMPAIGN, o.PO, o.MOCapping, o.RatioSend, o.RatioReceive, summary_date, o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Service, o.Adnet, o.CampaignId)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("(%s) Error %s when preparing SQL statement", SQL, err))

		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when update to table", SQL, err))

		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when finding rows affected", SQL, err))

		return err
	}

	r.Logs.Debug(fmt.Sprintf("SQL : %s, row affected : %d", SQL, rows))
	return nil
}

func (r *BaseModel) UpdateSummaryCampaign(summary_date string, o entity.DataConfig) error {

	SQL := fmt.Sprintf(UPDATESUMMARYCAMPAIGN, o.IsActive, summary_date, o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Service, o.Adnet, o.CampaignId)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("(%s) Error %s when preparing SQL statement", SQL, err))

		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when update to table", SQL, err))

		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when finding rows affected", SQL, err))

		return err
	}

	r.Logs.Debug(fmt.Sprintf("SQL : %s, row affected : %d", SQL, rows))
	return nil
}

func (r *BaseModel) SummaryCampaign(d entity.Summary) int {

	SQL := fmt.Sprintf(SUMMARYCAMPAIGN, d.IsActive, d.SummaryDate, d.CampaignId, d.CampaignName, d.Country, d.Partner, d.Operator, d.URLServiceKey, d.Aggregator, d.Service, d.Adnet, d.ShortCode, d.TotalTraffic, d.TotalLanding, d.TotalMOReceived, d.CRMO, d.CRPostback, d.TotalPostback, d.TotalFP, d.SuccessFP, d.BillRate, d.PO, d.SBAF, d.SAAF, d.CPA, d.Revenue, d.URLWarpLanding, d.URLLanding, d.MOCapping, d.RatioSend, d.RatioReceive, d.ClientType, d.CPCR, d.AgencyFee, d.TotalWakiAgencyFee, d.TDB, d.BudgetUsage, d.TotalTraffic, d.TotalLanding, d.TotalMOReceived, d.CRMO, d.CRPostback, d.TotalPostback, d.TotalFP, d.SuccessFP, d.BillRate, d.PO, d.SBAF, d.SAAF, d.CPA, d.Revenue, d.URLWarpLanding, d.URLLanding, d.MOCapping, d.RatioSend, d.RatioReceive, d.ClientType, d.CPCR, d.AgencyFee, d.TotalWakiAgencyFee, d.TDB, d.BudgetUsage, d.CampaignName)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("(%s) Error %s when preparing SQL statement", SQL, err))

		return 0
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when update to table", SQL, err))

		return 0
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when finding rows affected", SQL, err))

		return 0
	}

	r.Logs.Debug(fmt.Sprintf("SQL : %s, row affected : %d", SQL, rows))
	return int(rows)
}

func (r *BaseModel) DataTraffic(data map[string]string, o entity.DataCounter) int {

	traffic_added_time, _ := strconv.Atoi(data["traffic_added_time"])
	httpstatus, _ := strconv.Atoi(data["http_status"])

	SQL := fmt.Sprintf(INSERTTRAFFIC, data["traffic_time"], traffic_added_time, httpstatus, o.URLServiceKey, o.CampaignId, o.Country, o.Partner, o.Operator, o.Aggregator, o.Service, o.ShortCode, o.Adnet, o.Keyword, o.SubKeyword, o.IsBillable, o.Plan)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("(%s) Error %s when preparing SQL statement", SQL, err))

		return 0
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when update to table", SQL, err))

		return 0
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when finding rows affected", SQL, err))

		return 0
	}

	r.Logs.Debug(fmt.Sprintf("SQL : %s, row affected : %d", SQL, rows))
	return int(rows)
}

func (r *BaseModel) DataLanding(data map[string]string, o entity.DataCounter) int {

	landed_time, _ := strconv.Atoi(data["landed_time"])
	httpstatus, _ := strconv.Atoi(data["http_status"])

	SQL := fmt.Sprintf(INSERTLANDING, data["landing_time"], landed_time, httpstatus, o.URLServiceKey, o.CampaignId, o.Country, o.Partner, o.Operator, o.Aggregator, o.Service, o.ShortCode, o.Adnet, o.Keyword, o.SubKeyword, o.IsBillable, o.Plan)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("(%s) Error %s when preparing SQL statement", SQL, err))

		return 0
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when update to table", SQL, err))

		return 0
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when finding rows affected", SQL, err))

		return 0
	}

	r.Logs.Debug(fmt.Sprintf("SQL : %s, row affected : %d", SQL, rows))
	return int(rows)
}

func (r *BaseModel) DataClicked(data map[string]string, o entity.DataCounter) int {

	clicked_button_time, _ := strconv.Atoi(data["clicked_button_time"])
	httpstatus, _ := strconv.Atoi(data["http_status"])

	SQL := fmt.Sprintf(INSERTCLICKED, data["clicked_time"], clicked_button_time, httpstatus, o.URLServiceKey, o.CampaignId, o.Country, o.Partner, o.Operator, o.Aggregator, o.Service, o.ShortCode, o.Adnet, o.Keyword, o.SubKeyword, o.IsBillable, o.Plan)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("(%s) Error %s when preparing SQL statement", SQL, err))

		return 0
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when update to table", SQL, err))

		return 0
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when finding rows affected", SQL, err))

		return 0
	}

	r.Logs.Debug(fmt.Sprintf("SQL : %s, row affected : %d", SQL, rows))
	return int(rows)
}

func (r *BaseModel) DataRedirect(data map[string]string, o entity.DataCounter) int {

	redirect_added_time, _ := strconv.Atoi(data["redirect_added_time"])
	httpstatus, _ := strconv.Atoi(data["http_status"])

	SQL := fmt.Sprintf(INSERTREDIRECT, data["redirect_time"], redirect_added_time, httpstatus, o.URLServiceKey, o.CampaignId, o.Country, o.Partner, o.Operator, o.Aggregator, o.Service, o.ShortCode, o.Adnet, o.Keyword, o.SubKeyword, o.IsBillable, o.Plan)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("(%s) Error %s when preparing SQL statement", SQL, err))

		return 0
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when update to table", SQL, err))

		return 0
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when finding rows affected", SQL, err))

		return 0
	}

	r.Logs.Debug(fmt.Sprintf("SQL : %s, row affected : %d", SQL, rows))
	return int(rows)
}

func (r *BaseModel) UpdateCPAReportSummaryCampaign(summary_date string, o entity.DataConfig) error {

	SQL := fmt.Sprintf(UPDATECPAREPORTSUMMARYCAMPAIGN, summary_date, o.CPCR, o.AgencyFee, o.URLServiceKey, o.Country, o.Operator, o.Partner, o.Service, o.Adnet, o.CampaignId)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("(%s) Error %s when preparing SQL statement", SQL, err))

		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when update to table", SQL, err))

		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when finding rows affected", SQL, err))

		return err
	}

	r.Logs.Debug(fmt.Sprintf("SQL : %s, row affected : %d", SQL, rows))
	return nil
}

func (r *BaseModel) UpdateReportSummaryCampaignMonitoringBudget(summary_date string, o entity.DataConfig) error {

	SQL := fmt.Sprintf(UPDATEREPORTSUMMARYCAMPAIGNMONITORINGBUDGET, o.TargetDailyBudget, summary_date, o.Country, o.Operator)

	stmt, err := r.DBPostgre.PrepareContext(context.Background(), SQL)

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("(%s) Error %s when preparing SQL statement", SQL, err))

		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background())

	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when update to table", SQL, err))

		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {

		r.Logs.Debug(fmt.Sprintf("SQL : %s, Error %s when finding rows affected", SQL, err))

		return err
	}

	r.Logs.Debug(fmt.Sprintf("SQL : %s, row affected : %d", SQL, rows))
	return nil
} */


func (r *BaseModel) GetCampaignBudgetSummary(params entity.ParamsCampaignSummary, dateStart, dateEnd time.Time) ([]entity.BudgetDetailItem, []entity.BudgetSummaryItem, []entity.BudgetDetailItem, error) {
	base := r.DB.Table("summary_campaigns").
		Where("deleted_at IS NULL").
		Where("summary_date BETWEEN ? AND ?", dateStart, dateEnd)
	if params.Country != "" {
		base = base.Where("country = ?", params.Country)
	}
	if params.Operator != "" {
		base = base.Where("operator = ?", params.Operator)
	}
	if params.PartnerName != "" {
		base = base.Where("partner = ?", params.PartnerName)
	}
	if params.Service != "" {
		base = base.Where("service = ?", params.Service)
	}
	if params.Adnet != "" {
		base = base.Where("adnet = ?", params.Adnet)
	}

	var detail []entity.BudgetDetailItem
	err := base.Session(&gorm.Session{}).
		Select(`country, operator, partner, service, adnet,
			EXTRACT(YEAR FROM summary_date)::int AS year,
			EXTRACT(MONTH FROM summary_date)::int AS month,
			MAX(target_monthly_budget) AS budget,
			SUM(saaf) AS spending,
			CASE WHEN MAX(target_monthly_budget) = 0 THEN 0
			     ELSE SUM(saaf) / MAX(target_monthly_budget) * 100
			END AS budget_usage`).
		Group("country, operator, partner, service, adnet, EXTRACT(YEAR FROM summary_date), EXTRACT(MONTH FROM summary_date)").
		Unscoped().
		Find(&detail).Error
	if err != nil {
		return nil, nil, nil, err
	}

	type cKey struct {
		country      string
		year, month int
	}

	spendingMap := make(map[cKey]float64)
	for _, d := range detail {
		k := cKey{d.Country, d.Year, d.Month}
		spendingMap[k] += d.Spending
	}

	periodStart := time.Date(dateStart.Year(), dateStart.Month(), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	periodEnd := time.Date(dateEnd.Year(), dateEnd.Month()+1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1).Format("2006-01-02")

	type rawBudgetRow struct {
		Country string
		Year    int
		Month   int
		Budget  float64
	}

	var rawBudgets []rawBudgetRow
	tbQuery := r.DB.Table("target_budgets").
		Select("country, year, month, budget").
		Where("(year::text||'-'||lpad(month::text,2,'0')||'-02')::date BETWEEN ? AND ?", periodStart, periodEnd)
	if params.Country != "" {
		tbQuery = tbQuery.Where("country = ?", params.Country)
	}
	tbQuery.Scan(&rawBudgets)

	var summaryList []entity.BudgetSummaryItem
	for _, rb := range rawBudgets {
		k := cKey{rb.Country, rb.Year, rb.Month}
		summaryList = append(summaryList, entity.BudgetSummaryItem{
			Country:  rb.Country,
			Year:     rb.Year,
			Month:    rb.Month,
			Budget:   rb.Budget,
			Spending: spendingMap[k],
		})
	}

	if len(summaryList) == 0 {
		agg := make(map[cKey]*entity.BudgetSummaryItem)
		for _, d := range detail {
			k := cKey{d.Country, d.Year, d.Month}
			if agg[k] == nil {
				agg[k] = &entity.BudgetSummaryItem{Country: d.Country, Year: d.Year, Month: d.Month}
			}
			agg[k].Budget += d.Budget
			agg[k].Spending += d.Spending
		}
		for _, v := range agg {
			summaryList = append(summaryList, *v)
		}
	}

	var budgetSelf []entity.BudgetDetailItem
	tbdQuery := r.DB.Table("target_budget_details").
		Select("country, operator, partner, service, adnet, year, month, budget, 0 as spending, 0 as budget_usage").
		Where("adnet = ''").
		Where("(year::text||'-'||lpad(month::text,2,'0')||'-02')::date BETWEEN ? AND ?", periodStart, periodEnd)
	if params.Country != "" {
		tbdQuery = tbdQuery.Where("country = ?", params.Country)
	}
	if params.Operator != "" {
		tbdQuery = tbdQuery.Where("operator = ?", params.Operator)
	}
	if params.PartnerName != "" {
		tbdQuery = tbdQuery.Where("partner = ?", params.PartnerName)
	}
	if params.Service != "" {
		tbdQuery = tbdQuery.Where("service = ?", params.Service)
	}
	tbdQuery.Find(&budgetSelf)

	return detail, summaryList, budgetSelf, nil
}

func (r *BaseModel) UpdateTargetBudgetByLevel(req entity.EditTargetBudgetRequest) error {
	query := r.DB.Table("summary_campaigns").
		Where("EXTRACT(YEAR FROM summary_date) = ? AND EXTRACT(MONTH FROM summary_date) = ?", req.Year, req.Month).
		Where("country = ?", req.Country)
	switch req.Level {
	case "adnet":
		query = query.Where("operator = ? AND partner = ? AND service = ? AND adnet = ?", req.Operator, req.Partner, req.Service, req.Adnet)
	case "service":
		query = query.Where("operator = ? AND partner = ? AND service = ?", req.Operator, req.Partner, req.Service)
	case "partner":
		query = query.Where("operator = ? AND partner = ?", req.Operator, req.Partner)
	case "operator":
		query = query.Where("operator = ?", req.Operator)
	}
	days := float64(daysInMonth(req.Year, req.Month))
	result := query.Updates(map[string]interface{}{
		"target_daily_budget": req.Budget / days,
	})
	r.Logs.Debug(fmt.Sprintf("UpdateTargetBudgetByLevel affected: %d, error: %v", result.RowsAffected, result.Error))
	return result.Error
}

func budgetLockKey(country string, year, month int) int64 {
	h := fnv.New64a()
	h.Write([]byte(fmt.Sprintf("%s_%d_%d", country, year, month)))
	return int64(h.Sum64())
}

func (r *BaseModel) UpdateTargetBudgetBatch(reqs []entity.EditTargetBudgetRequest) error {
	if len(reqs) == 0 {
		return nil
	}
	type tbdEntry struct {
		country, operator, partner, service, adnet string
		year, month                                 int
		monthly, daily                              float64
	}
	type tbdKey struct {
		country, operator, partner, service, adnet string
		year, month                                 int
	}

	var scEntries []scBatchEntry
	tbdMap := make(map[tbdKey]tbdEntry)
	for _, req := range reqs {
		if req.Level == "country" {
			continue
		}
		days := float64(daysInMonth(req.Year, req.Month))
		daily := req.Budget / days
		var operator, partner, service, adnet string
		switch req.Level {
		case "adnet":
			operator, partner, service, adnet = strings.ToUpper(req.Operator), req.Partner, req.Service, req.Adnet
		case "service":
			operator, partner, service, adnet = strings.ToUpper(req.Operator), req.Partner, req.Service, ""
		case "partner":
			operator, partner, service, adnet = strings.ToUpper(req.Operator), req.Partner, "", ""
		case "operator":
			operator, partner, service, adnet = strings.ToUpper(req.Operator), "", "", ""
		}
		tbdMap[tbdKey{req.Country, operator, partner, service, adnet, req.Year, req.Month}] =
			tbdEntry{req.Country, operator, partner, service, adnet, req.Year, req.Month, req.Budget, daily}
		scEntries = append(scEntries, scBatchEntry{req.Level, req.Operator, req.Partner, req.Service, req.Adnet, daily, req.Budget})
	}

	lockKey := budgetLockKey(reqs[0].Country, reqs[0].Year, reqs[0].Month)
	if err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT pg_advisory_xact_lock(?)", lockKey).Error; err != nil {
			return err
		}
		for _, req := range reqs {
			if req.Level != "country" {
				continue
			}
			upsertSQL := `INSERT INTO target_budgets (country, year, month, budget, created_at, updated_at)
				VALUES (?, ?, ?, ?, NOW(), NOW())
				ON CONFLICT (country, year, month) DO UPDATE SET budget = ?, updated_at = NOW()`
			if err := tx.Exec(upsertSQL, req.Country, req.Year, req.Month, req.Budget, req.Budget).Error; err != nil {
				return err
			}
		}
		if len(tbdMap) > 0 {
			placeholders := make([]string, 0, len(tbdMap))
			args := make([]interface{}, 0, len(tbdMap)*9)
			for _, t := range tbdMap {
				placeholders = append(placeholders, "(?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())")
				args = append(args, t.country, t.year, t.month, t.operator, t.partner, t.service, t.adnet, t.monthly, t.daily)
			}
			tbdSQL := `INSERT INTO target_budget_details (country, year, month, operator, partner, service, adnet, budget, budget_per_day, created_at, updated_at)
				VALUES ` + strings.Join(placeholders, ", ") + `
				ON CONFLICT (country, year, month, operator, partner, service, adnet)
				DO UPDATE SET budget = EXCLUDED.budget, budget_per_day = EXCLUDED.budget_per_day, updated_at = NOW()`
			if err := tx.Exec(tbdSQL, args...).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	if len(scEntries) > 0 {
		year := reqs[0].Year
		month := reqs[0].Month
		country := reqs[0].Country
		startDate := fmt.Sprintf("%d-%02d-01", year, month)
		endDate := time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
		entries := scEntries
		go func() {
			for _, level := range []string{"operator", "partner", "service", "adnet"} {
				var levelEntries []scBatchEntry
				for _, e := range entries {
					if e.level == level {
						levelEntries = append(levelEntries, e)
					}
				}
				if len(levelEntries) == 0 {
					continue
				}
				if err := batchUpdateSummaryCampaigns(r.DB, level, country, startDate, endDate, levelEntries); err != nil {
					fmt.Printf("[UpdateTargetBudget] async SC update error level=%s: %v\n", level, err)
				}
			}
		}()
	}
	return nil
}

type scBatchEntry struct {
	level, operator, partner, service, adnet string
	daily, monthly                            float64
}

func batchUpdateSummaryCampaigns(tx *gorm.DB, level, country, startDate, endDate string, entries []scBatchEntry) error {
	var valRows []string
	args := make([]interface{}, 0)
	switch level {
	case "operator":
		for _, e := range entries {
			valRows = append(valRows, "(?, ?::float8, ?::float8)")
			args = append(args, strings.ToUpper(e.operator), e.daily, e.monthly)
		}
		SQL := fmt.Sprintf(`UPDATE summary_campaigns sc
			SET target_daily_budget = v.daily, target_monthly_budget = v.monthly
			FROM (VALUES %s) AS v(operator, daily, monthly)
			WHERE sc.country = ? AND sc.summary_date >= ? AND sc.summary_date < ?
			  AND UPPER(sc.operator) = v.operator`,
			strings.Join(valRows, ", "))
		args = append(args, country, startDate, endDate)
		return tx.Exec(SQL, args...).Error
	case "partner":
		for _, e := range entries {
			valRows = append(valRows, "(?, ?, ?::float8, ?::float8)")
			args = append(args, strings.ToUpper(e.operator), e.partner, e.daily, e.monthly)
		}
		SQL := fmt.Sprintf(`UPDATE summary_campaigns sc
			SET target_daily_budget = v.daily, target_monthly_budget = v.monthly
			FROM (VALUES %s) AS v(operator, partner, daily, monthly)
			WHERE sc.country = ? AND sc.summary_date >= ? AND sc.summary_date < ?
			  AND UPPER(sc.operator) = v.operator AND sc.partner = v.partner`,
			strings.Join(valRows, ", "))
		args = append(args, country, startDate, endDate)
		return tx.Exec(SQL, args...).Error
	case "service":
		for _, e := range entries {
			valRows = append(valRows, "(?, ?, ?, ?::float8, ?::float8)")
			args = append(args, strings.ToUpper(e.operator), e.partner, e.service, e.daily, e.monthly)
		}
		SQL := fmt.Sprintf(`UPDATE summary_campaigns sc
			SET target_daily_budget = v.daily, target_monthly_budget = v.monthly
			FROM (VALUES %s) AS v(operator, partner, service, daily, monthly)
			WHERE sc.country = ? AND sc.summary_date >= ? AND sc.summary_date < ?
			  AND UPPER(sc.operator) = v.operator AND sc.partner = v.partner AND sc.service = v.service`,
			strings.Join(valRows, ", "))
		args = append(args, country, startDate, endDate)
		return tx.Exec(SQL, args...).Error
	case "adnet":
		for _, e := range entries {
			valRows = append(valRows, "(?, ?, ?, ?, ?::float8, ?::float8)")
			args = append(args, strings.ToUpper(e.operator), e.partner, e.service, e.adnet, e.daily, e.monthly)
		}
		SQL := fmt.Sprintf(`UPDATE summary_campaigns sc
			SET target_daily_budget = v.daily, target_monthly_budget = v.monthly
			FROM (VALUES %s) AS v(operator, partner, service, adnet, daily, monthly)
			WHERE sc.country = ? AND sc.summary_date >= ? AND sc.summary_date < ?
			  AND UPPER(sc.operator) = v.operator AND sc.partner = v.partner
			  AND sc.service = v.service AND sc.adnet = v.adnet`,
			strings.Join(valRows, ", "))
		args = append(args, country, startDate, endDate)
		return tx.Exec(SQL, args...).Error
	}
	return nil
}

func daysInMonth(year, month int) int {
	return time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC).Day()
}
