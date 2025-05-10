package server

import (
	"fmt"
	"time"

	"github.com/donus-turkiye/backend/infra/postgres"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Server struct {
	App *fiber.App
}

func (s *Server) NewServer(repo *postgres.PgRepository) {
	s.App = fiber.New(fiber.Config{
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Concurrency:  256 * 1024,
	})

	// Initialize session store with repository
	initSessionStore(repo)

	// Add session middleware
	s.App.Use(SessionHeaderMiddleware())
	s.App.Use(SessionMiddleware())
	// Add validator middleware
	validate := validator.New()
	s.App.Use(ValidatorMiddleware(validate))

	// Register routes
	s.registerRoutes(repo)
}

func (s *Server) Start(port string) error {
	zap.L().Info("Starting server on port " + port)
	return s.App.Listen(fmt.Sprintf("0.0.0.0:%s", port))
}
