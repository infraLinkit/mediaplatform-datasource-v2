package persistence

import (
	"strconv"
	"strings"
	"time"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
)

func uniqueStrings(arr []string) []string {
	encountered := make(map[string]bool)
	result := []string{}

	for _, v := range arr {
		if !encountered[v] {
			encountered[v] = true
			result = append(result, v)
		}
	}
	return result
}

func GetDaysInMonth(year int, month time.Month) int {
	// 1. Find the first day of the *next* month.
	// We create a time.Time object for the 1st day of the specified month/year.
	firstOfCurrentMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)

	// Add one month to get the first day of the next month.
	firstOfNextMonth := firstOfCurrentMonth.AddDate(0, 1, 0)

	// 2. Subtract one day from the first day of the next month.
	// This gives us the last day of the current month.
	lastOfCurrentMonth := firstOfNextMonth.AddDate(0, 0, -1)

	// 3. The Day component of the last day is the number of days in the month.
	return lastOfCurrentMonth.Day()
}

func (r *BaseModel) CreateSummaryDashboard(s entity.SummaryCampaign) error {
	SQL := `
	INSERT INTO summary_dashboards(summary_date,adnet,company,total_mo,total_postback,total_cpa_mo,total_sms_mo,
		total_mainstream_mo,total_cpa_postback,total_sms_postback,total_mainstream_postback,
		total_spending,total_cpa_spending,total_sms_spending,total_mainstream_spending,
		total_saaf,total_cpa_saaf,total_sms_saaf,total_mainstream_saaf,updated_at)
		SELECT 
		summary_date as date,
		adnet,
		company,
		SUM(mo_received) as total_mo,
		SUM(postback) as total_postback,
		SUM(CASE WHEN campaign_objective='CPA' THEN mo_received ELSE 0 END) as total_cpa_mo,
		SUM(CASE WHEN campaign_objective='UPLOAD SMS' THEN mo_received ELSE 0 END) as total_sms_mo,
		SUM(CASE WHEN campaign_objective='MAINSTREAM' THEN mo_received ELSE 0 END) as total_mainstream_mo,

		SUM(CASE WHEN campaign_objective='CPA' THEN postback ELSE 0 END) as total_cpa_postback,
		SUM(CASE WHEN campaign_objective='UPLOAD SMS' THEN postback ELSE 0 END) as total_sms_postback,
		SUM(CASE WHEN campaign_objective='MAINSTREAM' THEN postback ELSE 0 END) as total_mainstream_postback,

		SUM(sbaf) as total_spending,
		SUM(CASE WHEN campaign_objective='CPA' THEN sbaf ELSE 0 END) as total_cpa_spending,
		SUM(CASE WHEN campaign_objective='UPLOAD SMS' THEN sbaf ELSE 0 END) as total_sms_spending,
		SUM(CASE WHEN campaign_objective='MAINSTREAM' THEN sbaf ELSE 0 END) as total_mainstream_spending,

		SUM(saaf) as total_saaf,
		SUM(CASE WHEN campaign_objective='CPA' THEN saaf ELSE 0 END) as total_cpa_saaf,
		SUM(CASE WHEN campaign_objective='UPLOAD SMS' THEN saaf ELSE 0 END) as total_sms_saaf,
		SUM(CASE WHEN campaign_objective='MAINSTREAM' THEN saaf ELSE 0 END) as total_mainstream_saaf,
		NOW() FROM summary_campaigns WHERE 
		DATE(summary_date) = DATE(?) AND  
		adnet = ? AND 
		company = ? 
		GROUP BY date, adnet, company
		ON CONFLICT(summary_date,adnet,company) DO UPDATE SET
		updated_at=NOW(),
		total_mo=EXCLUDED.total_mo,
		total_postback=EXCLUDED.total_postback,
		total_cpa_mo=EXCLUDED.total_cpa_mo,
		total_sms_mo=EXCLUDED.total_sms_mo,
		total_mainstream_mo=EXCLUDED.total_mainstream_mo,
		total_cpa_postback=EXCLUDED.total_cpa_postback,
		total_sms_postback=EXCLUDED.total_sms_postback,
		total_mainstream_postback=EXCLUDED.total_mainstream_postback,
		total_spending=EXCLUDED.total_spending,
		total_cpa_spending=EXCLUDED.total_cpa_spending,
		total_sms_spending=EXCLUDED.total_sms_spending,
		total_mainstream_spending=EXCLUDED.total_mainstream_spending,
		total_saaf=EXCLUDED.total_saaf,
		total_cpa_saaf=EXCLUDED.total_cpa_saaf,
		total_sms_saaf=EXCLUDED.total_sms_saaf,
		total_mainstream_saaf=EXCLUDED.total_mainstream_saaf
		`
	query := r.DB.Model(&entity.SummaryDashboard{})
	result := query.Exec(SQL, s.SummaryDate, s.Adnet, s.Company)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *BaseModel) GetReport(country string, operator string, client_type string, partner string, service string, campaign_objective string, date_range string, date_before string, date_after string, allowedAdnets []string, allowedCompanies []string) (entity.SummaryDashboardReport, error) {

	// API objective: query api_pin_reports table
	if campaign_objective == "API" {
		apiQ := r.DB.Model(&entity.ApiPinReport{})
		selectDateAPI := "DATE(date_send) as date, "
		switch date_range {
		case "TODAY":
			apiQ = apiQ.Where("date_send = CURRENT_DATE")
		case "YESTERDAY":
			apiQ = apiQ.Where("date_send = CURRENT_DATE - INTERVAL '1 DAY'")
		case "LAST7DAY":
			apiQ = apiQ.Where("date_send BETWEEN CURRENT_DATE - INTERVAL '7 DAY' AND CURRENT_DATE")
		case "LAST30DAY":
			apiQ = apiQ.Where("date_send BETWEEN CURRENT_DATE - INTERVAL '30 DAY' AND CURRENT_DATE")
		case "THISMONTH":
			apiQ = apiQ.Where("date_send >= DATE_TRUNC('month', CURRENT_DATE)")
		case "LASTMONTH":
			apiQ = apiQ.Where("date_send BETWEEN DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 MONTH') AND DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '1 DAY'")
		case "CUSTOMRANGE":
			apiQ = apiQ.Where("date_send BETWEEN ? AND ?", date_before, date_after)
		case "MONTHLY":
			apiQ = apiQ.Where("date_send BETWEEN TO_DATE(?, 'YYYY-MM') AND TO_DATE(?, 'YYYY-MM') + INTERVAL '1 month' - INTERVAL '1 day'", date_before, date_after)
			selectDateAPI = "TO_CHAR(date_send,'YYYY-MM') as date, "
		}
		apiRows, apiErr := apiQ.Select(selectDateAPI +
			"SUM(total_mo) as mo_received, " +
			"SUM(total_postback) as mo_sent, " +
			"SUM(sbaf) as spending_to_adnets, " +
			"SUM(saaf) as spending, " +
			"SUM(waki_revenue) as waki_revenue").Group("date").Order("date ASC").Rows()
		if apiErr != nil {
			return entity.SummaryDashboardReport{DateRange: date_range, Detail: []entity.SummaryDashboardReportDetail{}}, apiErr
		}
		defer apiRows.Close()
		var ss []entity.SummaryDashboardReportDetail
		var date_list []string
		for apiRows.Next() {
			var s entity.SummaryDashboardReportDetail
			r.DB.ScanRows(apiRows, &s)
			s.Date = strings.TrimSuffix(s.Date, "T00:00:00Z")
			ss = append(ss, s)
			date_list = append(date_list, s.Date)
		}
		format := "2006-01-02"
		if date_range == "MONTHLY" {
			format = "2006-01"
		}
		var start_date, end_date string
		if len(date_list) > 0 {
			start_date = date_list[0]
			end_date = date_list[len(date_list)-1]
		}
		start, _ := time.Parse(format, start_date)
		end, _ := time.Parse(format, end_date)
		var date_list_result []string
		current := start
		for !current.After(end) {
			date_list_result = append(date_list_result, current.Format(format))
			current = current.AddDate(0, 0, 1)
		}
		date_list_result = uniqueStrings(date_list_result)
		return entity.SummaryDashboardReport{DateRange: date_range, DateList: date_list_result, Detail: ss}, nil
	}

	query := r.DB.Model(&entity.SummaryCampaign{})

	select_date := " DATE(summary_date) as date, "

	switch date_range {
	case "TODAY":
		query = query.Where("summary_date = CURRENT_DATE")
	case "YESTERDAY":
		query = query.Where("summary_date = CURRENT_DATE - INTERVAL '1 DAY'")
	case "LAST7DAY":
		query = query.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '7 DAY' AND CURRENT_DATE")
	case "LAST30DAY":
		query = query.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '30 DAY' AND CURRENT_DATE")
	case "THISMONTH":
		query = query.Where("summary_date >= DATE_TRUNC('month', CURRENT_DATE)")
	case "LASTMONTH":
		query = query.Where("summary_date BETWEEN DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 MONTH') AND DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '1 DAY'")
	case "CUSTOMRANGE":
		query = query.Where("summary_date BETWEEN ? AND ?", date_before, date_after)
	case "MONTHLY":
		query = query.Where("summary_date BETWEEN TO_DATE(?, 'YYYY-MM') AND TO_DATE(?, 'YYYY-MM') + INTERVAL '1 month' - INTERVAL '1 day' ", date_before, date_after)
		select_date = "TO_CHAR(summary_date,'YYYY-MM') as date, "
	}

	if country != "" {
		query = query.Where("country = ?", country)
	}

	if operator != "" {
		query = query.Where("operator = ?", operator)
	}

	if partner != "" {
		query = query.Where("partner = ?", partner)
	}

	if service != "" {
		query = query.Where("service = ?", service)
	}

	if client_type != "" {
		query = query.Where("client_type = ?", client_type)
	}

	if campaign_objective != "" {
		query = query.Where("campaign_objective = ?", campaign_objective)
	}

	rows, err := query.Select(select_date +
		"SUM(mo_received) as mo_received, " +
		"SUM(postback) as mo_sent, " +
		"SUM(sbaf) as spending_to_adnets, " +
		"SUM(saaf) as spending, " +
		"SUM(saaf)-SUM(sbaf) as waki_revenue").Group("date").Order("date ASC").Rows()

	if err == nil {
		defer rows.Close()

		var ss []entity.SummaryDashboardReportDetail
		var date_list []string
		var date_list_result []string
		var start_date string
		var end_date string

		for rows.Next() {
			var s entity.SummaryDashboardReportDetail

			r.DB.ScanRows(rows, &s)

			s.Date = strings.TrimSuffix(s.Date, "T00:00:00Z")

			ss = append(ss, s)
			date_list = append(date_list, s.Date)
		}

		var SummaryDashboardReport entity.SummaryDashboardReport
		SummaryDashboardReport.DateRange = date_range
		SummaryDashboardReport.Detail = ss

		if len(date_list) == 0 {
			start_date = ""
			end_date = ""
		} else {
			start_date = date_list[0]
			end_date = date_list[len(date_list)-1]
		}

		format := "2006-01-02"

		if date_range == "MONTHLY" {
			format = "2006-01"
		}

		start, _ := time.Parse(format, start_date)
		end, _ := time.Parse(format, end_date)

		current := start
		for !current.After(end) {
			date_list_result = append(date_list_result, current.Format(format))
			current = current.AddDate(0, 0, 1)
		}

		date_list_result = uniqueStrings(date_list_result)

		SummaryDashboardReport.DateList = date_list_result

		return SummaryDashboardReport, nil
	}

	return entity.SummaryDashboardReport{DateRange: date_range, Detail: []entity.SummaryDashboardReportDetail{}}, nil
}

