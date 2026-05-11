package handler

import (
	"fmt"
	"math"
	"net/url"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/config"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
	"github.com/ledongthuc/goterators"
)

func (h *IncomingHandler) DisplayTrafficReport(c *fiber.Ctx) error {
	dataIndicators := extractQueryArrayTrafficReport(c, "data-indicators[]")
	if len(dataIndicators) == 0 {
		dataIndicators = append(dataIndicators, "landing", "mo_received", "cr_mo", "first_push")
	}

	dataType := c.Query("data-type")
	if dataType == "" {
		dataType = "daily_report"
	}

	params := entity.TrafficReportParams{
		DataType:             dataType,
		ReportType:           c.Query("report-type"),
		Country:              c.Query("country"),
		Operator:             c.Query("operator"),
		PartnerName:          c.Query("partner-name"),
		Service:              c.Query("service"),
		Adnet:                c.Query("adnet"),
		TypeData:             c.Query("type-data"),
		CampaignId:           c.Query("campaign_id"),
		CampaignName:         c.Query("campaign"),
		URLServiceKey:        c.Query("url_service_key"),
		DataIndicators:       dataIndicators,
		DataBasedOn:          c.Query("data-based-on"),
		DataBasedOnIndicator: c.Query("data-based-on-indicator"),
		DateRange:            c.Query("date-range"),
		DateStart:            c.Query("date-start"),
		DateEnd:              c.Query("date-end"),
		DateCustomRange:      c.Query("date-custom-range"),
		All:                  c.Query("all"),
	}

	if params.DataBasedOn == "" {
        params.DataBasedOn = "highest"
    }
    if params.DataBasedOnIndicator == "" {
        params.DataBasedOnIndicator = "landing"
    }

	r := h.GenerateTrafficReport(c, params)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) GenerateTrafficReport(c *fiber.Ctx, params entity.TrafficReportParams) entity.ReturnResponse {
	var (
		summaryCampaign []entity.SummaryTraffic
		startDate       time.Time
		endDate         time.Time
		err             error
	)

	var chartData interface{} // interface{} agar bisa assign berbagai tipe

	if params.DataType == "hourly_report" {
		summaryCampaign, startDate, endDate, err = h.DS.GetTrafficReportHourly(params)
		chartData = generateHourlyChartTrafficReport(summaryCampaign, params, startDate)
	} else {
		summaryCampaign, startDate, endDate, err = h.DS.GetTrafficReport(params)
	}

	if err != nil {
		return entity.ReturnResponse{
			HttpStatus: fiber.StatusNotFound,
			Rsp: entity.GlobalResponse{
				Code:    fiber.StatusNotFound,
				Message: "empty",
			},
		}
	}

	summary := formatSummaryDataValTrafficReport(summaryCampaign, params, startDate, endDate)
	sortedSummary := sortDataTrafficReport(summary, params.DataBasedOn, params.DataBasedOnIndicator)

	return entity.ReturnResponse{
		HttpStatus: fiber.StatusOK,
		Rsp: entity.GlobalResponseWithData{
			Code:    fiber.StatusOK,
			Message: config.OK_DESC,
			Data: map[string]interface{}{
				"summary":    sortedSummary,
				"chart_data": chartData,
			},
		},
	}
}

