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
	"github.com/kanta/backend-challenge/infrastructure"
	handlers "github.com/kanta/backend-challenge/internal/adapters/handlers/backend-handler"
	"github.com/kanta/backend-challenge/internal/adapters/repositories"
	"github.com/kanta/backend-challenge/internal/core/services"
	"github.com/kanta/backend-challenge/middlewares"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func newRouter(handler handlers.BackEndHandler) *fiber.App {
	app := fiber.New()
	app.Use(middlewares.Logger())
	docs.SwaggerInfo.Schemes = []string{"http"}
	swagHandler := swagger.New(swagger.Config{URL: "doc.json"})

	app.Get("/docs/*", func(c *fiber.Ctx) error {
		docs.SwaggerInfo.Schemes = []string{"http"}
		docs.SwaggerInfo.BasePath = "/api/v1"
		return swagHandler(c)
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("server is running")
	})

	v1 := app.Group("/api/v1")

	v1.Post("/register", handler.Register)
	v1.Post("/login", handler.Login)

	protected := v1.Group("", middlewares.JWTAuth(config.Get().JWT_Secret))
	protected.Post("/users", handler.CreateUser)
	protected.Get("/users/:id", handler.GetUserByID)
	protected.Get("/users", handler.ListUsers)
	protected.Put("/users/:id", handler.UpdateUser)
	protected.Delete("/users/:id", handler.DeleteUser)

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

	mongoClient := infrastructure.NewMongoClient(config.Get().Mongo.URI)
	defer infrastructure.MongoDisconnect(mongoClient)

	userRepo := repositories.NewUserRepository(mongoClient, config.Get().Mongo.DB)
	service := services.NewBackEndService(userRepo)

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
