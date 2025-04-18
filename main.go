package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/donus-turkiye/backend/infra/postgres"
	"github.com/donus-turkiye/backend/server"
	_ "github.com/lib/pq"

	_ "github.com/donus-turkiye/backend/docs" // swagger docs
	"github.com/donus-turkiye/backend/pkg/config"
	_ "github.com/donus-turkiye/backend/pkg/log"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func main() {
	// config
	appConfig := config.Read()

	defer zap.L().Sync()

	zap.L().Info("app starting...")
	zap.L().Info("app config", zap.Any("appConfig", appConfig))

	// Connect to DB
	repo, err := postgres.NewPgRepository(appConfig)
	if err != nil {
		fmt.Println(err)
		zap.L().Fatal("Failed to connect to DB:", zap.Error(err))
	}

	server := &server.Server{}
	server.NewServer(repo)

	// Start server in a goroutine
	go func() {
		if err := server.Start(appConfig.Port); err != nil {
			zap.L().Error("Failed to start server", zap.Error(err))
			os.Exit(1)
		}
	}()

	// Start session garbage collection in a separate goroutine
	go func() {
		ticker := time.NewTicker(6 * time.Hour)
		for range ticker.C {
			if err := repo.SessionStore.GC(); err != nil {
				zap.L().Error("Failed to clean expired sessions", zap.Error(err))
			}
		}
	}()

	zap.L().Info("Server started on port", zap.String("port", appConfig.Port))

	gracefulShutdown(server.App)
}

func gracefulShutdown(app *fiber.App) {
	// Create channel for shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for shutdown signal
	<-sigChan
	zap.L().Info("Shutting down server...")

	// Shutdown with 5 second timeout
	if err := app.ShutdownWithTimeout(5 * time.Second); err != nil {
		zap.L().Error("Error during server shutdown", zap.Error(err))
	}

	zap.L().Info("Server gracefully stopped")
}
