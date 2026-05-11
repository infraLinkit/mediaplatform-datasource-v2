package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/config"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
)

func (h *IncomingHandler) DisplayCPAMonitoring(c *fiber.Ctx) error {

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
	fe := entity.GlobalRequestFromDataTable{
		Page:     page,
		Action:   m["action"],
		Draw:     draw,
		PageSize: pageSize,
		Search:   m["search[value]"],
	}

	var (
		errResponse  error
		total_data   int64
		country_list []entity.Country
	)

	// key := "temp_key_api_country_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")

	// need to add redis mechanism here
	country_list, total_data, errResponse = h.DS.GetCountry(fe)

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
				Data:            country_list,
				Draw:            fe.Draw,
				RecordsTotal:    int(total_data),
				RecordsFiltered: int(total_data),
			},
		}

	}

	return c.Status(r.HttpStatus).JSON(r.Rsp)
}
