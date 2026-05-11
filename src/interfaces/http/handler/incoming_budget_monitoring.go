package handler

import (
	"fmt"
	"math"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/config"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
	"github.com/ledongthuc/goterators"
)

func (h *IncomingHandler) DisplayBudgetMonitoring(c *fiber.Ctx) error {
	dataIndicators := extractQueryArrayBudgetMonitoring(c, "data-indicators[]")
	if len(dataIndicators) == 0 {
		dataIndicators = append(dataIndicators, "target_daily_budget", "total_spending", "mo_received")
	}

	params := entity.ParamsCampaignSummary{
		DataType:       c.Query("data-type"),
		ReportType:     c.Query("report-type"),
		Country:        c.Query("country"),
		Operator:       c.Query("operator"),
		PartnerName:    c.Query("partner-name"),
		Adnet:          c.Query("adnet"),
		Service:        c.Query("service"),
		CampaignName:   c.Query("campaign_name"),
		UrlServiceKey:  c.Query("url_service_key"),
		DataIndicators: dataIndicators,
		DataBasedOn:    c.Query("data-based-on"),
		DateRange:      c.Query("date-range"),
		DateStart:      c.Query("date-start"),
		DateEnd:        c.Query("date-end"),
		All:            c.Query("all"),
	}

	r := h.GenerateCampaignSummaryforBudget(c, params)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) GenerateCampaignSummaryforBudget(c *fiber.Ctx, params entity.ParamsCampaignSummary) entity.ReturnResponse {
	summaryCampaign, startDate, endDate, err := h.DS.GetSummaryCampaignBudgetMonitoring(params)
	summary := formatBudgetMonitoringData(summaryCampaign, params, startDate, endDate)

	if err == nil {
		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithData{
				Code:    fiber.StatusOK,
				Message: config.OK_DESC,
				Data:    summary,
			},
		}
	} else {
		return entity.ReturnResponse{
			HttpStatus: fiber.StatusNotFound,
			Rsp: entity.GlobalResponse{
				Code:    fiber.StatusNotFound,
				Message: "empty",
			},
		}
	}
}

func formatBudgetMonitoringData(data []entity.CampaignSummaryMonitoring, params entity.ParamsCampaignSummary, startDate time.Time, endDate time.Time) []map[string]interface{} {
	var formattedData []map[string]interface{}

	uniqueData := make(map[string]map[string]interface{})

	if params.DataType == "cr" {
		if params.All == "true" {
			generatedSummary := generateBudget(data, params, startDate, endDate)
			placeHolder := map[string]any{
				"all":           "All Campaign",
				"campaign_id":   "All Campaign",
				"campaign_name": "ALL",
				"operator":      "ALL",
				"service":       "ALL",
				"adnet":         "ALL",
			}
			completeSummary := mergeMapsValue(generatedSummary, placeHolder)
			formattedData = append(formattedData, completeSummary)
		} else {
			groupedAdnet := goterators.Group(data, func(campaign entity.CampaignSummaryMonitoring) string {
				return campaign.UrlServiceKey
			})
			for _, campaignPerAdnet := range groupedAdnet {
				generatedSummary := generateBudget(campaignPerAdnet, params, startDate, endDate)

				placeHolder := map[string]interface{}{
					"level":         "campaign_id",
					"campaign_id":   campaignPerAdnet[0].UrlServiceKey,
					"campaign_name": campaignPerAdnet[0].CampaignName,
					"country":       campaignPerAdnet[0].Country,
					"operator":      campaignPerAdnet[0].Operator,
					"service":       campaignPerAdnet[0].Service,
					"adnet":         campaignPerAdnet[0].Adnet,
					"date":          campaignPerAdnet[0].SummaryDate,
				}
				uniqueKey := fmt.Sprintf("%s_%s_%s_%s_%s_%s", placeHolder["campaign_id"], placeHolder["campaign_name"], placeHolder["country"], placeHolder["operator"], placeHolder["service"], placeHolder["adnet"])
				if existing, exists := uniqueData[uniqueKey]; !exists {
					completeSummary := mergeMapsValue(generatedSummary, placeHolder)
					uniqueData[uniqueKey] = completeSummary
				} else {
					for key, value := range generatedSummary {
						if _, ok := existing[key]; ok {
							existing[key] = existing[key].(float64) + value.(float64)
						} else {
							existing[key] = value
						}
					}
					if campaignPerAdnet[0].SummaryDate.After(existing["date_latest"].(time.Time)) {
						existing["date_latest"] = campaignPerAdnet[0].SummaryDate
					}
				}
			}
		}
	} else {
		if params.All == "true" {
			generatedSummary := generateBudget(data, params, startDate, endDate)
			placeHolder := map[string]interface{}{
				"level":         "campaign_id",
				"country":       "All",
				"campaign_id":   "ALL",
				"campaign_name": "ALL",
				"operator":      "ALL",
				"service":       "ALL",
				"adnet":         "ALL",
			}
			completeSummary := mergeMapsValue(generatedSummary, placeHolder)
			formattedData = append(formattedData, completeSummary)
		} else {

			groupedCountry := goterators.Group(data, func(campaign entity.CampaignSummaryMonitoring) string {
				return campaign.UrlServiceKey
			})

			for _, campaignPerCountry := range groupedCountry {
				generatedCountrySummary := generateBudget(campaignPerCountry, params, startDate, endDate)

				// var children []any

				switch params.ReportType {
				case "campaign_summary":
					// children = groupPartner(campaignPerCountry, params, startDate, endDate)
				case "url_service_summary":
					// children = groupService(campaignPerCountry, params, startDate, endDate)
				case "adnet_summary":
					// children = groupAdnet(campaignPerCountry, params, startDate, endDate)
				default:
					// children = groupOperator(campaignPerCountry, params, startDate, endDate)
				}

				placeHolder := map[string]interface{}{
					"level":         "campaign_id",
					"country":       campaignPerCountry[0].Country,
					"campaign_name": campaignPerCountry[0].CampaignName,
					"date":          campaignPerCountry[0].SummaryDate,
					"operator":      campaignPerCountry[0].Operator,
					"service":       campaignPerCountry[0].Service,
					"adnet":         campaignPerCountry[0].Adnet,
					"campaign_id":   campaignPerCountry[0].UrlServiceKey,
					// "_children":     children,
				}
				completeSummary := mergeMapsValue(generatedCountrySummary, placeHolder)
				formattedData = append(formattedData, completeSummary)
			}

		}
	}

	for _, value := range uniqueData {
		formattedData = append(formattedData, value)
	}

	return formattedData
}

