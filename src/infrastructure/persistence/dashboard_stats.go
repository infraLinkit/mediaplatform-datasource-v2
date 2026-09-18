package persistence

import (
	"fmt"
	"sort"
	"strings"

	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
)

func topNKeys(m map[string]float64, n int) []string {
	type kv struct {
		k string
		v float64
	}
	kvs := make([]kv, 0, len(m))
	for k, v := range m {
		kvs = append(kvs, kv{k, v})
	}
	sort.Slice(kvs, func(i, j int) bool { return kvs[i].v > kvs[j].v })
	out := make([]string, 0, n)
	for i, p := range kvs {
		if i >= n {
			break
		}
		out = append(out, p.k)
	}
	return out
}

func containsStr(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func (r *BaseModel) GetCountryStats(date_range, date_before, date_after, country, service string, allowedAdnets, allowedCompanies []string) ([]entity.CountryStat, error) {
	query := r.DB.Model(&entity.SummaryCampaign{})
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
	default:
		query = query.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '7 DAY' AND CURRENT_DATE")
	}
	if country != "" {
		query = query.Where("country = ?", country)
	}
	if service != "" {
		query = query.Where("service = ?", service)
	}
	if len(allowedAdnets) > 0 {
		query = query.Where("adnet IN ?", allowedAdnets)
	}
	if len(allowedCompanies) > 0 {
		query = query.Where("company IN ?", allowedCompanies)
	}
	rows, err := query.Select("country, SUM(sbaf) as spend, SUM(saaf) as revenue, SUM(mo_received) as mo").
		Group("country").Order("spend DESC").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []entity.CountryStat
	var total float64
	for rows.Next() {
		var s entity.CountryStat
		r.DB.ScanRows(rows, &s)
		results = append(results, s)
		total += s.Spend
	}
	for i := range results {
		if results[i].Spend > 0 {
			results[i].ROAS = results[i].Revenue / results[i].Spend * 100
		}
		if total > 0 {
			results[i].Share = results[i].Spend / total * 100
		}
	}
	return results, nil
}

