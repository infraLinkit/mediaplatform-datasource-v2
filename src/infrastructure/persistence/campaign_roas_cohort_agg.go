package persistence

import (
	"fmt"
	"strings"

	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
)

// GetCampaignROASCohortAgg computes est_ltv = SUM(estimated_gross_revenue_full)
// / totalMO and roasAvg = est_ltv / cac — combining Mart's revenue projection
// with our own volume (totalMO) and cost (cac) data, rather than trusting
// Mart's own precomputed estimated_roas ratio. roasAvg is a plain ratio, the
// caller scales it to a percentage. roiMonthsAvg (roi_months_payback,
// averaged as-is) is already in months — NOT a ratio, do not scale. ok is
// false when no cohort rows match (or totalMO/cac are non-positive, making
// the formula undefined), so callers can fall back to the realized
// placeholders instead of showing a misleading zero.
func (r *BaseModel) GetCampaignROASCohortAgg(dateList []string, country, service string, totalMO int, cac float64) (roasAvg, roiMonthsAvg float64, ok bool) {

	query := r.DB.Model(&entity.CampaignROASCohort{}).Where("summary_date::date IN ?", dateList)

	if country != "" {
		query = query.Where("country = ?", country)
	}
	if service != "" {
		query = query.Where("service = ?", service)
	}

	// Mart's cohort table covers every campaign it tracks for a country,
	// which is a much broader set than what actually ran (has mo/spend) in
	// summary_campaigns for this filter — summing gross revenue over all of
	// them while totalMO/cac only reflect the campaigns that actually ran
	// would silently inflate est_roas. Restrict the sum to campaigns that
	// are actually present in summary_campaigns for the same window/filter,
	// so the numerator and the totalMO/cac denominator cover the same set.
	activeCampaigns := r.DB.Model(&entity.SummaryCampaign{}).
		Select("DISTINCT url_service_key").
		Where("summary_date::date IN ?", dateList)
	if country != "" {
		activeCampaigns = activeCampaigns.Where("country = ?", country)
	}
	if service != "" {
		activeCampaigns = activeCampaigns.Where("service = ?", service)
	}
	query = query.Where("url_service_key IN (?)", activeCampaigns)

	type aggRow struct {
		SumGrossRevenue float64 `gorm:"column:sum_gross_revenue"`
		AvgROIMonths    float64 `gorm:"column:avg_roi_months"`
		MatchCount      int64   `gorm:"column:match_count"`
	}
	var agg aggRow

	err := query.Select("SUM(estimated_gross_revenue_full) as sum_gross_revenue, AVG(roi_months_payback) as avg_roi_months, COUNT(*) as match_count").Scan(&agg).Error
	if err != nil {
		r.Logs.Error(fmt.Sprintf("GetCampaignROASCohortAgg query error: %v", err))
		return 0, 0, false
	}
	if agg.MatchCount == 0 || totalMO <= 0 || cac <= 0 {
		return 0, 0, false
	}

	estLTV := agg.SumGrossRevenue / float64(totalMO)
	roasAvg = estLTV / cac

	return roasAvg, agg.AvgROIMonths, true
}

// operatorKeywordOverrides maps a keyword found in the Mart API's raw
// operator string (campaign_roas_cohorts.operator, stored as-is from the
// API) to the canonical operator name used in summary_campaigns.operator.
// Mart's operator strings are a free-form "<company> <carrier>" label (e.g.
// "LINKIT TSEL DIRECT", "INDOSAT WAKI") that never equals our own operator
// value outright — only a specific, known abbreviation inside it does.
// Applied at query time (not storage time) so the raw Mart value stays
// intact in the table. Extend as more mismatches are confirmed; unmapped
// operators are left unresolved (raw) rather than guessed at.
var operatorKeywordOverrides = map[string]string{
	"TSEL":    "TELKOMSEL",
	"INDOSAT": "INDOSAT",
}

