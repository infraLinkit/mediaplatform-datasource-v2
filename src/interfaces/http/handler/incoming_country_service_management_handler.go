package handler

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/config"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/external"
)

func (h *IncomingHandler) CreateEmail(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	var email entity.Email

	if errForm := c.BodyParser(&email); errForm != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if errValidation := validate.Struct(email); errValidation != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation error",
			"errors":  errValidation.Error(),
		})
	}

	if errCreate := h.DS.CreateEmail(&email); errCreate != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create country",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})

}

func (h *IncomingHandler) UpdateEmail(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var email entity.Email

	if formErr := c.BodyParser(&email); formErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   formErr.Error(),
		})
	}

	email.ID = id

	if err := h.DS.UpdateEmail(&email); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update country",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

func (h *IncomingHandler) DeleteEmail(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.DS.DeleteEmail(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete country",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

func (h *IncomingHandler) DisplayEmail(c *fiber.Ctx) error {
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
		country_list []entity.Email
	)

	// key := "temp_key_api_country_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")

	// need to add redis mechanism here
	country_list, total_data, errResponse = h.DS.GetEmail(fe)

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

func (h *IncomingHandler) DisplayEmailByID(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	email, err := h.DS.GetEmailByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get email",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithData{Code: fiber.StatusOK, Message: config.OK_DESC, Data: email})
}

func (h *IncomingHandler) CreateCountry(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	var country entity.Country

	if errForm := c.BodyParser(&country); errForm != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if errValidation := validate.Struct(country); errValidation != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation error",
			"errors":  errValidation.Error(),
		})
	}

	if errCreate := h.DS.CreateCountry(&country); errCreate != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create country",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})

}

func (h *IncomingHandler) UpdateCountry(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var country entity.Country

	if formErr := c.BodyParser(&country); formErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   formErr.Error(),
		})
	}

	country.ID = uint(id)

	if err := h.DS.UpdateCountry(&country); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update country",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

func (h *IncomingHandler) DeleteCountry(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.DS.DeleteCountry(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete country",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

func (h *IncomingHandler) DisplayCountry(c *fiber.Ctx) error {

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

func (h *IncomingHandler) DisplayCountryInfo(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	country, _ := h.DS.GetCountryByCode(c.Params("code"))

	fmt.Println(country)
	return c.Status(200).JSON(country)

}

func (h *IncomingHandler) CreateCompany(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	var company entity.Company

	if errForm := c.BodyParser(&company); errForm != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if errValidation := validate.Struct(company); errValidation != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation error",
			"errors":  errValidation.Error(),
		})
	}

	if errCreate := h.DS.CreateCompany(&company); errCreate != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create company",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})

}

func (h *IncomingHandler) UpdateCompany(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var company entity.Company

	if formErr := c.BodyParser(&company); formErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   formErr.Error(),
		})
	}

	company.ID = uint(id)

	if err := h.DS.UpdateCompany(&company); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update company",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

func (h *IncomingHandler) DeleteCompany(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.DS.DeleteCompany(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete company",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

func (h *IncomingHandler) DisplayContinent(c *fiber.Ctx) error {

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
		errResponse    error
		total_data     int64
		continent_list []entity.Continent
	)

	// key := "temp_key_api_continent_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")

	// need to add redis mechanism here
	continent_list, total_data, errResponse = h.DS.GetContinent(fe)

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
				Data:            continent_list,
				Draw:            fe.Draw,
				RecordsTotal:    int(total_data),
				RecordsFiltered: int(total_data),
			},
		}

	}

	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayCompany(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	m := c.Queries()

	page, _ := strconv.Atoi(m["page"])
	pageSize, errRequest := strconv.Atoi(m["page_size"])
	if errRequest != nil {
		pageSize = 10
	}
	orderColumn := m["order_column"]
	orderDir := m["order_dir"]

	draw, _ := strconv.Atoi(m["draw"])
	fe := entity.GlobalRequestFromDataTableCompany{
		Page:        page,
		Action:      m["action"],
		Draw:        draw,
		PageSize:    pageSize,
		Search:      m["search[value]"],
		OrderColumn: orderColumn,
		OrderDir:    orderDir,
	}

	var (
		errResponse  error
		total_data   int64
		company_list []entity.Company
	)

	// key := "temp_key_api_company_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")

	// need to add redis mechanism here
	company_list, total_data, errResponse = h.DS.GetCompany(fe)

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
				Data:            company_list,
				Draw:            fe.Draw,
				RecordsTotal:    int(total_data),
				RecordsFiltered: int(total_data),
			},
		}

	}

	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) CreateDomain(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	var domain entity.Domain

	if errForm := c.BodyParser(&domain); errForm != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if errValidation := validate.Struct(domain); errValidation != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation error",
			"errors":  errValidation.Error(),
		})
	}

	if errCreate := h.DS.CreateDomain(&domain); errCreate != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create domain",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})

}

