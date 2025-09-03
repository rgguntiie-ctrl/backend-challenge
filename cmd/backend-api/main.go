package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/kanta/backend-challenge/config"
	"github.com/kanta/backend-challenge/docs"
	handlers "github.com/kanta/backend-challenge/internal/adapters/handlers/backend-handler"
	"github.com/kanta/backend-challenge/internal/core/services"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

func newRouter(handler handlers.BackEndHandler) *fiber.App {
	app := fiber.New()
	// TODO: register routes with handler

	docs.SwaggerInfo.Schemes = []string{"http"}
	swagHanlder := swagger.HandlerDefault
	swagHanlder = swagger.New(swagger.Config{URL: "doc.json"})

	app.Get("/docs/*", func(c *fiber.Ctx) error {
		// if !config.IsLocal() {
		docs.SwaggerInfo.Schemes = []string{"https"}
		docs.SwaggerInfo.BasePath = fmt.Sprintf("/%s/api/v1", c.Get("X-Forwarded-Prefix"))
		// }
		return swagHanlder(c)
	})
	// }

	app.Use(cors.New())
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("server is running")
	})

	v1 := app.Group("/api/v1")
	fmt.Print(v1)

	return app
}

func main() {
	config.Load()

	logrus.SetFormatter(&logrus.JSONFormatter{})

	var logger *zap.Logger
	logger, _ = zap.NewDevelopment()

	defer logger.Sync()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	service := services.NewBackEndService()
	handler := handlers.NewBackEndHandler(service)

	app := newRouter(handler)
	go func() {
		if err := app.Listen(fmt.Sprintf("%s:%d", config.Get().App.Host, config.Get().App.Port)); err != nil {
			zap.L().Sugar().Fatal(err)
		}
	}()

	gracefulShutdown(app)

}

func gracefulShutdown(app *fiber.App) {
	var (
		quit = make(chan os.Signal, 1)
		done = make(chan struct{}, 1)
	)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Sugar().Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		zap.L().Sugar().Fatal(err)
	}
	done <- struct{}{}

	select {
	case <-ctx.Done():
		zap.L().Sugar().Info("timeout exceeded, force shutdown")
	case <-done:
		zap.L().Sugar().Info("shutdown completed")
	}
}
