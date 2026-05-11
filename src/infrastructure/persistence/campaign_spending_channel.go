package persistence

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
)

var channelGroupMapModel = map[string]string{
	// TikTok
	"tiktok":                "Mainstream TikTok Ads",
	"mainstream_tiktok":     "Mainstream TikTok Ads",
	"mainstream_tiktok_ads": "Mainstream TikTok Ads",
	"tiktok ads":            "Mainstream TikTok Ads",
	// Google
	"google":                "Mainstream Google Ads",
	"mainstream_google":     "Mainstream Google Ads",
	"mainstream_google_ads": "Mainstream Google Ads",
	"google ads":            "Mainstream Google Ads",
	"google traffic":        "Mainstream Google Ads",
	// Meta
	"meta":                "Mainstream Meta Ads",
	"mainstream_meta":     "Mainstream Meta Ads",
	"mainstream_meta_ads": "Mainstream Meta Ads",
	"facebook":            "Mainstream Meta Ads",
	"fb":                  "Mainstream Meta Ads",
	"fbmeta":              "Mainstream Meta Ads",
	// Snack Video
	"snack video":          "Mainstream Snack Video Ads",
	"snack_video":          "Mainstream Snack Video Ads",
	"snackvideo":           "Mainstream Snack Video Ads",
	// Others
	"cpa":           "CPA",
	"dsp":           "DSP",
	"sms":           "SMS",
	"telco":         "Telco Channel",
	"telco_channel": "Telco Channel",
	"telco channel": "Telco Channel",
	"s2s":           "S2S",
	"api":           "API",
}

func reverseChannelLookupSQL(canonical string) []string {
	var keys []string
	for k, v := range channelGroupMapModel {
		if strings.EqualFold(v, canonical) || strings.EqualFold(k, canonical) {
			keys = append(keys, k)
		}
	}
	return keys
}

func resolveDateRange(params entity.ParamsCampaignSpendingChannel, today time.Time) (time.Time, time.Time) {
	var start, end time.Time
	switch strings.ToUpper(params.DateRange) {
	case "TODAY":
		start, end = today, today
	case "YESTERDAY":
		start = today.AddDate(0, 0, -1)
		end = today.AddDate(0, 0, -1)
	case "LAST_7_DAYS", "LAST_7_DAY":
		start = today.AddDate(0, 0, -6)
		end = today
	case "LAST_30_DAYS", "LAST_30_DAY":
		start = today.AddDate(0, -1, 0)
		end = today
	case "THIS_MONTH":
		start = time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
		end = today
	case "LAST_MONTH":
		lm := today.AddDate(0, -1, 0)
		start = time.Date(lm.Year(), lm.Month(), 1, 0, 0, 0, 0, today.Location())
		end = time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location()).AddDate(0, 0, -1)
	case "CUSTOM_RANGE":
		parts := strings.Split(params.DateCustomRange, " to ")
		var errS, errE error
		start, errS = time.Parse("2006-01-02", strings.TrimSpace(parts[0]))
		if len(parts) > 1 {
			end, errE = time.Parse("2006-01-02", strings.TrimSpace(parts[1]))
		}
		if errS != nil {
			start = today
		}
		if errE != nil {
			end = today
		}
	default:
		start = time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
		end = today
	}
	return start, end
}