func (h *IncomingHandler) UpdateDomain(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var domain entity.Domain

	if formErr := c.BodyParser(&domain); formErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   formErr.Error(),
		})
	}

	domain.ID = uint(id)

	if err := h.DS.UpdateDomain(&domain); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update domain",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

func (h *IncomingHandler) DeleteDomain(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.DS.DeleteDomain(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete domain",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

func (h *IncomingHandler) DisplayDomain(c *fiber.Ctx) error {

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
		errResponse error
		total_data  int64
		domain_list []entity.Domain
	)

	// key := "temp_key_api_domain_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")

	// need to add redis mechanism here
	domain_list, total_data, errResponse = h.DS.GetDomain(fe)

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
				Data:            domain_list,
				Draw:            fe.Draw,
				RecordsTotal:    int(total_data),
				RecordsFiltered: int(total_data),
			},
		}

	}

	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) CreateOperator(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	var operator entity.Operator

	if errForm := c.BodyParser(&operator); errForm != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	operator.Lastupdate = time.Now()

	if errValidation := validate.Struct(operator); errValidation != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation error",
			"errors":  errValidation.Error(),
		})
	}

	if errCreate := h.DS.CreateOperator(&operator); errCreate != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create operator",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})

}

func (h *IncomingHandler) UpdateOperator(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var operator entity.Operator

	if formErr := c.BodyParser(&operator); formErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   formErr.Error(),
		})
	}

	operator.ID = uint(id)

	if err := h.DS.UpdateOperator(&operator); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update operator",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

func (h *IncomingHandler) DeleteOperator(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.DS.DeleteOperator(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete operator",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

func (h *IncomingHandler) DisplayOperator(c *fiber.Ctx) error {

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
		errResponse   error
		total_data    int64
		operator_list []entity.Operator
	)

	// key := "temp_key_api_operator_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")

	// need to add redis mechanism here
	operator_list, total_data, errResponse = h.DS.GetOperator(fe)

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
				Data:            operator_list,
				Draw:            fe.Draw,
				RecordsTotal:    int(total_data),
				RecordsFiltered: int(total_data),
			},
		}

	}

	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) CreatePartner(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	var partner entity.Partner

	if errForm := c.BodyParser(&partner); errForm != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   errForm.Error(),
		})
	}

	partner.Lastupdate = time.Now()

	if errValidation := validate.Struct(partner); errValidation != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation error",
			"errors":  errValidation.Error(),
		})
	}
	fmt.Printf("Parsed Partner: %+v\n", partner)

	if errCreate := h.DS.CreatePartner(&partner); errCreate != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create partner",
			"error":   errCreate.Error(),
		})
	}

	filePath := "/app_data/data/mediaplatform/new_partner/partner.txt"

	if err := external.AppendPartnerToTxt(filePath, strings.ToUpper(partner.Name)); err != nil {
		h.Logs.Error("failed append partner file: ", err)
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{
		Code:    fiber.StatusOK,
		Message: config.OK_DESC,
	})
}

func (h *IncomingHandler) UpdatePartner(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var partner entity.Partner

	if formErr := c.BodyParser(&partner); formErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   formErr.Error(),
		})
	}

	partner.ID = uint(id)

	if err := h.DS.UpdatePartner(&partner); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update partner",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

