package server

import (
	"encoding/gob"
	"time"

	"github.com/donus-turkiye/backend/domain"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"go.uber.org/zap"
)

var store *session.Store

func init() {
	// Register types for session storage
	gob.Register(domain.UserData{})

	store = session.New(session.Config{
		KeyLookup:      "cookie:session_id",
		Expiration:     24 * time.Hour,
		CookieSecure:   true,
		CookieHTTPOnly: true,
		CookiePath:     "/",
		CookieDomain:   "",
		CookieSameSite: "Lax",
		Storage:        nil,
	})
}

// Middleware function to add validator
func ValidatorMiddleware(validate *validator.Validate) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("validator", validate)
		return c.Next()
	}
}

// Middleware function to handle sessions
func SessionMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {

		// Log incoming cookie for debugging
		sessionCookie := c.Cookies("session_id")
		if sessionCookie != "" {
			zap.L().Debug("Incoming request cookie",
				zap.String("session_id", sessionCookie),
				zap.String("path", c.Path()))
		} else {
			zap.L().Debug("Incoming request without session cookie",
				zap.String("path", c.Path()))
		}

		sess, err := store.Get(c)
		if err != nil {
			zap.L().Error("Failed to get/create session", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not initialize session",
			})
		}

		// Check if this is a new session
		if sess.Fresh() {
			zap.L().Info("New session created",
				zap.String("session_id", sess.ID()),
				zap.String("client_ip", c.IP()))
		} else {
			zap.L().Debug("Existing session found",
				zap.String("session_id", sess.ID()),
				zap.String("client_ip", c.IP()))
		}

		// Add session to context locals
		c.Locals("session", sess)

		// Continue stack
		return c.Next()
	}
}
