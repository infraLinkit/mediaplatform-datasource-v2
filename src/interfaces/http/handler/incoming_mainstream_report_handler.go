package handler

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/config"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
)

// const PAGESIZE int = 10

func (h *IncomingHandler) DisplayMainstreamReport(c *fiber.Ctx) error {
	c.Set("Content-type", "application/x-www-form-urlencoded")
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
		SummaryDate:   time.Time{},
		CampaignId:    m["campaign_id"],
		CampaignName:  m["campaign_name"],
		UrlServiceKey: m["url_service_key"],
		Country:       m["country"],
		ClientType:    m["client_type"],
		Company:       m["company"],
		Operator:      m["operator"],
		Partner:       m["partner"],
		Channel:       m["channel"],
		Agency:        m["agency"],
		Aggregator:    m["aggregator"],
		Adnets:        adnets,
		Service:       m["service"],
		DataBasedOn:   m["data_based_on"],
		Draw:          draw,
		Page:          page,
		PageSize:      pageSize,
		Action:        m["action"],
		DateRange:     m["date_range"],
		DateBefore:    m["date_before"],
		DateAfter:     m["date_after"],
		Reload:        m["reload"],
		OrderColumn:   m["order_column"],
		OrderDir:      m["order_dir"],
	}

	allowedCompanies, _ := c.Locals("companies").([]string)
	allowedAgencies, _ := c.Locals("agencies").([]string)

	r := h.DisplayMainstreamReportExtra(c, fe, allowedCompanies, allowedAgencies)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayMainstreamReportExtra(c *fiber.Ctx, fe entity.DisplayCPAReport, allowedCompanies []string, allowedAgencies []string) entity.ReturnResponse {

	var (
		err              error
		total_data       int64
		mainstreamreport []entity.SummaryCampaign
		TotalSummary     interface{}
	)

	if fe.Action != "" || fe.Reload == "true" {
		fmt.Println("-----", fe.Reload, "-----")
		mainstreamreport, total_data, TotalSummary, err = h.DS.GetDisplayMainstreamReport(fe, allowedCompanies, allowedAgencies)
	} else {
		mainstreamreport, total_data, TotalSummary, err = h.DS.GetDisplayMainstreamReport(fe, allowedCompanies, allowedAgencies)
	}

	if err == nil {

		if mainstreamreport == nil {
			mainstreamreport = []entity.SummaryCampaign{}
		}

		//TotalSummary := interface{}

		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithDataTable{
				Draw:            fe.Draw,
				Code:            fiber.StatusOK,
				Message:         config.OK_DESC,
				Data:            mainstreamreport,
				RecordsTotal:    int(total_data),
				RecordsFiltered: int(total_data),
				TotalSummary:    TotalSummary,
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
