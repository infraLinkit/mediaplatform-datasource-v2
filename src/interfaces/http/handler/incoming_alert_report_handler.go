package handler

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/config"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/external"
)

const PAGESIZE_alert int = 4

func (h *IncomingHandler) DisplayAlertReportAll(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("applicaton/x-ww-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-1959-1")

	m := c.Queries()
	v := c.Params("v")
	page, _ := strconv.Atoi(c.Query("page"))
	PAGESIZE_alert := PAGESIZE_alert
	if pageSizeQuery := c.Query("page_size"); pageSizeQuery != "" {
		PAGESIZE_alert, _ = strconv.Atoi(pageSizeQuery)
	}
	fe := entity.DisplayAlertReport{
		Action:       m["action"],
		Country:      m["country"],
		Operator:     m["operator"],
		CampaignName: m["campaign_name"],
		Service:      m["service"],
		Page:         page,
		DateRange:    m["date_range"],
		DateBefore:   m["date_before"],
		DateAfter:    m["date_after"],
		PageSize:     PAGESIZE_alert,
		ExportData:   m["export_data"],
	}

	r := h.DisplayAlertReportAllExtra(c, fe, v)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayAlertReportAllExtra(c *fiber.Ctx, fe entity.DisplayAlertReport, v string) entity.ReturnResponse {
	key := "temp_key_api_alert_report_" + v + "_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")

	var (
		err                   error
		x                     int
		isempty               bool
		allalert              []entity.SummaryAll
		displayalertreportall []entity.SummaryAll
	)

	if fe.Action != "" {
		allalert, err = h.DS.GetAlertReportAll(fe, v)
	} else {
		if allalert, isempty = h.DS.RGetAlertReportAll(key, "$"); isempty {
			allalert, err = h.DS.GetAlertReportAll(fe, v)
			s, _ := json.Marshal(allalert)
			h.DS.SetData(key, "$", string(s))
			h.DS.SetExpireData(key, 60)
		}
	}
	if fe.ExportData != "" {
		if err == nil {
			pagesize := fe.PageSize
			if pagesize == 0 {
				pagesize = PAGESIZE_alert
			}
			if fe.Page >= 2 {
				x = pagesize * (fe.Page - 1)
			} else {
				x = 0
			}

			for i := x; i < len(allalert) && i < x+pagesize; i++ {
				displayalertreportall = append(displayalertreportall, allalert[i])
			}
		} else {
			displayalertreportall = allalert
		}

		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithDataTable{
				Data:            displayalertreportall,
				Draw:            fe.Page,
				Code:            fiber.StatusOK,
				Message:         config.OK_DESC,
				RecordsTotal:    int(len(allalert)),
				RecordsFiltered: int(len(allalert)),
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

func (h *IncomingHandler) UpdateStatusAlert(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("applicaton/x-ww-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-1959-1")

	m := c.Queries()
	v := c.Params("v")
	time := m["time"]
	ID := m["id"]
	Status, _ := strconv.ParseBool(m["status"])

	r := h.UpdateStatusAlertExtra(c, ID, Status, time, v)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) UpdateStatusAlertExtra(c *fiber.Ctx, ID string, Status bool, time, v string) entity.ReturnResponse {

	err := h.DS.UpdateStatusAlert(ID, Status, time, v)
	if err == nil {
		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponse{
				Code:    fiber.StatusOK,
				Message: "Status updated successfully",
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
