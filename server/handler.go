package server

import (
	"context"
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"go.uber.org/zap"
)

type Request any
type Response any

// Define an interface for handlers
type HandlerInterface[R Request, Res Response] interface {
	Handle(ctx context.Context, req *R) (*Res, int, error)
}

// Update handle function to accept HandlerInterface instead of Handler function
func handle[R Request, Res Response](handler HandlerInterface[R, Res]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req R

		// Read session from context locals
		sessionData, ok := c.Locals("session").(*session.Session)
		if !ok {
			zap.L().Error("Failed to retrieve session from context locals")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Session not found",
			})
		}

		// Parse request body, parameters, query, and headers
		// parse request body if request method is POST, PUT, or DELETE
		if c.Method() == fiber.MethodPost || c.Method() == fiber.MethodPut || c.Method() == fiber.MethodDelete {
			if err := c.BodyParser(&req); err != nil && !errors.Is(err, fiber.ErrUnprocessableEntity) {
				zap.L().Error("Failed to parse request body", zap.Error(err))
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
			}
		}

		if err := c.ParamsParser(&req); err != nil {
			zap.L().Error("Failed to parse request parameters", zap.Error(err))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if err := c.QueryParser(&req); err != nil {
			zap.L().Error("Failed to parse request query", zap.Error(err))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if err := c.ReqHeaderParser(&req); err != nil {
			zap.L().Error("Failed to parse request headers", zap.Error(err))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		/*
			ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
			defer cancel()
		*/

		// Validate request
		validate := c.Locals("validator").(*validator.Validate)
		if err := validate.Struct(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Validation failed",
				"details": err.Error(),
			})
		}

		// Pass session in context
		ctx := context.WithValue(c.UserContext(), "session", sessionData)

		// log the client and request
		zap.L().Info("Request Received", zap.Any("details", map[string]interface{}{
			"method":            c.Method(),
			"endpoint":          c.Path(),
			"request":           req,
			"client_ip":         c.IP(),
			"client_user_agent": c.Get("User-Agent"),
		}))

		// Handle request
		res, status, err := handler.Handle(ctx, &req)
		if err != nil {
			zap.L().Error("Failed to handle request", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": getRootError(err),
			})
		}

		// Save session data
		if err := sessionData.Save(); err != nil {
			zap.L().Error("Failed to save session", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to save session",
			})
		}

		return c.Status(status).JSON(res)
	}
}

func getRootError(err error) string {
	cleanErrorMessage := err.Error()
	if idx := strings.Index(cleanErrorMessage, ":"); idx != -1 {
		cleanErrorMessage = cleanErrorMessage[:idx]
	}
	return cleanErrorMessage
}