// resolveOperator looks for a known keyword (case-insensitive) inside raw
// and returns the matching canonical operator name, or raw unchanged when
// nothing matches.
func resolveOperator(raw string) string {
	upper := strings.ToUpper(raw)
	for keyword, canonical := range operatorKeywordOverrides {
		if strings.Contains(upper, keyword) {
			return canonical
		}
	}
	return raw
}

// cohortGrossRevenueByGroup runs SUM(estimated_gross_revenue_full) from
// campaign_roas_cohorts grouped by groupByCols, filtered the same way
// GetCampaign/GetRollup/GetAdnetStats/GetHeatmap filter summary_campaigns
// (date_range BETWEEN-style, plus optional country/service). selectCols must
// list groupByCols followed by "SUM(estimated_gross_revenue_full) as
// sum_gross_revenue", and dest must be a pointer to a slice of a struct
// whose gorm column tags match selectCols exactly — callers build their own
// lookup map from the result.
func (r *BaseModel) cohortGrossRevenueByGroup(date_range, date_before, date_after, country, service, selectCols, groupByCols string, dest interface{}) error {
	query := r.DB.Model(&entity.CampaignROASCohort{})
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
		// Matches GetRollup/GetAdnetStats/GetHeatmap's own default branch —
		// without this, an empty/unrecognized date_range left the cohort
		// side unfiltered (full table history) while those realized queries
		// scoped to this month, silently inflating EstROAS.
		query = query.Where("summary_date >= DATE_TRUNC('month', CURRENT_DATE)")
	}
	if country != "" {
		query = query.Where("country = ?", country)
	}
	if service != "" {
		query = query.Where("service = ?", service)
	}

	// Same reasoning as GetCampaignROASCohortAgg: restrict to campaigns that
	// actually ran (have a row in summary_campaigns) for this window/filter,
	// so a group's cohort sum can't be inflated by campaigns Mart tracks but
	// that never ran here — those would otherwise silently pull in revenue
	// with no matching mo/spend to divide it by.
	activeCampaigns := r.DB.Model(&entity.SummaryCampaign{}).Select("DISTINCT url_service_key")
	switch date_range {
	case "TODAY":
		activeCampaigns = activeCampaigns.Where("summary_date = CURRENT_DATE")
	case "YESTERDAY":
		activeCampaigns = activeCampaigns.Where("summary_date = CURRENT_DATE - INTERVAL '1 DAY'")
	case "LAST7DAY":
		activeCampaigns = activeCampaigns.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '7 DAY' AND CURRENT_DATE")
	case "LAST30DAY":
		activeCampaigns = activeCampaigns.Where("summary_date BETWEEN CURRENT_DATE - INTERVAL '30 DAY' AND CURRENT_DATE")
	case "THISMONTH":
		activeCampaigns = activeCampaigns.Where("summary_date >= DATE_TRUNC('month', CURRENT_DATE)")
	case "LASTMONTH":
		activeCampaigns = activeCampaigns.Where("summary_date BETWEEN DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 MONTH') AND DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '1 DAY'")
	case "CUSTOMRANGE":
		activeCampaigns = activeCampaigns.Where("summary_date BETWEEN ? AND ?", date_before, date_after)
	default:
		activeCampaigns = activeCampaigns.Where("summary_date >= DATE_TRUNC('month', CURRENT_DATE)")
	}
	if country != "" {
		activeCampaigns = activeCampaigns.Where("country = ?", country)
	}
	if service != "" {
		activeCampaigns = activeCampaigns.Where("service = ?", service)
	}
	query = query.Where("url_service_key IN (?)", activeCampaigns)

	err := query.Select(selectCols).Group(groupByCols).Scan(dest).Error
	if err != nil {
		r.Logs.Error(fmt.Sprintf("cohortGrossRevenueByGroup(%s) query error: %v", groupByCols, err))
	}
	return err
}

