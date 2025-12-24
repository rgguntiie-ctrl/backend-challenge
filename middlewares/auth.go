package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/kanta/backend-challenge/infrastructure"
	"github.com/kanta/backend-challenge/internal/core/ports"
)

func JWTAuth(secret string, cache ports.CachePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization format",
			})
		}

		token := parts[1]

		userID, err := infrastructure.ValidateAccessTokenWithCache(c.Context(), token, secret, cache)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		claims, _ := infrastructure.ParseToken(token, secret)

		c.Locals("user_id", userID)
		c.Locals("claims", claims)

		return c.Next()
	}
}

func RefreshTokenAuth(secret string, cache ports.CachePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization format",
			})
		}

		token := parts[1]

		userID, err := infrastructure.ValidateRefreshTokenWithCache(c.Context(), token, secret, cache)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired refresh token",
			})
		}

		claims, _ := infrastructure.ParseToken(token, secret)

		c.Locals("user_id", userID)
		c.Locals("claims", claims)

		return c.Next()
	}
}

func OptionalAuth(secret string, cache ports.CachePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")

		if auth == "" {
			return c.Next()
		}

		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Next()
		}

		token := parts[1]

		userID, err := infrastructure.ValidateAccessTokenWithCache(c.Context(), token, secret, cache)
		if err == nil {
			claims, _ := infrastructure.ParseToken(token, secret)
			c.Locals("user_id", userID)
			c.Locals("claims", claims)
		}

		return c.Next()
	}
}
