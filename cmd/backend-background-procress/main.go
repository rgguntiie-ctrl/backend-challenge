package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kanta/backend-challenge/config"
	"github.com/kanta/backend-challenge/infrastructure"
	"github.com/kanta/backend-challenge/internal/adapters/repositories"
	"github.com/kanta/backend-challenge/internal/core/services"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

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

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("server is running"))
	})

	go func() {
		zap.L().Info(fmt.Sprintf("health check server running on port: %d", config.Get().App.Port))
		if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Get().App.Port), nil); err != nil {
			zap.L().Fatal("failed to start server", zap.Error(err))
		}
	}()

	service.RunUserCountLogger(10 * time.Second)

}