func (r *BaseModel) GetOpsStats(date_range, date_before, date_after, country, service string, allowedAdnets, allowedCompanies []string) (entity.OpsStats, error) {
	var stats entity.OpsStats
	var dspCodesRaw string
	adnetRows, adnetErr := r.DB.Model(&entity.AdnetList{}).Select("code, is_dsp").Rows()
	if adnetErr == nil {
		defer adnetRows.Close()
		var row struct {
			Code  string `gorm:"column:code"`
			IsDsp bool   `gorm:"column:is_dsp"`
		}
		for adnetRows.Next() {
			r.DB.ScanRows(adnetRows, &row)
			if row.IsDsp {
				dspCodesRaw += "'" + row.Code + "',"
			}
		}
	}
	dspIn := "false"
	if dspCodesRaw != "" {
		dspIn = "adnet IN (" + strings.TrimSuffix(dspCodesRaw, ",") + ")"
	}
	q := r.DB.Model(&entity.SummaryCampaign{})
	switch date_range {
	case "TODAY":
		q = q.Where("summary_date = CURRENT_DATE")
	case "YESTERDAY":
		q = q.Where("summary_date = CURRENT_DATE - INTERVAL '1 DAY'")
	case "LAST30DAY":
		q = q.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '30 DAY' AND CURRENT_DATE")
	case "THISMONTH":
		q = q.Where("summary_date >= DATE_TRUNC('month', CURRENT_DATE)")
	case "LASTMONTH":
		q = q.Where("summary_date BETWEEN DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 MONTH') AND DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '1 DAY'")
	case "CUSTOMRANGE":
		q = q.Where("summary_date BETWEEN ? AND ?", date_before, date_after)
	default:
		q = q.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '7 DAY' AND CURRENT_DATE")
	}
	q = q.Where("mo_received > 0")
	if country != "" {
		q = q.Where("country = ?", country)
	}
	if service != "" {
		q = q.Where("service = ?", service)
	}
	if len(allowedAdnets) > 0 {
		q = q.Where("adnet IN ?", allowedAdnets)
	}
	if len(allowedCompanies) > 0 {
		q = q.Where("company IN ?", allowedCompanies)
	}
	var res struct {
		Total int `gorm:"column:total"`
		S2S   int `gorm:"column:s2s"`
		MS    int `gorm:"column:ms"`
		DSP   int `gorm:"column:dsp"`
	}
	_ = q.Select(fmt.Sprintf(`
		COUNT(DISTINCT url_service_key) as total,
		COUNT(DISTINCT CASE WHEN LOWER(channel) IN('s2s','cpa','telco_channel') AND NOT(%s) THEN url_service_key END) as s2s,
		COUNT(DISTINCT CASE WHEN LOWER(channel) NOT IN('s2s','cpa','telco_channel','dsp') AND NOT(%s) THEN url_service_key END) as ms,
		COUNT(DISTINCT CASE WHEN %s THEN url_service_key END) as dsp
	`, dspIn, dspIn, dspIn)).Scan(&res).Error
	stats.TotalCampaigns = res.Total
	stats.S2SCampaigns = res.S2S
	stats.MSCampaigns = res.MS
	stats.DSPCampaigns = res.DSP
	var avgLoad float64
	lq := r.DB.Model(&entity.SummaryLanding{})
	switch date_range {
	case "TODAY":
		lq = lq.Where("summary_date_hour >= CURRENT_DATE AND summary_date_hour < CURRENT_DATE + INTERVAL '1 DAY'")
	case "YESTERDAY":
		lq = lq.Where("summary_date_hour >= CURRENT_DATE - INTERVAL '1 DAY' AND summary_date_hour < CURRENT_DATE")
	case "LAST30DAY":
		lq = lq.Where("summary_date_hour >= CURRENT_DATE - INTERVAL '30 DAY' AND summary_date_hour < CURRENT_DATE + INTERVAL '1 DAY'")
	case "THISMONTH":
		lq = lq.Where("summary_date_hour >= DATE_TRUNC('month', CURRENT_DATE)")
	case "LASTMONTH":
		lq = lq.Where("summary_date_hour >= DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 MONTH') AND summary_date_hour < DATE_TRUNC('month', CURRENT_DATE)")
	case "CUSTOMRANGE":
		lq = lq.Where("summary_date_hour >= ? AND summary_date_hour < ?::date + INTERVAL '1 DAY'", date_before, date_after)
	default:
		lq = lq.Where("summary_date_hour >= CURRENT_DATE - INTERVAL '7 DAY' AND summary_date_hour < CURRENT_DATE + INTERVAL '1 DAY'")
	}
	if country != "" {
		lq = lq.Where("country = ?", country)
	}
	if service != "" {
		lq = lq.Where("service = ?", service)
	}
	lq.Where("total_load_time > 0").Select("AVG(total_load_time)").Scan(&avgLoad)
	stats.AvgLoadTime = avgLoad
	return stats, nil
}