func generateBudget(data []entity.CampaignSummaryMonitoring, params entity.ParamsCampaignSummary, startDate time.Time, endDate time.Time) map[string]interface{} {
	days := map[string]map[string]map[string]interface{}{}
	totals := make(map[string]float64)

	for _, campaign := range data {
		date := campaign.SummaryDate.Format("2006-01-02")
		prevDate := campaign.SummaryDate.AddDate(0, 0, -1).Format("2006-01-02")

		if params.DataType == "monthly_report" {
			date = campaign.SummaryDate.Format("2006-01")
			prevDate = campaign.SummaryDate.AddDate(0, -1, 0).Format("2006-01")
		}

		// Initialize the day's indicator map if it doesn't exist
		if days[date] == nil {
			days[date] = make(map[string]map[string]interface{})
		}

		for _, indicator := range params.DataIndicators {
			indicatorValue := getIndicator(campaign, indicator)

			if days[date][indicator] == nil {
				days[date][indicator] = map[string]interface{}{
					"value":      0.0, // Initialize "value" to 0
					"percentage": 0.0, // Initialize "percentage" to 0
				}
			}

			prevValue := getPreviousVal(days[prevDate], indicator)

			currentValue, ok := days[date][indicator]["value"].(float64)
			if !ok {
				currentValue = 0.0
			}
			newValue := currentValue + indicatorValue

			days[date][indicator]["value"] = newValue
			days[date][indicator]["percentage"] = countPercentageValue(newValue, prevValue)

			totals[indicator] += indicatorValue
		}
	}

	// Final calculations
	tmoEnd := countTmoEndValue(totals, startDate, endDate)

	// Prepare summary data
	summaryData := map[string]interface{}{
		"data_indicators": params.DataIndicators,
		"total":           totals,
		"avg":             countAverageValue(totals, startDate, endDate),
		"t_mo_end":        tmoEnd,
	}

	// Merge with daily breakdowns
	completeSummary := mergeDaysValue(summaryData, days)

	return completeSummary
}

func getIndicator(item entity.CampaignSummaryMonitoring, key string) float64 {

	key = SnakeToCamelBudget(key)
	values := reflect.ValueOf(item)
	keyValues := values.FieldByName(key)

	if !keyValues.IsValid() {
		return 0
	}

	switch keyValues.Kind() {
	case reflect.Float64:
		return keyValues.Float()
	case reflect.Int, reflect.Int64:
		return float64(keyValues.Int())
	default:
		return 0
	}

}

func getPreviousVal(data map[string]map[string]interface{}, key string) float64 {
	value := 0.0
	if val, exists := data[key]; exists {
		value = val["value"].(float64)
	}
	return value
}

func countPercentageValue(now, prev float64) float64 {
	if prev == 0 {
		if now == 0 {
			return 0
		}
		return 100
	}
	if now == 0 {
		return -100
	}
	percentage := ((now - prev) / prev) * 100
	return percentage
}