func (r *BaseModel) GetCampaign(order_type string, order_by string, offset string, client_type string, date_range string, date_before string, date_after string, country, service string, allowedAdnets []string, allowedCompanies []string) ([]entity.TopCampaign, error) {
	query := r.DB.Model(&entity.SummaryCampaign{})
	field_order := "SUM(mo_received)"
	desc := "DESC"

	if order_type == "WORST" {
		desc = "ASC"
	}

	if order_type == "WORST" {
		query = query.Having("SUM(mo_received) > 0")
	}

	switch order_by {
	case "MO_RECEIVED":
		field_order = "SUM(mo_received)"
	case "SPENDING":
		field_order = "SUM(saaf)"
	case "REVENUE":
		field_order = "SUM(saaf) - SUM(sbaf)"
	case "PROFIT":
		field_order = "(SUM(saaf)-SUM(sbaf)-SUM(technical_fee))"
	case "ROAS":
		field_order = "SUM(saaf)/NULLIF(SUM(sbaf),0)"
	case "CR_MO":
		field_order = "SUM(mo_received)/NULLIF(SUM(landing),0)"
	case "CR_POSTBACK":
		field_order = "SUM(postback)/NULLIF(SUM(mo_received),0)"
	case "E_CPA":
		field_order = "SUM(sbaf)/NULLIF(SUM(mo_received),0)"
	}

	query.Order(field_order + " " + desc)
	limit, e := strconv.Atoi(offset)

	if e != nil {
		limit = 5
	}

	switch date_range {
	case "TODAY":
		query = query.Where("summary_date = CURRENT_DATE")
	case "YESTERDAY":
		query = query.Where("summary_date = CURRENT_DATE - INTERVAL '1 DAY'")
	case "LAST7DAY":
		query = query.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '7 DAY' AND CURRENT_DATE")
	case "LAST30DAY":
		query = query.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '30 DAY' AND CURRENT_DATE")
	case "THISMONTH":
		query = query.Where("summary_date >= DATE_TRUNC('month', CURRENT_DATE)")
	case "LASTMONTH":
		query = query.Where("summary_date BETWEEN DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 MONTH') AND DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '1 DAY'")
	case "CUSTOMRANGE":
		query = query.Where("summary_date BETWEEN ? AND ?", date_before, date_after)
	}

	if client_type != "" {
		query = query.Where("client_type = ?", client_type)
	}

	if country != "" {
		query = query.Where("country = ?", country)
	}
	if service != "" {
		query = query.Where("service = ?", service)
	}

	rows, err := query.Select(
		`url_service_key,
		MAX(campaign_id) as campaign_id,
		MAX(country) as country,
		MAX(operator) as operator,
		MAX(service) as service,
		MAX(adnet) as adnet,
		MAX(client_type) as client_type,
		SUM(mo_received) as mo,
		SUM(postback) as postback,
		SUM(saaf) as spend,
		SUM(sbaf) as spend_to_adnets,
		SUM(technical_fee) as technical_fee,
		SUM(saaf)-SUM(sbaf) as revenue`).Group("url_service_key").Limit(limit).Rows()

	if err == nil {
		defer rows.Close()

		var ss []entity.TopCampaign

		cohortSums, cohortErr := r.GetCampaignROASCohortSumByCampaign(date_range, date_before, date_after, country, service)
		cohortROI, cohortROIErr := r.GetCampaignROASCohortROIByCampaign(date_range, date_before, date_after, country, service)

		for rows.Next() {
			var s entity.TopCampaign

			r.DB.ScanRows(rows, &s)

			if s.SpendToAdnets > 0 {
				s.ROAS = s.Spend / s.SpendToAdnets * 100
				s.Profit = s.Spend - s.SpendToAdnets - s.TechnicalFee
			}

			s.EstROAS = s.ROAS
			if cohortErr == nil {
				sumGrossRevenue, hasCohort := cohortSums[s.URLServiceKey]
				cac := 0.0
				if s.MO > 0 {
					cac = s.SpendToAdnets / float64(s.MO)
				}
				s.EstROAS = estROASOrFallback(sumGrossRevenue, hasCohort, s.MO, cac, s.ROAS)
			}
			if cohortROIErr == nil {
				roiMonths, hasROI := cohortROI[s.URLServiceKey]
				if hasROI {
					s.ROIMonths = roiMonths
				}
			}

			c := r.DB.Model(&entity.Country{})
			_ = c.Select("name").Where("code=?", s.Country).Row().Scan(&s.CountryName)

			ss = append(ss, s)
		}

		return ss, nil
	}

	return []entity.TopCampaign{}, err
}

