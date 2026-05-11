package handler

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/config"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
)

func (h *IncomingHandler) DisplayGoogleTrafficReport(c *fiber.Ctx) error {
	c.Set("Content-type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	m := c.Queries()

	draw, _ := strconv.Atoi(m["draw"])
	page, _ := strconv.Atoi(m["page"])
	pageSize, err := strconv.Atoi(m["page_size"])
	if err != nil || pageSize == 0 {
		pageSize = PAGESIZE
	}

	fe := entity.DisplayGoogleTrafficReport{
		CampaignId:    m["campaign_id"],
		CampaignName:  m["campaign_name"],
		UrlServiceKey: m["url_service_key"],
		Country:       m["country"],
		Operator:      m["operator"],
		Partner:       m["partner"],
		Service:       m["service"],
		Company:       m["company"],
		Adnet:         m["adnet"],
		AdgroupID:     m["adgroup_id"],
		PeriodType:    m["period_type"],
		Month:         m["month"],
		Year:          m["year"],
		Week:          m["week"],
		DateFrom:      m["date_from"],
		DateTo:        m["date_to"],
		Action:        m["action"],
		Draw:          draw,
		Page:          page,
		PageSize:      pageSize,
	}

	r := h.DisplayGoogleTrafficReportExtra(c, fe)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayGoogleTrafficReportExtra(
	c *fiber.Ctx,
	fe entity.DisplayGoogleTrafficReport,
) entity.ReturnResponse {

	reportData, totalData, totalSummary, err := h.DS.GetDisplayGoogleTrafficReport(fe)

	if err != nil {
		fmt.Printf("[DisplayGoogleTrafficReportExtra] error: %v\n", err)
		return entity.ReturnResponse{
			HttpStatus: fiber.StatusNotFound,
			Rsp: entity.GlobalResponse{
				Code:    fiber.StatusNotFound,
				Message: err.Error(),
			},
		}
	}

	if reportData == nil {
		reportData = []entity.GoogleTrafficReportRow{}
	}

	return entity.ReturnResponse{
		HttpStatus: fiber.StatusOK,
		Rsp: entity.GlobalResponseWithDataTable{
			Draw:            fe.Draw,
			Code:            fiber.StatusOK,
			Message:         config.OK_DESC,
			Data:            reportData,
			RecordsTotal:    int(totalData),
			RecordsFiltered: int(totalData),
			TotalSummary:    totalSummary,
		},
	}
}