// estROASOrFallback applies the est_ltv/cac formula for one group's summed
// cohort revenue against that same group's own mo/cac (already computed by
// the realized-ROAS query), falling back to realizedROAS when there's no
// cohort match or mo/cac are non-positive — same convention as the KPI
// strip's EstROAS.
func estROASOrFallback(sumGrossRevenue float64, hasCohort bool, mo int, cac, realizedROAS float64) float64 {
	if !hasCohort || mo <= 0 || cac <= 0 {
		// 0 is a sentinel here, not a real estimate — no cohort data means
		// no cohort-based number to show. Callers (blade) render "—" for 0
		// rather than duplicating the realized ROAS as if it were a
		// separate estimate.
		return 0
	}
	estLTV := sumGrossRevenue / float64(mo)
	return estLTV / cac * 100
}

// GetCampaignROASCohortSumByCampaign sums estimated_gross_revenue_full per
// url_service_key, for GetCampaign's per-campaign EstROAS.
func (r *BaseModel) GetCampaignROASCohortSumByCampaign(date_range, date_before, date_after, country, service string) (map[string]float64, error) {
	type row struct {
		URLServiceKey   string  `gorm:"column:url_service_key"`
		SumGrossRevenue float64 `gorm:"column:sum_gross_revenue"`
	}
	var rows []row
	if err := r.cohortGrossRevenueByGroup(date_range, date_before, date_after, country, service,
		"url_service_key, SUM(estimated_gross_revenue_full) as sum_gross_revenue",
		"url_service_key", &rows); err != nil {
		return nil, err
	}
	out := make(map[string]float64, len(rows))
	for _, rr := range rows {
		out[rr.URLServiceKey] = rr.SumGrossRevenue
	}
	return out, nil
}

// GetCampaignROASCohortROIByCampaign averages roi_months_payback per
// url_service_key — same shape as GetCampaignROASCohortSumByCampaign, but
// for the "months to payback" ROI figure instead of gross revenue. Unlike
// EstROAS, ROI months is already in the right unit (months) as-is — callers
// use the raw average, no est_ltv/cac math needed.
func (r *BaseModel) GetCampaignROASCohortROIByCampaign(date_range, date_before, date_after, country, service string) (map[string]float64, error) {
	type row struct {
		URLServiceKey string  `gorm:"column:url_service_key"`
		AvgROIMonths  float64 `gorm:"column:avg_roi_months"`
	}
	var rows []row
	if err := r.cohortGrossRevenueByGroup(date_range, date_before, date_after, country, service,
		"url_service_key, AVG(roi_months_payback) as avg_roi_months",
		"url_service_key", &rows); err != nil {
		return nil, err
	}
	out := make(map[string]float64, len(rows))
	for _, rr := range rows {
		out[rr.URLServiceKey] = rr.AvgROIMonths
	}
	return out, nil
}

// GetCampaignROASCohortROIByRollup averages roi_months_payback per
// (country, operator, service). Key format: "<country>|<operator>|<service>".
func (r *BaseModel) GetCampaignROASCohortROIByRollup(date_range, date_before, date_after, country, service string) (map[string]float64, error) {
	type row struct {
		Country      string  `gorm:"column:country"`
		Operator     string  `gorm:"column:operator"`
		Service      string  `gorm:"column:service"`
		AvgROIMonths float64 `gorm:"column:avg_roi_months"`
	}
	var rows []row
	if err := r.cohortGrossRevenueByGroup(date_range, date_before, date_after, country, service,
		"country, operator, service, AVG(roi_months_payback) as avg_roi_months",
		"country, operator, service", &rows); err != nil {
		return nil, err
	}
	out := make(map[string]float64, len(rows))
	for _, rr := range rows {
		out[rr.Country+"|"+resolveOperator(rr.Operator)+"|"+rr.Service] = rr.AvgROIMonths
	}
	return out, nil
}