func (r *BaseModel) GetDisplayDashboard(date_range string, date_before string, date_after string, client_type string, country, service string, allowedAdnets []string, allowedCompanies []string) (entity.SummaryDashboardData, error) {

	query := r.DB.Model(&entity.SummaryCampaign{})
	query_last_month := r.DB.Model(&entity.SummaryCampaign{})

	var SummaryDashboard entity.SummaryDashboardData

	var date_list []string

	currentTime := time.Now()

	switch date_range {
	case "TODAY":
		date_list = append(date_list, currentTime.Format("2006-01-02"))
	case "YESTERDAY":
		date_list = append(date_list, currentTime.AddDate(0, 0, -1).Format("2006-01-02"))
	case "LAST7DAY":
		for i := 0; i < 7; i++ {
			date := currentTime.AddDate(0, 0, -i).Format("2006-01-02")
			date_list = append(date_list, date)
		}
	case "LAST30DAY":
		current_day := currentTime.Day()
		for i := current_day; i <= 30+current_day; i++ {
			newDate := time.Date(currentTime.Year(), currentTime.Month(), i, 0, 0, 0, 0, time.Local)
			date_list = append(date_list, newDate.Format("2006-01-02"))
		}
	case "THISMONTH":
		current_day := 1
		last_day := GetDaysInMonth(currentTime.Year(), currentTime.Month())
		for i := current_day; i <= last_day; i++ {
			newDate := time.Date(currentTime.Year(), currentTime.Month(), i, 0, 0, 0, 0, time.Local)
			date_list = append(date_list, newDate.Format("2006-01-02"))
		}
	case "LASTMONTH":
		lastMonth := currentTime.AddDate(0, -1, 0)
		current_day := 1
		last_day := GetDaysInMonth(lastMonth.Year(), lastMonth.Month())
		for i := current_day; i <= last_day; i++ {
			newDate := time.Date(lastMonth.Year(), lastMonth.Month(), i, 0, 0, 0, 0, time.Local)
			date_list = append(date_list, newDate.Format("2006-01-02"))
		}
	case "CUSTOMRANGE":
		start, _ := time.Parse("2006-01-02", date_before)
		end, _ := time.Parse("2006-01-02", date_after)
		current := start
		for !current.After(end) {
			date_list = append(date_list, current.Format("2006-01-02"))
			current = current.AddDate(0, 0, 1)
		}
	}

	query = query.Where("summary_date::date IN ?", date_list)

	if client_type != "" {
		query = query.Where("client_type = ?", client_type)
	}

	where := ""

	for _, date := range date_list {
		where +=
			`CASE WHEN EXTRACT(DAY FROM (DATE_TRUNC('month', DATE '` + date + `') - INTERVAL '1 day')) < EXTRACT(DAY FROM DATE '` + date + `')
        THEN '0001-01-01' ELSE DATE '` + date + `' - INTERVAL '1 month'
    	END ,`
	}

	query_last_month = query_last_month.Where("summary_date::date IN(" + strings.TrimSuffix(where, ",") + ")")

	if client_type != "" {
		query_last_month = query_last_month.Where("client_type = ?", client_type)
	}

	if country != "" {
		query = query.Where("country = ?", country)
		query_last_month = query_last_month.Where("country = ?", country)
	}
	if service != "" {
		query = query.Where("service = ?", service)
		query_last_month = query_last_month.Where("service = ?", service)
	}

	// Build DSP / non-DSP adnet lists for channel split
	var dspCodesRaw, nonDspCodesRaw string

	var dsp struct {
		Code  string `gorm:"column:code"`
		IsDsp bool   `gorm:"column:is_dsp"`
	}

	query_adnet := r.DB.Model(&entity.AdnetList{})
	adnetRows, adnetErr := query_adnet.Select("code,is_dsp").Rows()

	if adnetErr == nil {
		defer adnetRows.Close()
		for adnetRows.Next() {
			r.DB.ScanRows(adnetRows, &dsp)
			if dsp.IsDsp {
				dspCodesRaw += "'" + dsp.Code + "',"
			} else {
				nonDspCodesRaw += "'" + dsp.Code + "',"
			}
		}
	}

	dspIn := "false"
	if dspCodesRaw != "" {
		dspIn = "adnet IN (" + strings.TrimSuffix(dspCodesRaw, ",") + ")"
	}
	nonDspIn := "false"
	if nonDspCodesRaw != "" {
		nonDspIn = "adnet IN (" + strings.TrimSuffix(nonDspCodesRaw, ",") + ")"
	}

	var totalAdnetCount int64
	r.DB.Model(&entity.AdnetList{}).Count(&totalAdnetCount)
	SummaryDashboard.TotalAdnet = int(totalAdnetCount)

	var totalActiveAdnet int64

	activeAdnetQuery := r.DB.Model(&entity.SummaryCampaign{}).
		Where("DATE(summary_date) IN ?", date_list).
		Where("landing > 0")

	if client_type != "" {
		activeAdnetQuery = activeAdnetQuery.Where("client_type = ?", client_type)
	}

	if country != "" {
		activeAdnetQuery = activeAdnetQuery.Where("country = ?", country)
	}

	if service != "" {
		activeAdnetQuery = activeAdnetQuery.Where("service = ?", service)
	}

	type adnetCountRow struct{ Count int64 }
	var acr adnetCountRow
	activeAdnetQuery.Select("COUNT(DISTINCT adnet) as count").Scan(&acr)
	totalActiveAdnet = acr.Count

	SummaryDashboard.TotalActiveAdnet = int(totalActiveAdnet)

	selectSQL := `DATE(summary_date) as date,
		SUM(mo_received) as total_mo,
		SUM(saaf) as total_spending,
		SUM(saaf) as total_saaf,
		SUM(sbaf) as spending_to_adnets,
		SUM(CASE WHEN campaign_objective IN ('CPA','UPLOAD SMS','SINGLE URL S2S') AND mo_received > 0 THEN technical_fee ELSE 0 END) as total_technical_fee,
		SUM(CASE WHEN campaign_objective IN('CPA','UPLOAD SMS', 'SINGLE URL S2S') AND ` + nonDspIn + ` THEN saaf ELSE 0 END) as total_s2s_spending,
		0 as total_api_spending,
		SUM(CASE WHEN campaign_objective IN('MAINSTREAM', 'SINGLE URL MAINSTREAM') THEN saaf ELSE 0 END) as total_mainstream_spending,
		SUM(CASE WHEN campaign_objective IN('CPA','UPLOAD SMS', 'SINGLE URL S2S') AND ` + dspIn + ` THEN saaf ELSE 0 END) as total_dsp_spending,
		SUM(CASE WHEN client_type='internal' THEN saaf - sbaf ELSE 0 END) as internal_revenue,
		SUM(CASE WHEN client_type='external' THEN saaf - sbaf ELSE 0 END) as external_revenue,
		SUM(CASE WHEN client_type='internal' THEN saaf ELSE 0 END) as internal_spend,
		SUM(CASE WHEN client_type='external' THEN saaf ELSE 0 END) as external_spend,
		SUM(CASE WHEN campaign_objective IN('CPA','UPLOAD SMS', 'SINGLE URL S2S') AND ` + nonDspIn + ` THEN saaf - sbaf ELSE 0 END) as s2s_revenue,
		SUM(CASE WHEN campaign_objective IN('MAINSTREAM', 'SINGLE URL MAINSTREAM') THEN saaf - sbaf ELSE 0 END) as mainstream_revenue,
		SUM(CASE WHEN campaign_objective IN('CPA','UPLOAD SMS', 'SINGLE URL S2S') AND ` + dspIn + ` THEN saaf ELSE 0 END) as dsp_revenue,
		SUM(landing) as total_landing,
		SUM(clicked) as total_clicked,
		SUM(postback) as total_postback`

	// Pre-fetch all api_pin_reports for the period in one query — avoids N+1
	type apiDayRow struct {
		Date          string  `gorm:"column:date"`
		Spend         float64 `gorm:"column:spend"`
		SpendToAdnets float64 `gorm:"column:spend_to_adnets"`
		MO            int     `gorm:"column:mo"`
		WakiRevenue   float64 `gorm:"column:waki_revenue"`
	}
	var apiDayResults []apiDayRow
	apiPreQ := r.DB.Model(&entity.ApiPinReport{})
	switch date_range {
	case "TODAY":
		apiPreQ = apiPreQ.Where("date_send = CURRENT_DATE")
	case "YESTERDAY":
		apiPreQ = apiPreQ.Where("date_send = CURRENT_DATE - INTERVAL '1 DAY'")
	case "LAST7DAY":
		apiPreQ = apiPreQ.Where("date_send BETWEEN CURRENT_DATE - INTERVAL '7 DAY' AND CURRENT_DATE")
	case "LAST30DAY":
		apiPreQ = apiPreQ.Where("date_send BETWEEN CURRENT_DATE - INTERVAL '30 DAY' AND CURRENT_DATE")
	case "THISMONTH":
		apiPreQ = apiPreQ.Where("date_send >= DATE_TRUNC('month', CURRENT_DATE)")
	case "LASTMONTH":
		apiPreQ = apiPreQ.Where("date_send BETWEEN DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 MONTH') AND DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '1 DAY'")
	case "CUSTOMRANGE":
		apiPreQ = apiPreQ.Where("date_send BETWEEN ? AND ?", date_before, date_after)
	default:
		apiPreQ = apiPreQ.Where("date_send >= DATE_TRUNC('month', CURRENT_DATE)")
	}
	if client_type != "internal" {
		apiPreQ.Select("DATE(date_send) as date, SUM(saaf) as spend, SUM(sbaf) as spend_to_adnets, SUM(total_mo) as mo, SUM(waki_revenue) as waki_revenue").
			Group("DATE(date_send)").Scan(&apiDayResults)
	}
	apiByDate := make(map[string]apiDayRow, len(apiDayResults))
	for _, ar := range apiDayResults {
		apiByDate[strings.TrimSuffix(ar.Date, "T00:00:00Z")] = ar
	}

	rows, err := query.Select(selectSQL).Group("DATE(summary_date)").Order("date ASC").Rows()

	if err == nil {
		defer rows.Close()
		var ss []entity.SummaryDashboardDetail

		for rows.Next() {
			var s entity.SummaryDashboardDetail
			r.DB.ScanRows(rows, &s)
			s.Date = strings.TrimSuffix(s.Date, "T00:00:00Z")

			SummaryDashboard.TotalMO += s.TotalMO
			SummaryDashboard.TotalSpending += s.TotalSpending // saaf
			SummaryDashboard.SpendingToAdnets += s.SpendingToAdnets // sbaf
			SummaryDashboard.TotalTechnicalFee += s.TotalTechnicalFee
			SummaryDashboard.Revenue += s.TotalSaaf - s.SpendingToAdnets // waki_revenue
			SummaryDashboard.TotalS2SSpending += s.TotalS2sSpending
			SummaryDashboard.TotalMainstreamSpending += s.TotalMainstreamSpending
			SummaryDashboard.TotalDSPSpending += s.TotalDSPSpending
			SummaryDashboard.InternalRevenue += s.InternalRevenue
			SummaryDashboard.ExternalRevenue += s.ExternalRevenue
			SummaryDashboard.InternalSpend += s.InternalSpend
			SummaryDashboard.ExternalSpend += s.ExternalSpend
			SummaryDashboard.S2SRevenue += s.S2SRevenue
			SummaryDashboard.MainstreamRevenue += s.MainstreamRevenue
			SummaryDashboard.DSPRevenue += s.DSPRevenue
			SummaryDashboard.TotalLanding += s.TotalLanding
			SummaryDashboard.TotalClicked += s.TotalClicked
			SummaryDashboard.TotalPostback += s.TotalPostback

			if apiRow, ok := apiByDate[s.Date]; ok {
				SummaryDashboard.TotalAPISpending += apiRow.Spend
				SummaryDashboard.SpendingToAdnets += apiRow.SpendToAdnets
				SummaryDashboard.TotalSpending += apiRow.Spend
				SummaryDashboard.Revenue += apiRow.WakiRevenue
				SummaryDashboard.TotalMO += apiRow.MO
			}

			ss = append(ss, s)
		}

		// Add API data for dates not covered by any summary_campaigns row
		if len(apiByDate) > 0 {
			seenDates := make(map[string]bool, len(ss))
			for _, s := range ss {
				seenDates[s.Date] = true
			}
			for date, apiRow := range apiByDate {
				if !seenDates[date] {
					SummaryDashboard.TotalAPISpending += apiRow.Spend
					SummaryDashboard.SpendingToAdnets += apiRow.SpendToAdnets
					SummaryDashboard.TotalSpending += apiRow.Spend
					SummaryDashboard.Revenue += apiRow.WakiRevenue
					SummaryDashboard.TotalMO += apiRow.MO
				}
			}
		}

		// Compute derived metrics
		if SummaryDashboard.SpendingToAdnets > 0 {
			SummaryDashboard.ROAS = SummaryDashboard.TotalSpending / SummaryDashboard.SpendingToAdnets * 100
		}
		if SummaryDashboard.TotalMO > 0 {
			SummaryDashboard.ECPA = SummaryDashboard.SpendingToAdnets / float64(SummaryDashboard.TotalMO)
			SummaryDashboard.CAC = SummaryDashboard.ECPA
		}
		// The campaign_roas_cohorts table has no client_type column, so cohort
		// data can't be honestly split by segment. Only apply the cohort
		// override when client_type == "" (both segments combined), since
		// that's the only case where the realized ROAS above is also
		// unsegmented and the comparison stays apples-to-apples.
		if client_type == "" {
			// est_ltv = SUM(estimated_gross_revenue_full) / total_mo, est_roas = est_ltv / cac
			// — combines Mart's revenue projection with our own volume/cost
			// data, rather than trusting Mart's own precomputed ratio.
			cohortROAS, cohortROIMonths, cohortOK := r.GetCampaignROASCohortAgg(date_list, country, service, SummaryDashboard.TotalMO, SummaryDashboard.CAC)
			if cohortOK {
				SummaryDashboard.EstROAS = cohortROAS * 100
				// roi_months_payback is already in months, not a ratio — do not scale.
				SummaryDashboard.EstROI = cohortROIMonths
			} else {
				// 0 is a sentinel for "no cohort data" — the blade renders
				// "—" for 0 rather than duplicating the realized ROAS as if
				// it were a separate estimate.
				SummaryDashboard.EstROAS = 0
				SummaryDashboard.EstROI = 0
			}
		} else {
			SummaryDashboard.EstROAS = 0
			SummaryDashboard.EstROI = 0
		}
		SummaryDashboard.Profit = SummaryDashboard.TotalSpending - SummaryDashboard.SpendingToAdnets - SummaryDashboard.TotalTechnicalFee
		if SummaryDashboard.TotalSpending > 0 {
			SummaryDashboard.MarginPct = SummaryDashboard.Profit / SummaryDashboard.TotalSpending * 100
			// ROI realized == margin on spend, same formula the CMS blade used
			// to compute client-side from profit/total_spending.
			SummaryDashboard.ROI = SummaryDashboard.MarginPct
		}
		if SummaryDashboard.InternalSpend > 0 {
			SummaryDashboard.InternalROAS = SummaryDashboard.InternalRevenue / SummaryDashboard.InternalSpend * 100
		}
		if SummaryDashboard.ExternalSpend > 0 {
			SummaryDashboard.ExternalROAS = SummaryDashboard.ExternalRevenue / SummaryDashboard.ExternalSpend * 100
		}
		// ECPA/CAC computed earlier in this function, before the cohort block —
		// GetCampaignROASCohortAgg needs them.
		if SummaryDashboard.Revenue > 0 {
			SummaryDashboard.RecoveryDays = SummaryDashboard.SpendingToAdnets * 30.0 / SummaryDashboard.Revenue
		}

		// Forecast: Actual / running_days * days_in_month
		today := currentTime.Format("2006-01-02")
		running_days := 0
		for _, d := range date_list {
			if d <= today {
				running_days++
			}
		}
		days_in_month := GetDaysInMonth(currentTime.Year(), currentTime.Month())
		SummaryDashboard.RunningDays = running_days
		SummaryDashboard.DaysInMonth = days_in_month
		if running_days > 0 {
			ratio := float64(days_in_month) / float64(running_days)
			SummaryDashboard.ForecastMO = int(float64(SummaryDashboard.TotalMO) * ratio)
			SummaryDashboard.ForecastRevenue = SummaryDashboard.Revenue * ratio
			SummaryDashboard.ForecastSpending = SummaryDashboard.TotalSpending * ratio
			SummaryDashboard.ForecastProfit = SummaryDashboard.Profit * ratio
		}

		// Target budget from budget_ios
		currentMonth := currentTime.Format("2006-01")
		budgetQ := r.DB.Model(&entity.BudgetIO{}).Where("month = ?", currentMonth)
		if country != "" {
			budgetQ = budgetQ.Where("country = ?", country)
		}
		var targetBudget float64
		budgetQ.Select("COALESCE(SUM(io_target), 0)").Scan(&targetBudget)
		SummaryDashboard.TargetBudget = targetBudget

		// GET 1 MONTH PRIOR DATA
		priorData := make(map[string]map[string]float64)

		priorRows, _ := query_last_month.Select(selectSQL).Group("DATE(summary_date)").Order("date ASC").Rows()

		if priorRows != nil {
			defer priorRows.Close()
			var sl []entity.SummaryDashboardDetail

			for priorRows.Next() {
				var s entity.SummaryDashboardDetail
				r.DB.ScanRows(priorRows, &s)
				s.Date = strings.TrimSuffix(s.Date, "T00:00:00Z")
				sl = append(sl, s)

				priorData[s.Date] = make(map[string]float64)
				priorData[s.Date]["total_spending"] = s.TotalSpending // saaf
				priorData[s.Date]["total_mo"] = float64(s.TotalMO)
				priorData[s.Date]["total_revenue"] = s.TotalSaaf - s.SpendingToAdnets // waki_revenue
			}
		}

		SummaryDashboard.DateList = date_list
		var DetailChartData []entity.DetailChartData

		for _, date := range SummaryDashboard.DateList {
			var DetailChart entity.DetailChartData
			var last_date string

			t, _ := time.Parse("2006-01-02", date)

			m1 := t.AddDate(0, -1, 0).Month().String()
			m2, _ := time.Parse("2006-01-02", date)

			if m1 == m2.Month().String() {
				last_date = ""
			} else {
				last_date = t.AddDate(0, -1, 0).Format("2006-01-02")
			}

			DetailChart.Date = date
			DetailChart.LastMonthDate = last_date
			DetailChart.TotalMO = 0
			DetailChart.TotalSpending = 0
			DetailChart.TotalRevenue = 0
			DetailChart.LastMonthTotalMO = 0
			DetailChart.LastMonthTotalSpending = 0
			DetailChart.LastMonthTotalRevenue = 0

			if innerMap, ok := priorData[last_date]; ok {
				if val, exists := innerMap["total_mo"]; exists {
					DetailChart.LastMonthTotalMO = int(val)
				}
				if val, exists := innerMap["total_spending"]; exists {
					DetailChart.LastMonthTotalSpending = val
				}
				if val, exists := innerMap["total_revenue"]; exists {
					DetailChart.LastMonthTotalRevenue = val
				}
			}

			for _, detail := range ss {
				if date == detail.Date {
					DetailChart.TotalMO = detail.TotalMO
					DetailChart.TotalSpending = detail.TotalSpending // saaf
					DetailChart.TotalRevenue = detail.TotalSaaf - detail.SpendingToAdnets // waki_revenue
					DetailChart.TotalTechnicalFee = detail.TotalTechnicalFee
				}
			}

			// Add API data for this date (API = external only)
			if client_type != "internal" {
				var total_spending_api, total_waki_rev_api float64
				var total_mo_api int
				api_query := r.DB.Model(&entity.ApiPinReport{})
				api_query.Where("date_send = ? ", date)

				apiErr := api_query.Select("SUM(saaf),SUM(total_mo),SUM(waki_revenue)").Limit(1).Row().Scan(&total_spending_api, &total_mo_api, &total_waki_rev_api)

				if apiErr == nil {
					DetailChart.TotalMO += total_mo_api
					DetailChart.TotalSpending += total_spending_api
					DetailChart.TotalRevenue += total_waki_rev_api
				}

				if last_date != "" {
					var last_spending_api, last_waki_rev_api float64
					var last_mo_api int
					api_query = r.DB.Model(&entity.ApiPinReport{})
					api_query.Where("date_send = ? ", last_date)

					apiErr = api_query.Select("SUM(saaf),SUM(total_mo),SUM(waki_revenue)").Limit(1).Row().Scan(&last_spending_api, &last_mo_api, &last_waki_rev_api)

					if apiErr == nil {
						DetailChart.LastMonthTotalMO += last_mo_api
						DetailChart.LastMonthTotalSpending += last_spending_api
						DetailChart.LastMonthTotalRevenue += last_waki_rev_api
					}
				}
			}

			if DetailChart.TotalSpending > 0 {
				DetailChart.TotalROAS = DetailChart.TotalRevenue / DetailChart.TotalSpending * 100
			}

			DetailChartData = append(DetailChartData, DetailChart)
		}

		SummaryDashboard.DetailChartData = DetailChartData
		return SummaryDashboard, nil
	}

	SummaryDashboard.DateList = date_list
	return SummaryDashboard, nil
}