func (r *BaseModel) GetAlerts(country, service string, allowedAdnets, allowedCompanies []string) []entity.AlertItem {
	var alerts []entity.AlertItem
	// High eCPA (7d)
	type ecpaRow struct {
		CampaignId string  `gorm:"column:campaign_id"`
		Country    string  `gorm:"column:country"`
		Operator   string  `gorm:"column:operator"`
		Service    string  `gorm:"column:service"`
		Adnet      string  `gorm:"column:adnet"`
		ECPA       float64 `gorm:"column:ecpa"`
	}
	var highEcpa []ecpaRow
	q1 := r.DB.Model(&entity.SummaryCampaign{}).
		Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '7 DAY' AND CURRENT_DATE")
	if country != "" {
		q1 = q1.Where("country = ?", country)
	}
	if service != "" {
		q1 = q1.Where("service = ?", service)
	}
	if len(allowedAdnets) > 0 {
		q1 = q1.Where("adnet IN ?", allowedAdnets)
	}
	if len(allowedCompanies) > 0 {
		q1 = q1.Where("company IN ?", allowedCompanies)
	}
	q1.Select(`campaign_id, MAX(country) as country, MAX(operator) as operator,
		MAX(service) as service, MAX(adnet) as adnet,
		SUM(sbaf)/NULLIF(SUM(mo_received),0) as ecpa`).
		Group("campaign_id").
		Having("SUM(mo_received) > 0 AND SUM(sbaf)/NULLIF(SUM(mo_received),0) > 0.170").
		Order("ecpa DESC").Limit(3).Scan(&highEcpa)
	for _, h := range highEcpa {
		alerts = append(alerts, entity.AlertItem{
			Type:   "warn",
			Head:   fmt.Sprintf("eCPA above target on %s · %s", h.Country, h.Service),
			Meta:   fmt.Sprintf("$%.3f vs target $0.170 · adnet %s", h.ECPA, h.Adnet),
			Action: "Investigate publishers →",
		})
	}
	// Zero MO today
	type zeroRow struct {
		CampaignId string `gorm:"column:campaign_id"`
		Country    string `gorm:"column:country"`
		Adnet      string `gorm:"column:adnet"`
	}
	var zeroMo []zeroRow
	q2 := r.DB.Model(&entity.SummaryCampaign{}).Where("summary_date = CURRENT_DATE")
	if country != "" {
		q2 = q2.Where("country = ?", country)
	}
	if service != "" {
		q2 = q2.Where("service = ?", service)
	}
	if len(allowedAdnets) > 0 {
		q2 = q2.Where("adnet IN ?", allowedAdnets)
	}
	if len(allowedCompanies) > 0 {
		q2 = q2.Where("company IN ?", allowedCompanies)
	}
	q2.Select("campaign_id, MAX(country) as country, MAX(adnet) as adnet").
		Group("campaign_id").
		Having("SUM(mo_received) = 0 AND SUM(sbaf) > 0").Limit(10).Scan(&zeroMo)
	if len(zeroMo) > 0 {
		var ids, adnets []string
		for i, z := range zeroMo {
			if i < 3 {
				ids = append(ids, z.CampaignId)
				adnets = append(adnets, z.Adnet)
			}
		}
		alerts = append(alerts, entity.AlertItem{
			Type:   "neg",
			Head:   fmt.Sprintf("%d campaign(s) with 0 MO today", len(zeroMo)),
			Meta:   strings.Join(ids, ", ") + " · adnets: " + strings.Join(adnets, ", "),
			Action: "View affected campaigns →",
		})
	}
	// Outperforming (ROAS > 250%)
	type outRow struct {
		CampaignId string  `gorm:"column:campaign_id"`
		Country    string  `gorm:"column:country"`
		Adnet      string  `gorm:"column:adnet"`
		ROAS       float64 `gorm:"column:roas"`
		Spend      float64 `gorm:"column:spend"`
	}
	var outperf []outRow
	q3 := r.DB.Model(&entity.SummaryCampaign{}).
		Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '7 DAY' AND CURRENT_DATE")
	if country != "" {
		q3 = q3.Where("country = ?", country)
	}
	if service != "" {
		q3 = q3.Where("service = ?", service)
	}
	if len(allowedAdnets) > 0 {
		q3 = q3.Where("adnet IN ?", allowedAdnets)
	}
	if len(allowedCompanies) > 0 {
		q3 = q3.Where("company IN ?", allowedCompanies)
	}
	q3.Select(`campaign_id, MAX(country) as country, MAX(adnet) as adnet,
		SUM(saaf)/NULLIF(SUM(sbaf),0)*100 as roas, SUM(sbaf) as spend`).
		Group("campaign_id").
		Having("SUM(sbaf) > 0 AND SUM(saaf)/NULLIF(SUM(sbaf),0)*100 > 250").
		Order("roas DESC").Limit(1).Scan(&outperf)
	for _, o := range outperf {
		alerts = append(alerts, entity.AlertItem{
			Type:   "pos",
			Head:   fmt.Sprintf("%s-%s outperforming · scale opportunity", o.Country, o.CampaignId),
			Meta:   fmt.Sprintf("%s on %s · ROAS %.0f%% · 7d spend $%.0f", o.Country, o.Adnet, o.ROAS, o.Spend),
			Action: "Scale by +30% →",
		})
	}
	return alerts
}

