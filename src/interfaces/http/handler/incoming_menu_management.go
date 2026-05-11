package handler

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/config"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/external"
)

var validate = validator.New()

func (h *IncomingHandler) CreateMenu(c *fiber.Ctx) error {
	var req entity.Menu
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation error",
			"errors":  err.Error(),
		})
	}

	existingMenu, _ := h.DS.GetAllMenus()
	for _, menu := range existingMenu {
		if menu.Code == req.Code || menu.Sort == req.Sort {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Code or Sort already exists",
			})
		}
	}

	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	if err := h.DS.CreateMenu(&req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create menu",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

func (h *IncomingHandler) GetAllMenus(c *fiber.Ctx) error {
	key := "temp_key_api_menu_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")

	var (
		err     error
		isempty bool
		menus   []entity.Menu
	)

	if menus, isempty = h.DS.RGetMenu(key, "$"); isempty {
		menus, err = h.DS.GetAllMenus()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to fetch menus",
			})
		}
		s, _ := json.Marshal(menus)
		h.DS.SetData(key, "$", string(s))
		h.DS.SetExpireData(key, 60)
	}

	totalRecords := len(menus)
	pagesize := PAGESIZE
	page, _ := strconv.Atoi(c.Query("page", "1"))
	start := (page - 1) * pagesize
	end := start + pagesize

	var displayMenus []entity.Menu
	if start < totalRecords {
		if end > totalRecords {
			end = totalRecords
		}
		displayMenus = menus[start:end]
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithDataTable{
		Draw:            page,
		Code:            fiber.StatusOK,
		Message:         config.OK_DESC,
		Data:            displayMenus,
		RecordsTotal:    totalRecords,
		RecordsFiltered: totalRecords,
	})
}

func (h *IncomingHandler) GetMenuByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponseWithData{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid ID format",
			Data:    []entity.Menu{},
		})
	}

	menu, err := h.DS.GetMenuByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(entity.GlobalResponseWithData{
			Code:    fiber.StatusNotFound,
			Message: "Menu not found",
			Data:    []entity.Menu{},
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithData{
		Code:    fiber.StatusOK,
		Message: config.OK_DESC,
		Data:    menu,
	})
}


func (h *IncomingHandler) UpdateMenu(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var req entity.Menu

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	req.ID = id
	req.UpdatedAt = time.Now()

	if err := h.DS.UpdateMenu(&req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update menu",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

func (h *IncomingHandler) DeleteMenu(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	if err := h.DS.DeleteMenu(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete menu",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}
