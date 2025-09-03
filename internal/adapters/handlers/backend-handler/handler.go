package handlers

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/kanta/backend-challenge/config"
	jwt "github.com/kanta/backend-challenge/infrastructure"
	"github.com/kanta/backend-challenge/internal/core/domain"
	"github.com/kanta/backend-challenge/internal/core/ports"
)

// @title Example API
// @version 1.0
// @description Example Fiber JWT API
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
type BackEndHandler interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	CreateUser(c *fiber.Ctx) error
	GetUserByID(c *fiber.Ctx) error
	ListUsers(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
}

type backEndHandler struct {
	service ports.Service
}

func NewBackEndHandler(
	service ports.Service,
) BackEndHandler {
	return &backEndHandler{
		service,
	}
}

// Register godoc
// @Summary Register new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body domain.User true "User info"
// @Router /register [post]
func (h *backEndHandler) Register(c *fiber.Ctx) error {
	var req domain.User
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	err := h.service.Register(req.Name, req.Email, req.Password)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "registered"})
}

// Login godoc
// @Summary Login and get JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body domain.Login true "Login info"
// @Router /login [post]
func (h *backEndHandler) Login(c *fiber.Ctx) error {
	var req domain.User
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	user, err := h.service.Authenticate(req.Email, req.Password)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	token, _ := jwt.GenerateToken(user.ID, config.Get().JWT_Secret)
	return c.JSON(fiber.Map{"token": token})
}

// CreateUser godoc
// @Summary Create user (admin)
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body domain.User true "User info"
// @Router /users [post]
func (h *backEndHandler) CreateUser(c *fiber.Ctx) error {
	var req domain.User
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	err := h.service.CreateUser(&req)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "created"})
}

// GetUserByID godoc
// @Summary Get user by ID
// @Tags User
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Router /users/{id} [get]
func (h *backEndHandler) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.service.GetUserByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}
	return c.JSON(user)
}

// ListUsers godoc
// @Summary List users (paginated)
// @Tags User
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Router /users [get]
func (h *backEndHandler) ListUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	users, err := h.service.ListUsers(page, limit)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}

// UpdateUser godoc
// @Summary Update user's name or email
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param body body domain.User true "User info"
// @Router /users/{id} [put]
func (h *backEndHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var req domain.User
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	err := h.service.UpdateUser(id, req.Name, req.Email)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "updated"})
}

// DeleteUser godoc
// @Summary Delete user
// @Tags User
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Router /users/{id} [delete]
func (h *backEndHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.service.DeleteUser(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "deleted"})
}