func (h *IncomingHandler) DeletePartner(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid partner id",
		})
	}

	// partner, err := h.DS.GetPartnerByID(uint(id))
	// if err != nil {
	// 	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
	// 		"message": "Partner not found",
	// 	})
	// }

	if err := h.DS.DeletePartner(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete partner",
		})
	}

	// cmd := exec.Command(
	// 	"/bin/bash",
	// 	"/app/mediaplatform/resources/dockerver/yaml/staging/bin/delete_partner_ratio.sh",
	// 	"staging",
	// 	strings.ToUpper(partner.Name),
	// )

	// if err := cmd.Run(); err != nil {
	// 	h.Logs.Error("delete partner infra failed:", err)
	// }

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{
		Code:    fiber.StatusOK,
		Message: config.OK_DESC,
	})
}

func (h *IncomingHandler) DisplayPartner(c *fiber.Ctx) error {

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
		partner_list []entity.Partner
	)

	// key := "temp_key_api_partner_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")

	// need to add redis mechanism here
	partner_list, total_data, errResponse = h.DS.GetPartner(fe)

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
				Data:            partner_list,
				Draw:            fe.Draw,
				RecordsTotal:    int(total_data),
				RecordsFiltered: int(total_data),
			},
		}

	}

	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) CreateService(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	var service entity.Service

	if errForm := c.BodyParser(&service); errForm != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if errValidation := validate.Struct(service); errValidation != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation error",
			"errors":  errValidation.Error(),
		})
	}

	if errCreate := h.DS.CreateService(&service); errCreate != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create service",
			"errors":  errCreate.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})

}

func (h *IncomingHandler) UpdateService(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var service entity.Service

	if formErr := c.BodyParser(&service); formErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   formErr.Error(),
		})
	}

	service.ID = uint(id)

	if err := h.DS.UpdateService(&service); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update service",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

func (h *IncomingHandler) DeleteService(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.DS.DeleteService(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete service",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

func (h *IncomingHandler) DisplayService(c *fiber.Ctx) error {

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
		service_list []entity.Service
	)

	// key := "temp_key_api_service_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")

	// need to add redis mechanism here
	service_list, total_data, errResponse = h.DS.GetService(fe)

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
				Data:            service_list,
				Draw:            fe.Draw,
				RecordsTotal:    int(total_data),
				RecordsFiltered: int(total_data),
			},
		}

	}

	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) CreateAdnetList(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	var adnet_list entity.AdnetList

	if errForm := c.BodyParser(&adnet_list); errForm != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   errForm.Error(),
		})
	}

	if errValidation := validate.Struct(adnet_list); errValidation != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation error",
			"errors":  errValidation.Error(),
		})
	}

	if errCreate := h.DS.CreateAdnetList(&adnet_list); errCreate != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create adnet_list",
			"error":   errCreate.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})

}

func (h *IncomingHandler) UpdateAdnetList(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var adnet_list entity.AdnetList

	if formErr := c.BodyParser(&adnet_list); formErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   formErr.Error(),
		})
	}

	adnet_list.ID = uint(id)

	if err := h.DS.UpdateAdnetList(&adnet_list); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update adnet_list",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

func (h *IncomingHandler) DeleteAdnetList(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.DS.DeleteAdnetList(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete adnet_list",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

func (h *IncomingHandler) DisplayAdnetList(c *fiber.Ctx) error {

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
		errResponse     error
		total_data      int64
		adnet_list_list []entity.AdnetList
	)

	// key := "temp_key_api_adnet_list_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")

	// need to add redis mechanism here
	adnet_list_list, total_data, errResponse = h.DS.GetAdnetList(fe)

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
				Data:            adnet_list_list,
				Draw:            fe.Draw,
				RecordsTotal:    int(total_data),
				RecordsFiltered: int(total_data),
			},
		}

	}

	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayAPIAdnetList(c *fiber.Ctx) error {

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
		errResponse    error
		total_data     int64
		api_adnet_list []entity.ApiPinReport
	)

	// key := "temp_key_api_adnet_list_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")

	// need to add redis mechanism here
	api_adnet_list, total_data, errResponse = h.DS.GetAPIAdnetList(fe)

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
				Data:            api_adnet_list,
				Draw:            fe.Draw,
				RecordsTotal:    int(total_data),
				RecordsFiltered: int(total_data),
			},
		}

	}

	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayAPIOperatorList(c *fiber.Ctx) error {

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
		errResponse     error
		total_data      int64
		api_adnet_list []entity.ApiPinReport
	)

	// key := "temp_key_api_adnet_list_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")

	// need to add redis mechanism here
	api_adnet_list, total_data, errResponse = h.DS.GetAPIOperatorList(fe)

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
				Data:            api_adnet_list,
				Draw:            fe.Draw,
				RecordsTotal:    int(total_data),
				RecordsFiltered: int(total_data),
			},
		}

	}

	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayAPIServiceList(c *fiber.Ctx) error {

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
		errResponse     error
		total_data      int64
		api_adnet_list []entity.ApiPinReport
	)

	// key := "temp_key_api_adnet_list_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")

	// need to add redis mechanism here
	api_adnet_list, total_data, errResponse = h.DS.GetAPIServiceList(fe)

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
				Data:            api_adnet_list,
				Draw:            fe.Draw,
				RecordsTotal:    int(total_data),
				RecordsFiltered: int(total_data),
			},
		}

	}

	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) CreateAgency(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	var agency entity.Agency

	if errForm := c.BodyParser(&agency); errForm != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if errValidation := validate.Struct(agency); errValidation != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation error",
			"errors":  errValidation.Error(),
		})
	}

	if errCreate := h.DS.CreateAgency(&agency); errCreate != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create agency",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})

}

