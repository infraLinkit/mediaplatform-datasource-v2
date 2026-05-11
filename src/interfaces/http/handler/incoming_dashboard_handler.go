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

	allowedAdnets, _ := c.Locals("adnets").([]string)
	allowedCompanies, _ := c.Locals("companies").([]string)

	// m["type"] // WORST or BEST

	TopCampaign, _ = h.DS.GetCampaign(order_type, order_by, limit, date_range, date_before, date_after, allowedAdnets, allowedCompanies)

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

	allowedAdnets, _ := c.Locals("adnets").([]string)
	allowedCompanies, _ := c.Locals("companies").([]string)

	SummaryDashboard, _ = h.DS.GetDisplayDashboard(date_range, date_before, date_after, allowedAdnets, allowedCompanies)

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
