package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/donus-turkiye/backend/app/auth"
	"github.com/donus-turkiye/backend/app/healthcheck"
	"github.com/donus-turkiye/backend/infra/postgres"
	_ "github.com/lib/pq"

	"github.com/donus-turkiye/backend/pkg/config"
	_ "github.com/donus-turkiye/backend/pkg/log"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Request any
type Response any

// Define an interface for handlers
type HandlerInterface[R Request, Res Response] interface {
	Handle(ctx context.Context, req *R) (*Res, error)
}

// Update handle function to accept HandlerInterface instead of Handler function
func handle[R Request, Res Response](handler HandlerInterface[R, Res]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req R

		if err := c.BodyParser(&req); err != nil && !errors.Is(err, fiber.ErrUnprocessableEntity) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if err := c.ParamsParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if err := c.QueryParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if err := c.ReqHeaderParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		/*
			ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
			defer cancel()
		*/

		ctx := c.UserContext()

		res, err := handler.Handle(ctx, &req)
		if err != nil {
			zap.L().Error("Failed to handle request", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(res)
	}
}

func main() {
	// config
	appConfig := config.Read()

	defer zap.L().Sync()

	zap.L().Info("app starting...")
	zap.L().Info("app config", zap.Any("appConfig", appConfig))
	fmt.Println(appConfig)

	app := fiber.New(fiber.Config{
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Concurrency:  256 * 1024,
	})

	// Connect to DB
	repo, err := postgres.NewPgRepository(appConfig)
	if err != nil {
		fmt.Println(err)
		zap.L().Fatal("Failed to connect to DB:", zap.Error(err))
	}

	// Define handlers
	healthcheckHandler := healthcheck.NewHealthCheckHandler()
	registerHandler := auth.NewRegisterHandler(repo)

	app.Get("/healthcheck", handle[healthcheck.HealthCheckRequest, healthcheck.HealthCheckResponse](healthcheckHandler))
	app.Post("/users", handle[auth.RegisterRequest, auth.RegisterResponse](registerHandler))

	// Start server in a goroutine
	go func() {
		if err = app.Listen(fmt.Sprintf("0.0.0.0:%s", appConfig.Port)); err != nil {
			zap.L().Error("Failed to start server", zap.Error(err))
			os.Exit(1)
		}
	}()

	zap.L().Info("Server started on port", zap.String("port", appConfig.Port))

	gracefulShutdown(app)
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