func (h *IncomingHandler) UpdateAgency(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var agency entity.Agency

	if formErr := c.BodyParser(&agency); formErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   formErr.Error(),
		})
	}

	agency.ID = uint(id)

	if err := h.DS.UpdateAgency(&agency); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update agency",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

func (h *IncomingHandler) DeleteAgency(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.DS.DeleteAgency(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete agency",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

func (h *IncomingHandler) DisplayAgency(c *fiber.Ctx) error {

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
		errResponse error
		total_data  int64
		agency_list []entity.Agency
	)

	// key := "temp_key_api_agency_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")

	// need to add redis mechanism here
	agency_list, total_data, errResponse = h.DS.GetAgency(fe)

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
				Data:            agency_list,
				Draw:            fe.Draw,
				RecordsTotal:    int(total_data),
				RecordsFiltered: int(total_data),
			},
		}

	}

	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) CreateChannel(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	var channel entity.Channel

	if errForm := c.BodyParser(&channel); errForm != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if errValidation := validate.Struct(channel); errValidation != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation error",
			"errors":  errValidation.Error(),
		})
	}

	if errCreate := h.DS.CreateChannel(&channel); errCreate != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create channel",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})

}

func (h *IncomingHandler) UpdateChannel(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var channel entity.Channel

	if formErr := c.BodyParser(&channel); formErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   formErr.Error(),
		})
	}

	channel.ID = uint(id)
	fmt.Println("CHANNEL APIKEY: ", channel.ApiKey)
	if err := h.DS.UpdateChannel(&channel); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update channel",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

func (h *IncomingHandler) DeleteChannel(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.DS.DeleteChannel(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete channel",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

func (h *IncomingHandler) DisplayChannel(c *fiber.Ctx) error {

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
		channel_list []entity.Channel
	)

	// key := "temp_key_api_channel_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")

	// need to add redis mechanism here
	channel_list, total_data, errResponse = h.DS.GetChannel(fe)

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
				Data:            channel_list,
				Draw:            fe.Draw,
				RecordsTotal:    int(total_data),
				RecordsFiltered: int(total_data),
			},
		}

	}

	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func formatDomainKey(domain string) string {
	domain = strings.ToLower(domain)
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "www.")
	domain = strings.TrimSuffix(domain, "/")
	return strings.ReplaceAll(domain, ".", "_")
}

