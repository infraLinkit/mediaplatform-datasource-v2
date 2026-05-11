package entity

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/config"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/external"
	"github.com/sirupsen/logrus"
)

type (
	DisplayPinReport struct { // api report
		ID          int       `form:"id" json:"id"`
		DateSend    time.Time `form:"date_send" json:"date_send"`
		CampaignId  string    `form:"campaign_id" json:"campaign_id"`
		Channel     string    `form:"channel" json:"channel"`
		Country     string    `form:"country" json:"country"`
		Operator    string    `form:"operator" json:"operator"`
		Partner     string    `form:"partner" json:"partner"`
		Agency      string    `form:"agency" json:"agency"`
		Aggregator  string    `form:"aggregator" json:"aggregator"`
		Adnets      []string  `form:"adnet" json:"adnet"`
		Service     string    `form:"service" json:"service"`
		Company     string    `form:"company" json:"company"`
		PageSize    int       `form:"page_size" json:"page_size"`
		Page        int       `form:"page" json:"page"`
		Action      string    `form:"action" json:"action"`
		DateRange   string    `form:"date_range" json:"date_range"`
		DateBefore  string    `form:"date_before" json:"date_before"`
		DateAfter   string    `form:"date_after" json:"date_after"`
		Draw        int       `form:"draw" json:"draw"`
		Reload      string    `form:"draw" json:"reload"`
		OrderColumn string    `form:"order_column" json:"order_column"`
		OrderDir    string    `form:"order_dir" json:"order_dir"`
	}

	DisplayPinPerformanceReport struct {
		Country    string   `form:"country" json:"country"`
		Adnet      string   `form:"adnet" json:"adnet"`
		Operator   string   `form:"operator" json:"operator"`
		Service    string   `form:"service" json:"service"`
		Adnets     []string `form:"adnet" json:"adnet"`
		Page       int      `form:"page" json:"page"`
		DateRange  string   `form:"date_range" json:"date_range"`
		DateBefore string   `form:"date_before" json:"date_before"`
		DateAfter  string   `form:"date_after" json:"date_after"`
		Action     string   `form:"action" json:"action"`
		Draw       int      `form:"draw" json:"draw"`
		PageSize   int      `form:"page_size" json:"page_size"`
	}

	DisplayConversionLogReport struct {
		Country        string `form:"country" json:"country"`
		Adnet          string `form:"adnet" json:"adnet"`
		Agency         string `form:"agency" json:"agency"`
		Operator       string `form:"operator" json:"operator"`
		CampaignType   string `form:"campaign_type" json:"campaign_type"`
		CampaignId     string `form:"campaign_id" json:"campaign_id"`
		CampaignName   string `form:"campaign_name" json:"campaign_name"`
		StatusPostback string `form:"status_postback" json:"status_postback"`
		Pixel          string `form:"pixel" json:"pixel"`
		Page           int    `form:"page" json:"page"`
		DateRange      string `form:"date_range" json:"date_range"`
		DateStart      string `form:"date_start" json:"date_start"`
		DateEnd        string `form:"date_end" json:"date_end"`
		Action         string `form:"action" json:"action"`
		Draw           int    `form:"draw" json:"draw"`
		PageSize       int    `form:"page_size" json:"page_size"`
		Order          string `form:"order" json:"order"`
	}

	DisplayCPAReport struct { // cpa
		ID                int       `form:"id" json:"id"`
		SummaryDate       time.Time `form:"summary_date" json:"summary_date"`
		CampaignId        string    `form:"campaign_id" json:"campaign_id"`
		UrlServiceKey     string    `form:"url_service_key" json:"url_service_key"`
		Channel           string    `form:"channel" json:"channel"`
		CampaignName      string    `form:"campaign_name" json:"campaign_name"`
		Country           string    `form:"country" json:"country"`
		Operator          string    `form:"operator" json:"operator"`
		Partner           string    `form:"partner" json:"partner"`
		Agency            string    `form:"agency" json:"agency"`
		Aggregator        string    `form:"aggregator" json:"aggregator"`
		Adnets            []string  `form:"adnet" json:"adnet"`
		Service           string    `form:"service" json:"service"`
		DataBasedOn       string    `form:"data_based_on" json:"data_based_on"`
		Cost              float64   `form:"cost" json:"cost"`
		SBAF              float64   `form:"sbaf" json:"sbaf"`
		SAAF              float64   `form:"saaf" json:"saaf"`
		RatioSend         int       `form:"ratio_send" json:"ratio_send"`
		RatioReceive      int       `form:"ratio_receive" json:"ratio_receive"`
		Company           string    `form:"company" json:"company"`
		ClientType        string    `form:"client_type" json:"client_type"`
		CostPerConversion float64   `form:"cost_per_conversion" json:"cost_per_conversion"`
		AgencyFee         float64   `form:"agency_fee" json:"agency_fee"`
		PageSize          int       `form:"page_size" json:"page_size"`
		Page              int       `form:"page" json:"page"`
		Action            string    `form:"action" json:"action"`
		DateRange         string    `form:"date_range" json:"date_range"`
		DateBefore        string    `form:"date_before" json:"date_before"`
		DateAfter         string    `form:"date_after" json:"date_after"`
		Draw              int       `form:"draw" json:"draw"`
		Reload            string    `form:"draw" json:"reload"`
		OrderColumn       string    `form:"order_column" json:"order_column"`
		OrderDir          string    `form:"order_dir" json:"order_dir"`
		CampaignObjective string    `form:"order_dir" json:"campaign_objective"`
	}

	TotalSummaryCampaign struct {
		MoReceived         int     `json:"mo_received"`
		Postback           int     `json:"postback"`
		PO                 float64 `json:"po"`
		CostPerConversion  float64 `json:"cost_per_conversion"`
		TotalWakiAgencyFee float64 `json:"total_waki_agency_fee"`
		SBAF               float64 `json:"sbaf"`
		SAAF               float64 `json:"saaf"`
		TechnicalFee       float64 `json:"technical_fee"`
		Landing            int     `json:"landing"`
		Clicked            int     `json:"clicked"`
		CrMO               float64 `json:"cr_mo"`
		CrPostback         float64 `json:"cr_postback"`
		WakiRevenue        float64 `json:"waki_revenue"`
		ECPA               float64 `json:"e_cpa"`
		PricePerMO         float64 `json:"price_per_mo"`
	}

	// CostReport is the data row returned by all cost report queries.
	CostReport struct {
		SummaryDate time.Time `json:"summary_date"`
		Adnet       string    `json:"adnet"`
		Country     string    `json:"country"`
		Operator    string    `json:"operator"`
		ChannelType string    `json:"channel_type"` // NEW: mapped from campaign_objective
		Landing     float64   `json:"landing"`
		CrPostback  float64   `json:"cr_postback"`
		ShortCode   string    `json:"short_code"`
		UrlAfter    string    `json:"url_after"`
		Conversion1 float64   `json:"conversion1"`
		Cost1       float64   `json:"cost1"`
		Conversion2 float64   `json:"conversion2"`
		Cost2       float64   `json:"cost2"`
		Saaf1       float64   `json:"saaf1"`
		Saaf2       float64   `json:"saaf2"`
	}

	// DisplayCostReport holds all filter parameters from the frontend.
	DisplayCostReport struct {
		SummaryDate   time.Time `json:"summary_date"`
		Adnet         string    `json:"adnet"`
		Adnets        []string  `json:"adnets"`
		Country       string    `json:"country"`
		Operator      string    `json:"operator"`
		ChannelType   string    `json:"channel_type"`   // NEW: filter by channel type
		GroupBy       string    `json:"group_by"`        // NEW: "adnet" | "country"
		DataIndicator string    `json:"data_indicator"`  // NEW: "s2s" | "api" | ""
		Landing       float64   `json:"landing"`
		CrPostback    float64   `json:"cr_postback"`
		ShortCode     string    `json:"short_code"`
		UrlAfter      string    `json:"url_after"`
		Conversion1   float64   `json:"conversion1"`
		Cost1         float64   `json:"cost1"`
		Conversion2   float64   `json:"conversion2"`
		Cost2         float64   `json:"cost2"`
		Action        string    `json:"action"`
		CampaignType  string    `json:"campaign_type"`
		DateRange     string    `json:"date_range"`
		DateBefore    string    `json:"date_before"`
		DateAfter     string    `json:"date_after"`
		PageSize      int       `json:"page_size"`
		Page          int       `json:"page"`
		Draw          int       `json:"draw"`
		DataBasedOn   string    `json:"data_based_on"`
		FromChannel   bool 
	}

	DisplayAlertReport struct {
		Action       string `json:"action"`
		Country      string `json:"country"`
		Operator     string `json:"operator"`
		CampaignName string `json:"campaign_name"`
		Service      string `json:"service"`
		DateRange    string `json:"date_range"`
		DateBefore   string `json:"date_before"`
		DateAfter    string `json:"date_after"`
		Page         int    `json:"page"`
		Draw         int    `json:"draw"`
		PageSize     int    `json:"page_size"`
		ExportData   string `json:"export_data"`
	}

	PerformaceReportParams struct {
		Country      string `form:"country" json:"country"`
		Company      string `form:"company" json:"company"`
		ClientType   string `form:"client_type" json:"client_type"`
		Operator     string `form:"operator" json:"operator"`
		CampaignName string `form:"campaign_name" json:"campaign_name"`
		CampaignType string `form:"campaign_type" json:"campaign_type"`
		Publisher    string `form:"publisher" json:"publisher"`
		Service      string `form:"service" json:"service"`
		CampaignId   string `form:"campaign_id" json:"campaign_id"`
		Partner      string `form:"partner" json:"partner"`
		PageSize     int    `form:"page_size" json:"page_size"`
		Page         int    `form:"page" json:"page"`
		Action       string `form:"action" json:"action"`
		DateStart    string `form:"date_before" json:"date_start"`
		DateEnd      string `form:"date_after" json:"date_end"`
		DateRange    string `form:"date_range" json:"date_range"`
		DateBefore   string `form:"date_before" json:"date_before"`
		DateAfter    string `form:"date_after" json:"date_after"`
		Draw         int    `form:"draw" json:"draw"`
	}

	PerformanceReport struct {
		Country            string  `json:"country"`
		Company            string  `json:"company"`
		ClientType         string  `json:"client_type"`
		CampaignName       string  `json:"campaign_name"`
		Partner            string  `json:"partner"`
		Operator           string  `json:"operator"`
		Service            string  `json:"service"`
		Adnet              string  `json:"adnet"`
		PixelReceived      int     `json:"pixel_received"`
		PixelSend          int     `json:"pixel_send"`
		CRPostback         int     `json:"cr_postback"`
		CRMo               int     `json:"cr_mo"`
		Landing            int     `json:"landing"`
		RatioSend          int     `json:"ratio_send"`
		RatioReceive       int     `json:"ratio_receive"`
		PricePerPostback   float64 `json:"price_per_postback"`
		CostPerConversion  float64 `json:"cost_per_conversion"`
		AgencyFee          float64 `json:"agency_fee"`
		SpendingToAdnets   float64 `json:"spending_to_adnets"`
		TotalWakiAgencyFee float64 `json:"total_waki_agency_fee"`
		TotalSpending      float64 `json:"total_spending"`
		TotalFP            float64 `json:"total_fp"`
		SuccessFP          float64 `json:"success_fp"`
		ECPA               float64 `json:"e_cpa"`
		ARPUROI            float64 `json:"arpu_roi"`
		ARPU90             float64 `json:"arpu_90"`
		BillrateFP         float64 `json:"billrate_fp"`
	}

	ArpuParams struct {
		Country  string `form:"country" json:"country"`
		Operator string `form:"operator" json:"operator"`
		Service  string `form:"service" json:"service"`
		From     string `form:"from" json:"from"`
		To       string `form:"to" json:"to"`
	}

	ARPUResponse struct {
		Status  int          `json:"status"`
		Message string       `json:"message"`
		Data    *ARPUDataSet `json:"data"`
	}

	ARPUDataSet struct {
		DateHit   string         `json:"date_hit"`
		Country   string         `json:"country"`
		Operator  string         `json:"operator"`
		Service   string         `json:"service"`
		Keyword   string         `json:"keyword"`
		Publisher string         `json:"publisher"`
		Data      []ARPUDataItem `json:"data"`
	}

	ARPUDataItem struct {
		Adnet        string  `json:"adnet"`
		Arpu90       float64 `json:"arpu90"`
		Arpu90Net    float64 `json:"arpu90_net"`
		Arpu90USD    float64 `json:"arpu90_usd"`
		Arpu90USDNet float64 `json:"arpu90_usd_net"`
		Service      string  `json:"service"`
	}

	WakiCallbackParams struct {
		Date          string
		Publisher     string
		Adnet         string
		Operator      string
		Adn           string
		Client        string
		Aggregator    string
		Country       string
		Service       string
		MoReceived    string
		MoPostback    string
		TotalMo       string
		TotalPostback string
		Landing       string
		CrMoReceived  string
		CrMoPostback  string
		UrlCampaign   string
		UrlService    string
		Sbaf          string
		Saaf          string
		Spending      string
		Campaign      string
		Payout        string
		PricePerMo    string
	}

	SuccessRateData struct {
		Date        string `json:"date"`
		Operator    string `json:"operator"`
		Service     string `json:"service"`
		SuccessRate string `json:"success_rate"`
	}

	SuccessRateResponse struct {
		Status string          `json:"status"`
		Code   int             `json:"code"`
		Data   SuccessRateData `json:"data"`
	}
)

