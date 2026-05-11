package handler

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/config"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/external"
)

func (h *IncomingHandler) CreateUser(c *fiber.Ctx) error {
	var req struct {
		Name     string              `json:"name" binding:"required"`
		Username string              `json:"username" binding:"required"`
		Email    string              `json:"email" binding:"required,email"`
		Password string              `json:"password" binding:"required"`
		RoleID   int                 `json:"role_id" binding:"required"`
		Services []entity.DetailUser `json:"services"`
		Adnets   []entity.UserAdnet  `json:"adnets"`
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

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to hash password",
		})
	}

	newUser := entity.User{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword), // Store hashed password
		RoleID:   req.RoleID,
		IsVerify: true,
		Status:   true,
	}

	createdUser, err := h.DS.CreateUser(&newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Insert Services
	for _, service := range req.Services {
		service.UserID = createdUser.ID
		if _, err := h.DS.CreateOrUpdateDetailUser(&service); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to update services",
			})
		}
	}

	// Insert Adnets
	for _, adnet := range req.Adnets {
		adnet.UserID = createdUser.ID
		if _, err := h.DS.CreateOrUpdateUserAdnet(&adnet); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to update adnets",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

func (h *IncomingHandler) GetUserTable(c *fiber.Ctx) error {
	key := "temp_key_api_user_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")

	var (
		err     error
		isempty bool
		user    []entity.UserManagementData
	)

	if user, isempty = h.DS.RGetUser(key, "$"); isempty {
		user, err = h.DS.GetAllUserWithRelation()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to fetch user",
			})
		}
		s, _ := json.Marshal(user)
		h.DS.SetData(key, "$", string(s))
		h.DS.SetExpireData(key, 60)
	}

	totalRecords := len(user)
	pagesize := PAGESIZE
	page, _ := strconv.Atoi(c.Query("page", "1"))
	start := (page - 1) * pagesize
	end := start + pagesize

	var displayUser []entity.UserManagementData
	if start < totalRecords {
		if end > totalRecords {
			end = totalRecords
		}
		displayUser = user[start:end]
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithDataTable{
		Draw:            page,
		Code:            fiber.StatusOK,
		Message:         config.OK_DESC,
		Data:            displayUser,
		RecordsTotal:    totalRecords,
		RecordsFiltered: totalRecords,
	})
}

func (h *IncomingHandler) GetUserCounts(c *fiber.Ctx) error {
	counts, err := h.DS.GetUserCounts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get user counts",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(counts)
}

func (h *IncomingHandler) UpdateUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid user ID",
		})
	}

	var req struct {
		Name     string              `json:"name" binding:"required"`
		Username string              `json:"username" binding:"required"`
		Email    string              `json:"email" binding:"required,email"`
		Password string              `json:"password,omitempty"`
		RoleID   int                 `json:"role_id" binding:"required"`
		Services []entity.DetailUser `json:"services"`
		Adnets   []entity.UserAdnet  `json:"adnets"`
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

	// Fetch existing user
	existingUser, err := h.DS.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	// Update user data
	existingUser.Name = req.Name
	existingUser.Username = req.Username
	existingUser.Email = req.Email
	existingUser.RoleID = req.RoleID

	// Hash password only if a new one is provided
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to hash password",
			})
		}
		existingUser.Password = string(hashedPassword)
	}

	updatedUser, err := h.DS.UpdateUser(userID, existingUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Update Services
	for _, service := range req.Services {
		service.UserID = updatedUser.ID
		if _, err := h.DS.CreateOrUpdateDetailUser(&service); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to update services",
			})
		}
	}

	// Update Adnets
	for _, adnet := range req.Adnets {
		adnet.UserID = updatedUser.ID
		if _, err := h.DS.CreateOrUpdateUserAdnet(&adnet); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to update adnets",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

func (h *IncomingHandler) AssignService(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid user ID",
		})
	}

	var req struct {
		Services []entity.DetailUser `json:"services"`
		Adnets   []entity.UserAdnet  `json:"adnets"`
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

	// Fetch existing user
	existingUser, err := h.DS.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	// Update Services
	for _, service := range req.Services {
		service.UserID = existingUser.ID
		if _, err := h.DS.CreateOrUpdateDetailUser(&service); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to update services",
			})
		}
	}

	// Update Adnets
	for _, adnet := range req.Adnets {
		adnet.UserID = existingUser.ID
		if _, err := h.DS.CreateOrUpdateUserAdnet(&adnet); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to update adnets",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{Code: fiber.StatusOK, Message: config.OK_DESC})
}

func (h *IncomingHandler) UpdateUserStatus(c *fiber.Ctx) error {
	idParam := c.Params("id")

	userID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid user ID",
		})
	}

	var req struct {
		Status bool `json:"status" binding:"required"`
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

	_, err = h.DS.UpdateUserStatus(userID, req.Status)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	statusMessage := "User is now inactive"
	if req.Status {
		statusMessage = "User is now active"
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": statusMessage,
	})
}

func (h *IncomingHandler) DeleteUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid user ID",
		})
	}

	user, err := h.DS.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	if err := h.DS.DeleteUser(user.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

func (h *IncomingHandler) GetUserApplovalRequestTable(c *fiber.Ctx) error {
	key := "temp_key_api_user_approval_request_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")

	var (
		err     error
		isempty bool
		user    []entity.UserApprovalRequestData
	)

	if user, isempty = h.DS.RGetUserApprovalRequest(key, "$"); isempty {
		user, err = h.DS.GetAllUserApprovalRequest()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to fetch user",
			})
		}
		s, _ := json.Marshal(user)
		h.DS.SetData(key, "$", string(s))
		h.DS.SetExpireData(key, 60)
	}

	totalRecords := len(user)
	pagesize := PAGESIZE
	page, _ := strconv.Atoi(c.Query("page", "1"))
	start := (page - 1) * pagesize
	end := start + pagesize

	var displayUser []entity.UserApprovalRequestData
	if start < totalRecords {
		if end > totalRecords {
			end = totalRecords
		}
		displayUser = user[start:end]
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithDataTable{
		Draw:            page,
		Code:            fiber.StatusOK,
		Message:         config.OK_DESC,
		Data:            displayUser,
		RecordsTotal:    totalRecords,
		RecordsFiltered: totalRecords,
	})
}

func (h *IncomingHandler) ApproveUser(c *fiber.Ctx) error {
	idParam := c.Params("id")

	userID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid user ID"})
	}

	var req struct {
		RoleID   int    `json:"role_id"`
		VerifyBy string `json:"verify_by"`
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

	err = h.DS.ApproveUser(userID, req.RoleID, req.VerifyBy)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User approved successfully"})
}
