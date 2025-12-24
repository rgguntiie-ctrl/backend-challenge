package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		zap.L().Info("request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Duration("duration", time.Since(start)),
		)
		return err
	}
}