func NewInstanceTrxPinReport(c *fiber.Ctx, cfg *config.Cfg) *ApiPinReport {

	m := c.Queries()

	mo, _ := strconv.Atoi(m["mo"])
	postback, _ := strconv.Atoi(m["postback"])

	sbaf, _ := strconv.ParseFloat(m["sbaf"], 64)
	saaf, _ := strconv.ParseFloat(m["saaf"], 64)
	payoutAdn, _ := strconv.ParseFloat(m["payout_adn"], 64)
	payoutAF, _ := strconv.ParseFloat(m["payout_af"], 64)

	country := external.NormalizeCountries(m["country"])

	operator := strings.ToUpper(m["operator"])
	if operator == "" && m["telco"] != "" {
		operator = strings.ToUpper(m["telco"])
	}

	var dateSend time.Time
	if m["date_send"] == "" {
		dateSend = time.Now().In(cfg.TZ)
	} else {
		dateSend = external.ParseVendorDateSend(cfg.TZ, m["date_send"])
	}

	pin := ApiPinReport{
		CampaignId:    strings.ToUpper(m["campaign_id"]),
		Country:       country,
		Company:       strings.ToUpper(m["company"]),
		Adnet:         strings.ToUpper(m["adnet"]),
		Service:       strings.ToUpper(m["service"]),
		Operator:      operator,
		DateSend:      dateSend,
		PayoutAdn:     payoutAdn,
		PayoutAF:      payoutAF,
		TotalMO:       mo,
		TotalPostback: postback,
		SBAF:          sbaf,
		SAAF:          saaf,
	}

	return &pin
}

