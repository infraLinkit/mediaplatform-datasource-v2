package handler

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/config"
)

func (h *IncomingHandler) CreateBudgetIO(c *fiber.Ctx) error {
	var req entity.BudgetIORequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	if len(req.Data) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Data rows cannot be empty",
		})
	}

	for _, row := range req.Data {

		budget := entity.BudgetIO{
			CPName:            req.CPName,
			PICName:           req.PICName,
			ContactEmail:      req.ContactEmail,
			CPBusinessPICName: req.CPBusinessPICName,
			Signature:         req.Signature,
			SubmittedBy: 	   req.SubmittedBy,	

			IOID:               row.IOID,
			CampaignType:       row.CampaignType,
			Month:              row.Month,
			Country:            row.Country,
			CountryName:        row.CountryName,
			Continent:          row.Continent,
			CompanyGroupName:   row.CompanyGroupName,
			Company:            row.Company,
			Partner:            row.Partner,
			Service:            row.Service,
			TargetCAC:          row.TargetCAC,
			TargetROI:          row.TargetROI,
			MonthlyMOTarget:    row.MonthlyMOTarget,
			MonthlySpendTarget: row.MonthlySpendTarget,
		}

		if err := h.DS.DB.Create(&budget).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to save data",
				"error":   err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

func (h *IncomingHandler) DisplayBudgetIO(c *fiber.Ctx) error {

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
	
	fe := entity.DisplayBudgetIO{
		Continent:    m["continent"],
		Country:      m["country"],
		CampaignType: m["campaign_type"],
		Company:      m["company"],
		Partner:      m["partner"],
		Draw:         draw,
		Page:         page,
		PageSize:     pageSize,
		Action:       m["action"],
		DateRange:    m["date_range"],
		DateBefore:   m["date_before"],
		DateAfter:    m["date_after"],
		Reload:       m["reload"],
		OrderColumn:  m["order_column"],
		OrderDir:     m["order_dir"],
	}

	allowedCompanies, _ := c.Locals("companies").([]string)

	r := h.DisplayBudgetIOExtra(c, fe, allowedCompanies)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayBudgetIOExtra(c *fiber.Ctx, fe entity.DisplayBudgetIO, allowedCompanies []string) entity.ReturnResponse {
	// key := "temp_key_api_cpa_report_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")

	var (
		err        error
		total_data int64
		// isempty    bool
		budgetio []entity.BudgetIO
		// displaybudgetio []entity.BudgetIO
	)

	if fe.Action != "" || fe.Reload == "true" {
		budgetio, total_data, err = h.DS.GetDisplayBudgetIO(fe, allowedCompanies)
	} else {
		budgetio, total_data, err = h.DS.GetDisplayBudgetIO(fe, allowedCompanies)
	}

	if err == nil {

		if budgetio == nil {
			budgetio = []entity.BudgetIO{}
		}

		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithDataTable{
				Draw:            fe.Draw,
				Code:            fiber.StatusOK,
				Message:         config.OK_DESC,
				Data:            budgetio,
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

func (h *IncomingHandler) DisplayBudgetIOAll(c *fiber.Ctx) error {

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
	
	fe := entity.DisplayBudgetIO{
		Continent:    m["continent"],
		Country:      m["country"],
		CampaignType: m["campaign_type"],
		Company:      m["company"],
		Partner:      m["partner"],
		Draw:         draw,
		Page:         page,
		PageSize:     pageSize,
		Action:       m["action"],
		DateRange:    m["date_range"],
		DateBefore:   m["date_before"],
		DateAfter:    m["date_after"],
		Reload:       m["reload"],
		OrderColumn:  m["order_column"],
		OrderDir:     m["order_dir"],
	}


	r := h.DisplayBudgetIOExtraAll(c, fe)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayBudgetIOExtraAll(c *fiber.Ctx, fe entity.DisplayBudgetIO) entity.ReturnResponse {
	// key := "temp_key_api_cpa_report_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")

	var (
		err        error
		total_data int64
		// isempty    bool
		budgetio []entity.BudgetIO
		// displaybudgetio []entity.BudgetIO
	)

	if fe.Action != "" || fe.Reload == "true" {
		budgetio, total_data, err = h.DS.GetDisplayBudgetIOAll(fe)
	} else {
		budgetio, total_data, err = h.DS.GetDisplayBudgetIOAll(fe)
	}

	if err == nil {

		if budgetio == nil {
			budgetio = []entity.BudgetIO{}
		}

		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithDataTable{
				Draw:            fe.Draw,
				Code:            fiber.StatusOK,
				Message:         config.OK_DESC,
				Data:            budgetio,
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

func (h *IncomingHandler) DisplayBudgetIOApproved(c *fiber.Ctx) error {

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
	
	fe := entity.DisplayBudgetIO{
		Continent:    m["continent"],
		Country:      m["country"],
		CampaignType: m["campaign_type"],
		Company:      m["company"],
		Partner:      m["partner"],
		Draw:         draw,
		Page:         page,
		PageSize:     pageSize,
		Action:       m["action"],
		DateRange:    m["date_range"],
		DateBefore:   m["date_before"],
		DateAfter:    m["date_after"],
		Reload:       m["reload"],
		OrderColumn:  m["order_column"],
		OrderDir:     m["order_dir"],
	}

	allowedCompanies, _ := c.Locals("companies").([]string)

	r := h.DisplayBudgetIOExtra(c, fe, allowedCompanies)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayBudgetIOApprovedExtra(c *fiber.Ctx, fe entity.DisplayBudgetIO, allowedCompanies []string) entity.ReturnResponse {
	// key := "temp_key_api_cpa_report_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")

	var (
		err        error
		total_data int64
		// isempty    bool
		budgetio []entity.BudgetIO
		// displaybudgetio []entity.BudgetIO
	)

	if fe.Action != "" || fe.Reload == "true" {
		budgetio, total_data, err = h.DS.GetDisplayBudgetIOApproved(fe, allowedCompanies)
	} else {
		budgetio, total_data, err = h.DS.GetDisplayBudgetIOApproved(fe, allowedCompanies)
	}

	if err == nil {

		if budgetio == nil {
			budgetio = []entity.BudgetIO{}
		}

		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithDataTable{
				Draw:            fe.Draw,
				Code:            fiber.StatusOK,
				Message:         config.OK_DESC,
				Data:            budgetio,
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

func (h *IncomingHandler) DisplayBudgetIOApprovedAll(c *fiber.Ctx) error {

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
	
	fe := entity.DisplayBudgetIO{
		Continent:    m["continent"],
		Country:      m["country"],
		CampaignType: m["campaign_type"],
		Company:      m["company"],
		Partner:      m["partner"],
		Draw:         draw,
		Page:         page,
		PageSize:     pageSize,
		Keyword: 	  m["keyword"],
		Action:       m["action"],
		DateRange:    m["date_range"],
		DateBefore:   m["date_before"],
		DateAfter:    m["date_after"],
		Reload:       m["reload"],
		OrderColumn:  m["order_column"],
		OrderDir:     m["order_dir"],
	}


	r := h.DisplayBudgetIOExtraApproved(c, fe)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayBudgetIOExtraApproved(c *fiber.Ctx, fe entity.DisplayBudgetIO) entity.ReturnResponse {
	// key := "temp_key_api_cpa_report_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")

	var (
		err        error
		total_data int64
		// isempty    bool
		budgetio []entity.BudgetIO
		// displaybudgetio []entity.BudgetIO
	)

	if fe.Action != "" || fe.Reload == "true" {
		budgetio, total_data, err = h.DS.GetDisplayBudgetIOApprovedAll(fe)
	} else {
		budgetio, total_data, err = h.DS.GetDisplayBudgetIOApprovedAll(fe)
	}

	if err == nil {

		if budgetio == nil {
			budgetio = []entity.BudgetIO{}
		}

		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithDataTable{
				Draw:            fe.Draw,
				Code:            fiber.StatusOK,
				Message:         config.OK_DESC,
				Data:            budgetio,
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

func (h *IncomingHandler) DisplaySummaryBudgetIO(c *fiber.Ctx) error {
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

	fe := entity.DisplaySummaryBudgetIO{
		Continent:    m["continent"],
		Country:      m["country"],
		Company:      m["company"],
		Partner:      m["partner"],
		Draw:         draw,
		Page:         page,
		PageSize:     pageSize,
		Action:       m["action"],
		DateRange:    m["date_range"],
		DateBefore:   m["date_before"],
		DateAfter:    m["date_after"],
		Reload:       m["reload"],
		OrderColumn:  m["order_column"],
		OrderDir:     m["order_dir"],
	}

	r := h.DisplaySummaryBudgetIOExtra(c, fe)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplaySummaryBudgetIOExtra(c *fiber.Ctx, fe entity.DisplaySummaryBudgetIO) entity.ReturnResponse {
	var (
		err    error
		result []entity.SummaryBudgetIOAgg
	)

	result, _, err = h.DS.GetDisplaySummaryBudgetIO(fe)

	if err != nil {
		return entity.ReturnResponse{
			HttpStatus: fiber.StatusNotFound,
			Rsp: entity.GlobalResponse{
				Code:    fiber.StatusNotFound,
				Message: "empty",
			},
		}
	}

	if result == nil {
		result = []entity.SummaryBudgetIOAgg{}
	}

	rows := MapSummaryBudgetIOToReportRows(result)

	return entity.ReturnResponse{
		HttpStatus: fiber.StatusOK,
		Rsp: entity.GlobalResponseWithDataTable{
			Draw:            fe.Draw,
			Code:            fiber.StatusOK,
			Message:         config.OK_DESC,
			Data:            rows,
			RecordsTotal:    len(rows),
			RecordsFiltered: len(rows),
		},
	}
}

// MapSummaryBudgetIOToReportRows ports datasource's flat per-dimension IO
// report shape (week1-5 actual cost + MO breakdown vs target), replacing this
// repo's earlier continent-rollup mapping which had no datasource counterpart
// and couldn't represent MO-per-week at all.
func MapSummaryBudgetIOToReportRows(data []entity.SummaryBudgetIOAgg) []entity.IOReportRow {
	rows := make([]entity.IOReportRow, 0, len(data))
	for _, d := range data {
		recordedDay := 0
		now := time.Now()
		if d.Month == now.Format("2006-01") {
			recordedDay = now.Day()
		}

		rows = append(rows, entity.IOReportRow{
			BudgetIOID:  d.BudgetIOID,
			Region:      d.Continent,
			Country:     d.Country,
			Company:     d.Company,
			Partner:     d.Partner,
			Operator:    d.Operator,
			Channel:     d.Channel,
			ClientType:  d.ClientType,
			Service:     d.Service,
			Month:       d.Month,
			MOWeek1:     d.MOWeek1,
			MOWeek2:     d.MOWeek2,
			MOWeek3:     d.MOWeek3,
			MOWeek4:     d.MOWeek4,
			MOWeek5:     d.MOWeek5,
			CostWeek1:   d.ActualWeek1,
			CostWeek2:   d.ActualWeek2,
			CostWeek3:   d.ActualWeek3,
			CostWeek4:   d.ActualWeek4,
			CostWeek5:   d.ActualWeek5,
			IOTarget:    d.IOTarget,
			MOTarget:    d.MOTarget,
			TargetCAC:   d.TargetCAC,
			EstLTV:      d.LTV,
			EstROAS:     d.ROAS,
			ROI:         d.ROI,
			RecordedDay: recordedDay,
		})
	}
	return rows
}

func (h *IncomingHandler) UpdateSummaryBudgetIO(c *fiber.Ctx) error {
	var req entity.UpdateSummaryBudgetIORequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "NOK", "error": "invalid request body"})
	}
	if err := h.DS.UpdateSummaryBudgetIO(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "NOK", "error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "OK", "error": ""})
}

