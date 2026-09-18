package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/config"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
)

func (h *IncomingHandler) DisplayDashboardReport(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	var SummaryDashboardReport entity.SummaryDashboardReport
	m := c.Queries()

	date_range := m["date_range"]
	date_before := m["date_before"]
	date_after := m["date_after"]
	country := m["country"]
	operator := m["operator"]
	client_type := m["client_type"]
	partner := m["partner"]
	service := m["service"]
	campaign_objective := m["campaign_objective"]

	allowedAdnets, _ := c.Locals("adnets").([]string)
	allowedCompanies, _ := c.Locals("companies").([]string)

	SummaryDashboardReport, _ = h.DS.GetReport(country, operator, client_type, partner, service, campaign_objective, date_range, date_before, date_after, allowedAdnets, allowedCompanies)
	Response := entity.ReturnResponse{
		HttpStatus: fiber.StatusOK,
		Rsp: entity.GlobalResponseWithDataTable{
			Draw:            0,
			Code:            fiber.StatusOK,
			Message:         config.OK_DESC,
			Data:            SummaryDashboardReport,
			RecordsTotal:    1,
			RecordsFiltered: 1,
		},
	}

	return c.Status(Response.HttpStatus).JSON(Response.Rsp)
}

func (h *IncomingHandler) DisplayDashboardTopCampaign(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	var TopCampaign []entity.TopCampaign

	m := c.Queries()
	date_range := m["date_range"]
	date_before := m["date_before"]
	date_after := m["date_after"]
	limit := m["limit"]
	order_type := m["order_type"] // BEST WORST
	order_by := m["order_by"]
	client_type := m["client_type"]
	country := m["country"]
	service := m["service"]

	allowedAdnets, _ := c.Locals("adnets").([]string)
	allowedCompanies, _ := c.Locals("companies").([]string)

	// m["type"] // WORST or BEST

	TopCampaign, _ = h.DS.GetCampaign(order_type, order_by, limit, client_type, date_range, date_before, date_after, country, service, allowedAdnets, allowedCompanies)

	Response := entity.ReturnResponse{
		HttpStatus: fiber.StatusOK,
		Rsp: entity.GlobalResponseWithDataTable{
			Draw:            0,
			Code:            fiber.StatusOK,
			Message:         config.OK_DESC,
			Data:            TopCampaign,
			RecordsTotal:    1,
			RecordsFiltered: 1,
		},
	}

	return c.Status(Response.HttpStatus).JSON(Response.Rsp)

}
func (h *IncomingHandler) DisplayDashboardData(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	/*
	   date_range: TODAY,YESTERDAY,LAST7DAY,LAST30DAY,THISMONTH,LASTMONTH,CUSTOMRANGE
	   if CUSTOMRANGE => date_before,date_after
	*/

	var SummaryDashboard entity.SummaryDashboardData

	m := c.Queries()

	date_range := m["date_range"]
	date_before := m["date_before"]
	date_after := m["date_after"]
	client_type := m["client_type"]
	country := m["country"]
	service := m["service"]

	allowedAdnets, _ := c.Locals("adnets").([]string)
	allowedCompanies, _ := c.Locals("companies").([]string)

	SummaryDashboard, _ = h.DS.GetDisplayDashboard(date_range, date_before, date_after, client_type, country, service, allowedAdnets, allowedCompanies)

	Response := entity.ReturnResponse{
		HttpStatus: fiber.StatusOK,
		Rsp: entity.GlobalResponseWithDataTable{
			Draw:            0,
			Code:            fiber.StatusOK,
			Message:         config.OK_DESC,
			Data:            SummaryDashboard,
			RecordsTotal:    1,
			RecordsFiltered: 1,
		},
	}

	return c.Status(Response.HttpStatus).JSON(Response.Rsp)
}

func (h *IncomingHandler) DisplayCountryStats(c *fiber.Ctx) error {
	m := c.Queries()
	country := m["country"]
	service := m["service"]
	allowedAdnets, _ := c.Locals("adnets").([]string)
	allowedCompanies, _ := c.Locals("companies").([]string)
	data, _ := h.DS.GetCountryStats(m["date_range"], m["date_before"], m["date_after"], country, service, allowedAdnets, allowedCompanies)
	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithDataTable{
		Draw: 0, Code: fiber.StatusOK, Message: config.OK_DESC,
		Data: data, RecordsTotal: 1, RecordsFiltered: 1,
	})
}

