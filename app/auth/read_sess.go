package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/donus-turkiye/backend/domain"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type SessionRequest struct{}
type SessionResponse struct {
	UserData domain.UserData `json:"user_data"`
}

type SessionHandler struct {
}

func NewSessionHandler() *SessionHandler {
	return &SessionHandler{}
}
func (h *SessionHandler) Handle(ctx context.Context, req *SessionRequest) (*SessionResponse, int, error) {
	// Get session from context
	sess, ok := ctx.Value("session").(*session.Session)
	if !ok {
		return nil, http.StatusInternalServerError, errors.New("session not found")
	}

	value := sess.Get(string(domain.UserDataKey))
	if value == nil {
		return nil, http.StatusUnauthorized, errors.New("no user data in session")
	}

	userData, ok := value.(domain.UserData)
	if !ok {
		return nil, http.StatusInternalServerError, errors.New("invalid user data format in session")
	}

	return &SessionResponse{
		UserData: userData,
	}, http.StatusOK, nil
}
