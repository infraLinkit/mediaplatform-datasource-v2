package persistence

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
)

func (r *BaseModel) GetDisplayCPAReport(o entity.DisplayCPAReport, allowedCompanies []string, allowedAdnets []string) ([]entity.SummaryCampaign, int64, entity.TotalSummaryCampaign, error) {
	var rows *sql.Rows
	var err error
	var total_rows int64
	var TotalSummaryCampaign entity.TotalSummaryCampaign

	t_query := r.DB.Model(&entity.SummaryCampaign{})
	//fmt.Println("Company: ", o.Company)
	query := r.DB.Model(&entity.SummaryCampaign{}).Select(`
		summary_campaigns.*,
		CASE WHEN campaign_objective='UPLOAD SMS' THEN saaf
			 ELSE
			 (
				CASE 
					WHEN LOWER(client_type) = 'external' 
						THEN (mo_received * poaf) 
					ELSE (total_waki_agency_fee + (po * postback) + technical_fee) 
				END
			)	 
		END AS saaf,
		CASE WHEN campaign_objective='UPLOAD SMS' THEN sbaf
			 ELSE (po * postback)
		END AS sbaf,
		CASE WHEN campaign_objective='UPLOAD SMS' THEN saaf-sbaf
			 ELSE
			 (
				CASE 
					WHEN LOWER(client_type) = 'external' 
						THEN (mo_received * poaf) 
					ELSE (total_waki_agency_fee + (po * postback) + technical_fee) 
				END - (po * postback)
			 )
		END AS revenue
	`).Where("mo_received > 0").Where("company IN ?", allowedCompanies).Where("adnet IN ?", allowedAdnets)
	t_query.Where("mo_received > 0").Where("company IN ?", allowedCompanies).Where("adnet IN ?", allowedAdnets)

	if o.CampaignObjective != "" {
		query.Where("campaign_objective = ? ", o.CampaignObjective)
		t_query.Where("campaign_objective = ? ", o.CampaignObjective)
	} else {
		query.Where("campaign_objective IN ?", []string{"CPA", "UPLOAD SMS", "SINGLE URL S2S"})
		t_query.Where("campaign_objective IN ?", []string{"CPA", "UPLOAD SMS", "SINGLE URL S2S"})
	}

	//fmt.Println("o.CampaignObjective: ", o.CampaignObjective)

	if o.Action == "Search" {
		if o.CampaignId != "" {
			query = query.Where("campaign_id = ?", o.CampaignId)
			t_query = t_query.Where("campaign_id = ?", o.CampaignId)
		}
		if o.UrlServiceKey != "" {
			query = query.Where("url_service_key = ?", o.UrlServiceKey)
			t_query = t_query.Where("url_service_key = ?", o.UrlServiceKey)
		}
		if o.Country != "" {
			query = query.Where("country = ?", o.Country)
			t_query = t_query.Where("country = ?", o.Country)
		}
		if o.Company != "" {
			query = query.Where("company = ?", o.Company)
			t_query = t_query.Where("company = ?", o.Company)
		}
		if o.ClientType != "" {
			query = query.Where("client_type = ?", o.ClientType)
			t_query = t_query.Where("client_type = ?", o.ClientType)
		}
		if o.Operator != "" {
			query = query.Where("operator = ?", o.Operator)
			t_query = t_query.Where("operator = ?", o.Operator)
		}
		if o.CampaignName != "" {
			query = query.Where("campaign_name = ?", o.CampaignName)
			t_query = t_query.Where("campaign_name = ?", o.CampaignName)
		}
		if o.Partner != "" {
			query = query.Where("partner = ?", o.Partner)
			t_query = t_query.Where("partner = ?", o.Partner)
		}
		if len(o.Adnets) > 0 {
			query = query.Where("adnet IN ?", o.Adnets)
			t_query = t_query.Where("adnet IN ?", o.Adnets)
		}
		if o.Service != "" {
			query = query.Where("service = ?", o.Service)
			t_query = t_query.Where("service = ?", o.Service)
		}

		if o.DateRange != "" {
			switch strings.ToUpper(o.DateRange) {
			case "TODAY":
				query = query.Where("summary_date = CURRENT_DATE")
				t_query = t_query.Where("summary_date = CURRENT_DATE")
			case "YESTERDAY":
				query = query.Where("summary_date = CURRENT_DATE - INTERVAL '1 DAY'")
				t_query = t_query.Where("summary_date = CURRENT_DATE - INTERVAL '1 DAY'")
			case "LAST7DAY":
				query = query.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '7 DAY' AND CURRENT_DATE")
				t_query = t_query.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '7 DAY' AND CURRENT_DATE")
			case "LAST30DAY":
				query = query.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '30 DAY' AND CURRENT_DATE")
				t_query = t_query.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '30 DAY' AND CURRENT_DATE")
			case "THISMONTH":
				query = query.Where("summary_date >= DATE_TRUNC('month', CURRENT_DATE)")
				t_query = t_query.Where("summary_date >= DATE_TRUNC('month', CURRENT_DATE)")
			case "LASTMONTH":
				query = query.Where("summary_date BETWEEN DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 MONTH') AND DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '1 DAY'")
				t_query = t_query.Where("summary_date BETWEEN DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 MONTH') AND DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '1 DAY'")
			case "CUSTOMRANGE":
				query = query.Where("summary_date BETWEEN ? AND ?", o.DateBefore, o.DateAfter)
				t_query = t_query.Where("summary_date BETWEEN ? AND ?", o.DateBefore, o.DateAfter)
			default:
				query = query.Where("summary_date = ?", o.DateRange)
				t_query = t_query.Where("summary_date = ?", o.DateRange)
			}
		} else {
			query = query.Where("summary_date = CURRENT_DATE")
			t_query = t_query.Where("summary_date = CURRENT_DATE")
		}
	} else {
		t_query = t_query.Where("summary_date = CURRENT_DATE")
	}

	if o.OrderColumn != "" {
		dir := "ASC"
		if strings.ToUpper(o.OrderDir) == "DESC" {
			dir = "DESC"
		}

		switch o.OrderColumn {
		case "saaf":
			query = query.Order(fmt.Sprintf(`
				CASE 
					WHEN LOWER(client_type) = 'external' 
						THEN (mo_received * poaf) 
					ELSE (total_waki_agency_fee + (po * postback) + technical_fee) 
				END %s
			`, dir))
		case "sbaf":
			query = query.Order(fmt.Sprintf("(po * postback) %s", dir))
		case "revenue":
			query = query.Order(fmt.Sprintf(`
				(
					CASE 
						WHEN LOWER(client_type) = 'external' 
							THEN (mo_received * poaf) 
						ELSE (total_waki_agency_fee + (po * postback) + technical_fee) 
					END - (po * postback)
				) %s
			`, dir))
		default:
			query = query.Order(fmt.Sprintf("%s %s", o.OrderColumn, dir))
		}
	} else {
		query = query.Order("summary_date DESC").Order("id DESC")
	}

	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, err = query_limit.Rows()
	if err != nil {
		return []entity.SummaryCampaign{}, 0, entity.TotalSummaryCampaign{}, err
	}
	defer rows.Close()

	var ss []entity.SummaryCampaign
	for rows.Next() {
		var s entity.SummaryCampaign
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	// GET TOTAL HEADER
	if total_rows > 0 {
		// COUNT THE SUMMARIZE
		_ = t_query.Select(
			`SUM(landing) as landing,
			 SUM(mo_received) as mo_received,
			 SUM(postback) as postback,
			 AVG(po) as price_per_postback,
			 SUM(cost_per_conversion) as cost_per_conversion,
			 SUM(total_waki_agency_fee) as total_waki_agency_fee,
			 SUM(CASE WHEN campaign_objective='UPLOAD SMS' THEN sbaf 
			 		  ELSE po * postback
				 END) as spending_to_adnet, -- sbaf
			 SUM(CASE WHEN campaign_objective='UPLOAD SMS' THEN saaf
			 		  ELSE
						CASE
							WHEN LOWER(client_type) = 'external'
								THEN (mo_received * poaf)
							ELSE (total_waki_agency_fee + (po * postback) + technical_fee)
						END
				 END) as spending, --saaf
			 SUM(technical_fee) as technical_fee,
			 SUM(
			     CASE WHEN campaign_objective='UPLOAD SMS' THEN saaf-sbaf
			     ELSE
					CASE
						WHEN LOWER(client_type) = 'external'
							THEN (mo_received * poaf)
						ELSE (total_waki_agency_fee + (po * postback) + technical_fee)
					END
				 END) - 
				 SUM(CASE WHEN campaign_objective='UPLOAD SMS' THEN 0
				    	 ELSE po * postback
					 END) as waki_revenue,
			 CASE WHEN SUM(landing)>0 THEN ROUND(SUM(mo_received)/SUM(landing)::numeric,5) ELSE 0 END as cr_mo,
			 CASE WHEN SUM(landing)>0 THEN ROUND(SUM(postback)/SUM(landing)::numeric,5) ELSE 0 END as cr_postback,
			 ROUND(AVG(cpa)::numeric,5) as avg_cpa`).Row().Scan(
			&TotalSummaryCampaign.Landing,
			&TotalSummaryCampaign.MoReceived,
			&TotalSummaryCampaign.Postback,
			&TotalSummaryCampaign.PO,
			&TotalSummaryCampaign.CostPerConversion,
			&TotalSummaryCampaign.TotalWakiAgencyFee,
			&TotalSummaryCampaign.SBAF,
			&TotalSummaryCampaign.SAAF,
			&TotalSummaryCampaign.TechnicalFee,
			&TotalSummaryCampaign.WakiRevenue,
			&TotalSummaryCampaign.CrMO,
			&TotalSummaryCampaign.CrPostback,
			&TotalSummaryCampaign.ECPA)
	}

	r.Logs.Debug(fmt.Sprintf("Total data : %d ... \n", len(ss)))

	return ss, total_rows, TotalSummaryCampaign, rows.Err()
}

func (r *BaseModel) CreateCpaReport(s entity.SummaryCampaign) error {
	s.SuccessFP = 0
	s.PO = 0
	s.AgencyFee = 0
	s.TotalWakiAgencyFee = 0
	s.TechnicalFee = 0
	s.CPA = 0
	s.Traffic = 0
	s.CrPostback = 0
	s.CrMO = 0
	return r.DB.Create(&s).Error
}

func (r *BaseModel) UpdateCpaReport(s entity.SummaryCampaign) error {
	if err := r.DB.Model(&entity.SummaryCampaign{}).Where("id = ?", s.ID).Updates(&s).Error; err != nil {
		return err
	}

	// Override field default (misalnya 0) agar ikut terupdate
	resetFields := map[string]interface{}{
		"success_fp":            0,
		"po":                    0,
		"agency_fee":            0,
		"total_waki_agency_fee": 0,
		"technical_fee":         0,
		"cpa":                   0,
		"traffic":               0,
		"cr_postback":           0,
		"cr_mo":                 0,
	}

	return r.DB.Model(&entity.SummaryCampaign{}).Where("id = ?", s.ID).Updates(resetFields).Error
}

func (r *BaseModel) FindSummaryCampaignByUniqueKey(
	summaryDate *time.Time,
	campaignId, country, operator, partner, service, adnet, urlServiceKey, campaignObjective string,
) (entity.SummaryCampaign, error) {
	db := r.DB.Model(&entity.SummaryCampaign{})
	if summaryDate != nil {
		db = db.Where("summary_date = ?", *summaryDate)
	}
	if campaignId != "" {
		db = db.Where("LOWER(campaign_id) = LOWER(?)", campaignId)
	}
	if country != "" {
		db = db.Where("LOWER(country) = LOWER(?)", country)
	}
	if operator != "" {
		db = db.Where("LOWER(operator) = LOWER(?)", operator)
	}
	if partner != "" {
		db = db.Where("LOWER(partner) = LOWER(?)", partner)
	}
	if service != "" {
		db = db.Where("LOWER(service) = LOWER(?)", service)
	}
	if adnet != "" {
		db = db.Where("LOWER(adnet) = LOWER(?)", adnet)
	}
	if urlServiceKey != "" {
		db = db.Where("LOWER(url_service_key) = LOWER(?)", urlServiceKey)
	}
	if campaignObjective != "" {
		db = db.Where("LOWER(campaign_objective) = LOWER(?)", campaignObjective)
	}

	var s entity.SummaryCampaign
	result := db.First(&s)
	return s, result.Error
}

func (r *BaseModel) FindLatestSummaryCampaignByUniqueKey(service, adnet, operator string, partner string) (entity.SummaryCampaign, error) {
	var s entity.SummaryCampaign
	result := r.DB.Where("LOWER(service) = LOWER(?) AND LOWER(adnet) = LOWER(?) AND LOWER(operator) = LOWER(?) AND LOWER(partner)=LOWER(?)", service, adnet, operator, partner).
		Order("summary_date DESC").
		First(&s)
	return s, result.Error
}

func (r *BaseModel) GetDisplayMainstreamReport(o entity.DisplayCPAReport, allowedCompanies []string, allowedAgencies []string) ([]entity.SummaryCampaign, int64, interface{}, error) {
	var rows *sql.Rows
	var err error
	var total_rows int64
	var total_summary entity.TotalSummaryCampaign

	/*

	 */
	t_query := r.DB.Model(&entity.SummaryCampaign{}).Where("campaign_objective LIKE ?", "%MAINSTREAM%").
		Where("mo_received > 0").
		Where("company IN ?", allowedCompanies).
		Where("adnet IN ?", allowedAgencies)

	query := r.DB.Model(&entity.SummaryCampaign{}).Select(`
		summary_campaigns.*,
		saaf,
		sbaf,
		price_per_mo,
		revenue
	`).Where("campaign_objective LIKE ?", "%MAINSTREAM%").
		Where("mo_received > 0").
		Where("company IN ?", allowedCompanies).
		Where("adnet IN ?", allowedAgencies)

	if o.Action == "Search" {
		if o.CampaignId != "" {
			query = query.Where("LOWER(campaign_id) = LOWER(?)", o.CampaignId)
			t_query = t_query.Where("LOWER(campaign_id) = LOWER(?)", o.CampaignId)
		}
		if o.Agency != "" {
			query = query.Where("LOWER(adnet) = LOWER(?)", o.Agency)
			t_query = t_query.Where("LOWER(adnet) = LOWER(?)", o.Agency)
		}
		if o.UrlServiceKey != "" {
			keys := strings.Split(o.UrlServiceKey, ",")
			lowerKeys := make([]string, len(keys))
			for i, k := range keys {
				lowerKeys[i] = strings.ToLower(strings.TrimSpace(k))
			}
			query = query.Where("LOWER(url_service_key) IN ?", lowerKeys)
			t_query = t_query.Where("LOWER(url_service_key) IN ?", lowerKeys)
		}
		if o.Country != "" {
			query = query.Where("LOWER(country) = LOWER(?)", o.Country)
			t_query = t_query.Where("LOWER(country) = LOWER(?)", o.Country)
		}
		if o.Company != "" {
			query = query.Where("LOWER(company) = LOWER(?)", o.Company)
			t_query = t_query.Where("LOWER(company) = LOWER(?)", o.Company)
		}
		if o.ClientType != "" {
			query = query.Where("LOWER(client_type) = LOWER(?)", o.ClientType)
			t_query = t_query.Where("LOWER(client_type) = LOWER(?)", o.ClientType)
		}
		if o.Operator != "" {
			query = query.Where("LOWER(operator) = LOWER(?)", o.Operator)
			t_query = t_query.Where("LOWER(operator) = LOWER(?)", o.Operator)
		}
		if o.CampaignName != "" {
			query = query.Where("LOWER(campaign_name) = LOWER(?)", o.CampaignName)
			t_query = t_query.Where("LOWER(campaign_name) = LOWER(?)", o.CampaignName)
		}
		if o.Channel != "" {
			query = query.Where("LOWER(channel) = LOWER(?)", o.Channel)
			t_query = t_query.Where("LOWER(channel) = LOWER(?)", o.Channel)
		}
		if o.Partner != "" {
			query = query.Where("LOWER(partner) = LOWER(?)", o.Partner)
			t_query = t_query.Where("LOWER(partner) = LOWER(?)", o.Partner)
		}
		if len(o.Adnets) > 0 {
			query = query.Where("adnet IN ?", o.Adnets)
			t_query = t_query.Where("adnet IN ?", o.Adnets)
		}
		if o.Service != "" {
			query = query.Where("LOWER(service) = LOWER(?)", o.Service)
			t_query = t_query.Where("LOWER(service) = LOWER(?)", o.Service)
		}

		// Date range filter
		if o.DateRange != "" {
			switch strings.ToUpper(o.DateRange) {
			case "TODAY":
				query = query.Where("summary_date = CURRENT_DATE")
				t_query = t_query.Where("summary_date = CURRENT_DATE")
			case "YESTERDAY":
				query = query.Where("summary_date = CURRENT_DATE - INTERVAL '1 DAY'")
				t_query = t_query.Where("summary_date = CURRENT_DATE - INTERVAL '1 DAY'")
			case "LAST7DAY":
				query = query.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '7 DAY' AND CURRENT_DATE")
				t_query = t_query.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '7 DAY' AND CURRENT_DATE")
			case "LAST30DAY":
				query = query.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '30 DAY' AND CURRENT_DATE")
				t_query = t_query.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '30 DAY' AND CURRENT_DATE")
			case "THISMONTH":
				query = query.Where("summary_date >= DATE_TRUNC('month', CURRENT_DATE)")
				t_query = t_query.Where("summary_date >= DATE_TRUNC('month', CURRENT_DATE)")
			case "LASTMONTH":
				query = query.Where("summary_date BETWEEN DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 MONTH') AND DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '1 DAY'")
				t_query = t_query.Where("summary_date BETWEEN DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 MONTH') AND DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '1 DAY'")
			case "CUSTOMRANGE":
				query = query.Where("summary_date BETWEEN ? AND ?", o.DateBefore, o.DateAfter)
				t_query = t_query.Where("summary_date BETWEEN ? AND ?", o.DateBefore, o.DateAfter)
			case "ALLDATERANGE":
				// no filter
			default:
				query = query.Where("summary_date = ?", o.DateRange)
				t_query = t_query.Where("summary_date = ?", o.DateRange)
			}
		} else {
			query = query.Where("summary_date = CURRENT_DATE")
			t_query = t_query.Where("summary_date = CURRENT_DATE")
		}
	}

	if o.OrderColumn != "" {
		dir := "ASC"
		if strings.ToUpper(o.OrderDir) == "DESC" {
			dir = "DESC"
		}

		switch o.OrderColumn {
		case "saaf":
			query = query.Order(fmt.Sprintf("saaf %s", dir))
		case "sbaf":
			query = query.Order(fmt.Sprintf("sbaf %s", dir))
		case "price_per_mo":
			query = query.Order(fmt.Sprintf("price_per_mo %s", dir))
		case "revenue":
			query = query.Order(fmt.Sprintf("revenue %s", dir))
		default:
			query = query.Order(fmt.Sprintf("%s %s", o.OrderColumn, dir))
		}
	} else {
		query = query.Order("summary_date DESC").Order("id DESC")
	}

	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, err = query_limit.Rows()
	if err != nil {
		return []entity.SummaryCampaign{}, 0, total_summary, err
	}
	defer rows.Close()

	var ss []entity.SummaryCampaign
	for rows.Next() {
		var s entity.SummaryCampaign
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	if total_rows > 0 {
		// COUNT THE SUMMARIZE
		_ = t_query.Select(
			`SUM(mo_received) as mo_received,
			 SUM(postback) as postback,
			 SUM(landing) as landing,
			 SUM(clicked) as clicked,
			 SUM(saaf) as saaf,
			 SUM(sbaf) as sbaf,
			 SUM(price_per_mo) as price_per_mo,
		     SUM(revenue) as revenue,
			 SUM(po) as po`).Row().Scan(
			&total_summary.MoReceived,
			&total_summary.Postback,
			&total_summary.Landing,
			&total_summary.Clicked,
			&total_summary.SAAF,
			&total_summary.SBAF,
			&total_summary.PricePerMO,
			&total_summary.WakiRevenue,
			&total_summary.PO,
		)
	}

	r.Logs.Debug(fmt.Sprintf("Total data : %d ... \n", len(ss)))
	return ss, total_rows, total_summary, rows.Err()
}

func (r *BaseModel) FindSummaryCampaign(id int) (entity.SummaryCampaign, error) {
	var s entity.SummaryCampaign
	result := r.DB.First(&s, id)
	r.Logs.Debug(fmt.Sprintf("Total data : %d ... \n", result.RowsAffected))
	return s, result.Error
}

func (r *BaseModel) UpdateRatioModel(o entity.SummaryCampaign, id int) error {
	result := r.DB.Exec("UPDATE summary_campaigns SET ratio_send = ?, ratio_receive = ? WHERE ID = ?", o.RatioSend, o.RatioReceive, id)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) UpdatePostbackModel(o entity.SummaryCampaign, id int) error {
	result := r.DB.Exec("UPDATE summary_campaigns SET po = ? WHERE ID = ?", o.Postback, id)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error
}

func (r *BaseModel) UpdateAgencyCostModel(o entity.SummaryCampaign) error {
	result := r.DB.Exec("UPDATE summary_campaigns SET agency_fee = COALESCE(?, agency_fee), cost_per_conversion = COALESCE(?, cost_per_conversion) WHERE DATE_TRUNC('month', summary_date) = DATE_TRUNC('month', CURRENT_DATE) AND summary_date >= DATE_TRUNC('month', CURRENT_DATE) + INTERVAL '0 DAY'", o.AgencyFee, o.CostPerConversion)

	r.Logs.Debug(fmt.Sprintf("affected: %d, is error : %#v", result.RowsAffected, result.Error))

	return result.Error

}

func buildChannelCaseSQL(column string) string {

	var caseSQL strings.Builder
	caseSQL.WriteString("CASE ")

	for k, v := range channelGroupMapModel {
		caseSQL.WriteString(fmt.Sprintf(
			"WHEN LOWER(%s) = '%s' THEN '%s' ",
			column,
			strings.ToLower(k),
			v,
		))
	}

	caseSQL.WriteString(fmt.Sprintf("ELSE %s END", column))

	return caseSQL.String()
}


func applyDateFilter(query *gorm.DB, dateRange, dateBefore, dateAfter string) *gorm.DB {
	switch strings.ToUpper(dateRange) {
	case "TODAY":
		return query.Where("summary_date = CURRENT_DATE")
	case "YESTERDAY":
		return query.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '1 DAY' AND CURRENT_DATE")
	case "LAST7DAY":
		return query.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '7 DAY' AND CURRENT_DATE")
	case "LAST30DAY":
		return query.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '30 DAY' AND CURRENT_DATE")
	case "THISMONTH":
		return query.Where("summary_date >= DATE_TRUNC('month', CURRENT_DATE)")
	case "LASTMONTH":
		return query.Where("summary_date BETWEEN DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 MONTH') AND DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '1 DAY'")
	case "CUSTOMRANGE":
		return query.Where("summary_date BETWEEN ? AND ?", dateBefore, dateAfter)
	case "ALLDATERANGE":
		return query
	default:
		if dateRange != "" {
			return query.Where("summary_date = ?", dateRange)
		}
		return query
	}
}

func applyDataBasedOnOrder(query *gorm.DB, dataBasedOn string) *gorm.DB {
	switch strings.ToUpper(dataBasedOn) {
	case "HIGHEST CONVERSION S2S":
		return query.Order("SUM(CASE WHEN type = 's2s' THEN conversion ELSE 0 END) DESC")
	case "LOWEST CONVERSION S2S":
		return query.Order("SUM(CASE WHEN type = 's2s' THEN conversion ELSE 0 END) ASC")
	case "HIGHEST COST S2S":
		return query.Order("SUM(CASE WHEN type = 's2s' THEN cost ELSE 0 END) DESC")
	case "LOWEST COST S2S":
		return query.Order("SUM(CASE WHEN type = 's2s' THEN cost ELSE 0 END) ASC")
	case "HIGHEST CONVERSION API":
		return query.Order("SUM(CASE WHEN type = 'api' THEN conversion ELSE 0 END) DESC")
	case "LOWEST CONVERSION API":
		return query.Order("SUM(CASE WHEN type = 'api' THEN conversion ELSE 0 END) ASC")
	case "HIGHEST COST API":
		return query.Order("SUM(CASE WHEN type = 'api' THEN cost ELSE 0 END) DESC")
	case "LOWEST COST API":
		return query.Order("SUM(CASE WHEN type = 'api' THEN cost ELSE 0 END) ASC")
	default:
		return query.Order("SUM(CASE WHEN type = 's2s' THEN conversion ELSE 0 END) DESC")
	}
}

func applyDataBasedOnOrderByCountry(query *gorm.DB, dataBasedOn string) *gorm.DB {
	switch strings.ToUpper(dataBasedOn) {

	case "HIGHEST CONVERSION S2S":
		return query.Order("parent_conversion1 DESC, country, operator, channel_type, group_rank")

	case "LOWEST CONVERSION S2S":
		return query.Order("parent_conversion1 ASC, country, operator, channel_type, group_rank")

	case "HIGHEST COST S2S":
		return query.Order("parent_cost1 DESC, country, operator, channel_type, group_rank")

	case "LOWEST COST S2S":
		return query.Order("parent_cost1 ASC, country, operator, channel_type, group_rank")

	case "HIGHEST CONVERSION API":
		return query.Order("parent_conversion2 DESC, country, operator, channel_type, group_rank")

	case "LOWEST CONVERSION API":
		return query.Order("parent_conversion2 ASC, country, operator, channel_type, group_rank")

	case "HIGHEST COST API":
		return query.Order("parent_cost2 DESC, country, operator, channel_type, group_rank")

	case "LOWEST COST API":
		return query.Order("parent_cost2 ASC, country, operator, channel_type, group_rank")

	default:
		return query.Order("parent_conversion1 DESC, country, operator, channel_type, group_rank")
	}
}

// applyDataBasedOnOrderDetail adds ORDER BY for the detail query.
func applyDataBasedOnOrderDetail(query *gorm.DB, dataBasedOn string) *gorm.DB {
	switch strings.ToUpper(dataBasedOn) {
	case "HIGHEST CONVERSION S2S":
		return query.Order("SUM(postback) DESC")
	case "LOWEST CONVERSION S2S":
		return query.Order("SUM(postback) ASC")
	case "HIGHEST COST S2S":
		return query.Order("SUM(sbaf) DESC")
	case "LOWEST COST S2S":
		return query.Order("SUM(sbaf) ASC")
	default:
		return query.Order("SUM(postback) DESC")
	}
}

func applyPagination(query *gorm.DB, page, pageSize int) *gorm.DB {
	if pageSize <= 0 {
		return query // show all
	}
	if page > 1 {
		query = query.Offset((page - 1) * pageSize)
	}
	return query.Limit(pageSize)
}

// scanCostRows iterates *sql.Rows into []entity.CostReport.
func (r *BaseModel) scanCostRows(rows *sql.Rows) ([]entity.CostReport, error) {
	defer rows.Close()
	var results []entity.CostReport
	for rows.Next() {
		var item entity.CostReport
		if err := r.DB.ScanRows(rows, &item); err != nil {
			r.Logs.Error(fmt.Sprintf("Error scanning row: %v", err))
			continue
		}
		results = append(results, item)
	}
	return results, rows.Err()
}

func (r *BaseModel) GetDisplayCostReport(o entity.DisplayCostReport, allowedAdnets []string) ([]entity.CostReport, int64, error) {
	var total_rows int64

	var apiAdnets []string
	_ = r.DB.Model(&entity.ApiPinReport{}).Distinct("adnet").Pluck("adnet", &apiAdnets)
	
	var agencies []string
	_ = r.DB.Model(&entity.Agency{}).
		Distinct("name").
		Pluck("name", &agencies)

	for i, a := range agencies {
		agencies[i] = strings.ToUpper(a)
	}

	allowedAdnets = append(allowedAdnets, apiAdnets...)
	allowedAdnets = append(allowedAdnets, agencies...)

	baseSQL := `
			SELECT
				date_send                           AS summary_date,
				adnet,
				SUM(total_mo)                       AS mo_received,
				SUM(total_postback)                 AS conversion,
				SUM(sbaf)     AS cost,
				SUM(saaf)                           AS saaf,
				'api'                               AS type,
				''                                  AS country,
				''                                  AS operator,
				'API'                               AS channel_type
			FROM api_pin_reports
			WHERE total_mo > 0
			GROUP BY adnet, date_send

			UNION ALL

			SELECT
				summary_date,
				adnet,
				SUM(mo_received)                    AS mo_received,
				SUM(postback)                       AS conversion,
				SUM(CASE WHEN campaign_objective = 'UPLOAD SMS' THEN sbaf
				         ELSE po * postback END)    AS cost,
				SUM(CASE WHEN campaign_objective = 'UPLOAD SMS' THEN saaf
				         WHEN LOWER(client_type) = 'external'   THEN (mo_received * poaf)
				         ELSE (total_waki_agency_fee + (po * postback) + technical_fee)
				    END)                            AS saaf,
				's2s'                               AS type,
				country,
				operator,
				channel                             AS channel_type
			FROM summary_campaigns
			WHERE mo_received > 0
			  AND deleted_at IS NULL
			GROUP BY summary_date, adnet, country, operator, channel
		`

	query := r.DB.Table("(" + baseSQL + ") as t")

	// Adnet allow-list
	if len(o.Adnets) > 0 {
		query = query.Where("adnet IN ?", o.Adnets)
	} else {
		query = query.Where("adnet IN ?", allowedAdnets)
	}

	// Filters
	if o.Adnet != "" {
		query = query.Where("adnet = ?", o.Adnet)
	}
	if o.Country != "" {
		query = query.Where("country = ?", o.Country)
	}
	if o.Operator != "" {
		query = query.Where("operator = ?", o.Operator)
	}
	if o.ChannelType != "" {
		query = query.Where("channel_type = ?", o.ChannelType)
	}
	if o.DateRange != "" {
		query = applyDateFilter(query, o.DateRange, o.DateBefore, o.DateAfter)
	}

	selectClause := `
		adnet,
		SUM(CASE WHEN type = 's2s' THEN conversion ELSE 0 END) AS conversion1,
		SUM(CASE WHEN type = 's2s' THEN cost       ELSE 0 END) AS cost1,
		SUM(CASE WHEN type = 'api' THEN conversion ELSE 0 END) AS conversion2,
		SUM(CASE WHEN type = 'api' THEN cost       ELSE 0 END) AS cost2,
		SUM(CASE WHEN type = 's2s' THEN saaf       ELSE 0 END) AS saaf1,
		SUM(CASE WHEN type = 'api' THEN saaf       ELSE 0 END) AS saaf2`

	// Count distinct adnets matching filters
	countQuery := query.
	Select("adnet").
	Group("adnet")

	r.DB.Table("(?) AS counted", countQuery).
		Count(&total_rows)

	// Order + paginate + fetch
	finalQuery := applyPagination(
		applyDataBasedOnOrder(query, o.DataBasedOn),
		o.Page, o.PageSize,
	)

	rows, err := finalQuery.Select(selectClause).Group("adnet").Rows()
	if err != nil {
		return nil, 0, err
	}
	if rows == nil {
		return []entity.CostReport{}, 0, nil
	}

	results, err := r.scanCostRows(rows)
	r.Logs.Debug(fmt.Sprintf("GetDisplayCostReport total: %d\n", len(results)))
	return results, total_rows, err
}

func (r *BaseModel) GetDisplayCostReportByCountry(o entity.DisplayCostReport, allowedAdnets []string,) ([]entity.CostReport, int64, error) {

	var total_rows int64

	channelCase := buildChannelCaseSQL("channel")

	var apiAdnets []string
	_ = r.DB.Model(&entity.ApiPinReport{}).
		Distinct("adnet").
		Pluck("adnet", &apiAdnets)

	var agencies []string
	_ = r.DB.Model(&entity.Agency{}).
		Distinct("name").
		Pluck("name", &agencies)

	for i, a := range agencies {
		agencies[i] = strings.ToUpper(a)
	}

	allowedAdnets = append(allowedAdnets, apiAdnets...)
	allowedAdnets = append(allowedAdnets, agencies...)

	baseSQL := fmt.Sprintf(`
	SELECT
		date_send AS summary_date,
		adnet,
		SUM(total_mo) AS mo_received,
		SUM(total_postback) AS conversion,
		SUM(sbaf) AS cost,
		SUM(saaf) AS saaf,
		'api' AS type,
		country,
		operator,
		'API' AS channel_type
	FROM api_pin_reports
	WHERE total_mo > 0
	GROUP BY adnet, date_send, country, operator

	UNION ALL

	SELECT
		summary_date,
		adnet,
		SUM(mo_received) AS mo_received,
		SUM(postback) AS conversion,
		SUM(
			CASE
				WHEN campaign_objective = 'UPLOAD SMS'
					THEN sbaf
				ELSE po * postback
			END
		) AS cost,
		SUM(
			CASE
				WHEN campaign_objective = 'UPLOAD SMS'
					THEN saaf
				WHEN LOWER(client_type) = 'external'
					THEN mo_received * poaf
				ELSE total_waki_agency_fee + (po * postback) + technical_fee
			END
		) AS saaf,
		's2s' AS type,
		country,
		operator,
		%s AS channel_type
	FROM summary_campaigns
	WHERE mo_received > 0
	AND deleted_at IS NULL
	GROUP BY summary_date, adnet, country, operator, channel
	`, channelCase)

	query := r.DB.Table("(" + baseSQL + ") as t")

	if len(o.Adnets) > 0 {
		query = query.Where("adnet IN ?", o.Adnets)
	} else {
		query = query.Where("adnet IN ?", allowedAdnets)
	}

	// filters
	if o.Adnet != "" {
		query = query.Where("adnet = ?", o.Adnet)
	}

	if o.Country != "" {
		query = query.Where("country = ?", o.Country)
	}

	if o.Operator != "" {
		query = query.Where("operator = ?", o.Operator)
	}

	if o.ChannelType != "" {
		query = query.Where("channel_type = ?", o.ChannelType)
	}

	if o.DataIndicator != "" {
		query = query.Where("type = ?", strings.ToLower(o.DataIndicator))
	}

	if o.DateRange != "" {
		query = applyDateFilter(query, o.DateRange, o.DateBefore, o.DateAfter)
	}

	selectClause := `
	country,
	operator,
	channel_type,
	adnet,

	SUM(CASE WHEN type='s2s' THEN conversion ELSE 0 END) AS conversion1,
	SUM(CASE WHEN type='s2s' THEN cost ELSE 0 END) AS cost1,

	SUM(CASE WHEN type='api' THEN conversion ELSE 0 END) AS conversion2,
	SUM(CASE WHEN type='api' THEN cost ELSE 0 END) AS cost2,

	SUM(CASE WHEN type='s2s' THEN saaf ELSE 0 END) AS saaf1,
	SUM(CASE WHEN type='api' THEN saaf ELSE 0 END) AS saaf2,

	SUM(SUM(CASE WHEN type='s2s' THEN conversion ELSE 0 END))
	OVER (PARTITION BY country, operator, channel_type) AS parent_conversion1,

	SUM(SUM(CASE WHEN type='s2s' THEN cost ELSE 0 END))
	OVER (PARTITION BY country, operator, channel_type) AS parent_cost1,

	ROW_NUMBER() OVER (
		PARTITION BY country, operator, channel_type
		ORDER BY SUM(CASE WHEN type='s2s' THEN conversion ELSE 0 END) DESC
	) AS group_rank
	`

	groupClause := "country, operator, channel_type, adnet"

	r.DB.Table("(?) AS counted",
		query.Session(&gorm.Session{}).
			Select(groupClause).
			Group(groupClause),
	).Count(&total_rows)

	finalQuery := applyPagination(
		applyDataBasedOnOrderByCountry(query, o.DataBasedOn),
		o.Page,
		o.PageSize,
	)

	rows, err := finalQuery.
		Session(&gorm.Session{}).
		Select(selectClause).
		Group(groupClause).
		Rows()

	if err != nil {
		return nil, 0, err
	}

	results, err := r.scanCostRows(rows)

	return results, total_rows, err
}

func (r *BaseModel) GetDisplayCostReportDetail(o entity.DisplayCostReport, allowedAdnets []string) ([]entity.CostReport, int64, error) {

	var total_rows int64

	channelCase := buildChannelCaseSQL("channel")

	var apiAdnets []string
	_ = r.DB.Model(&entity.ApiPinReport{}).
		Distinct("adnet").
		Pluck("adnet", &apiAdnets)

	var agencies []string
	_ = r.DB.Model(&entity.Agency{}).
		Distinct("name").
		Pluck("name", &agencies)

	for i, a := range agencies {
		agencies[i] = strings.ToUpper(a)
	}

	allowedAdnets = append(allowedAdnets, apiAdnets...)
	allowedAdnets = append(allowedAdnets, agencies...)

	baseSQL := fmt.Sprintf(`

	SELECT
		adnet,
		date_send AS summary_date,
		country,
		operator,
		'API' AS channel_type,
		'' AS short_code,
		'' AS url_after,

		SUM(total_postback) AS landing,
		SUM(total_postback) AS conversion,
		SUM(total_mo) AS mo_received,
		SUM(sbaf) AS cost,
		SUM(saaf) AS saaf,

		'api' AS type

	FROM api_pin_reports
	WHERE total_mo > 0
	GROUP BY adnet, date_send, country, operator

	UNION ALL

	SELECT
		adnet,
		summary_date,
		country,
		operator,
		%s AS channel_type,
		short_code,
		url_after,

		SUM(landing) AS landing,
		SUM(postback) AS conversion,
		SUM(mo_received) AS mo_received,

		SUM(
			CASE
				WHEN campaign_objective = 'UPLOAD SMS'
				THEN sbaf
				ELSE po * postback
			END
		) AS cost,

		SUM(
			CASE
				WHEN campaign_objective = 'UPLOAD SMS'
				THEN saaf
				WHEN LOWER(client_type) = 'external'
				THEN mo_received * poaf
				ELSE total_waki_agency_fee + (po * postback) + technical_fee
			END
		) AS saaf,

		's2s' AS type

	FROM summary_campaigns
	WHERE mo_received > 0
	AND deleted_at IS NULL
	GROUP BY summary_date, adnet, country, operator, channel, short_code, url_after

	`, channelCase)

	query := r.DB.Table("(" + baseSQL + ") as t")

	// adnet filter
	if len(o.Adnets) > 0 {
		query = query.Where("adnet IN ?", o.Adnets)
	} else if len(allowedAdnets) > 0 {
		query = query.Where("adnet IN ?", allowedAdnets)
	}

	if o.Adnet != "" {
		query = query.Where("adnet = ?", o.Adnet)
	}

	if o.Country != "" {
		query = query.Where("country = ?", o.Country)
	}

	if o.Operator != "" {
		query = query.Where("operator = ?", o.Operator)
	}

	if o.ChannelType != "" {
		query = query.Where("channel_type = ?", o.ChannelType)
	}

	if o.DataIndicator != "" {
		query = query.Where("type = ?", strings.ToLower(o.DataIndicator))
	}

	if o.DateRange != "" {
		query = applyDateFilter(query, o.DateRange, o.DateBefore, o.DateAfter)
	}

	selectClause := `
	adnet,
	summary_date,
	country,
	operator,
	channel_type,
	short_code,
	url_after,

	SUM(landing) AS landing,

	CASE WHEN SUM(landing) = 0 THEN 0
		ELSE ROUND((SUM(conversion)::numeric / SUM(landing)::numeric) * 100, 2)
	END AS cr_postback,


	SUM(CASE WHEN type='s2s' THEN conversion ELSE 0 END) AS conversion1,
	SUM(CASE WHEN type='s2s' THEN cost ELSE 0 END) AS cost1,

	SUM(CASE WHEN type='api' THEN conversion ELSE 0 END) AS conversion2,
	SUM(CASE WHEN type='api' THEN cost ELSE 0 END) AS cost2,

	SUM(CASE WHEN type='s2s' THEN saaf ELSE 0 END) AS saaf1,
	SUM(CASE WHEN type='api' THEN saaf ELSE 0 END) AS saaf2
	`

	groupClause := `
	summary_date,
	adnet,
	country,
	operator,
	channel_type,
	short_code,
	url_after
	`

	r.DB.Table("(?) AS counted",
		query.Session(&gorm.Session{}).
			Select(groupClause).
			Group(groupClause),
	).Count(&total_rows)

	if o.DataBasedOn != "" {
		query = applyDataBasedOnOrderDetail(query, o.DataBasedOn)
	} else {
		query = query.Order("summary_date ASC")
	}

	finalQuery := applyPagination(query, o.Page, o.PageSize)

	rows, err := finalQuery.
		Session(&gorm.Session{}).
		Select(selectClause).
		Group(groupClause).
		Rows()

	if err != nil {
		return nil, 0, err
	}

	if rows == nil {
		return []entity.CostReport{}, 0, nil
	}

	results, err := r.scanCostRows(rows)

	r.Logs.Debug(fmt.Sprintf("GetDisplayCostReportDetail total: %d\n", len(results)))

	return results, total_rows, err
}

func (r *BaseModel) GetSummaryReportById(id []string) ([]entity.SummaryCampaign, error) {
	query := r.DB.Model(&entity.SummaryCampaign{}).Where("mo_received > 0 AND id IN ?", id)
	rows, err := query.Rows()

	if err != nil {
		return []entity.SummaryCampaign{}, err
	}

	var results []entity.SummaryCampaign
	for rows.Next() {
		var s entity.SummaryCampaign
		r.DB.ScanRows(rows, &s)
		results = append(results, s)
	}

	return results, nil
}

func (r *BaseModel) AddSMSReport(s entity.SummaryCampaign) error {
	SQL := `
	INSERT INTO summary_campaigns 
	(created_at, updated_at, status, summary_date, url_service_key, campaign_id, campaign_name, country, operator,
	partner, aggregator, adnet, service, short_code, traffic, landing, mo_received, cr, postback,
	total_fp, success_fp, billrate, roi, po, first_push, cost, sbaf, saaf, cpa, revenue,
	url_after, url_before, mo_limit, ratio_send, ratio_receive, company, client_type,
	cost_per_conversion, agency_fee, target_daily_budget, cr_mo, cr_postback, total_waki_agency_fee,
	budget_usage, target_daily_budget_changes, technical_fee, campaign_objective, 
	channel, price_per_mo, target_monthly_budget, poaf)
	SELECT NOW(), NOW(), true, ?,
	cd.url_service_key, cd.campaign_id, cp.name, cd.country, 
	?, ?, cd.aggregator, ?, ?, ?, 0, 0, ?, -- mo_received
	0, ?, -- postback
	0, 0, 0, 0, ?, -- po
	0, -- first_push
	0, -- cost
	?, -- sbaf
	?, -- saaf
	?, -- cpa
	0, cd.url_landing, cd.url_landing, ?, ?,
	?, -- ratio_receive 
	pt.company, pt.client_type, 0, 0, 0, 0, 0, 0, 0, 0, 0, 'UPLOAD SMS',
	'NA', 0, 0, 0
	from campaign_details cd 
	left join partners as pt on pt.name=cd.partner 
	left join campaigns as cp on cp.id = cd.campaign_id::INTEGER where 
	cd.url_service_key = ?
	ON CONFLICT (summary_date,url_service_key,campaign_id,country,operator,partner,adnet,service,campaign_objective) 
	DO UPDATE SET 
		updated_at=NOW(),
		po = EXCLUDED.po,
		cpa = EXCLUDED.cpa,
		sbaf = EXCLUDED.sbaf,
		saaf = EXCLUDED.saaf,
		ratio_send = EXCLUDED.ratio_send,
		ratio_receive = EXCLUDED.ratio_receive,
		mo_received = EXCLUDED.mo_received,
		postback = EXCLUDED.postback `

	q := r.DB.Exec(SQL, s.SummaryDate, s.Operator, s.Partner, s.Adnet,
		s.Service, s.ShortCode, s.MoReceived, s.Postback,
		s.PO, s.SBAF, s.SAAF, s.CPA, 500, s.RatioSend, s.RatioReceive,
		s.URLServiceKey)
	return q.Error
}