func BuildPinReportCalculation(o *ApiPinReport) {
	o.SBAF = float64(o.TotalPostback) * o.PayoutAdn
	o.SAAF = float64(o.TotalPostback) * o.PayoutAF

	if o.TotalMO > 0 {
		o.PricePerMO = o.SAAF / float64(o.TotalMO)
	} else {
		o.PricePerMO = 0
	}

	o.WakiRevenue = o.SAAF - o.SBAF
}

func BuildAPIPerformanceReportCalculation(o *ApiPinPerformance) {
	if o.PinOkRatio > 0 {
		o.TotalSpending = o.CPA * float64(o.PinOkRatio)
		o.TotalSpendingAfterWaki = o.CPAWaki * float64(o.PinOkRatio)
		o.CAC = o.TotalSpendingAfterWaki / float64(o.PinOK)
		if o.ChargedMO > 0 {
			o.PaidCAC = o.TotalSpendingAfterWaki / o.ChargedMO
		}
	}
}

func (t *ApiPinReport) ValidateParams(Logs *logrus.Logger) ReturnResponse {

	if t.Adnet == "" {

		return ReturnResponse{HttpStatus: fiber.StatusBadRequest, Rsp: GlobalResponse{Code: fiber.StatusBadRequest, Message: "Parameter Adnet is mandatory"}}

	} else if t.Country == "" {

		return ReturnResponse{HttpStatus: fiber.StatusBadRequest, Rsp: GlobalResponse{Code: fiber.StatusBadRequest, Message: "Parameter Country is mandatory"}}

	} else if t.Service == "" {

		return ReturnResponse{HttpStatus: fiber.StatusBadRequest, Rsp: GlobalResponse{Code: fiber.StatusBadRequest, Message: "Parameter Service is mandatory"}}

	} else if t.Operator == "" {

		return ReturnResponse{HttpStatus: fiber.StatusBadRequest, Rsp: GlobalResponse{Code: fiber.StatusBadRequest, Message: "Parameter Operator is mandatory"}}

	} else {

		return ReturnResponse{HttpStatus: fiber.StatusOK, Rsp: GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC}}

	}
}