func (r *BaseModel) GetRollup(date_range, date_before, date_after, client_type, country, service string, allowedAdnets, allowedCompanies []string) ([]entity.RollupRow, error) {
	query := r.DB.Model(&entity.SummaryCampaign{})
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
	default:
		query = query.Where("summary_date >= DATE_TRUNC('month', CURRENT_DATE)")
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
	if len(allowedAdnets) > 0 {
		query = query.Where("adnet IN ?", allowedAdnets)
	}
	if len(allowedCompanies) > 0 {
		query = query.Where("company IN ?", allowedCompanies)
	}
	rows, err := query.Select(`country, operator, service,
		MAX(client_type) as client_type,
		SUM(mo_received) as mo,
		SUM(sbaf) as spend,
		SUM(saaf) as revenue`).
		Group("country, operator, service").Order("spend DESC").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []entity.RollupRow

	cohortSums, cohortErr := r.GetCampaignROASCohortSumByRollup(date_range, date_before, date_after, country, service)
	cohortROI, cohortROIErr := r.GetCampaignROASCohortROIByRollup(date_range, date_before, date_after, country, service)

	for rows.Next() {
		var row entity.RollupRow
		r.DB.ScanRows(rows, &row)
		if row.Spend > 0 {
			row.ROAS = row.Revenue / row.Spend * 100
		}
		if row.Revenue > 0 {
			row.MarginPct = (row.Revenue - row.Spend) / row.Revenue * 100
			row.RecoveryDays = row.Spend * 30.0 / row.Revenue
		}
		if row.MO > 0 {
			row.CAC = row.Spend / float64(row.MO)
		}

		row.EstROAS = row.ROAS
		if cohortErr == nil {
			sumGrossRevenue, hasCohort := cohortSums[row.Country+"|"+row.Operator+"|"+row.Service]
			row.EstROAS = estROASOrFallback(sumGrossRevenue, hasCohort, row.MO, row.CAC, row.ROAS)
		}
		if cohortROIErr == nil {
			roiMonths, hasROI := cohortROI[row.Country+"|"+row.Operator+"|"+row.Service]
			if hasROI {
				row.ROIMonths = roiMonths
				row.HasROI = true
			}
		}

		results = append(results, row)
	}
	return results, nil
}

