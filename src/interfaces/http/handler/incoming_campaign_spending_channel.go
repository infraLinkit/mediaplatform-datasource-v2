package handler

import (
	"encoding/json"
	"log"
	"math"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/config"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/external"
)


func (h *IncomingHandler) DisplayCampaignSpendingChannel(c *fiber.Ctx) error {
	params := buildCampaignSpendingChannelParams(c)
	viewType := c.Query("view-type", "country")

	data, startDate, endDate, err := h.DS.GetSpendingChannelMonitoring(params)
	if err != nil {
		log.Printf("[SpendingChannel] GetSpendingChannelMonitoring error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(entity.GlobalResponse{
			Code: fiber.StatusInternalServerError, Message: err.Error(),
		})
	}

	log.Printf("[SpendingChannel] Handler got %d rows, all=%s, viewType=%s, dataType=%s",
		len(data), params.All, viewType, params.DataType)

	var result []map[string]interface{}

	if params.All == "true" {
		row := generateSpendingValue(data, params, startDate, endDate)
		if viewType == "channel" {
			row["channel"] = "All Data"
		} else {
			row["country"] = "All Data"
		}
		row["level"] = "all"
		result = []map[string]interface{}{row}
	} else if viewType == "channel" {
		result = formatSpendingByChannel(data, params, startDate, endDate)
	} else {
		result = formatSpendingByCountry(data, params, startDate, endDate)
	}

	sorted := sortBySpending(result, params.DataBasedOn)

	if len(sorted) > 0 {
		if b, jerr := json.Marshal(sorted[0]); jerr == nil {
			log.Printf("[SpendingChannel] First result item: %s", string(b))
		}
	} else {
		log.Printf("[SpendingChannel] WARN: sorted result is empty! input rows=%d", len(data))
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithData{
		Code:    fiber.StatusOK,
		Message: config.OK_DESC,
		Data:    sorted,
	})
}

func (h *IncomingHandler) DisplayCampaignSpendingChannelCountryChildren(c *fiber.Ctx) error {
	params := buildCampaignSpendingChannelParams(c)
	viewMode := c.Query("view-mode", "operator")
	country := c.Query("country")

	raw, startDate, endDate, err := h.DS.GetSpendingChannelMonitoring(params)
	if err != nil || len(raw) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(entity.GlobalResponse{
			Code: fiber.StatusNotFound, Message: "empty",
		})
	}

	var filtered []entity.CampaignSpendingChannelMonitoring
	for _, r := range raw {
		if r.Country == country {
			filtered = append(filtered, r)
		}
	}

	var children []interface{}
	if viewMode == "channel" {
		children = groupChannelChildrenCSC(filtered, params, startDate, endDate)
	} else {
		children = groupOperatorChildrenCSC(filtered, params, startDate, endDate)
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithData{
		Code: fiber.StatusOK, Message: config.OK_DESC, Data: children,
	})
}

// ─── By-Channel Formatter ─────────────────────────────────────────────────────

func formatSpendingByChannel(
	data []entity.CampaignSpendingChannelMonitoring,
	params entity.ParamsCampaignSpendingChannel,
	startDate, endDate time.Time,
) []map[string]interface{} {

	grouped := make(map[string][]entity.CampaignSpendingChannelMonitoring)
	for _, c := range data {
		grp := external.ResolveChannelGroup(c.Channel)
		grouped[grp] = append(grouped[grp], c)
	}

	var result []map[string]interface{}
	for channelName, campaigns := range grouped {
		row := generateSpendingValue(campaigns, params, startDate, endDate)
		row["channel"] = channelName
		row["level"] = "channel"
		row["_children"] = groupCountryChildrenInChannel(campaigns, params, startDate, endDate)
		result = append(result, row)
	}
	return result
}

func groupCountryChildrenInChannel(
	campaigns []entity.CampaignSpendingChannelMonitoring,
	params entity.ParamsCampaignSpendingChannel,
	startDate, endDate time.Time,
) []interface{} {
	var out []interface{}

	grouped := make(map[string][]entity.CampaignSpendingChannelMonitoring)
	for _, c := range campaigns {
		grouped[c.Country] = append(grouped[c.Country], c)
	}

	for country, cc := range grouped {
		row := generateSpendingValue(cc, params, startDate, endDate)
		row["channel"] = country
		row["level"] = "country"
		out = append(out, row)
	}
	return out
}

