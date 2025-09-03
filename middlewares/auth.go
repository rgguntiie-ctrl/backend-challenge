package middlewares

import (
	"github.com/gofiber/fiber/v2"
	jwt "github.com/kanta/backend-challenge/interfrastucture"
)

func JWTAuth(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if len(auth) < 8 || auth[:7] != "Bearer " {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing or invalid token"})
		}
		userID, err := jwt.ParseToken(auth[7:], secret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}
		c.Locals("user_id", userID)
		return c.Next()
	}
}
