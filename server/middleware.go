package server

import (
	"encoding/gob"
	"encoding/json"
	"time"

	"github.com/donus-turkiye/backend/domain"
	"github.com/donus-turkiye/backend/infra/postgres"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"go.uber.org/zap"
)

var store *session.Store

func initSessionStore(repo *postgres.PgRepository) {
	// Register types for session storage
	gob.Register(domain.UserData{})

	store = session.New(session.Config{
		Storage:        repo.SessionStore,
		KeyLookup:      "header:X-Session-ID", // Changed from cookie to header
		Expiration:     24 * time.Hour,
		CookieSecure:   false, // Not needed for header based sessions
		CookieHTTPOnly: false, // Not needed for header based sessions
		CookiePath:     "",    // Not needed for header based sessions
		CookieDomain:   "",    // Not needed for header based sessions
		CookieSameSite: "",    // Not needed for header based sessions
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
		// Log incoming session header for debugging
		sessionID := c.Get("X-Session-ID")
		if sessionID != "" {
			zap.L().Debug("Incoming request session",
				zap.String("session_id", sessionID),
				zap.String("path", c.Path()))
		} else {
			zap.L().Debug("Incoming request without session",
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

func SessionHeaderMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Process request
		err := c.Next()
		if err != nil {
			return err
		}

		// Get session ID from response body
		var response map[string]interface{}
		if err := json.Unmarshal(c.Response().Body(), &response); err != nil {
			return err
		}

		// If session ID exists in response, set it in next request header
		if sessionID, ok := response["session_id"].(string); ok {
			c.Set("X-Session-ID", sessionID)
		}

		return nil
	}
}
