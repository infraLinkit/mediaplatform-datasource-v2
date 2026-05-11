package persistence

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
	"gorm.io/gorm/clause"
)

func (r *BaseModel) GetRedirectionTime(params entity.RedirectionTimeParams) ([]entity.SummaryLanding, time.Time, time.Time, entity.RedirectionKPIStats, error) {
	var rows *sql.Rows

	query := r.DB.Model(&entity.SummaryLanding{})

	// Set default values
	if params.DataType == "" {
		params.DataType = "daily_report"
	}

	// Select base fields
	selectedFields := []string{"country", "campaign_id", "url_service_key", "campaign_name", "partner", "operator", "service", "adnet", "url_campaign"}

	switch params.DataType {
	case "monthly_report":
		selectedFields = append(selectedFields, "DATE_TRUNC('month', summary_date_hour) as summary_date_hour")
	case "daily_report":
		selectedFields = append(selectedFields, "DATE(summary_date_hour) as summary_date_hour")
	case "hourly_report":
		selectedFields = append(selectedFields, "summary_date_hour")
	}

	formattedIndicators := formatQueryIndicatorsRedirection(params.DataIndicators, params.DataType)
	selectedFields = append(selectedFields, formattedIndicators...)

	query.Select(selectedFields)

	// Apply filters
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
		query.Where("url_service_key = ?", params.CampaignId)
	}

	// Date range logic
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

	// Normalize to full day
	dateStart = time.Date(dateStart.Year(), dateStart.Month(), dateStart.Day(), 0, 0, 0, 0, dateStart.Location())
	dateEnd = time.Date(dateEnd.Year(), dateEnd.Month(), dateEnd.Day(), 23, 59, 59, 999999999, dateEnd.Location())

	query.Where("summary_date_hour BETWEEN ? AND ?", dateStart, dateEnd)

	// Grouping fields
	groupFields := []string{"country", "campaign_id", "url_service_key", "campaign_name", "partner", "operator", "service", "adnet", "url_campaign"}

	switch params.DataType {
	case "monthly_report":
		groupFields = append(groupFields, "DATE_TRUNC('month', summary_date_hour)")
	case "daily_report":
		groupFields = append(groupFields, "DATE(summary_date_hour)")
	case "hourly_report":
		groupFields = append(groupFields, "summary_date_hour")
	}

	query.Group(strings.Join(groupFields, ", "))

	rows, err := query.Unscoped().Rows()
	if err != nil {
		return nil, dateStart, dateEnd, entity.RedirectionKPIStats{}, err
	}
	defer rows.Close()

	var results []entity.SummaryLanding
	var kpiStats entity.RedirectionKPIStats

	for rows.Next() {
		var row entity.SummaryLanding
		r.DB.ScanRows(rows, &row)
		results = append(results, row)

		// KPI Check
		if row.TotalLoadTime > 3.0 {
			kpiStats.ExceedLoadTimeCount++
		}
		if row.ResponseTime > 0.3 {
			kpiStats.ExceedResponseTimeCount++
		}
		if row.SuccessRate < 95 {
			kpiStats.BelowSuccessRateCount++
		}
	}

	return results, dateStart, dateEnd, kpiStats, rows.Err()
}