func formatSummaryDataValTrafficReport(data []entity.SummaryTraffic, params entity.TrafficReportParams, startDate, endDate time.Time) []map[string]interface{} {
	var formattedData []map[string]interface{}

	if params.DataType == "hourly_report" {
		var filtered []entity.SummaryTraffic
		for _, item := range data {
			if item.URLServiceKey == params.URLServiceKey && sameDay(item.SummaryDateHour, startDate) {
				filtered = append(filtered, item)
			}
		}

		if len(filtered) > 0 {
			generatedSummary := generateHourlyTrafficReport(filtered, params, startDate)

			placeHolder := map[string]interface{}{
				"level":            "detail",
				"campaign_id":      filtered[0].URLServiceKey,
				"campaign_name":    filtered[0].CampaignName,
				"country":          filtered[0].Country,
				"operator":         filtered[0].Operator,
				"service":          filtered[0].Service,
				"adnet":            filtered[0].Adnet,
				"url_warp_landing": filtered[0].URLWarpLanding,
				"date":             startDate.Format("2006-01-02"),
			}
			completeSummary := mergeMapsTrafficReport(generatedSummary, placeHolder)
			formattedData = append(formattedData, completeSummary)
		}
	} else if params.DataType == "daily_report" {
		if params.All == "true" {
			generatedSummary := generateSummaryTrafficReport(data, params, startDate, endDate)
			placeHolder := map[string]interface{}{
				"all": "All Campaign",
			}
			completeSummary := mergeMapsTrafficReport(generatedSummary, placeHolder)
			formattedData = append(formattedData, completeSummary)
		} else {
			groupedAdnet := goterators.Group(data, func(campaign entity.SummaryTraffic) string {
				return campaign.URLServiceKey
			})
			for _, campaignPerAdnet := range groupedAdnet {
				generatedSummary := generateSummaryTrafficReport(campaignPerAdnet, params, startDate, endDate)
				placeHolder := map[string]interface{}{
					"level":            "country",
					"campaign_id":      campaignPerAdnet[0].URLServiceKey,
					"campaign_name":    campaignPerAdnet[0].CampaignName,
					"country":          campaignPerAdnet[0].Country,
					"operator":         campaignPerAdnet[0].Operator,
					"service":          campaignPerAdnet[0].Service,
					"adnet":            campaignPerAdnet[0].Adnet,
					"url_warp_landing": campaignPerAdnet[0].URLWarpLanding,
					"date":             campaignPerAdnet[0].SummaryDateHour,
				}
				completeSummary := mergeMapsTrafficReport(generatedSummary, placeHolder)
				formattedData = append(formattedData, completeSummary)
			}
		}
	} else if params.DataType == "monthly_report" {
		grouped := goterators.Group(data, func(campaign entity.SummaryTraffic) string {
			return campaign.SummaryDateHour.Format("2006-01")
		})
		for _, campaignPerMonth := range grouped {
			generatedSummary := generateSummaryTrafficReport(campaignPerMonth, params, startDate, endDate)
			placeHolder := map[string]interface{}{
				"level":            "month",
				"campaign_id":      campaignPerMonth[0].URLServiceKey,
				"campaign_name":    campaignPerMonth[0].CampaignName,
				"country":          campaignPerMonth[0].Country,
				"operator":         campaignPerMonth[0].Operator,
				"service":          campaignPerMonth[0].Service,
				"adnet":            campaignPerMonth[0].Adnet,
				"url_warp_landing": campaignPerMonth[0].URLWarpLanding,
				"date":             campaignPerMonth[0].SummaryDateHour.Format("2006-01"),
			}
			completeSummary := mergeMapsTrafficReport(generatedSummary, placeHolder)
			formattedData = append(formattedData, completeSummary)
		}
	}
	

	return formattedData
}

func generateSummaryTrafficReport(data []entity.SummaryTraffic, params entity.TrafficReportParams, startDate time.Time, endDate time.Time) map[string]interface{} {
	days := map[string]map[string]map[string]interface{}{}
	totals := make(map[string]float64)
	dayLanding := make(map[string]float64)
	dayMoReceived := make(map[string]float64)
	totalLanding := 0.0
	totalMoReceived := 0.0

	for _, campaign := range data {
		date := campaign.SummaryDateHour.Format("2006-01-02")
		prevDate := campaign.SummaryDateHour.AddDate(0, 0, -1).Format("2006-01-02")

		if params.DataType == "monthly_report" {
			date = campaign.SummaryDateHour.Format("2006-01")
			prevDate = campaign.SummaryDateHour.AddDate(0, -1, 0).Format("2006-01")
		}

		if days[date] == nil {
			days[date] = make(map[string]map[string]interface{})
		}

		for _, indicator := range params.DataIndicators {
			val := getIndicatorValueTrafficReport(campaign, indicator)

			if days[date][indicator] == nil {
				days[date][indicator] = map[string]interface{}{
					"value":      0.0,
					"percentage": 0.0,
				}
			}

			prevValue := getPreviousValueTrafficReport(days[prevDate], indicator)

			curVal := days[date][indicator]["value"].(float64)
			newVal := curVal + val
			days[date][indicator]["value"] = newVal
			days[date][indicator]["percentage"] = countPercentageTrafficReport(newVal, prevValue)

			// simpan untuk hitung cr_mo nanti
			if indicator == "landing" {
				dayLanding[date] += val
				totalLanding += val
			}
			if indicator == "mo_received" {
				dayMoReceived[date] += val
				totalMoReceived += val
			}

			// simpan totals untuk indikator lain
			if indicator != "cr_mo" {
				totals[indicator] += val
			}
		}
	}

	// hitung cr_mo harian
	for date := range days {
		if contains(params.DataIndicators, "cr_mo") {
			var cr, prevCr float64

			if dayLanding[date] > 0 {
				cr = (dayMoReceived[date] / dayLanding[date]) * 100
			}

			// cari prevDate sesuai type
			var prevDate string
			if params.DataType == "monthly_report" {
				parts, _ := time.Parse("2006-01", date)
				prevDate = parts.AddDate(0, -1, 0).Format("2006-01")
			} else {
				parts, _ := time.Parse("2006-01-02", date)
				prevDate = parts.AddDate(0, 0, -1).Format("2006-01-02")
			}

			if prev, ok := days[prevDate]; ok {
				if pv, ok2 := prev["cr_mo"]["value"].(float64); ok2 {
					prevCr = pv
				}
			}

			days[date]["cr_mo"] = map[string]interface{}{
				"value":      cr,
				"percentage": countPercentageTrafficReport(cr, prevCr),
			}
		}
	}

	// hitung cr_mo total
	if contains(params.DataIndicators, "cr_mo") {
		if totalLanding > 0 {
			totals["cr_mo"] = (totalMoReceived / totalLanding) * 100
		} else {
			totals["cr_mo"] = 0
		}
	}

	// Final calculations
	tmoEnd := countTmoEndTrafficReport(totals, startDate, endDate)

	summaryData := map[string]interface{}{
		"data_indicators": params.DataIndicators,
		"total":           totals,
		"avg":             countAverageTrafficReport(totals, days),
		"t_mo_end":        tmoEnd,
	}

	return mergeDaysTrafficReport(summaryData, days)
}

