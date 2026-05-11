package persistence

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
)

func (r *BaseModel) GetTrafficReport(params entity.TrafficReportParams) ([]entity.SummaryTraffic, time.Time, time.Time, error) {
	var rows *sql.Rows

	query := r.DB.Model(&entity.SummaryTraffic{})

	if params.DataType == "" {
		params.DataType = "daily_report"
	}

	selectedFields := []string{"country", "campaign_id", "url_service_key", "campaign_name", "partner", "operator", "service", "adnet", "url_warp_landing"}

	switch params.DataType {
	case "monthly_report":
		selectedFields = append(selectedFields, "DATE_TRUNC('month', summary_date_hour) as summary_date_hour")
	case "daily_report":
		selectedFields = append(selectedFields, "DATE(summary_date_hour) as summary_date_hour")
	case "hourly_report":
		selectedFields = append(selectedFields, "summary_date_hour")
	}

	formattedIndicators := formatQueryIndicatorsTrafficReport(params.DataIndicators, params.DataType)
	selectedFields = append(selectedFields, formattedIndicators...)

	query.Select(selectedFields)

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

	dateStart = time.Date(dateStart.Year(), dateStart.Month(), dateStart.Day(), 0, 0, 0, 0, dateStart.Location())
	dateEnd = time.Date(dateEnd.Year(), dateEnd.Month(), dateEnd.Day(), 23, 59, 59, 999999999, dateEnd.Location())

	query.Where("summary_date_hour BETWEEN ? AND ?", dateStart, dateEnd)

	groupFields := []string{"country", "campaign_id", "url_service_key", "campaign_name", "partner", "operator", "service", "adnet", "url_warp_landing"}

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
		return nil, dateStart, dateEnd, err
	}
	defer rows.Close()

	var results []entity.SummaryTraffic

	for rows.Next() {
		var row entity.SummaryTraffic
		r.DB.ScanRows(rows, &row)
		results = append(results, row)

	}

	return results, dateStart, dateEnd, rows.Err()
}



func (r *BaseModel) GetTrafficReportHourly(params entity.TrafficReportParams) ([]entity.SummaryTraffic, time.Time, time.Time, error) {
	var rows *sql.Rows

	query := r.DB.Model(&entity.SummaryTraffic{})

	selectedFields := []string{
		"country", "campaign_id", "url_service_key", "campaign_name",
		"partner", "operator", "service", "adnet", "url_warp_landing",
		"DATE_TRUNC('hour', summary_date_hour) as summary_date_hour",
	}

	formattedIndicators := formatQueryIndicatorsTrafficReport(params.DataIndicators, "hourly_report")
	selectedFields = append(selectedFields, formattedIndicators...)

	query.Select(selectedFields)

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

	var dateStart, dateEnd time.Time
	loc := time.FixedZone("Asia/Jakarta", 7*3600)

	if strings.ToUpper(params.DateRange) == "CUSTOM_RANGE" {
		splitTime := strings.Split(params.DateCustomRange, "to")
		dateStart, _ = time.ParseInLocation("2006-01-02", strings.TrimSpace(splitTime[0]), loc)
		dateEnd, _ = time.ParseInLocation("2006-01-02", strings.TrimSpace(splitTime[1]), loc)
	} else {
		today := time.Now().In(loc)
		dateStart = today
		dateEnd = today
	}

	dateStart = time.Date(dateStart.Year(), dateStart.Month(), dateStart.Day(), 0, 0, 0, 0, loc)
	dateEnd = time.Date(dateEnd.Year(), dateEnd.Month(), dateEnd.Day(), 23, 59, 59, 999999999, loc)

	query.Where("summary_date_hour BETWEEN ? AND ?", dateStart, dateEnd)

	groupFields := []string{
		"country", "campaign_id", "url_service_key", "campaign_name",
		"partner", "operator", "service", "adnet", "url_warp_landing",
		"DATE_TRUNC('hour', summary_date_hour)",
	}
	query.Group(strings.Join(groupFields, ", "))

	rows, err := query.Unscoped().Rows()
	if err != nil {
		return nil, dateStart, dateEnd, err
	}
	defer rows.Close()

	var results []entity.SummaryTraffic
	for rows.Next() {
		var row entity.SummaryTraffic
		err := r.DB.ScanRows(rows, &row)
		if err != nil {
			return nil, dateStart, dateEnd, err
		}
		results = append(results, row)
	}

	return results, dateStart, dateEnd, rows.Err()
}


func formatQueryIndicatorsTrafficReport(selects []string, dataType string) []string {
	var formatted []string

	for _, value := range selects {
		switch value {
		case "landing", "mo_received":
			formatted = append(formatted, fmt.Sprintf("SUM(%s) AS %s", value, value))
		case "cr_mo", "first_push":
			formatted = append(formatted, fmt.Sprintf("AVG(%s) * 100 AS %s", value, value))		
		default:
			formatted = append(formatted, fmt.Sprintf("SUM(%s) AS %s", value, value))
		}
	}

	return formatted
}