func countTmoEndValue(totals map[string]float64, startDate time.Time, endDate time.Time) map[string]float64 {

	tmoEnd := map[string]float64{}
	totalDaysRunning := int(math.Ceil(endDate.Sub(startDate).Hours() / 24))
	if totalDaysRunning < 1 {
		totalDaysRunning = 1
	}

	// Calculate total days in the last month
	lastMonthEnd := endDate.AddDate(0, 0, -endDate.Day())
	lastMonthStart := lastMonthEnd.AddDate(0, 0, -lastMonthEnd.Day())
	totalDaysLastMonth := int(math.Ceil(lastMonthEnd.Sub(lastMonthStart).Hours()/24)) + 1
	for key, value := range totals {
		result := (value / float64(totalDaysRunning)) * float64(totalDaysLastMonth)
		tmoEnd[key] = result
	}

	return tmoEnd
}

func countAverageValue(totals map[string]float64, startDate, endDate time.Time) map[string]float64 {
	averages := map[string]float64{}
	totalDaysRunning := int(endDate.Sub(startDate).Hours() / 24)
	if totalDaysRunning < 1 {
		totalDaysRunning = 1
	}

	for key, value := range totals {
		averages[key] = value / float64(totalDaysRunning)
	}

	return averages
}

func extractQueryArrayBudgetMonitoring(c *fiber.Ctx, key string) []string {
	rawQuery, _ := url.QueryUnescape(string(c.Request().URI().QueryString()))
	values := []string{}

	// Adjust logic if the input `rawQuery` needs different parsing
	params := strings.Split(rawQuery, "&")

	for _, param := range params {
		parts := strings.SplitN(param, "=", 2)
		if len(parts) == 2 && parts[0] == key {
			values = append(values, parts[1])
		}
	}
	return values
}

func SnakeToCamelBudget(snake string) string {
	words := strings.Split(snake, "_")
	for i, word := range words {
		words[i] = strings.ToLower(word)
		if len(word) > 0 {
			words[i] = strings.ToUpper(words[i][:1]) + words[i][1:]
		}
	}
	return strings.Join(words, "")
}

func mergeDaysValue(summaryData map[string]interface{}, days map[string]map[string]map[string]interface{}) map[string]interface{} {
	for key, value := range days {
		summaryData[key] = value
	}
	return summaryData
}

func mergeMapsValue(map1, map2 map[string]interface{}) map[string]interface{} {
	mergedMap := make(map[string]interface{})
	for key, value := range map1 {
		mergedMap[key] = value
	}
	for key, value := range map2 {
		mergedMap[key] = value
	}
	return mergedMap
}

// func (h *IncomingHandler) DisplayBudgetMonitoringChart(c *fiber.Ctx) error {
// 	c.Set("Content-Type", "application/x-www-form-urlencoded")
// 	c.Accepts("application/x-www-form-urlencoded")
// 	c.AcceptsCharsets("utf-8", "iso-1959-1")

// 	currentDate := time.Now()
// 	currentMonth := currentDate.Format("2006-01")
// 	currentYear, _, _ := currentDate.Date()

// 	var summaryCampaigns []entity.SummaryCampaign
// 	_, err := h.DS.DB.Raw("SELECT country, SUM(budget_usage) as used_budget, SUM(target_daily_budget_changes) as target_daily_budget, SUM(target_daily_budget_changes) as target_daily_budget_changes FROM `summary_campaigns` WHERE `summary_date` >= ? AND `summary_date` < ? GROUP BY country", currentYear+"-"+currentMonth+"-01", currentDate.Format("2006-01-02")).Scan(&summaryCampaigns)
// 	if err != nil {
// 		h.Logs.Error(err.Error())
// 		return c.Status(fiber.StatusInternalServerError).JSON(entity.GlobalResponse{
// 			Code:    fiber.StatusInternalServerError,
// 			Message: "Failed execute query",
// 		})
// 	}

// 	underBudget := []entity.SummaryCampaign{}
// 	overBudget := []entity.SummaryCampaign{}
// 	for _, summaryCampaign := range summaryCampaigns {
// 		summaryCampaign.BudgetUnused = summaryCampaign.UsedBudget - summaryCampaign.TargetDailyBudgetChanges
// 		if summaryCampaign.UsedBudget < summaryCampaign.TargetDailyBudget {
// 			underBudget = append(underBudget, summaryCampaign)
// 		} else {
// 			overBudget = append(overBudget, summaryCampaign)
// 		}
// 	}

// 	return c.Status(fiber.StatusOK).JSON(entity.DisplayBudgetMonitoringChart{
// 		BudgetUsed:           budgetUsed,
// 		BudgetUnused:         budgetUnused,
// 		OverBudgetCountries:  overBudget,
// 		UnderBudgetCountries: underBudget,
// 	})
// }