// GetCampaignHierarchy returns one row per campaign, each carrying its full
// country/operator/service/adnet path — the FE builds the 5-level
// Country → Operator → Service → Adnet → Campaign accordion by grouping
// these rows client-side (parent levels are pure aggregation of their
// children, same convention as GetRollup's 3-level tree).
//
// Rows come from two sources, tagged by Source:
//   - "summary": summary_campaigns, keyed by url_service_key. CR/CPA/Payout
//     are computed here (landing/postback/saaf/sbaf are all present).
//   - "api": api_pin_reports, keyed by campaign_id — covers API-objective
//     traffic that never lands in summary_campaigns (mirrors GetReport's own
//     api_pin_reports merge). It has no landing/clicked column, so CR is
//     left at 0 (FE renders "—"); Status is unknown for the same reason.
//
// EstROAS is resolved per campaign from campaign_roas_cohorts and reported
// with an explicit HasEstROAS flag — unlike GetRollup/GetCampaign, this
// endpoint does NOT fall back to realized ROAS when no cohort row matches,
// since the caller wants a bare "no data" signal rather than a duplicate
// number.
func (r *BaseModel) GetCampaignHierarchy(date_range, date_before, date_after, client_type, country, service string, allowedAdnets, allowedCompanies []string) ([]entity.HierarchyCampaignRow, error) {
	query := r.DB.Model(&entity.SummaryCampaign{})
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
	default:
		query = query.Where("summary_date >= DATE_TRUNC('month', CURRENT_DATE)")
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
	if len(allowedAdnets) > 0 {
		query = query.Where("adnet IN ?", allowedAdnets)
	}
	if len(allowedCompanies) > 0 {
		query = query.Where("company IN ?", allowedCompanies)
	}

	type summaryAgg struct {
		Country    string  `gorm:"column:country"`
		Operator   string  `gorm:"column:operator"`
		Service    string  `gorm:"column:service"`
		Adnet      string  `gorm:"column:adnet"`
		CampaignID string  `gorm:"column:campaign_id"`
		ClientType string  `gorm:"column:client_type"`
		MO         int     `gorm:"column:mo"`
		Landing    int     `gorm:"column:landing"`
		Postback   int     `gorm:"column:postback"`
		Spend      float64 `gorm:"column:spend"`
		Revenue    float64 `gorm:"column:revenue"`
		ActiveN    int64   `gorm:"column:active_n"`
	}
	var summaryRows []summaryAgg
	err := query.Select(`country, operator, service, adnet,
		url_service_key as campaign_id,
		MAX(client_type) as client_type,
		SUM(mo_received) as mo,
		SUM(landing) as landing,
		SUM(postback) as postback,
		SUM(sbaf) as spend,
		SUM(saaf) as revenue,
		SUM(CASE WHEN status THEN 1 ELSE 0 END) as active_n`).
		Group("country, operator, service, adnet, url_service_key").Order("spend DESC").Scan(&summaryRows).Error
	if err != nil {
		return nil, err
	}

	cohortSums, cohortErr := r.GetCampaignROASCohortSumByCampaign(date_range, date_before, date_after, country, service)
	cohortROI, cohortROIErr := r.GetCampaignROASCohortROIByCampaign(date_range, date_before, date_after, country, service)

	var results []entity.HierarchyCampaignRow
	for _, s := range summaryRows {
		row := entity.HierarchyCampaignRow{
			Country: s.Country, Operator: s.Operator, Service: s.Service, Adnet: s.Adnet,
			CampaignID: s.CampaignID, ClientType: s.ClientType,
			MO: s.MO, Spend: s.Spend, Revenue: s.Revenue, Source: "summary",
		}
		if row.Spend > 0 {
			row.ROAS = row.Revenue / row.Spend * 100
		}
		if row.Revenue > 0 {
			row.RecoveryDays = row.Spend * 30.0 / row.Revenue
		}
		if row.MO > 0 {
			row.CPA = row.Revenue / float64(row.MO)  // CPA = SAAF / MO, same formula as FormulaCPA
			row.Payout = row.Spend / float64(row.MO) // Payout = SBAF / MO — adnet payout rate per MO
		}
		if s.Landing > 0 {
			row.CR = float64(s.MO) / float64(s.Landing) * 100
		}
		row.Status = "paused"
		if s.ActiveN > 0 {
			row.Status = "active"
		}

		if cohortErr == nil {
			cac := 0.0
			if row.MO > 0 {
				cac = row.Spend / float64(row.MO)
			}
			sumGrossRevenue, hasCohort := cohortSums[row.CampaignID]
			if hasCohort && row.MO > 0 && cac > 0 {
				row.EstROAS = (sumGrossRevenue / float64(row.MO)) / cac * 100
				row.HasEstROAS = true
			}
		}
		if cohortROIErr == nil {
			roiMonths, hasROI := cohortROI[row.CampaignID]
			if hasROI {
				row.ROIMonths = roiMonths
				row.HasROI = true
			}
		}

		results = append(results, row)
	}

	// Merge api_pin_reports — API-objective traffic that never lands in
	// summary_campaigns (same reasoning as GetReport's own API merge).
	// Skipped for internal-only scope, since api_pin_reports has no
	// client_type column and is always external-facing traffic.
	if client_type != "internal" {
		apiQ := r.DB.Model(&entity.ApiPinReport{})
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
		default:
			apiQ = apiQ.Where("date_send >= DATE_TRUNC('month', CURRENT_DATE)")
		}
		if country != "" {
			apiQ = apiQ.Where("country = ?", country)
		}
		if service != "" {
			apiQ = apiQ.Where("service = ?", service)
		}
		if len(allowedAdnets) > 0 {
			apiQ = apiQ.Where("adnet IN ?", allowedAdnets)
		}
		if len(allowedCompanies) > 0 {
			apiQ = apiQ.Where("company IN ?", allowedCompanies)
		}

		type apiAgg struct {
			Country    string  `gorm:"column:country"`
			Operator   string  `gorm:"column:operator"`
			Service    string  `gorm:"column:service"`
			Adnet      string  `gorm:"column:adnet"`
			CampaignID string  `gorm:"column:campaign_id"`
			MO         int     `gorm:"column:mo"`
			Spend      float64 `gorm:"column:spend"`
			Revenue    float64 `gorm:"column:revenue"`
		}
		var apiRows []apiAgg
		apiErr := apiQ.Select(`country, operator, service, adnet, campaign_id,
			SUM(total_mo) as mo,
			SUM(sbaf) as spend,
			SUM(saaf) as revenue`).
			Group("country, operator, service, adnet, campaign_id").Scan(&apiRows).Error
		if apiErr == nil {
			for _, a := range apiRows {
				if a.CampaignID == "" {
					continue
				}
				row := entity.HierarchyCampaignRow{
					Country: a.Country, Operator: a.Operator, Service: a.Service, Adnet: a.Adnet,
					CampaignID: a.CampaignID, ClientType: "external",
					MO: a.MO, Spend: a.Spend, Revenue: a.Revenue, Source: "api",
				}
				if row.Spend > 0 {
					row.ROAS = row.Revenue / row.Spend * 100
				}
				if row.Revenue > 0 {
					row.RecoveryDays = row.Spend * 30.0 / row.Revenue
				}
				if row.MO > 0 {
					row.CPA = row.Revenue / float64(row.MO)
					row.Payout = row.Spend / float64(row.MO)
				}
				// api_pin_reports has no landing/clicked column (CR stays
				// 0 → FE renders "—") and no status column (Status stays
				// "" → FE renders "—").
				if cohortErr == nil {
					cac := 0.0
					if row.MO > 0 {
						cac = row.Spend / float64(row.MO)
					}
					sumGrossRevenue, hasCohort := cohortSums[row.CampaignID]
					if hasCohort && row.MO > 0 && cac > 0 {
						row.EstROAS = (sumGrossRevenue / float64(row.MO)) / cac * 100
						row.HasEstROAS = true
					}
				}
				if cohortROIErr == nil {
					roiMonths, hasROI := cohortROI[row.CampaignID]
					if hasROI {
						row.ROIMonths = roiMonths
						row.HasROI = true
					}
				}
				results = append(results, row)
			}
		}
	}

	return results, nil
}