func (h *IncomingHandler) DisplayOpsStats(c *fiber.Ctx) error {
	m := c.Queries()
	date_range := m["date_range"]
	date_before := m["date_before"]
	date_after := m["date_after"]
	country := m["country"]
	service := m["service"]
	allowedAdnets, _ := c.Locals("adnets").([]string)
	allowedCompanies, _ := c.Locals("companies").([]string)
	data, _ := h.DS.GetOpsStats(date_range, date_before, date_after, country, service, allowedAdnets, allowedCompanies)
	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithDataTable{
		Draw: 0, Code: fiber.StatusOK, Message: config.OK_DESC,
		Data: data, RecordsTotal: 1, RecordsFiltered: 1,
	})
}

func (h *IncomingHandler) DisplayAlerts(c *fiber.Ctx) error {
	m := c.Queries()
	country := m["country"]
	service := m["service"]
	allowedAdnets, _ := c.Locals("adnets").([]string)
	allowedCompanies, _ := c.Locals("companies").([]string)
	data := h.DS.GetAlerts(country, service, allowedAdnets, allowedCompanies)
	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithDataTable{
		Draw: 0, Code: fiber.StatusOK, Message: config.OK_DESC,
		Data: data, RecordsTotal: 1, RecordsFiltered: 1,
	})
}

func (h *IncomingHandler) DisplayRollup(c *fiber.Ctx) error {
	m := c.Queries()
	country := m["country"]
	service := m["service"]
	allowedAdnets, _ := c.Locals("adnets").([]string)
	allowedCompanies, _ := c.Locals("companies").([]string)
	data, _ := h.DS.GetRollup(m["date_range"], m["date_before"], m["date_after"], m["client_type"], country, service, allowedAdnets, allowedCompanies)
	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithDataTable{
		Draw: 0, Code: fiber.StatusOK, Message: config.OK_DESC,
		Data: data, RecordsTotal: 1, RecordsFiltered: 1,
	})
}

func (h *IncomingHandler) DisplayAdnetStats(c *fiber.Ctx) error {
	m := c.Queries()
	country := m["country"]
	service := m["service"]
	allowedAdnets, _ := c.Locals("adnets").([]string)
	allowedCompanies, _ := c.Locals("companies").([]string)
	data, _ := h.DS.GetAdnetStats(m["date_range"], m["date_before"], m["date_after"], m["client_type"], country, service, allowedAdnets, allowedCompanies)
	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithDataTable{
		Draw: 0, Code: fiber.StatusOK, Message: config.OK_DESC,
		Data: data, RecordsTotal: 1, RecordsFiltered: 1,
	})
}

func (h *IncomingHandler) DisplayHeatmap(c *fiber.Ctx) error {
	m := c.Queries()
	country := m["country"]
	service := m["service"]
	allowedAdnets, _ := c.Locals("adnets").([]string)
	allowedCompanies, _ := c.Locals("companies").([]string)
	data, _ := h.DS.GetHeatmap(m["date_range"], m["date_before"], m["date_after"], country, service, allowedAdnets, allowedCompanies)
	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithDataTable{
		Draw: 0, Code: fiber.StatusOK, Message: config.OK_DESC,
		Data: data, RecordsTotal: 1, RecordsFiltered: 1,
	})
}

func (h *IncomingHandler) DisplayCampaignDaily(c *fiber.Ctx) error {
	m := c.Queries()
	data, _ := h.DS.GetCampaignDaily(m["campaign_id"], m["date_range"], m["date_before"], m["date_after"])
	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithDataTable{
		Draw: 0, Code: fiber.StatusOK, Message: config.OK_DESC,
		Data: data, RecordsTotal: 1, RecordsFiltered: 1,
	})
}

func (h *IncomingHandler) DisplayServiceDaily(c *fiber.Ctx) error {
	m := c.Queries()
	data, _ := h.DS.GetServiceDaily(m["country"], m["operator"], m["service"], m["date_range"], m["date_before"], m["date_after"])
	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithDataTable{
		Draw: 0, Code: fiber.StatusOK, Message: config.OK_DESC,
		Data: data, RecordsTotal: 1, RecordsFiltered: 1,
	})
}

func (h *IncomingHandler) DisplayFilterOptions(c *fiber.Ctx) error {
	allowedAdnets, _ := c.Locals("adnets").([]string)
	allowedCompanies, _ := c.Locals("companies").([]string)
	country := c.Query("country")
	data := h.DS.GetFilterOptions(country, allowedAdnets, allowedCompanies)
	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithDataTable{
		Draw: 0, Code: fiber.StatusOK, Message: config.OK_DESC,
		Data: data, RecordsTotal: 1, RecordsFiltered: 1,
	})
}
