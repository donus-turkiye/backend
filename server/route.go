package server

import (
	"github.com/donus-turkiye/backend/app/auth"
	"github.com/donus-turkiye/backend/app/healthcheck"
	"github.com/donus-turkiye/backend/infra/postgres"
	"github.com/gofiber/swagger"
)

func (s *Server) registerRoutes(repo *postgres.PgRepository) {
	// Swagger route
	s.App.Get("/swagger/*", swagger.HandlerDefault)

	// Define handlers
	healthcheckHandler := healthcheck.NewHealthCheckHandler()
	registerHandler := auth.NewRegisterHandler(repo)
	sessionHandler := auth.NewSessionHandler()

	// Register routes
	s.App.Get("/healthcheck", handle[healthcheck.HealthCheckRequest, healthcheck.HealthCheckResponse](healthcheckHandler))
	s.App.Post("/user", handle[auth.RegisterRequest, auth.RegisterResponse](registerHandler))
	s.App.Get("/session", handle[auth.SessionRequest, auth.SessionResponse](sessionHandler))
}