func (h *IncomingHandler) CreateMainstreamGroup(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	var mainstreamGroup entity.MainstreamGroup

	if errForm := c.BodyParser(&mainstreamGroup); errForm != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if errValidation := validate.Struct(mainstreamGroup); errValidation != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation error",
			"errors":  errValidation.Error(),
		})
	}

	if errCreate := h.DS.CreateMainstreamGroup(&mainstreamGroup); errCreate != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create mainstreamGroup",
		})
	}

	redisKey := "domain_services"
	path := "$"
	domainKey := formatDomainKey(mainstreamGroup.UniqueDomain)
	renderName := strings.ToLower(mainstreamGroup.DomainService)

	cfgDomain, _ := h.DS.GetDomainServices(redisKey, path)

	newItem := entity.DomainServices{
		Domain:           domainKey,
		Render:           renderName,
		Country:          mainstreamGroup.Country,
		Operator:         mainstreamGroup.Operator,
		Service:          mainstreamGroup.Service,
		Company:          mainstreamGroup.Company,
		CompanyLegalName: mainstreamGroup.CompanyLegalName,
		CompanyAddress:   mainstreamGroup.CompanyAddress,
		CompanyEmail:     mainstreamGroup.CompanyEmail,
		CompanyPhone:     mainstreamGroup.CompanyPhone,
		ServiceCurrency:  mainstreamGroup.ServiceCurrency,
		ServicePrice:     mainstreamGroup.ServicePrice,
		PortalURL:        mainstreamGroup.PortalURL,
	}

	cfgDomain = append(cfgDomain, newItem)
	cfgData := map[string]interface{}{
		"data": cfgDomain,
	}

	jsonData, _ := json.Marshal(cfgData)
	h.DS.SetData(redisKey, path, string(jsonData))

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})

}

func (h *IncomingHandler) UpdateMainstreamGroup(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var mainstreamGroup entity.MainstreamGroup

	if formErr := c.BodyParser(&mainstreamGroup); formErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   formErr.Error(),
		})
	}

	mainstreamGroup.ID = uint(id)

	if err := h.DS.UpdateMainstreamGroup(&mainstreamGroup); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update mainstreamGroup",
		})
	}

	redisKey := "domain_services"
	path := "$"
	domainKey := formatDomainKey(mainstreamGroup.UniqueDomain)
	renderName := strings.ToLower(mainstreamGroup.DomainService)

	cfgDomain, _ := h.DS.GetDomainServices(redisKey, path)

	updated := false

	for i, v := range cfgDomain {
		if v.Domain == domainKey {
			cfgDomain[i] = entity.DomainServices{
				Domain:           domainKey,
				Render:           renderName,
				Country:          mainstreamGroup.Country,
				Operator:         mainstreamGroup.Operator,
				Service:          mainstreamGroup.Service,
				Company:          mainstreamGroup.Company,
				CompanyLegalName: mainstreamGroup.CompanyLegalName,
				CompanyAddress:   mainstreamGroup.CompanyAddress,
				CompanyEmail:     mainstreamGroup.CompanyEmail,
				CompanyPhone:     mainstreamGroup.CompanyPhone,
				ServiceCurrency:  mainstreamGroup.ServiceCurrency,
				ServicePrice:     mainstreamGroup.ServicePrice,
				PortalURL:        mainstreamGroup.PortalURL,
			}
			updated = true
			break
		}
	}

	if !updated {
		cfgDomain = append(cfgDomain, entity.DomainServices{
			Domain:           domainKey,
			Render:           renderName,
			Country:          mainstreamGroup.Country,
			Operator:         mainstreamGroup.Operator,
			Service:          mainstreamGroup.Service,
			Company:          mainstreamGroup.Company,
			CompanyLegalName: mainstreamGroup.CompanyLegalName,
			CompanyAddress:   mainstreamGroup.CompanyAddress,
			CompanyEmail:     mainstreamGroup.CompanyEmail,
			CompanyPhone:     mainstreamGroup.CompanyPhone,
			ServiceCurrency:  mainstreamGroup.ServiceCurrency,
			ServicePrice:     mainstreamGroup.ServicePrice,
			PortalURL:        mainstreamGroup.PortalURL,
		})
	}

	cfgData := map[string]interface{}{
		"data": cfgDomain,
	}
	jsonData, _ := json.Marshal(cfgData)
	h.DS.SetData(redisKey, path, string(jsonData))

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