func (r *BaseModel) GetSpendingChannelMonitoring(
	params entity.ParamsCampaignSpendingChannel,
) ([]entity.CampaignSpendingChannelMonitoring, time.Time, time.Time, error) {

	today := time.Now().Truncate(24 * time.Hour)
	startDate, endDate := resolveDateRange(params, today)
	if endDate.After(today) {
		endDate = today
	}

	dateExpr := buildDateExpr(params.DataType, "summary_date")
	groupBy := buildGroupBy(params.DataType, "summary_date")

	scSQL := fmt.Sprintf(`
		SELECT
			country,
			operator,
			partner,
			service,
			adnet,
			COALESCE(channel, '') AS channel,
			url_service_key,
			campaign_name,
			campaign_id,
			%s AS summary_date,
			SUM(sbaf) AS sbaf
		FROM summary_campaigns
		WHERE deleted_at IS NULL
		  AND summary_date BETWEEN '%s' AND '%s'
		  %s
		GROUP BY %s, country, operator, partner, service, adnet, channel,
		         url_service_key, campaign_name, campaign_id
	`,
		dateExpr,
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"),
		buildSCWhereClause(params),
		groupBy,
	)

	includeAPI := params.ChannelCampaign == "" ||
		strings.EqualFold(params.ChannelCampaign, "api") ||
		strings.EqualFold(params.ChannelCampaign, "API")

	var finalSQL string

	if includeAPI {
		apiDateExpr := buildDateExpr(params.DataType, "date_send")
		apiGroupBy := buildGroupBy(params.DataType, "date_send")

		apiSQL := fmt.Sprintf(`
			SELECT
				country,
				operator,
				''           AS partner,
				service,
				'api'        AS adnet,
				'API'        AS channel,
				''           AS url_service_key,
				COALESCE(campaign_id::varchar, '') AS campaign_name,
				COALESCE(campaign_id::varchar, '') AS campaign_id,
				%s           AS summary_date,
				SUM(sbaf)    AS sbaf
			FROM api_pin_reports
			WHERE deleted_at IS NULL
			  AND date_send BETWEEN '%s' AND '%s'
			  %s
			GROUP BY %s, country, operator, service, campaign_id
		`,
			apiDateExpr,
			startDate.Format("2006-01-02"),
			endDate.Format("2006-01-02"),
			buildAPIWhereClause(params),
			apiGroupBy,
		)

		finalSQL = fmt.Sprintf("(%s) UNION ALL (%s)", scSQL, apiSQL)
	} else {
		finalSQL = scSQL
	}

	if params.ChannelCampaign != "" && !strings.EqualFold(params.ChannelCampaign, "api") {
		channelKeys := reverseChannelLookupSQL(params.ChannelCampaign)
		if len(channelKeys) > 0 {
			quoted := make([]string, len(channelKeys))
			for i, k := range channelKeys {
				quoted[i] = fmt.Sprintf("'%s'", esc(k))
			}
			inList := strings.Join(quoted, ", ")
			finalSQL = fmt.Sprintf(
				`SELECT * FROM (%s) _ch
				 WHERE LOWER(COALESCE(channel,'')) IN (%s)`,
				finalSQL, inList,
			)
		}
	}

	sqlRows, err := r.DB.Raw(finalSQL).Rows()
	if err != nil {
		log.Printf("[SpendingChannel] DB.Raw error: %v", err)
		return nil, startDate, endDate, err
	}
	defer sqlRows.Close()

	cols, colErr := sqlRows.Columns()
	if colErr != nil {
		log.Printf("[SpendingChannel] Columns() error: %v", colErr)
		return nil, startDate, endDate, colErr
	}
	log.Printf("[SpendingChannel] Query columns (%d): %v", len(cols), cols)

	var result []entity.CampaignSpendingChannelMonitoring
	rowCount := 0

	for sqlRows.Next() {
		rowCount++
		var (
			country       string
			operator      string
			partner       string
			service       string
			adnet         string
			channel       string
			urlServiceKey string
			campaignName  string
			campaignId    string
			summaryDate   time.Time
			sbaf          float64
		)
		if scanErr := sqlRows.Scan(
			&country, &operator, &partner, &service,
			&adnet, &channel, &urlServiceKey,
			&campaignName, &campaignId,
			&summaryDate, &sbaf,
		); scanErr != nil {
			log.Printf("[SpendingChannel] Scan error on row %d: %v", rowCount, scanErr)
			return nil, startDate, endDate, scanErr
		}
		result = append(result, entity.CampaignSpendingChannelMonitoring{
			Country:       country,
			Operator:      operator,
			Partner:       partner,
			Service:       service,
			Adnet:         adnet,
			Channel:       channel,
			UrlServiceKey: urlServiceKey,
			CampaignName:  campaignName,
			CampaignId:    campaignId,
			SummaryDate:   summaryDate,
			SBAF:          sbaf,
		})
	}

	if rowErr := sqlRows.Err(); rowErr != nil {
		log.Printf("[SpendingChannel] rows.Err(): %v", rowErr)
		return nil, startDate, endDate, rowErr
	}

	log.Printf("[SpendingChannel] Total rows scanned: %d, startDate: %s, endDate: %s",
		rowCount, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	return result, startDate, endDate, nil
}

// ─── SQL Helpers ──────────────────────────────────────────────────────────────

func buildDateExpr(dataType, col string) string {
	switch dataType {
	case "monthly_report":
		return fmt.Sprintf("DATE_TRUNC('month', %s)", col)
	case "weekly_report":
		return fmt.Sprintf("DATE_TRUNC('week', %s)", col)
	default:
		return col
	}
}

func buildGroupBy(dataType, col string) string {
	return buildDateExpr(dataType, col)
}

func buildSCWhereClause(params entity.ParamsCampaignSpendingChannel) string {
	var clauses []string
	if params.Country != "" {
		clauses = append(clauses, fmt.Sprintf("AND country = '%s'", esc(params.Country)))
	}
	if params.Operator != "" {
		clauses = append(clauses, fmt.Sprintf("AND operator = '%s'", esc(params.Operator)))
	}
	if params.PartnerName != "" {
		clauses = append(clauses, fmt.Sprintf("AND partner = '%s'", esc(params.PartnerName)))
	}
	if params.Service != "" {
		clauses = append(clauses, fmt.Sprintf("AND service = '%s'", esc(params.Service)))
	}
	return strings.Join(clauses, "\n  ")
}

func buildAPIWhereClause(params entity.ParamsCampaignSpendingChannel) string {
	var clauses []string
	if params.Country != "" {
		clauses = append(clauses, fmt.Sprintf("AND country = '%s'", esc(params.Country)))
	}
	if params.Operator != "" {
		clauses = append(clauses, fmt.Sprintf("AND operator = '%s'", esc(params.Operator)))
	}
	if params.Service != "" {
		clauses = append(clauses, fmt.Sprintf("AND service = '%s'", esc(params.Service)))
	}
	return strings.Join(clauses, "\n  ")
}

func esc(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}