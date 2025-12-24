package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/kanta/backend-challenge/config"
	jwt "github.com/kanta/backend-challenge/infrastructure"
	"github.com/kanta/backend-challenge/internal/core/domain"
	"github.com/kanta/backend-challenge/internal/core/ports"
	"github.com/kanta/backend-challenge/middlewares/meta"
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
	RefreshToken(c *fiber.Ctx) error
	GetMyProfile(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
}

type backEndHandler struct {
	service ports.Service
	cache   ports.CachePort
}

func NewBackEndHandler(
	service ports.Service,
	cache ports.CachePort,
) BackEndHandler {
	return &backEndHandler{
		service,
		cache,
	}
}

// Register godoc
// @Summary Register new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body domain.User true "User info"
// @Router /auth/register [post]
func (h *backEndHandler) Register(c *fiber.Ctx) error {
	var req domain.User
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(meta.NewMetaError(http.StatusBadRequest, "invalid request"))
	}
	err := h.service.Register(req.Name, req.Email, req.Password)
	if err != nil {
		return c.JSON(meta.NewMetaError(http.StatusBadRequest, err.Error()))
	}
	ok := meta.NewMetaOK("register successfully", map[string]interface{}{"message": "registered"})
	return c.JSON(ok)
}

// Login godoc
// @Summary Login and get JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body domain.Login true "Login info"
// @Router /auth/login [post]
func (h *backEndHandler) Login(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		req domain.User
	)
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(meta.NewMetaError(http.StatusBadRequest, "invalid request"))
	}
	user, err := h.service.Authenticate(req.Email, req.Password)
	if err != nil {

		return c.JSON(meta.NewMetaError(http.StatusUnauthorized, "invalid credentials"))
	}

	tokenPair, _ := jwt.GenerateTokenPairWithCache(ctx, user.ID, config.Get().JWT_Secret, h.cache)
	ok := meta.NewMetaOK("login successfully", map[string]interface{}{
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
		"user": fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})

	return c.JSON(ok)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get a new access token using refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param refresh body domain.RefreshTokenRequest true "Refresh token"
// @Router /auth/refresh [post]
func (h *backEndHandler) RefreshToken(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		req domain.RefreshTokenRequest
	)
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(meta.NewMetaError(http.StatusBadRequest, "invalid request"))

	}

	if req.RefreshToken == "" {
		return c.JSON(meta.NewMetaError(http.StatusBadRequest, "refresh token is required"))

	}

	accessToken, err := jwt.RefreshAccessTokenWithCache(ctx, req.RefreshToken, config.Get().JWT_Secret, h.cache)
	if err != nil {
		return c.JSON(meta.NewMetaError(http.StatusUnauthorized, "invalid or expired refresh token"))

	}
	ok := meta.NewMetaOK("refresh token successfully", map[string]interface{}{"access_token": accessToken})
	return c.JSON(ok)
}

// GetMyProfile godoc
// @Summary Get current user profile
// @Description Get authenticated user's profile information
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /users/me [get]
func (h *backEndHandler) GetMyProfile(c *fiber.Ctx) error {

	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.JSON(meta.NewMetaError(http.StatusUnauthorized, "unauthorized"))
	}

	user, err := h.service.GetUserByID(userID)
	if err != nil {
		return c.JSON(meta.NewMetaError(http.StatusNotFound, "user not found"))

	}
	res := domain.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	resOk := meta.NewMetaOK("get user successfully", res)
	return c.JSON(resOk)
}

// Logout godoc
// @Summary Logout user
// @Description Invalidate the current user session by revoking tokens from Redis
// @Tags Logout
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /auth/logout [post]
func (h *backEndHandler) Logout(c *fiber.Ctx) error {
	ctx := c.Context()
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.JSON(meta.NewMetaError(http.StatusUnauthorized, "unauthorized"))

	}

	if err := jwt.RevokeToken(ctx, userID, h.cache); err != nil {
		return c.JSON(meta.NewMetaError(http.StatusInternalServerError, "unauthorized"))

	}

	resOk := meta.NewMetaOK("logged out successfully", nil)
	return c.JSON(resOk)
}
