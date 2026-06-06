package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/config"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/external"
)

const PAGESIZE int = 10

func (h *IncomingHandler) DisplayPinReport(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	m := c.Queries()

	page, _ := strconv.Atoi(m["page"])
	pageSize, err := strconv.Atoi(m["page_size"])
	if err != nil {
		pageSize = PAGESIZE
	}
	draw, _ := strconv.Atoi(m["draw"])
	var adnets []string
	for k, v := range m {
		if strings.HasPrefix(k, "adnet[") {
			adnets = append(adnets, v)
		}
	}
	fe := entity.DisplayPinReport{
		DateSend:    time.Time{},
		CampaignId:  m["campaign_id"],
		Country:     m["country"],
		Company:     m["company"],
		Operator:    m["operator"],
		Partner:     m["partner"],
		Aggregator:  m["aggregator"],
		Adnets:      adnets,
		Service:     m["service"],
		Draw:        draw,
		Page:        page,
		PageSize:    pageSize,
		Action:      m["action"],
		DateRange:   m["date_range"],
		DateBefore:  m["date_before"],
		DateAfter:   m["date_after"],
		Reload:      m["reload"],
		OrderColumn: m["order_column"],
		OrderDir:    m["order_dir"],
	}

	r := h.DisplayPinReportExtra(c, fe)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayPinReportExtra(c *fiber.Ctx, fe entity.DisplayPinReport) entity.ReturnResponse {
	var (
		err        error
		total_data int64
		apireport  []entity.ApiPinReport
	)

	if fe.Action != "" || fe.Reload == "true" {
		fmt.Println("-----", fe.Reload, "-----")
		apireport, total_data, err = h.DS.GetDisplayPinReport(fe)
	} else {

		apireport, total_data, err = h.DS.GetDisplayPinReport(fe)
	}

	if err == nil {

		if apireport == nil {
			apireport = []entity.ApiPinReport{}
		}

		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithDataTable{
				Draw:            fe.Draw,
				Code:            fiber.StatusOK,
				Message:         config.OK_DESC,
				Data:            apireport,
				RecordsTotal:    int(total_data),
				RecordsFiltered: int(total_data),
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

func (h *IncomingHandler) EditPayoutAPIReport(c *fiber.Ctx) error {

	o := new(entity.ApiPinReport)

	if err := c.BodyParser(o); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	h.Logs.Debug(fmt.Sprintf("payload : %#v", o))

	if err := h.DS.EditPayoutAPIReport(*o); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).SendString("OK")
}

func (h *IncomingHandler) DisplayPinPerformanceReport(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	m := c.Queries()

	page, _ := strconv.Atoi(m["page"])
	pageSize, err := strconv.Atoi(m["page_size"])
	if err != nil {
		pageSize = 10
	}

	draw, _ := strconv.Atoi(m["draw"])
	var adnets []string
	for k, v := range m {
		if strings.HasPrefix(k, "adnet[") {
			adnets = append(adnets, v)
		}
	}
	fe := entity.DisplayPinPerformanceReport{
		Adnet:      m["adnet"],
		Country:    m["country"],
		Service:    m["service"],
		Operator:   m["operator"],
		Adnets:     adnets,
		DateRange:  m["date_range"],
		DateBefore: m["date_before"],
		DateAfter:  m["date_after"],
		Page:       page,
		Action:     m["action"],
		Draw:       draw,
		PageSize:   pageSize,
	}

	r := h.DisplayPinPerformanceReportExtra(c, fe)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayPinPerformanceReportExtra(c *fiber.Ctx, fe entity.DisplayPinPerformanceReport) entity.ReturnResponse {

	key := "temp_key_api_pin_performance_report_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")

	var (
		err                  error
		total_data           int64
		isempty              bool
		pinperformancereport []entity.ApiPinPerformance
	)

	if fe.Action != "" {
		pinperformancereport, total_data, err = h.DS.GetApiPinPerformanceReport(fe)
	} else {
		if pinperformancereport, isempty = h.DS.RGetApiPinPerformanceReport(key, "$"); isempty {

			pinperformancereport, total_data, err = h.DS.GetApiPinPerformanceReport(fe)

			s, _ := json.Marshal(pinperformancereport)

			h.DS.SetData(key, "$", string(s))
			h.DS.SetExpireData(key, 60)
		}
	}

	if err == nil {

		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithDataTable{
				Code:            fiber.StatusOK,
				Message:         config.OK_DESC,
				Data:            pinperformancereport,
				Draw:            fe.Draw,
				RecordsTotal:    int(total_data),
				RecordsFiltered: int(total_data),
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

func (h *IncomingHandler) EditCpaAPIPerformanceReport(c *fiber.Ctx) error {

	o := new(entity.ApiPinPerformance)

	if err := c.BodyParser(o); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	h.Logs.Debug(fmt.Sprintf("payload : %#v", o))

	if err := h.DS.EditCpaAPIPerformanceReport(*o); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).SendString("OK")
}

func (h *IncomingHandler) EditArpuAPIPerformanceReport(c *fiber.Ctx) error {

	o := new(entity.ApiPinPerformance)

	if err := c.BodyParser(o); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	o.Country = strings.ToUpper(o.Country)
	o.Operator = strings.ToUpper(o.Operator)
	o.Service = strings.ToUpper(o.Service)
	o.Adnet = strings.ToUpper(o.Adnet)

	h.Logs.Debug(fmt.Sprintf("payload : %#v", o))

	if err := h.DS.EditArpuAPIPerformanceReport(*o); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).SendString("OK")
}

func (h *IncomingHandler) DisplayConversionLogReport(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	m := c.Queries()

	page, errPage := strconv.Atoi(m["page"])
	pageSize, err := strconv.Atoi(m["page_size"])
	if err != nil {
		pageSize = 10
	}
	if errPage != nil {
		page = 1
	}

	draw, _ := strconv.Atoi(m["draw"])
	fe := entity.DisplayConversionLogReport{
		Adnet:          m["adnet"],
		Agency:         m["agency"],
		Country:        m["country"],
		Operator:       m["operator"],
		Pixel:          m["pixel"],
		CampaignType:   m["campaign_type"],
		StatusPostback: m["status_postback"],
		CampaignId:     m["campaign_id"],
		CampaignName:   m["campaign_name"],
		DateRange:      m["date_range"],
		DateStart:      m["date_start"],
		DateEnd:        m["date_end"],
		Page:           page,
		Action:         m["action"],
		Draw:           draw,
		PageSize:       pageSize,
		Order:          m["order"],
	}

	r := h.DisplayConversionLogReportExtra(c, fe)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayPerformanceReport(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	m := c.Queries()

	page, _ := strconv.Atoi(m["page"])
	pageSize, errRequest := strconv.Atoi(m["page_size"])
	if errRequest != nil {
		pageSize = 10
	}
	draw, _ := strconv.Atoi(m["draw"])
	params := entity.PerformaceReportParams{
		Page:         page,
		Action:       m["action"],
		Draw:         draw,
		PageSize:     pageSize,
		Country:      m["country"],
		Company:      m["company"],
		ClientType:   m["client_type"],
		Operator:     m["operator"],
		CampaignName: m["campaign_name"],
		CampaignType: m["campaign_type"],
		Publisher:    m["publisher"],
		Service:      m["service"],
		DateStart:    m["date_start"],
		DateEnd:      m["date_end"],
		DateRange:    m["date_range"],
		DateBefore:   m["date_before"],
		DateAfter:    m["date_after"],
	}

	var (
		errResponse             error
		total_data              int64
		performance_report_list []entity.PerformanceReport
	)

	// key := "temp_key_api_company_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")

	// need to add redis mechanism here

	performance_report_list, total_data, errResponse = h.DS.GetPerformanceReport(params)

	r := entity.ReturnResponse{
		HttpStatus: fiber.StatusNotFound,
		Rsp: entity.GlobalResponse{
			Code:    fiber.StatusNotFound,
			Message: "empty",
		},
	}

	if errResponse == nil {
		r = entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithDataTable{
				Code:            fiber.StatusOK,
				Message:         config.OK_DESC,
				Data:            performance_report_list,
				Draw:            params.Draw,
				RecordsTotal:    int(total_data),
				RecordsFiltered: int(total_data),
			},
		}

	}

	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayConversionLogReportExtra(c *fiber.Ctx, fe entity.DisplayConversionLogReport) entity.ReturnResponse {

	// key := "temp_key_api_conversion_log_report_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")

	var (
		err        error
		total_data int64
		// isempty               bool
		conversion_log_report []entity.PixelStorage
	)

	// if fe.Action != "" {
	conversion_log_report, total_data, err = h.DS.GetConversionLogReport(fe)
	// } else {
	// 	if conversion_log_report, isempty = h.DS.RGetConversionLogReport(key, "$"); isempty {

	// 		conversion_log_report, total_data, err = h.DS.GetConversionLogReport(fe)

	// 		s, _ := json.Marshal(conversion_log_report)

	// 		h.DS.SetData(key, "$", string(s))
	// 		h.DS.SetExpireData(key, 60)
	// 	}
	// }

	if err == nil {

		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithDataTable{
				Code:            fiber.StatusOK,
				Message:         config.OK_DESC,
				Data:            conversion_log_report,
				Draw:            fe.Draw,
				RecordsTotal:    int(total_data),
				RecordsFiltered: int(total_data),
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

func (h *IncomingHandler) DisplayCPAReport(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	m := c.Queries()

	page, _ := strconv.Atoi(m["page"])
	pageSize, err := strconv.Atoi(m["page_size"])
	if err != nil {
		pageSize = PAGESIZE
	}
	draw, _ := strconv.Atoi(m["draw"])
	var adnets []string
	for k, v := range m {
		if strings.HasPrefix(k, "adnet[") {
			adnets = append(adnets, v)
		}
	}
	fe := entity.DisplayCPAReport{
		SummaryDate:       time.Time{},
		CampaignId:        m["campaign_id"],
		CampaignName:      m["campaign_name"],
		Country:           m["country"],
		ClientType:        m["client_type"],
		Company:           m["company"],
		Operator:          m["operator"],
		Partner:           m["partner"],
		Aggregator:        m["aggregator"],
		Adnets:            adnets,
		Service:           m["service"],
		Draw:              draw,
		Page:              page,
		PageSize:          pageSize,
		Action:            m["action"],
		DateRange:         m["date_range"],
		DateBefore:        m["date_before"],
		DateAfter:         m["date_after"],
		Reload:            m["reload"],
		OrderColumn:       m["order_column"],
		OrderDir:          m["order_dir"],
		CampaignObjective: m["campaign_objective"],
	}

	allowedCompanies, _ := c.Locals("companies").([]string)
	allowedAdnets, _ := c.Locals("adnets").([]string)

	r := h.DisplayCPAReportExtra(c, fe, allowedCompanies, allowedAdnets)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayCPAReportExtra(c *fiber.Ctx, fe entity.DisplayCPAReport, allowedCompanies []string, allowedAdnets []string) entity.ReturnResponse {
	// key := "temp_key_api_cpa_report_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")

	var (
		err        error
		total_data int64
		// isempty    bool
		cpareport            []entity.SummaryCampaign
		TotalSummaryCampaign entity.TotalSummaryCampaign
		// displaycpareport []entity.SummaryCampaign
	)

	if fe.Action != "" || fe.Reload == "true" {
		cpareport, total_data, TotalSummaryCampaign, err = h.DS.GetDisplayCPAReport(fe, allowedCompanies, allowedAdnets)
	} else {
		cpareport, total_data, TotalSummaryCampaign, err = h.DS.GetDisplayCPAReport(fe, allowedCompanies, allowedAdnets)
	}

	if err == nil {

		if cpareport == nil {
			cpareport = []entity.SummaryCampaign{}
		}

		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithDataTable{
				Draw:            fe.Draw,
				Code:            fiber.StatusOK,
				Message:         config.OK_DESC,
				Data:            cpareport,
				TotalSummary:    TotalSummaryCampaign,
				RecordsTotal:    int(total_data),
				RecordsFiltered: int(total_data),
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

func (h *IncomingHandler) ExportCpaButton(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	m := c.Queries()

	page, _ := strconv.Atoi(m["page"])
	var adnets []string
	for k, v := range m {
		if strings.HasPrefix(k, "adnet[") {
			adnets = append(adnets, v)
		}
	}
	fe := entity.DisplayCPAReport{
		SummaryDate:  time.Time{},
		CampaignId:   m["campaign_id"],
		CampaignName: m["campaign_name"],
		Country:      m["country"],
		Operator:     m["operator"],
		Partner:      m["partner"],
		Aggregator:   m["aggregator"],
		Adnets:       adnets,
		Service:      m["service"],
		Page:         page,
		Action:       m["action"],
		DateRange:    m["date_range"],
		DateBefore:   m["date_before"],
		DateAfter:    m["date_after"],
		OrderColumn:  m["order_column"],
		OrderDir:     m["order_dir"],
	}

	export_cpa := m["export_cpa"]

	if export_cpa == "true" {

		r := h.ExportCpaReportExtraNoLimit(c, fe)
		return c.Status(r.HttpStatus).JSON(r.Rsp)
	}

	return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponse{
		Code:    fiber.StatusBadRequest,
		Message: config.BAD_REQUEST_DESC,
	})
}

func (h *IncomingHandler) ExportCpaReportExtraNoLimit(c *fiber.Ctx, fe entity.DisplayCPAReport) entity.ReturnResponse {
	key := "temp_key_api_cpa_report_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")

	var (
		err                  error
		cpareport            []entity.SummaryCampaign
		isempty              bool
		TotalSummaryCampaign entity.TotalSummaryCampaign
		// displaycpareport []entity.SummaryCampaign
	)

	allowedCompanies, _ := c.Locals("companies").([]string)
	allowedAdnets, _ := c.Locals("adnets").([]string)

	if fe.Action != "" {
		cpareport, _, TotalSummaryCampaign, err = h.DS.GetDisplayCPAReport(fe, allowedCompanies, allowedAdnets)
	} else {

		if cpareport, isempty = h.DS.RGetDisplayCPAReport(key, "$"); isempty {

			cpareport, _, TotalSummaryCampaign, err = h.DS.GetDisplayCPAReport(fe, allowedCompanies, allowedAdnets)

			s, _ := json.Marshal(cpareport)

			h.DS.SetData(key, "$", string(s))
			h.DS.SetExpireData(key, 60)
		}
	}

	if err == nil {
		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithData{
				Code:         fiber.StatusOK,
				Message:      config.OK_DESC,
				Data:         cpareport,
				TotalSummary: TotalSummaryCampaign,
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

func (h *IncomingHandler) DisplayCostReport(c *fiber.Ctx) error {
	c.Set("Content-type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")
 
	m := c.Queries()
 
	page, errPage := strconv.Atoi(m["page"])
	pageSize, err := strconv.Atoi(m["page_size"])
	if err != nil {
		pageSize = 10
	}
	if errPage != nil {
		page = 10
	}
 
	draw, _ := strconv.Atoi(m["draw"])
	v := c.Params("v")
 
	var adnets []string
	for k, val := range m {
		if strings.HasPrefix(k, "adnets[") {
			adnets = append(adnets, val)
		}
	}
 
	var adnetFilter []string
	for k, val := range m {
		if strings.HasPrefix(k, "adnet[") {
			adnetFilter = append(adnetFilter, val)
		}
	}
	if len(adnets) == 0 && len(adnetFilter) > 0 {
		adnets = adnetFilter
	}
 
	fromChannel := m["from_channel"] == "1"
 
	fe := entity.DisplayCostReport{
		Adnet:       m["adnet"],
		Adnets:      adnets,
		Country:     m["country"],
		Operator:    m["operator"],
		ChannelType: m["channel_type"],
		GroupBy:     m["group_by"],
		DataIndicator: m["data_indicator"],
		Page:        page,
		Action:      m["action"],
		DateRange:   m["date_range"],
		DateBefore:  m["date_before"],
		DateAfter:   m["date_after"],
		DataBasedOn: m["data_based_on"],
		PageSize:    pageSize,
		Draw:        draw,
		FromChannel: fromChannel,
	}
 
	allowedAdnets, _ := c.Locals("adnets").([]string)
 
	r := h.DisplayCostReportExtra(c, fe, v, allowedAdnets)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

 
func (h *IncomingHandler) DisplayCostReportExtra(
	c *fiber.Ctx,
	fe entity.DisplayCostReport,
	v string,
	allowedAdnets []string,
) entity.ReturnResponse {
 
	key       := "temp_key_api_cost_report_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")
	keydetail := "temp_key_api_cost_report_detail_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")
 
	var (
		err        error
		isempty    bool
		total_data int64
		costreport []entity.CostReport
	)
 
	if v == "list" {
		if fe.GroupBy == "country" {
			if fe.Action != "" {
				costreport, total_data, err = h.DS.GetDisplayCostReportByCountry(fe, allowedAdnets)
			} else {
				keyCountry := key + "_country"
				if costreport, isempty = h.DS.RGetDisplayCostReport(keyCountry, "$"); isempty {
					costreport, total_data, err = h.DS.GetDisplayCostReportByCountry(fe, allowedAdnets)
					s, _ := json.Marshal(costreport)
					h.DS.SetData(keyCountry, "$", string(s))
					h.DS.SetExpireData(keyCountry, 60)
				}
			}
		} else {
			if fe.Action != "" {
				costreport, total_data, err = h.DS.GetDisplayCostReport(fe, allowedAdnets)
			} else {
				if costreport, isempty = h.DS.RGetDisplayCostReport(key, "$"); isempty {
					costreport, total_data, err = h.DS.GetDisplayCostReport(fe, allowedAdnets)
					s, _ := json.Marshal(costreport)
					h.DS.SetData(key, "$", string(s))
					h.DS.SetExpireData(key, 60)
				}
			}
		}
	} else if v == "detail" {
		if fe.Action != "" {
			costreport, total_data, err = h.DS.GetDisplayCostReportDetail(fe, allowedAdnets)
		} else {
			if costreport, isempty = h.DS.RGetDisplayCostReportDetail(keydetail, "$"); isempty {
				costreport, total_data, err = h.DS.GetDisplayCostReportDetail(fe, allowedAdnets)
				s, _ := json.Marshal(costreport)
				h.DS.SetData(keydetail, "$", string(s))
				h.DS.SetExpireData(keydetail, 60)
			}
		}
	}
 
	if err == nil {
		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithDataTable{
				Draw:            fe.Draw,
				Code:            fiber.StatusOK,
				Message:         config.OK_DESC,
				Data:            costreport,
				RecordsTotal:    int(total_data),
				RecordsFiltered: int(total_data),
			},
		}
	}
 
	return entity.ReturnResponse{
		HttpStatus: fiber.StatusNotFound,
		Rsp: entity.GlobalResponse{
			Code:    fiber.StatusNotFound,
			Message: "empty",
		},
	}
}

func (h *IncomingHandler) ExportCostButton(c *fiber.Ctx) error {
	c.Set("Content-type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	m := c.Queries()

	page, _ := strconv.Atoi(m["page"])
	draw, _ := strconv.Atoi(m["draw"])
	fe := entity.DisplayCostReport{
		Adnet:       m["adnet"],
		Country:     m["country"],
		Operator:    m["operator"],
		ChannelType: m["channel_type"],
		GroupBy:     m["group_by"],
		Page:        page,
		Action:      m["action"],
		DateRange:   m["date_range"],
		DateBefore:  m["date_before"],
		DateAfter:   m["date_after"],
		DataBasedOn: m["data_based_on"],
		Draw:        draw,
	}
	export_cost := m["export_cost"]
	if export_cost == "true" {
		allowedAdnets, _ := c.Locals("adnets").([]string)
		r := h.ExportCostReportExtraNoLimit(c, fe, allowedAdnets)
		return c.Status(r.HttpStatus).JSON(r.Rsp)
	}

	return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponse{
		Code:    fiber.StatusBadRequest,
		Message: config.BAD_REQUEST_DESC,
	})
}

func (h *IncomingHandler) ExportCostReportExtraNoLimit(c *fiber.Ctx, fe entity.DisplayCostReport, allowedAdnets []string) entity.ReturnResponse {
	key := "temp_key_api_cost_report_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")

	var (
		err        error
		costreport []entity.CostReport
		isempty    bool
		total_data int64
	)

	if fe.GroupBy == "country" {
		if fe.Action != "" {
			costreport, total_data, err = h.DS.GetDisplayCostReportByCountry(fe, allowedAdnets)
		} else {
			keyCountry := key + "_country"
			if costreport, isempty = h.DS.RGetDisplayCostReport(keyCountry, "$"); isempty {
				costreport, total_data, err = h.DS.GetDisplayCostReportByCountry(fe, allowedAdnets)
				s, _ := json.Marshal(costreport)
				h.DS.SetData(keyCountry, "$", string(s))
				h.DS.SetExpireData(keyCountry, 60)
			}
		}
	} else {
		if fe.Action != "" {
			costreport, total_data, err = h.DS.GetDisplayCostReport(fe, allowedAdnets)
		} else {
			if costreport, isempty = h.DS.RGetDisplayCostReport(key, "$"); isempty {
				costreport, total_data, err = h.DS.GetDisplayCostReport(fe, allowedAdnets)
				s, _ := json.Marshal(costreport)
				h.DS.SetData(key, "$", string(s))
				h.DS.SetExpireData(key, 60)
			}
		}
	}

	if err == nil {
		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithTable{
				Code:            fiber.StatusOK,
				Message:         config.OK_DESC,
				Data:            costreport,
				Draw:            fe.Draw,
				RecordsTotal:    int(total_data),
				RecordsFiltered: int(total_data),
			},
		}
	}

	return entity.ReturnResponse{
		HttpStatus: fiber.StatusNotFound,
		Rsp: entity.GlobalResponse{
			Code:    fiber.StatusNotFound,
			Message: "empty",
		},
	}
}

func (h *IncomingHandler) ExportCostDetailButton(c *fiber.Ctx) error {
	c.Set("Content-type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	m := c.Queries()

	page, _ := strconv.Atoi(m["page"])
	draw, _ := strconv.Atoi(m["draw"])
	fe := entity.DisplayCostReport{
		Adnet:       m["adnet"],
		Country:     m["country"],
		Operator:    m["operator"],
		ChannelType: m["channel_type"],
		Page:        page,
		Action:      m["action"],
		DateRange:   m["date_range"],
		DateBefore:  m["date_before"],
		DateAfter:   m["date_after"],
		DataBasedOn: m["data_based_on"],
		Draw:        draw,
	}
	export_cost := m["export_cost"]
	if export_cost == "true" {
		r := h.ExportCostReportDetailExtraNoLimit(c, fe)
		return c.Status(r.HttpStatus).JSON(r.Rsp)
	}

	return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponse{
		Code:    fiber.StatusBadRequest,
		Message: config.BAD_REQUEST_DESC,
	})
}

func (h *IncomingHandler) ExportCostReportDetailExtraNoLimit(c *fiber.Ctx, fe entity.DisplayCostReport) entity.ReturnResponse {
	key := "temp_key_api_cost_report_detail_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")
	allowedAdnets, _ := c.Locals("adnets").([]string)

	var (
		err        error
		costreport []entity.CostReport
		isempty    bool
		total_data int64
	)

	if fe.Action != "" {
		costreport, total_data, err = h.DS.GetDisplayCostReportDetail(fe, allowedAdnets)
	} else {
		if costreport, isempty = h.DS.RGetDisplayCostReportDetail(key, "$"); isempty {
			costreport, total_data, err = h.DS.GetDisplayCostReportDetail(fe, allowedAdnets)
			s, _ := json.Marshal(costreport)
			h.DS.SetData(key, "$", string(s))
			h.DS.SetExpireData(key, 60)
		}
	}

	if err == nil {
		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithTable{
				Code:            fiber.StatusOK,
				Message:         config.OK_DESC,
				Data:            costreport,
				Draw:            fe.Draw,
				RecordsTotal:    int(total_data),
				RecordsFiltered: int(total_data),
			},
		}
	}

	return entity.ReturnResponse{
		HttpStatus: fiber.StatusNotFound,
		Rsp: entity.GlobalResponse{
			Code:    fiber.StatusNotFound,
			Message: "empty",
		},
	}
}

func (h *IncomingHandler) DisplayDefaultInput(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	gs, err := h.DS.GetDataConfig("global_setting", "$")
	if err != nil {
		h.Logs.Error(fmt.Sprintf("Failed to get global settings: %v", err))
		return c.Status(fiber.StatusInternalServerError).JSON(entity.GlobalResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to get global settings",
		})
	}

	// Convert string values to float64
	costPerConversion, _ := strconv.ParseFloat(gs.CPCR, 64)
	agencyFee, _ := strconv.ParseFloat(gs.AgencyFee, 64)
	technicalFee, _ := strconv.ParseFloat(gs.TechnicalFee, 64)
	targetDailyBudget, _ := strconv.ParseFloat(gs.TargetDailyBudget, 64)

	// Validasi jika error, jadi default agency 5, cost 0.06, tech fee 5
	if err != nil {
		costPerConversion = 0.06
		agencyFee = 5
		technicalFee = 5
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithData{
		Code:    fiber.StatusOK,
		Message: config.OK_DESC,
		Data: entity.DefaultInput{
			CostPerConversion: costPerConversion,
			AgencyFee:         agencyFee,
			TechnicalFee:      technicalFee,
			TargetDailyBudget: targetDailyBudget,
		},
	})
}

func (h *IncomingHandler) ResendData(c *fiber.Ctx) error {

	var ids []string
	total, _ := strconv.Atoi(c.FormValue("total"))

	for i := 0; i < total; i++ {
		id := strings.TrimSpace(c.FormValue("id[" + strconv.Itoa(i) + "]"))
		if id == "" {
			continue
		}
		ids = append(ids, id)
	}

	baseURL := h.Config.APILINKITDashboard

	reports, _ := h.DS.GetSummaryReportById(ids)
	errorCounter := 0

	for _, sc := range reports {
		channelGroup := external.ResolveChannelGroup(sc.Channel)
		q := url.Values{
			"date":           {sc.SummaryDate.Format("2006-01-02")},
			"campaign_id":    {sc.URLServiceKey},
			"publisher":      {sc.Adnet},
			"adnet":          {sc.Adnet},
			"operator":       {sc.Partner},
			"adn":            {sc.ShortCode},
			"client":         {sc.Partner},
			"aggregator":     {sc.Aggregator},
			"country":        {sc.Country},
			"service":        {sc.Service},
			"channel":        {channelGroup},
			"mo_received":    {strconv.Itoa(sc.MoReceived)},
			"mo_postback":    {strconv.Itoa(sc.Postback)},
			"total_mo":       {strconv.Itoa(sc.MoReceived)},
			"total_postback": {strconv.Itoa(sc.Postback)},
			"landing":        {strconv.Itoa(sc.Traffic)},
			"cr_mo_received": {strconv.FormatFloat(sc.CrMO, 'f', 2, 64)},
			"cr_mo_postback": {strconv.FormatFloat(sc.CrPostback, 'f', 2, 64)},
			"url_campaign":   {sc.URLAfter},
			"url_service":    {sc.URLBefore},
			"sbaf":           {strconv.FormatFloat(sc.SBAF, 'f', 2, 64)},
			"saaf":           {strconv.FormatFloat(sc.SAAF, 'f', 2, 64)},
			"spending":       {strconv.FormatFloat(sc.SAAF, 'f', 2, 64)},
			"campaign":       {sc.CampaignObjective},
			"payout":         {strconv.FormatFloat(sc.PO, 'f', 2, 64)},
			"price_per_mo":   {strconv.FormatFloat(sc.PricePerMO, 'f', 2, 64)},
		}

		fullURL := fmt.Sprintf("%s?%s", baseURL, q.Encode())
		message := `{"url":"` + fullURL + `"}`

		pubCtx, pubCancel := context.WithTimeout(c.UserContext(), time.Duration(h.Config.RabbitMQCtxTimeout)*time.Second)
		defer pubCancel()
		if err := h.RM.PublishWithRetry(pubCtx, "E_RESENDCAMPAIGNDATA", "Q_RESENDCAMPAIGNDATA", []byte(message), ""); err != nil {
			errorCounter++
			h.Logs.Debug(fmt.Sprintf("[x] Failed published: Data: %s ...", message))
		} else {
			h.Logs.Debug(fmt.Sprintf("[v] Published: Data: %s ...", message))
		}

	}

	if errorCounter > 0 {
		return c.Status(fiber.StatusOK).SendString(`{"status":"NOK","error":"Some data not published"}`)
	}

	return c.Status(fiber.StatusOK).SendString(`{"status":"OK","error":""}`)
}

func ResolveOperatorAlias(
	operator string,
	service string,
	country string,
	aliases []entity.OperatorAlias,
) string {

	operator = strings.ToLower(strings.TrimSpace(operator))
	service = strings.ToLower(strings.TrimSpace(service))
	country = strings.ToLower(strings.TrimSpace(country))

	countries := external.NormalizeCountry(country)
	for i := range countries {
		countries[i] = strings.ToLower(strings.TrimSpace(countries[i]))
	}

	for _, a := range aliases {
		if strings.ToUpper(a.Type) != "API" {
			continue
		}

		aliasOperator := strings.ToLower(strings.TrimSpace(a.Operator))
		aliasService := strings.ToLower(strings.TrimSpace(a.Service))
		aliasCountry := strings.ToLower(strings.TrimSpace(a.Country))

		if aliasOperator == operator &&
			aliasService != "" &&
			strings.Contains(service, aliasService) &&
			contains(countries, aliasCountry) {
			return strings.ToLower(strings.TrimSpace(a.Alias))
		}
	}

	for _, a := range aliases {
		if strings.ToUpper(a.Type) != "API" {
			continue
		}

		aliasOperator := strings.ToLower(strings.TrimSpace(a.Operator))
		aliasCountry := strings.ToLower(strings.TrimSpace(a.Country))

		if aliasOperator == operator &&
			a.Service == "" &&
			contains(countries, aliasCountry) {
			return strings.ToLower(strings.TrimSpace(a.Alias))
		}
	}

	return operator
}

func (h *IncomingHandler) ResendDataAPIReport(c *fiber.Ctx) error {

	var ids []string
	total, _ := strconv.Atoi(c.FormValue("total"))

	for i := 0; i < total; i++ {
		id := strings.TrimSpace(c.FormValue("id[" + strconv.Itoa(i) + "]"))
		if id == "" {
			continue
		}
		ids = append(ids, id)
	}

	baseURL := h.Config.APILINKITDashboard

	reports, err := h.DS.GetAPIReportById(ids)
	if err != nil {
		return c.Status(500).SendString(`{"status":"NOK","error":"failed load report"}`)
	}

	aliases, err := h.DS.GetOperatorAliases()
	if err != nil {
		return c.Status(500).SendString(`{"status":"NOK","error":"failed load operator alias"}`)
	}

	errorCounter := 0

	for _, sc := range reports {

		operatorAlias := ResolveOperatorAlias(
			sc.Operator,
			sc.Service,
			sc.Country,
			aliases,
		)

		q := url.Values{
			"date":           {sc.DateSend.Format("2006-01-02")},
			"publisher":      {strings.ToLower(sc.Adnet)},
			"adnet":          {strings.ToLower(sc.Adnet)},
			"operator":       {operatorAlias},
			"client":         {operatorAlias},
			"aggregator":     {""},
			"country":        {strings.ToLower(sc.Country)},
			"service":        {strings.ToLower(sc.Service)},
			"channel": 		  {"API"},
			"total_mo":       {strconv.Itoa(sc.TotalMO)},
			"total_postback": {strconv.Itoa(sc.TotalPostback)},
			"landing":        {""},
			"cr_mo_received": {""},
			"cr_mo_postback": {""},
			"url_campaign":   {""},
			"url_service":    {""},
			"sbaf":           {""},
			"spending":       {strconv.FormatFloat(sc.SAAF, 'f', 2, 64)},
			"payout":         {strconv.FormatFloat(sc.PayoutAdn, 'f', 2, 64)},
			"price_per_mo":   {strconv.FormatFloat(sc.PricePerMO, 'f', 2, 64)},
		}

		fullURL := fmt.Sprintf("%s?%s", baseURL, q.Encode())
		message := fmt.Sprintf(`{"url":"%s"}`, fullURL)

		pubCtx2, pubCancel2 := context.WithTimeout(c.UserContext(), time.Duration(h.Config.RabbitMQCtxTimeout)*time.Second)
		defer pubCancel2()
		if err := h.RM.PublishWithRetry(pubCtx2, "E_RESENDCAMPAIGNDATA", "Q_RESENDCAMPAIGNDATA", []byte(message), ""); err != nil {
			errorCounter++
			h.Logs.Debug(fmt.Sprintf("[x] Failed published: %s", message))
		} else {
			h.Logs.Debug(fmt.Sprintf("[v] Published: %s", message))
		}
	}

	if errorCounter > 0 {
		return c.Status(fiber.StatusOK).
			SendString(`{"status":"NOK","error":"Some data not published"}`)
	}

	return c.Status(fiber.StatusOK).
		SendString(`{"status":"OK","error":""}`)
}