func generateHourlyTrafficReport(data []entity.SummaryTraffic, params entity.TrafficReportParams, date time.Time) map[string]interface{} {
	hours := make(map[string]map[string]map[string]interface{})
	totals := make(map[string]float64)

	for i := 0; i < 24; i++ {
		hourStr := fmt.Sprintf("%02d", i)
		hours[hourStr] = make(map[string]map[string]interface{})
		for _, indicator := range params.DataIndicators {
			hours[hourStr][indicator] = map[string]interface{}{
				"value":      0.0,
				"percentage": 0.0,
			}
		}
	}

	for _, campaign := range data {
		if !sameDay(campaign.SummaryDateHour, date) {
			continue
		}
		hour := campaign.SummaryDateHour.Format("15")

		for _, indicator := range params.DataIndicators {
			val := getIndicatorValueTrafficReport(campaign, indicator)
			current := hours[hour][indicator]["value"].(float64)
			hours[hour][indicator]["value"] = current + val
			totals[indicator] += val
		}
	}

	for i := 0; i < 24; i++ {
		hourStr := fmt.Sprintf("%02d", i)
		for _, indicator := range params.DataIndicators {
			current := hours[hourStr][indicator]["value"].(float64)

			previous := 0.0
			for j := i - 1; j >= 0; j-- {
				prevHour := fmt.Sprintf("%02d", j)
				if valMap, ok := hours[prevHour]; ok {
					previous = valMap[indicator]["value"].(float64)
					break
				}
			}

			hours[hourStr][indicator]["percentage"] = countPercentageTrafficReport(current, previous)
		}
	}

	if sameDay(date, time.Now()) {
		currHour := time.Now().Hour()
		for i := currHour + 1; i < 24; i++ {
			hourStr := fmt.Sprintf("%02d", i)
			delete(hours, hourStr)
		}
	}

	summaryData := map[string]interface{}{
		"data_indicators": params.DataIndicators,
		"total":           totals,
		"avg":             countAverageTrafficReport(totals, hours),
	}

	return mergeHoursTrafficReport(summaryData, hours)
}


