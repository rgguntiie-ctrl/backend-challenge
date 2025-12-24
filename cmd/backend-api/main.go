package main

import (
	"context"
	"fmt"
	"log"
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
	cache "github.com/kanta/backend-challenge/internal/adapters/cache"
	handlers "github.com/kanta/backend-challenge/internal/adapters/handlers/backend-handler"
	"github.com/kanta/backend-challenge/internal/adapters/repositories"
	"github.com/kanta/backend-challenge/internal/core/ports"
	"github.com/kanta/backend-challenge/internal/core/services"
	"github.com/kanta/backend-challenge/middlewares"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func newRouter(handler handlers.BackEndHandler, cache ports.CachePort) *fiber.App {
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

	v1.Post("/auth/register", handler.Register)
	v1.Post("/auth/login", handler.Login)
	v1.Post("/auth/refresh", handler.RefreshToken)

	protected := v1.Group("", middlewares.JWTAuth(config.Get().JWT_Secret, cache))
	protected.Get("/users/me", handler.GetMyProfile)
	protected.Post("/auth/logout", handler.Logout)
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

	db := infrastructure.NewPostgresClient(
		config.Get().Psql.Host,
		config.Get().Psql.User,
		config.Get().Psql.Pass,
		config.Get().Psql.DB,
		config.Get().Psql.Port)

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
	}
	defer sqlDB.Close()

	redisConfig := config.Get().Redis
	redisClient := infrastructure.NewRedisClient(redisConfig.Addr, redisConfig.Password, redisConfig.DB)

	userRepo := repositories.NewUserRepository(db)
	tokenCache := cache.NewTokenCache(redisClient)

	service := services.NewBackEndService(userRepo)
	handler := handlers.NewBackEndHandler(service, tokenCache)

	app := newRouter(handler, tokenCache)
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