// GetCampaignROASCohortROIByAdnet averages roi_months_payback per adnet.
func (r *BaseModel) GetCampaignROASCohortROIByAdnet(date_range, date_before, date_after, country, service string) (map[string]float64, error) {
	type row struct {
		Adnet        string  `gorm:"column:adnet"`
		AvgROIMonths float64 `gorm:"column:avg_roi_months"`
	}
	var rows []row
	if err := r.cohortGrossRevenueByGroup(date_range, date_before, date_after, country, service,
		"adnet, AVG(roi_months_payback) as avg_roi_months",
		"adnet", &rows); err != nil {
		return nil, err
	}
	out := make(map[string]float64, len(rows))
	for _, rr := range rows {
		out[rr.Adnet] = rr.AvgROIMonths
	}
	return out, nil
}

// GetCampaignROASCohortSumByRollup sums estimated_gross_revenue_full per
// (country, operator, service), for GetRollup's per-row EstROAS. Key format:
// "<country>|<operator>|<service>".
func (r *BaseModel) GetCampaignROASCohortSumByRollup(date_range, date_before, date_after, country, service string) (map[string]float64, error) {
	type row struct {
		Country         string  `gorm:"column:country"`
		Operator        string  `gorm:"column:operator"`
		Service         string  `gorm:"column:service"`
		SumGrossRevenue float64 `gorm:"column:sum_gross_revenue"`
	}
	var rows []row
	if err := r.cohortGrossRevenueByGroup(date_range, date_before, date_after, country, service,
		"country, operator, service, SUM(estimated_gross_revenue_full) as sum_gross_revenue",
		"country, operator, service", &rows); err != nil {
		return nil, err
	}
	out := make(map[string]float64, len(rows))
	for _, rr := range rows {
		out[rr.Country+"|"+resolveOperator(rr.Operator)+"|"+rr.Service] = rr.SumGrossRevenue
	}
	return out, nil
}

// GetCampaignROASCohortSumByAdnet sums estimated_gross_revenue_full per
// adnet, for GetAdnetStats' per-row EstROAS.
func (r *BaseModel) GetCampaignROASCohortSumByAdnet(date_range, date_before, date_after, country, service string) (map[string]float64, error) {
	type row struct {
		Adnet           string  `gorm:"column:adnet"`
		SumGrossRevenue float64 `gorm:"column:sum_gross_revenue"`
	}
	var rows []row
	if err := r.cohortGrossRevenueByGroup(date_range, date_before, date_after, country, service,
		"adnet, SUM(estimated_gross_revenue_full) as sum_gross_revenue",
		"adnet", &rows); err != nil {
		return nil, err
	}
	out := make(map[string]float64, len(rows))
	for _, rr := range rows {
		out[rr.Adnet] = rr.SumGrossRevenue
	}
	return out, nil
}

// GetCampaignROASCohortSumByHeatmapCell sums estimated_gross_revenue_full
// per (url_service_key, adnet, service), for GetHeatmap's per-cell EstROAS.
// Key format: "<url_service_key>|<adnet>|<service>".
func (r *BaseModel) GetCampaignROASCohortSumByHeatmapCell(date_range, date_before, date_after, country, service string) (map[string]float64, error) {
	type row struct {
		URLServiceKey   string  `gorm:"column:url_service_key"`
		Adnet           string  `gorm:"column:adnet"`
		Service         string  `gorm:"column:service"`
		SumGrossRevenue float64 `gorm:"column:sum_gross_revenue"`
	}
	var rows []row
	if err := r.cohortGrossRevenueByGroup(date_range, date_before, date_after, country, service,
		"url_service_key, adnet, service, SUM(estimated_gross_revenue_full) as sum_gross_revenue",
		"url_service_key, adnet, service", &rows); err != nil {
		return nil, err
	}
	out := make(map[string]float64, len(rows))
	for _, rr := range rows {
		out[rr.URLServiceKey+"|"+rr.Adnet+"|"+rr.Service] = rr.SumGrossRevenue
	}
	return out, nil
}