func getIndicatorValueTrafficReport(item entity.SummaryTraffic, key string) float64 {

	key = SnakeToCamelTrafficReport(key)
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


func countAverageTrafficReport(totals map[string]float64, days map[string]map[string]map[string]interface{}) map[string]float64 {
	averages := map[string]float64{}
	counts := map[string]int{}

	// hitung jumlah hari yang punya data untuk tiap indikator
	for _, indicators := range days {
		for key, val := range indicators {
			if v, ok := val["value"].(float64); ok && v > 0 {
				counts[key]++
			}
		}
	}

	for key, total := range totals {
		if cnt := counts[key]; cnt > 0 {
			avg := total / float64(cnt)

			// untuk mo_received & landing dibulatkan ke integer
			if key == "mo_received" || key == "landing" {
				avg = float64(int(avg + 0.5))
			}
			averages[key] = avg
		} else {
			averages[key] = 0
		}
	}

	return averages
}

func countTmoEndTrafficReport(totals map[string]float64, startDate time.Time, endDate time.Time) map[string]float64 {
	tmoEnd := map[string]float64{}
	totalDaysRunning := int(math.Ceil(endDate.Sub(startDate).Hours() / 24))
	if totalDaysRunning < 1 {
		totalDaysRunning = 1
	}

	lastMonthEnd := endDate.AddDate(0, 0, -endDate.Day())
	lastMonthStart := lastMonthEnd.AddDate(0, 0, -lastMonthEnd.Day())
	totalDaysLastMonth := int(math.Ceil(lastMonthEnd.Sub(lastMonthStart).Hours()/24)) + 1

	for key, value := range totals {
		result := (value / float64(totalDaysRunning)) * float64(totalDaysLastMonth)

		if key == "mo_received" || key == "landing" {
			result = float64(int(result + 0.5))
		}

		tmoEnd[key] = result
	}

	return tmoEnd
}

func mergeDaysTrafficReport(summaryData map[string]interface{}, days map[string]map[string]map[string]interface{}) map[string]interface{} {
	for key, value := range days {
		summaryData[key] = value
	}
	return summaryData
}

func mergeHoursTrafficReport(summaryData map[string]interface{}, hours map[string]map[string]map[string]interface{}) map[string]interface{} {
	for key, value := range hours {
		summaryData[key] = value
	}
	return summaryData
}

func getPreviousValueTrafficReport(data map[string]map[string]interface{}, key string) float64 {
	value := 0.0
	if val, exists := data[key]; exists {
		value = val["value"].(float64)
	}
	return value
}

func countPercentageTrafficReport(now, prev float64) float64 {
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

func extractQueryArrayTrafficReport(c *fiber.Ctx, key string) []string {
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

func SnakeToCamelTrafficReport(snake string) string {
	if snake == "cr_mo" {
		return "CrMO"
	}

	words := strings.Split(snake, "_")
	for i, word := range words {
		words[i] = strings.ToLower(word)
		if len(word) > 0 {
			words[i] = strings.ToUpper(words[i][:1]) + words[i][1:]
		}
	}
	return strings.Join(words, "")
}

func mergeMapsTrafficReport(map1, map2 map[string]interface{}) map[string]interface{} {
	mergedMap := make(map[string]interface{})
	for key, value := range map1 {
		mergedMap[key] = value
	}
	for key, value := range map2 {
		mergedMap[key] = value
	}
	return mergedMap
}

func sortDataTrafficReport(data []map[string]interface{}, dataBasedOn string, dataBasedOnIndicator string) []map[string]interface{} {
	sortedData := make([]map[string]interface{}, len(data))
	copy(sortedData, data)

	sort.Slice(sortedData, func(i, j int) bool {
		totalI, okI := sortedData[i]["total"].(map[string]float64)
		totalJ, okJ := sortedData[j]["total"].(map[string]float64)

		if !okI || !okJ {
			return i < j
		}

		valI := totalI[dataBasedOnIndicator]
		valJ := totalJ[dataBasedOnIndicator]

		if dataBasedOn == "highest" {
			return valI > valJ
		} else {
			return valI < valJ
		}
	})

	return sortedData
}

func generateHourlyChartTrafficReport(data []entity.SummaryTraffic, params entity.TrafficReportParams, date time.Time) map[string][]map[string]interface{} {
	chart := make(map[string][]map[string]interface{})

	for _, indicator := range params.DataIndicators {
		chart[indicator] = []map[string]interface{}{}
	}

	hourlyMap := make(map[string]map[string]float64)
	hoursToUse := []string{}
	for i := 0; i < 24; i++ {
		hourStr := fmt.Sprintf("%02d", i)
		hoursToUse = append(hoursToUse, hourStr)
		hourlyMap[hourStr] = make(map[string]float64)
		for _, indicator := range params.DataIndicators {
			hourlyMap[hourStr][indicator] = 0.0
		}
	}

	// Fill hourly values
	countsPerHour := make(map[string]map[string]int)
	for _, item := range data {
		if !sameDay(item.SummaryDateHour, date) {
			continue
		}
		if params.URLServiceKey != "" && item.URLServiceKey != params.URLServiceKey {
			continue
		}

		hour := item.SummaryDateHour.Format("15")
		if countsPerHour[hour] == nil {
			countsPerHour[hour] = make(map[string]int)
		}

		for _, indicator := range params.DataIndicators {
			val := getIndicatorValueTrafficReport(item, indicator)
			hourlyMap[hour][indicator] += val
			countsPerHour[hour][indicator]++
		}
	}

	// Average if "All" param
	if params.All == "true" {
		for hour, indicators := range countsPerHour {
			for indicator, cnt := range indicators {
				if cnt > 0 {
					hourlyMap[hour][indicator] = hourlyMap[hour][indicator] / float64(cnt)
				} else {
					hourlyMap[hour][indicator] = 0.0
				}
			}
		}
	}

	// Build chart with percentage
	for i, hourStr := range hoursToUse {
		prevHourStr := fmt.Sprintf("%02d", i-1)
		for _, indicator := range params.DataIndicators {
			current := hourlyMap[hourStr][indicator]
			previous := 0.0
			if i > 0 {
				previous = hourlyMap[prevHourStr][indicator]
			}
			percentage := countPercentageTrafficReport(current, previous)

			chart[indicator] = append(chart[indicator], map[string]interface{}{
				"time":       hourStr + ":00",
				"value":      current,
				"percentage": percentage,
			})
		}
	}

	// Trim future hours if today
	if sameDay(date, time.Now()) {
		currentHour := time.Now().Hour()
		for _, indicator := range params.DataIndicators {
			if currentHour+1 < len(chart[indicator]) {
				chart[indicator] = chart[indicator][:currentHour+1]
			}
		}
	}

	return chart
}


func (h *IncomingHandler) GetTrafficReportChart(c *fiber.Ctx) error {
    dataType := c.Query("data-type")
	if dataType == "" {
		dataType = "daily_report"
	}
    dataIndicators := extractQueryArrayTrafficReport(c, "data-indicators[]")
    if len(dataIndicators) == 0 {
        dataIndicators = []string{"landing", "mo_received", "cr_mo", "first_push"}
    }

    params := entity.TrafficReportParams{
        DataType:        dataType,
        DataIndicators:  dataIndicators,
        DateRange:       c.Query("date-range"),
        DateCustomRange: c.Query("date-custom-range"),
        Country:         c.Query("country"),
        Operator:        c.Query("operator"),
        PartnerName:     c.Query("partner-name"),
        Service:         c.Query("service"),
        Adnet:           c.Query("adnet"),
        CampaignId:      c.Query("campaign_id"),
        CampaignName:    c.Query("campaign"),
        URLServiceKey:   c.Query("url_service_key"),
    }

    var summaryData []entity.SummaryTraffic
    var err error

    if dataType == "daily_report" {
        summaryData, _, _, err = h.DS.GetTrafficReport(params)
    } else if dataType == "monthly_report" {
        summaryData, _, _, err = h.DS.GetTrafficReport(params)
    }

    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "code":    fiber.StatusNotFound,
            "message": "empty",
        })
    }

    chart := generateChartData(summaryData, params, dataType)
    var response entity.ReturnResponse
	if err == nil {
		response = entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithData{
				Code:    fiber.StatusOK,
				Message: config.OK_DESC,
				Data:    chart,
			},
		}
	} else {
		response = entity.ReturnResponse{
			HttpStatus: fiber.StatusNotFound,
			Rsp: entity.GlobalResponse{
				Code:    fiber.StatusNotFound,
				Message: "empty",
			},
		}
	}

	r := response
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func generateChartData(data []entity.SummaryTraffic, params entity.TrafficReportParams, chartType string) []map[string]interface{} {
	chart := []map[string]interface{}{}
	grouped := map[string][]entity.SummaryTraffic{}

	for _, item := range data {
		var key string
		if chartType == "daily_report" {
			key = item.SummaryDateHour.Format("2006-01-02")
		} else if chartType == "monthly_report" {
			key = item.SummaryDateHour.Format("2006-01")
		}
		grouped[key] = append(grouped[key], item)
	}

	for date, items := range grouped {
		record := map[string]interface{}{
			"summary_date": date,
		}
	
		landingSum := 0.0
		moReceivedSum := 0.0
	
		for _, indicator := range params.DataIndicators {
			sum := 0.0
			for _, item := range items {
				val := getIndicatorValueTrafficReport(item, indicator)
				sum += val
			}
	
			if indicator == "landing" {
				landingSum = sum
			}
			if indicator == "mo_received" {
				moReceivedSum = sum
			}
	
			if indicator != "cr_mo" {
				record[indicator] = sum
			}
		}
	
		if contains(params.DataIndicators, "cr_mo") {
			if landingSum > 0 {
				record["cr_mo"] = (moReceivedSum / landingSum) * 100
			} else {
				record["cr_mo"] = 0
			}
		}
	
		chart = append(chart, record)
	}
	

	sort.Slice(chart, func(i, j int) bool {
		return chart[i]["summary_date"].(string) < chart[j]["summary_date"].(string)
	})

	return chart
}

func contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}