// ─── By-Country Formatter ─────────────────────────────────────────────────────

func formatSpendingByCountry(
	data []entity.CampaignSpendingChannelMonitoring,
	params entity.ParamsCampaignSpendingChannel,
	startDate, endDate time.Time,
) []map[string]interface{} {

	grouped := make(map[string][]entity.CampaignSpendingChannelMonitoring)
	for _, c := range data {
		grouped[c.Country] = append(grouped[c.Country], c)
	}

	var result []map[string]interface{}
	for country, perCountry := range grouped {
		row := generateSpendingValue(perCountry, params, startDate, endDate)
		row["country"] = country
		row["level"] = "country"
		row["_children"] = groupOperatorChildrenCSC(perCountry, params, startDate, endDate)
		result = append(result, row)
	}
	return result
}

// ─── Children Grouping ────────────────────────────────────────────────────────

func groupOperatorChildrenCSC(
	campaigns []entity.CampaignSpendingChannelMonitoring,
	params entity.ParamsCampaignSpendingChannel,
	startDate, endDate time.Time,
) []interface{} {
	var out []interface{}

	grouped := make(map[string][]entity.CampaignSpendingChannelMonitoring)
	for _, c := range campaigns {
		grouped[c.Operator] = append(grouped[c.Operator], c)
	}

	for operator, cc := range grouped {
		row := generateSpendingValue(cc, params, startDate, endDate)
		row["country"] = operator
		row["level"] = "operator"
		out = append(out, row)
	}
	return out
}

func groupChannelChildrenCSC(
	campaigns []entity.CampaignSpendingChannelMonitoring,
	params entity.ParamsCampaignSpendingChannel,
	startDate, endDate time.Time,
) []interface{} {
	var out []interface{}

	grouped := make(map[string][]entity.CampaignSpendingChannelMonitoring)
	for _, c := range campaigns {
		grp := external.ResolveChannelGroup(c.Channel)
		grouped[grp] = append(grouped[grp], c)
	}

	for channelName, cc := range grouped {
		row := generateSpendingValue(cc, params, startDate, endDate)
		row["country"] = channelName
		row["level"] = "channel"
		out = append(out, row)
	}
	return out
}

// ─── Core Value Generator ─────────────────────────────────────────────────────

func generateSpendingValue(
	data []entity.CampaignSpendingChannelMonitoring,
	params entity.ParamsCampaignSpendingChannel,
	startDate, endDate time.Time,
) map[string]interface{} {

	startDate = normalisePeriodStart(startDate.In(time.Local), params.DataType)

	days := map[string]map[string]interface{}{}
	cur := startDate
	for !cur.After(endDate) {
		key := formatDateKey(cur, params.DataType)
		if _, exists := days[key]; !exists {
			days[key] = map[string]interface{}{"value": 0.0, "percentage": 0.0}
		}
		cur = incrementDate(cur, params.DataType)
	}

	total := 0.0

	for _, campaign := range data {
		key := formatDateKey(campaign.SummaryDate.In(time.Local), params.DataType)
		if slot, ok := days[key]; ok {
			slot["value"] = slot["value"].(float64) + campaign.SBAF
			days[key] = slot
		} else {
			log.Printf("[SpendingChannel] WARN: date key %q not found in slots (SummaryDate=%s)",
				key, campaign.SummaryDate.Format(time.RFC3339))
		}
		total += campaign.SBAF
	}

	cur = startDate
	for !cur.After(endDate) {
		key := formatDateKey(cur, params.DataType)
		prev := formatPreviousDateKey(cur, params.DataType)
		curr := days[key]["value"].(float64)
		prevVal := 0.0
		if s, ok := days[prev]; ok {
			prevVal = s["value"].(float64)
		}
		days[key]["percentage"] = safePercentage(curr, prevVal)
		cur = incrementDate(cur, params.DataType)
	}

	periodCount := countPeriods(startDate, endDate, params.DataType)

	result := map[string]interface{}{
		"total":    total,
		"avg":      total / float64(periodCount),
		"t_mo_end": calcTmoEnd(total, startDate, endDate),
	}
	for k, v := range days {
		result[k] = v
	}
	return result
}