func (r *BaseModel) GetAdnetStats(date_range, date_before, date_after, client_type, country, service string, allowedAdnets, allowedCompanies []string) ([]entity.AdnetStat, error) {
	query := r.DB.Model(&entity.SummaryCampaign{})
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
	default:
		query = query.Where("summary_date >= DATE_TRUNC('month', CURRENT_DATE)")
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
	if len(allowedAdnets) > 0 {
		query = query.Where("adnet IN ?", allowedAdnets)
	}
	if len(allowedCompanies) > 0 {
		query = query.Where("company IN ?", allowedCompanies)
	}
	rows, err := query.Select(`adnet,
		SUM(sbaf) as spend, SUM(saaf) as revenue,
		SUM(mo_received) as mo,
		COUNT(DISTINCT campaign_id) as campaigns`).
		Group("adnet").Order("spend DESC").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []entity.AdnetStat

	cohortSums, cohortErr := r.GetCampaignROASCohortSumByAdnet(date_range, date_before, date_after, country, service)
	cohortROI, cohortROIErr := r.GetCampaignROASCohortROIByAdnet(date_range, date_before, date_after, country, service)

	for rows.Next() {
		var row entity.AdnetStat
		r.DB.ScanRows(rows, &row)
		if row.Spend > 0 {
			row.ROAS = row.Revenue / row.Spend * 100
		}
		if row.Revenue > 0 {
			row.RecoveryDays = row.Spend * 30.0 / row.Revenue
		}
		if row.MO > 0 {
			row.CAC = row.Spend / float64(row.MO)
		}

		row.EstROAS = row.ROAS
		if cohortErr == nil {
			sumGrossRevenue, hasCohort := cohortSums[row.Adnet]
			row.EstROAS = estROASOrFallback(sumGrossRevenue, hasCohort, row.MO, row.CAC, row.ROAS)
		}
		if cohortROIErr == nil {
			roiMonths, hasROI := cohortROI[row.Adnet]
			if hasROI {
				row.ROIMonths = roiMonths
				row.HasROI = true
			}
		}

		results = append(results, row)
	}
	return results, nil
}

