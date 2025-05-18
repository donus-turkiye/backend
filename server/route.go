package server

import (
	"github.com/donus-turkiye/backend/app/auth"
	"github.com/donus-turkiye/backend/app/healthcheck"
	"github.com/donus-turkiye/backend/app/user"
	"github.com/donus-turkiye/backend/infra/postgres"
	"github.com/gofiber/swagger"
)

func (s *Server) registerRoutes(repo *postgres.PgRepository) {
	// Swagger route
	s.App.Get("/swagger/*", swagger.HandlerDefault)

	// Define handlers
	healthcheckHandler := healthcheck.NewHealthCheckHandler()
	registerHandler := auth.NewRegisterHandler(repo)
	loginHandler := auth.NewLoginHandler(repo)
	sessionHandler := auth.NewSessionHandler()
	userDataHandler := user.NewUserDataHandler(repo)

	// Register routes
	s.App.Get("/healthcheck", handle[healthcheck.HealthCheckRequest, healthcheck.HealthCheckResponse](healthcheckHandler))
	s.App.Post("/user", handle[auth.RegisterRequest, auth.RegisterResponse](registerHandler))
	s.App.Post("/login", handle[auth.LoginRequest, auth.LoginResponse](loginHandler))
	s.App.Get("/session", handle[auth.SessionRequest, auth.SessionResponse](sessionHandler))
	s.App.Get("/user", handle[user.UserDataRequest, user.UserDataResponse](userDataHandler))
}