func normalisePeriodStart(d time.Time, dataType string) time.Time {
	switch dataType {
	case "monthly_report":
		return time.Date(d.Year(), d.Month(), 1, 0, 0, 0, 0, d.Location())
	case "weekly_report":
		wd := int(d.Weekday())
		if wd == 0 {
			wd = 7
		}
		return d.AddDate(0, 0, -(wd - 1))
	default:
		return d
	}
}

// ─── Stats Helpers ────────────────────────────────────────────────────────────

func countPeriods(startDate, endDate time.Time, dataType string) int {
	n := 0
	cur := startDate
	for !cur.After(endDate) {
		n++
		cur = incrementDate(cur, dataType)
	}
	if n < 1 {
		return 1
	}
	return n
}

func calcTmoEnd(total float64, startDate, endDate time.Time) float64 {
	daysRunning := int(math.Ceil(endDate.Sub(startDate).Hours()/24)) + 1
	if daysRunning < 1 {
		daysRunning = 1
	}
	now := time.Now()
	daysInMonth := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, now.Location()).Day()
	return (total / float64(daysRunning)) * float64(daysInMonth)
}

// ─── Sort ─────────────────────────────────────────────────────────────────────

func sortBySpending(data []map[string]interface{}, dataBasedOn string) []map[string]interface{} {
	sorted := make([]map[string]interface{}, len(data))
	copy(sorted, data)
	sort.Slice(sorted, func(i, j int) bool {
		vi, _ := sorted[i]["total"].(float64)
		vj, _ := sorted[j]["total"].(float64)
		if dataBasedOn == "lowest_cost" || dataBasedOn == "lowest_traffic" || dataBasedOn == "lowest" {
			return vi < vj
		}
		return vi > vj
	})
	return sorted
}

// ─── Param Builder ────────────────────────────────────────────────────────────

func buildCampaignSpendingChannelParams(c *fiber.Ctx) entity.ParamsCampaignSpendingChannel {
	return entity.ParamsCampaignSpendingChannel{
		DataType:        c.Query("data-type", "daily_report"),
		ReportType:      c.Query("report-type", "operator_summary"),
		Country:         c.Query("country"),
		Operator:        c.Query("operator"),
		PartnerName:     c.Query("partner_name"),
		Service:         c.Query("service"),
		ChannelCampaign: c.Query("channel-campaign"),
		DataBasedOn:     c.Query("data-based-on", "highest_cost"),
		DateRange:       c.Query("date-range", "this_month"),
		DateStart:       c.Query("date-start"),
		DateEnd:         c.Query("date-end"),
		DateCustomRange: c.Query("date-custom-range"),
		All:             c.Query("all"),
	}
}

// ─── Date Helpers ─────────────────────────────────────────────────────────────

func formatDateKey(date time.Time, dataType string) string {
	d := date.In(time.Local)
	switch dataType {
	case "monthly_report":
		return d.Format("2006-01")
	case "weekly_report":
		wd := int(d.Weekday())
		if wd == 0 {
			wd = 7
		}
		monday := d.AddDate(0, 0, -(wd - 1))
		return monday.Format("2006-01-02")
	default:
		return d.Format("2006-01-02")
	}
}

func formatPreviousDateKey(date time.Time, dataType string) string {
	d := date.In(time.Local)
	switch dataType {
	case "monthly_report":
		return d.AddDate(0, -1, 0).Format("2006-01")
	case "weekly_report":
		return d.AddDate(0, 0, -7).Format("2006-01-02")
	default:
		return d.AddDate(0, 0, -1).Format("2006-01-02")
	}
}

func incrementDate(date time.Time, dataType string) time.Time {
	switch dataType {
	case "monthly_report":
		return date.AddDate(0, 1, 0)
	case "weekly_report":
		return date.AddDate(0, 0, 7)
	default:
		return date.AddDate(0, 0, 1)
	}
}