func (r *BaseModel) GetHeatmap(date_range, date_before, date_after, country, service string, allowedAdnets, allowedCompanies []string) (entity.HeatmapData, error) {
	query := r.DB.Model(&entity.SummaryCampaign{})
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
	default:
		query = query.Where("summary_date >= DATE_TRUNC('month', CURRENT_DATE)")
	}
	if len(allowedAdnets) > 0 {
		query = query.Where("adnet IN ?", allowedAdnets)
	}
	if len(allowedCompanies) > 0 {
		query = query.Where("company IN ?", allowedCompanies)
	}
	if country != "" {
		query = query.Where("country = ?", country)
	}
	if service != "" {
		query = query.Where("service = ?", service)
	}
	type cellRow struct {
		UrlServiceKey string  `gorm:"column:url_service_key"`
		Adnet         string  `gorm:"column:adnet"`
		Service       string  `gorm:"column:service"`
		ROAS          float64 `gorm:"column:roas"`
		Spend         float64 `gorm:"column:spend"`
		MO            int     `gorm:"column:mo"`
	}
	var raw []cellRow
	err := query.Select(`url_service_key, adnet, service,
		SUM(saaf)/NULLIF(SUM(sbaf),0)*100 as roas,
		SUM(sbaf) as spend,
		SUM(mo_received) as mo`).
		Group("url_service_key, adnet, service").Order("spend DESC").Scan(&raw).Error
	if err != nil {
		return entity.HeatmapData{}, err
	}
	campSpend := make(map[string]float64)
	adnetSpend := make(map[string]float64)
	for _, c := range raw {
		campSpend[c.UrlServiceKey] += c.Spend
		adnetSpend[c.Adnet] += c.Spend
	}
	topCamps := topNKeys(campSpend, 8)
	topAdnets := topNKeys(adnetSpend, 8)

	cohortSums, cohortErr := r.GetCampaignROASCohortSumByHeatmapCell(date_range, date_before, date_after, country, service)

	var cells []entity.HeatmapCell
	for _, c := range raw {
		estROAS := c.ROAS
		if cohortErr == nil {
			cac := 0.0
			if c.MO > 0 {
				cac = c.Spend / float64(c.MO)
			}
			sumGrossRevenue, hasCohort := cohortSums[c.UrlServiceKey+"|"+c.Adnet+"|"+c.Service]
			estROAS = estROASOrFallback(sumGrossRevenue, hasCohort, c.MO, cac, c.ROAS)
		}
		cells = append(cells, entity.HeatmapCell{
			Campaign: c.UrlServiceKey,
			Adnet:    c.Adnet,
			Service:  c.Service,
			ROAS:     c.ROAS,
			EstROAS:  estROAS,
			Spend:    c.Spend,
			MO:       c.MO,
		})
	}
	return entity.HeatmapData{
		Campaigns: topCamps,
		Adnets:    topAdnets,
		Cells:     cells,
	}, nil
}

