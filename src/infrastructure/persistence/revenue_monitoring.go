package persistence

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
)

func (r *BaseModel) GetRevenueMonitoring(params entity.ParamsRevenueMonitoring) ([]entity.CampaignSummaryMonitoring, time.Time, time.Time, error) {
	var (
		rows *sql.Rows
	)
	query := r.DB.Model(&entity.CampaignSummaryMonitoring{})

	// Apply Indicator Selection
	selectedFields := []string{"country", "campaign_id", "url_service_key", "campaign_name", "partner", "operator", "service", "adnet"}
	if params.DataType == "monthly_report" {
		selectedFields = append(selectedFields, "DATE_TRUNC('month', summary_date) as summary_date")
	} else {
		selectedFields = append(selectedFields, "summary_date")
	}

	formattedIndicators := formatQueryIndicatorsRevenue(params.DataIndicators, params.DataType)
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
	if params.CampaignId != "" {
		query.Where("campaign_id = ?", params.CampaignId)
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
	case "LAST_7_DAY":
		dateStart, dateEnd = today.AddDate(0, 0, -6), today
	case "LAST_30_DAY":
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
		query.Group("DATE_TRUNC('month', summary_date),country, campaign_name,campaign_id, partner, operator, service, adnet")
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

func (r *BaseModel) GetRevenueChart(params entity.ParamsRevenueMonitoring) (entity.DataRevenue, time.Time, time.Time, error) {
	var (
		rows           *sql.Rows
		selectedFields []string
	)
	query := r.DB.Model(&entity.CampaignSummaryMonitoring{})

	if params.ChartType == "spending" {
		selectedFields = []string{"SUM(saaf) as spending,SUM(mo_received) as mo, SUM(revenue) as revenue,summary_date"}
	} else if params.ChartType == "cr" {
		selectedFields = []string{"SUM(saaf) as spending,SUM(mo_received) as mo, SUM(revenue) as revenue,SUM(cr) as cr, summary_date"}
	} else {
		selectedFields = []string{"SUM(cr) as cr,SUM(saaf) as spending,SUM(mo_received) as mo, SUM(revenue) as revenue, summary_date"}
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
	if params.CampaignId != "" {
		query.Where("campaign_id = ?", params.CampaignId)
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
	case "LAST_7_DAY":
		dateStart, dateEnd = today.AddDate(0, 0, -6), today
	case "LAST_30_DAY":
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
		ss            []entity.RevenueMonitoringChart
		total_revenue float64
	)

	for rows.Next() {
		var s entity.RevenueMonitoringChart
		// ScanRows scans a row into a struct
		r.DB.ScanRows(rows, &s)
		sDate, err := time.Parse(time.RFC3339, s.SummaryDate)
		if err == nil {
			s.SummaryDate = sDate.Format("2006-01-02")
		}
		ss = append(ss, s)
		total_revenue += s.Revenue
	}

	data := entity.DataRevenue{
		RevenueMonitoringChart: ss,
		TimeRevenue:            total_revenue,
		TimeInternalRevenue:    total_revenue,
		TimeExternalRevenue:    total_revenue,
	}

	return data, dateStart, dateEnd, rows.Err()

}

// helper
func formatQueryIndicatorsRevenue(selects []string, dataType string) []string {
	var formattedSelects []string

	for _, value := range selects {
		var formattedValue string

		if dataType == "monthly_report" {
			switch value {
			case "waki_revenue":
				formattedValue = "SUM(saaf - sbaf) AS waki_revenue"
			case "budget_usage":
				formattedValue = "SUM(CASE WHEN target_daily_budget = 0 THEN 0 ELSE (sbaf / target_daily_budget * 100) END) AS budget_usage"
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
				formattedValue = "CASE WHEN target_daily_budget = 0 THEN NULL ELSE (sbaf / target_daily_budget * 100) END AS budget_usage, sbaf AS sbaf_t, target_daily_budget AS target_daily_budget_t"
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
			case "traffic":
				formattedValue = "landing AS traffic"
			case "mo":
				formattedValue = "mo_received AS mo"
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
