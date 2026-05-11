package handler

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/config"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
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
        err      error
        budgetio []entity.SummaryBudgetIO
    )

    budgetio, _, err = h.DS.GetDisplaySummaryBudgetIO(fe)

    if err != nil {
        return entity.ReturnResponse{
            HttpStatus: fiber.StatusNotFound,
            Rsp: entity.GlobalResponse{
                Code:    fiber.StatusNotFound,
                Message: "empty",
            },
        }
    }

    if budgetio == nil {
        budgetio = []entity.SummaryBudgetIO{}
    }

    mapped := MapSummaryBudgetIOToContinentReport(budgetio)

    return entity.ReturnResponse{
        HttpStatus: fiber.StatusOK,
        Rsp: entity.GlobalResponseWithDataTable{
            Draw:    fe.Draw,
            Code:    fiber.StatusOK,
            Message: config.OK_DESC,
            Data:    mapped,
            RecordsTotal: len(mapped),
            RecordsFiltered: len(mapped),
        },
    }
}

func MapSummaryBudgetIOToContinentReport(data []entity.SummaryBudgetIO) []entity.ContinentReport {
	continentMap := make(map[string]map[string]map[string]*entity.SummaryBudgetIO)

	for _, d := range data {
		month := d.Month
		continent := d.Continent
		country := d.Country

		if continentMap[month] == nil {
			continentMap[month] = make(map[string]map[string]*entity.SummaryBudgetIO)
		}
		if continentMap[month][continent] == nil {
			continentMap[month][continent] = make(map[string]*entity.SummaryBudgetIO)
		}

		if existing, ok := continentMap[month][continent][country]; ok {
			existing.TotalMonthlySpendTarget += d.TotalMonthlySpendTarget
			existing.ActualWeek1 += d.ActualWeek1
			existing.ActualWeek2 += d.ActualWeek2
			existing.ActualWeek3 += d.ActualWeek3
			existing.ActualWeek4 += d.ActualWeek4
		} else {
			copy := d
			continentMap[month][continent][country] = &copy
		}
	}

	var reports []entity.ContinentReport

	for month, continents := range continentMap {
		for continent, countries := range continents {
			var contRep entity.ContinentReport
			contRep.Continent = continent
			contRep.Month = month
			contRep.TotalIOTargetContinent = 0
			contRep.TotalActualCostContinent = 0

			for _, c := range countries {
				totalActual := c.ActualWeek1 + c.ActualWeek2 + c.ActualWeek3 + c.ActualWeek4
				kpi := func(a float64) float64 {
					if c.TotalMonthlySpendTarget == 0 {
						return 0
					}
					return float64(math.Round(a / ((c.TotalMonthlySpendTarget / 30) * 7) * 100))
				}

				countryRep := entity.CountryReport{
					Country:                c.Country,
					ActualCostWeek1Country: float64(c.ActualWeek1),
					KPIWeek1Country:        kpi(c.ActualWeek1),
					ActualCostWeek2Country: float64(c.ActualWeek2),
					KPIWeek2Country:        kpi(c.ActualWeek2),
					ActualCostWeek3Country: float64(c.ActualWeek3),
					KPIWeek3Country:        kpi(c.ActualWeek3),
					ActualCostWeek4Country: float64(c.ActualWeek4),
					KPIWeek4Country:        kpi(c.ActualWeek4),
					TotalActualCostCountry: float64(totalActual),
					TotalIOTargetCountry:   float64(c.TotalMonthlySpendTarget),
					BudgetUsageCountry:     float64(math.Round(totalActual / c.TotalMonthlySpendTarget * 100)),
				}

				contRep.Countries = append(contRep.Countries, countryRep)

				contRep.ActualCostWeek1Continent += float64(c.ActualWeek1)
				contRep.ActualCostWeek2Continent += float64(c.ActualWeek2)
				contRep.ActualCostWeek3Continent += float64(c.ActualWeek3)
				contRep.ActualCostWeek4Continent += float64(c.ActualWeek4)
				contRep.TotalIOTargetContinent += float64(c.TotalMonthlySpendTarget)
			}

			kpiCont := func(a float64) float64 {
				if contRep.TotalIOTargetContinent == 0 {
					return 0
				}
				return float64(math.Round(float64(a) / ((float64(contRep.TotalIOTargetContinent) / 30) * 7) * 100))
			}

			contRep.KPIWeek1Continent = kpiCont(contRep.ActualCostWeek1Continent)
			contRep.KPIWeek2Continent = kpiCont(contRep.ActualCostWeek2Continent)
			contRep.KPIWeek3Continent = kpiCont(contRep.ActualCostWeek3Continent)
			contRep.KPIWeek4Continent = kpiCont(contRep.ActualCostWeek4Continent)

			contRep.TotalActualCostContinent = contRep.ActualCostWeek1Continent + contRep.ActualCostWeek2Continent +
				contRep.ActualCostWeek3Continent + contRep.ActualCostWeek4Continent

			if contRep.TotalIOTargetContinent != 0 {
				contRep.BudgetUsageContinent = float64(math.Round(float64(contRep.TotalActualCostContinent) / float64(contRep.TotalIOTargetContinent) * 100))
			}

			reports = append(reports, contRep)
		}
	}

	return reports
}

