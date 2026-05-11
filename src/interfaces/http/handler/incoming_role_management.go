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

func (h *IncomingHandler) CreateRole(c *fiber.Ctx) error {
	var req struct {
		Code        string              `json:"code"`
		Name        string              `json:"name"`
		Permissions []entity.Permission `json:"permissions"`
	}

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

	existingRoles, _ := h.DS.GetAllRoles()
	for _, role := range existingRoles {
		if role.Code == req.Code {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Role code already exists",
			})
		}
	}

	newRole := entity.Role{
		Code: req.Code,
		Name: req.Name,
	}

	if err := h.DS.CreateRole(&newRole); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create role",
		})
	}

	for _, perm := range req.Permissions {
		perm.RoleID = newRole.ID
		if err := h.DS.CreateOrUpdatePermission(&perm); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to update permissions",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

func (h *IncomingHandler) GetRoleTable(c *fiber.Ctx) error {
	key := "temp_key_api_role_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")

	var (
		err     error
		isempty bool
		roles   []entity.RoleManagementData
	)

	if roles, isempty = h.DS.RGetRole(key, "$"); isempty {
		roles, err = h.DS.GetAllRolesWithPermission()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to fetch roles",
			})
		}
		s, _ := json.Marshal(roles)
		h.DS.SetData(key, "$", string(s))
		h.DS.SetExpireData(key, 60)
	}

	totalRecords := len(roles)
	pagesize := PAGESIZE
	page, _ := strconv.Atoi(c.Query("page", "1"))
	start := (page - 1) * pagesize
	end := start + pagesize

	var displayRoles []entity.RoleManagementData
	if start < totalRecords {
		if end > totalRecords {
			end = totalRecords
		}
		displayRoles = roles[start:end]
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithDataTable{
		Draw:            page,
		Code:            fiber.StatusOK,
		Message:         config.OK_DESC,
		Data:            displayRoles,
		RecordsTotal:    totalRecords,
		RecordsFiltered: totalRecords,
	})
}

func (h *IncomingHandler) UpdateRole(c *fiber.Ctx) error {
	idParam := c.Params("id")
	roleID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid role ID",
		})
	}

	var req struct {
		Code        string              `json:"code"`
		Name        string              `json:"name"`
		Permissions []entity.Permission `json:"permissions"`
	}

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

	existingRole, err := h.DS.GetRoleByID(roleID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Role not found",
		})
	}

	duplicateRole, _ := h.DS.GetRoleByCode(req.Code)
	if duplicateRole != nil && duplicateRole.ID != existingRole.ID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Role code already exists",
		})
	}

	existingRole.Code = req.Code
	existingRole.Name = req.Name

	if err := h.DS.UpdateRole(existingRole); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update role",
		})
	}

	for _, perm := range req.Permissions {
		perm.RoleID = existingRole.ID
		if err := h.DS.CreateOrUpdatePermission(&perm); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to update permissions",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

func (h *IncomingHandler) DeleteRole(c *fiber.Ctx) error {
	idParam := c.Params("id")
	roleID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid role ID",
		})
	}

	role, err := h.DS.GetRoleByID(roleID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Role not found",
		})
	}

	if err := h.DS.DeletePermissionsByRoleID(role.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete role permissions",
		})
	}

	if err := h.DS.DeleteRole(role.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete role",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Role deleted successfully",
	})
}