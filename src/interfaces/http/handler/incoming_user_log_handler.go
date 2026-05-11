package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/config"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/external"
)

func (h *IncomingHandler) DisplayUserLogList(c *fiber.Ctx) error {

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
		country_list []entity.DisplayUserLogList
	)

	// key := "temp_key_api_country_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")

	// need to add redis mechanism here
	country_list, total_data, errResponse = h.DS.GetUserLogList(fe)

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

func (h *IncomingHandler) DisplayUserLogHistory(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponseWithData{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid ID format",
			Data:    []entity.Menu{},
		})
	}

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
		country_list []entity.DisplayUserLogList
	)

	// key := "temp_key_api_country_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")

	// need to add redis mechanism here
	country_list, total_data, errResponse = h.DS.GetUserLogHistory(fe, id)

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

func (h *IncomingHandler) CreateUserLog(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	var audit entity.Audit

	if errForm := c.BodyParser(&audit); errForm != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   errForm.Error(),
		})
	}

	if errValidation := validate.Struct(audit); errValidation != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation error",
			"errors":  errValidation.Error(),
		})
	}

	if errCreate := h.DS.CreateUserLog(&audit); errCreate != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create log",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})

}

func (h *IncomingHandler) GetIpAddress(c *fiber.Ctx) error {
    ip := external.GetIpAddress(c)

    return c.JSON(entity.ReturnResponse{
        HttpStatus: fiber.StatusOK,
        Rsp: entity.GlobalResponseWithData{
            Data:    ip, // Pass the IP address here
            Code:    fiber.StatusOK,
            Message: config.OK_DESC,
        },
    })
}