func (h *IncomingHandler) DeleteMainstreamGroup(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	mainstreamGroup, err := h.DS.FindMainstreamGroupByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "MainstreamGroup not found",
		})
	}

	if err := h.DS.DeleteMainstreamGroup(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete mainstreamGroup",
		})
	}

	redisKey := "domain_services"
	path := "$"
	domainKey := formatDomainKey(mainstreamGroup.UniqueDomain)

	cfgDomain, _ := h.DS.GetDomainServices(redisKey, path)
	if len(cfgDomain) > 0 {
		filtered := make([]entity.DomainServices, 0)
		for _, v := range cfgDomain {
			if v.Domain != domainKey {
				filtered = append(filtered, v)
			}
		}

		cfgData := map[string]interface{}{
			"data": filtered,
		}
		jsonData, _ := json.Marshal(cfgData)
		h.DS.SetData(redisKey, path, string(jsonData))
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

func (h *IncomingHandler) DisplayMainstreamGroup(c *fiber.Ctx) error {

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
		errResponse          error
		total_data           int64
		mainstreamGroup_list []entity.MainstreamGroup
	)

	// key := "temp_key_api_mainstreamGroup_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")

	// need to add redis mechanism here
	mainstreamGroup_list, total_data, errResponse = h.DS.GetMainstreamGroup(fe)

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
				Data:            mainstreamGroup_list,
				Draw:            fe.Draw,
				RecordsTotal:    int(total_data),
				RecordsFiltered: int(total_data),
			},
		}

	}

	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayDomainService(c *fiber.Ctx) error {

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
		errResponse         error
		total_data          int64
		domain_service_list []entity.DomainService
	)

	// key := "temp_key_api_channel_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")

	// need to add redis mechanism here
	domain_service_list, total_data, errResponse = h.DS.GetDomainService(fe)

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
				Data:            domain_service_list,
				Draw:            fe.Draw,
				RecordsTotal:    int(total_data),
				RecordsFiltered: int(total_data),
			},
		}

	}

	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) UpdateDSPAdnetStatus(c *fiber.Ctx) error {

	var params struct {
		ID     string "id"
		Status string "status"
	}

	var adnet_list entity.AdnetList

	if formErr := c.BodyParser(&params); formErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   formErr.Error(),
		})
	}

	fmt.Println("PARAMS ", params)

	adnet_list, _ = h.DS.GetAdnet(params.ID)

	if params.Status == "1" || params.Status == "true" {
		adnet_list.IsDsp = "true"
	} else {
		adnet_list.IsDsp = "false"
	}

	fmt.Println("adnet_list ", adnet_list)

	if err := h.DS.UpdateAdnetList(&adnet_list); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update adnet_list",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

func (h *IncomingHandler) CreateCompanyGroup(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	var companyGroup entity.CompanyGroup

	if errForm := c.BodyParser(&companyGroup); errForm != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if errValidation := validate.Struct(companyGroup); errValidation != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation error",
			"errors":  errValidation.Error(),
		})
	}

	if errCreate := h.DS.CreateCompanyGroup(&companyGroup); errCreate != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create companyGroup",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})

}

func (h *IncomingHandler) UpdateCompanyGroup(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var companyGroup entity.CompanyGroup

	if formErr := c.BodyParser(&companyGroup); formErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   formErr.Error(),
		})
	}

	companyGroup.ID = uint(id)

	if err := h.DS.UpdateCompanyGroup(&companyGroup); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update companyGroup",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

func (h *IncomingHandler) DeleteCompanyGroup(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.DS.DeleteCompanyGroup(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete companyGroup",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

func (h *IncomingHandler) DisplayCompanyGroup(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	m := c.Queries()

	page, _ := strconv.Atoi(m["page"])
	pageSize, errRequest := strconv.Atoi(m["page_size"])
	if errRequest != nil {
		pageSize = 10
	}
	orderColumn := m["order_column"]
	orderDir := m["order_dir"]

	draw, _ := strconv.Atoi(m["draw"])
	fe := entity.GlobalRequestFromDataTableCompany{
		Page:        page,
		Action:      m["action"],
		Draw:        draw,
		PageSize:    pageSize,
		Search:      m["search[value]"],
		OrderColumn: orderColumn,
		OrderDir:    orderDir,
	}

	var (
		errResponse        error
		total_data         int64
		company_group_list []entity.CompanyGroup
	)

	// key := "temp_key_api_companyGroup_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")

	// need to add redis mechanism here
	company_group_list, total_data, errResponse = h.DS.GetCompanyGroup(fe)

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
				Data:            company_group_list,
				Draw:            fe.Draw,
				RecordsTotal:    int(total_data),
				RecordsFiltered: int(total_data),
			},
		}

	}

	return c.Status(r.HttpStatus).JSON(r.Rsp)
}