func NewInstanceTrxPinPerfonrmanceReport(c *fiber.Ctx, cfg *config.Cfg) *ApiPinPerformance {
	m := c.Queries()

	var dateSend time.Time
	if m["date_send"] == "" {
		dateSend = time.Now().In(cfg.TZ)
	} else {
		dateSend = external.ParseVendorDateSend(cfg.TZ, m["date_send"])
	}

	pinRequest, _ := strconv.Atoi(m["pin_request"])
	uniquePinRequest, _ := strconv.Atoi(m["pin_request_unique"])
	pinSent, _ := strconv.Atoi(m["pin_success"])
	pinFailed, _ := strconv.Atoi(m["pin_failed"])
	verifyRequest, _ := strconv.Atoi(m["pin_verify_request"])
	verifyRequestUnique, _ := strconv.Atoi(m["pin_verify_request_unique"])
	pinOK, _ := strconv.Atoi(m["pin_ok"])
	pinNotOK, _ := strconv.Atoi(m["pin_notok"])
	pinOkRatio, _ := strconv.Atoi(m["pin_ok_ratio"])
	chargedMO, _ := strconv.Atoi(m["charged_mo"])

	country := external.NormalizeCountries(m["country"])

	operator := strings.ToUpper(m["operator"])
	if operator == "" && m["telco"] != "" {
		operator = strings.ToUpper(m["telco"])
	}

	pin := ApiPinPerformance{
		Adnet:               strings.ToUpper(m["adnet"]),
		Country:             country,
		Company:             strings.ToUpper(m["company"]),
		Service:             strings.ToUpper(m["service"]),
		Operator:            operator,
		DateSend:            dateSend,
		PinRequest:          pinRequest,
		UniquePinRequest:    uniquePinRequest,
		PinSuccess:          pinSent,
		PinFailed:           pinFailed,
		PinVerifyRequest:       verifyRequest,
		PinVerifyRequestUnique: verifyRequestUnique,
		PinOK:                  pinOK,
		PinNotOK:               pinNotOK,
		PinOkRatio:             pinOkRatio,
		ChargedMO:              float64(chargedMO),
	}

	return &pin
}