func (r *BaseModel) GetRedirectionTimeHourly(params entity.RedirectionTimeParams) ([]entity.SummaryLanding, time.Time, time.Time, error) {
	var rows *sql.Rows

	query := r.DB.Model(&entity.SummaryLanding{})

	// Select field utama
	selectedFields := []string{
		"country", "campaign_id", "url_service_key", "campaign_name",
		"partner", "operator", "service", "adnet", "url_campaign",
		"DATE_TRUNC('hour', summary_date_hour) as summary_date_hour",
	}

	formattedIndicators := formatQueryIndicatorsRedirection(params.DataIndicators, "hourly_report")
	selectedFields = append(selectedFields, formattedIndicators...)

	query.Select(selectedFields)

	// Filter umum
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
		query.Where("url_service_key = ?", params.CampaignId)
	}

	// Tentukan tanggal
	var dateStart, dateEnd time.Time
	loc := time.FixedZone("Asia/Jakarta", 7*3600) // UTC+7 sesuai data DB

	if strings.ToUpper(params.DateRange) == "CUSTOM_RANGE" {
		splitTime := strings.Split(params.DateCustomRange, "to")
		dateStart, _ = time.ParseInLocation("2006-01-02", strings.TrimSpace(splitTime[0]), loc)
		dateEnd, _ = time.ParseInLocation("2006-01-02", strings.TrimSpace(splitTime[1]), loc)
	} else {
		today := time.Now().In(loc)
		dateStart = today
		dateEnd = today
	}

	// Normalize range untuk 1 hari penuh
	dateStart = time.Date(dateStart.Year(), dateStart.Month(), dateStart.Day(), 0, 0, 0, 0, loc)
	dateEnd = time.Date(dateEnd.Year(), dateEnd.Month(), dateEnd.Day(), 23, 59, 59, 999999999, loc)

	// Filter waktu
	query.Where("summary_date_hour BETWEEN ? AND ?", dateStart, dateEnd)

	// Group by
	groupFields := []string{
		"country", "campaign_id", "url_service_key", "campaign_name",
		"partner", "operator", "service", "adnet", "url_campaign",
		"DATE_TRUNC('hour', summary_date_hour)",
	}
	query.Group(strings.Join(groupFields, ", "))

	// Eksekusi
	rows, err := query.Unscoped().Rows()
	if err != nil {
		return nil, dateStart, dateEnd, err
	}
	defer rows.Close()

	var results []entity.SummaryLanding
	for rows.Next() {
		var row entity.SummaryLanding
		err := r.DB.ScanRows(rows, &row)
		if err != nil {
			return nil, dateStart, dateEnd, err
		}
		results = append(results, row)
	}

	return results, dateStart, dateEnd, rows.Err()
}


// helper
func formatQueryIndicatorsRedirection(selects []string, dataType string) []string {
	var formatted []string

	for _, value := range selects {
		switch value {
		case "landing":
			formatted = append(formatted, fmt.Sprintf("SUM(%s) AS %s", value, value))
		case "click_ios":
			formatted = append(formatted, fmt.Sprintf("SUM(%s) AS %s", value, value))
		case "click_android":
			formatted = append(formatted, fmt.Sprintf("SUM(%s) AS %s", value, value))
		case "click_operator":
			formatted = append(formatted, fmt.Sprintf("SUM(%s) AS %s", value, value))
		case "click_non_operator":
			formatted = append(formatted, fmt.Sprintf("SUM(%s) AS %s", value, value))
		default:
			formatted = append(formatted, fmt.Sprintf("AVG(%s) AS %s", value, value))
		}
	}

	formatted = append(formatted,
		"(AVG(response_time) + AVG(response_url_service_time)) AS total_load_time")

	return formatted
}

func (r *BaseModel) GetDataLandingsByCreatedAtRange(start, end time.Time) ([]entity.PixelStorage, error) {
	var landings []entity.PixelStorage
	err := r.DB.
		Where("created_at >= ? AND created_at < ?", start, end).
		Find(&landings).Error

	if err != nil {
		r.Logs.Error(fmt.Sprintf("Failed to get data_landings: %v", err))
		return nil, err
	}

	return landings, nil
}

func (r *BaseModel) SummaryLanding(o entity.SummaryLanding) int {
	result := r.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "summary_date_hour"},
			{Name: "url_service_key"},
			{Name: "campaign_id"},
			{Name: "country"},
			{Name: "operator"},
			{Name: "partner"},
			{Name: "aggregator"},
			{Name: "adnet"},
			{Name: "service"},
		},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"landing":            o.Landing,
			"response_time":      o.ResponseTime,
			"success_rate":       o.SuccessRate,
			"click_ios":          o.ClickIOS,
			"click_android":      o.ClickAndroid,
			"click_operator":     o.ClickOperator,
			"click_non_operator": o.ClickNonOperator,
			"updated_at":         time.Now(),
		}),
	}).Create(&o)

	r.Logs.Debug(fmt.Sprintf("SummaryLanding affected: %d, is error: %#v", result.RowsAffected, result.Error))

	return o.ID
}