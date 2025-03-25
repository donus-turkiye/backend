package server

import (
	"context"
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
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

		// Validate request
		validate := c.Locals("validator").(*validator.Validate)
		if err := validate.Struct(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Validation failed",
				"details": err.Error(),
			})
		}

		ctx := c.UserContext()

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