func (t *ApiPinPerformance) ValidateParams(Logs *logrus.Logger) ReturnResponse {

	if t.Adnet == "" {

		return ReturnResponse{HttpStatus: fiber.StatusBadRequest, Rsp: GlobalResponse{Code: fiber.StatusBadRequest, Message: "Parameter Adnet is mandatory"}}

	} else if t.Country == "" {

		return ReturnResponse{HttpStatus: fiber.StatusBadRequest, Rsp: GlobalResponse{Code: fiber.StatusBadRequest, Message: "Parameter Country is mandatory"}}

	} else if t.Service == "" {

		return ReturnResponse{HttpStatus: fiber.StatusBadRequest, Rsp: GlobalResponse{Code: fiber.StatusBadRequest, Message: "Parameter Service is mandatory"}}

	} else if t.Operator == "" {

		return ReturnResponse{HttpStatus: fiber.StatusBadRequest, Rsp: GlobalResponse{Code: fiber.StatusBadRequest, Message: "Parameter Operator is mandatory"}}

	} else {

		return ReturnResponse{HttpStatus: fiber.StatusOK, Rsp: GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC}}

	}
}

func NewInstancePinPerformance(c *fiber.Ctx, cfg *config.Cfg) *ApiPinPerformance {
	m := c.Queries()

	toInt := func(key string) int {
		val, _ := strconv.Atoi(m[key])
		return val
	}

	toFloat := func(key string) float64 {
		val, _ := strconv.ParseFloat(m[key], 64)
		return val
	}

	pin := ApiPinPerformance{
		DateSend:               external.GetCurrentTime(cfg.TZ, time.RFC3339),
		Country:                m["country"],
		Company:                m["company"],
		Adnet:                  m["adnet"],
		Operator:               m["operator"],
		Service:                m["service"],
		PinRequest:             toInt("pin_request"),
		UniquePinRequest:       toInt("unique_pin_request"),
		PinSuccess:             toInt("pin_success"),
		PinFailed:              toInt("pin_failed"),
		PinVerifyRequest:       toInt("pin_verify_request"),
		PinVerifyRequestUnique: toInt("pin_verify_request_unique"),
		PinOK:                  toInt("pin_ok"),
		PinNotOK:               toInt("pin_not_ok"),
		PinOkRatio:             toInt("pin_ok_ratio"),
		CPA:                    toFloat("cpa"),
		CPAWaki:                toFloat("cpa_waki"),
		EstimatedARPU:          toFloat("estimated_arpu"),
		SBAF:                   toFloat("sbaf"),
		SAAF:                   toFloat("saaf"),
		ChargedMO:              toFloat("charged_mo"),
		SubsCR:                 toFloat("subs_cr"),
		AdnetCR:                toFloat("adnet_cr"),
		CAC:                    toFloat("cac"),
		PaidCAC:                toFloat("paid_cac"),
		TotalWakiAgencyFee:     toFloat("total_waki_agency_fee"),
		TotalSpending:          toFloat("total_spending"),
		// ClientType:          m["client_type"],
		// CampaignName: m["campaign_name"],
	}

	return &pin
}
