package server

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Middleware function to add validator
func ValidatorMiddleware(validate *validator.Validate) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("validator", validate)
		return c.Next()
	}
}