func (r *BaseModel) GetFilterOptions(country string, allowedAdnets, allowedCompanies []string) entity.FilterOptions {
	var opts entity.FilterOptions

	// Country list — always unfiltered (show all countries)
	q1 := r.DB.Model(&entity.SummaryCampaign{}).
		Where("summary_date >= CURRENT_DATE - INTERVAL '90 DAY'")
	if len(allowedAdnets) > 0 {
		q1 = q1.Where("adnet IN ?", allowedAdnets)
	}
	if len(allowedCompanies) > 0 {
		q1 = q1.Where("company IN ?", allowedCompanies)
	}
	var countries []string
	q1.Distinct("country").Order("country").Pluck("country", &countries)
	opts.Countries = countries

	// Service list — optionally filtered by country
	q2 := r.DB.Model(&entity.SummaryCampaign{}).
		Where("summary_date >= CURRENT_DATE - INTERVAL '90 DAY'")
	if len(allowedAdnets) > 0 {
		q2 = q2.Where("adnet IN ?", allowedAdnets)
	}
	if len(allowedCompanies) > 0 {
		q2 = q2.Where("company IN ?", allowedCompanies)
	}
	if country != "" {
		q2 = q2.Where("country = ?", country)
	}
	var services []string
	q2.Distinct("service").Order("service").Pluck("service", &services)
	opts.Services = services

	return opts
}

func (r *BaseModel) GetCampaignDaily(campaign_id, date_range, date_before, date_after string) ([]entity.CampaignDailyStat, error) {
	query := r.DB.Model(&entity.SummaryCampaign{}).Where("url_service_key = ? OR campaign_id = ?", campaign_id, campaign_id)
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
	default:
		query = query.Where("summary_date >= DATE_TRUNC('month', CURRENT_DATE)")
	}
	rows, err := query.Select(`DATE(summary_date) as date,
		SUM(mo_received) as mo,
		SUM(sbaf) as spend,
		SUM(saaf) as revenue,
		SUM(saaf)/NULLIF(SUM(sbaf),0)*100 as roas`).
		Group("DATE(summary_date)").Order("date ASC").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []entity.CampaignDailyStat
	for rows.Next() {
		var row entity.CampaignDailyStat
		r.DB.ScanRows(rows, &row)
		row.Date = strings.TrimSuffix(row.Date, "T00:00:00Z")
		results = append(results, row)
	}
	return results, nil
}

func (r *BaseModel) GetServiceDaily(country, operator, service, date_range, date_before, date_after string) ([]entity.CampaignDailyStat, error) {
	query := r.DB.Model(&entity.SummaryCampaign{}).
		Where("country = ? AND operator = ? AND service = ?", country, operator, service)
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
	default:
		query = query.Where("summary_date >= DATE_TRUNC('month', CURRENT_DATE)")
	}
	rows, err := query.Select(`DATE(summary_date) as date,
		SUM(mo_received) as mo,
		SUM(sbaf) as spend,
		SUM(saaf) as revenue,
		SUM(saaf)/NULLIF(SUM(sbaf),0)*100 as roas`).
		Group("DATE(summary_date)").Order("date ASC").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []entity.CampaignDailyStat
	for rows.Next() {
		var row entity.CampaignDailyStat
		r.DB.ScanRows(rows, &row)
		row.Date = strings.TrimSuffix(row.Date, "T00:00:00Z")
		results = append(results, row)
	}
	return results, nil
}
