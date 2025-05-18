package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/donus-turkiye/backend/app"
	"github.com/donus-turkiye/backend/app/session"
	"github.com/donus-turkiye/backend/domain"
	"go.uber.org/zap"
)

type LoginRequest struct {
	Email     string `json:"email" validate:"omitempty,email,excluded_with=TelNumber"`
	Password  string `json:"password" validate:"required,min=6"`
	TelNumber string `json:"tel_number" validate:"omitempty,excluded_with=Email"`
}

type LoginResponse struct {
	UserData domain.UserData `json:"user_data"`
}

type LoginHandler struct {
	repository app.Repository
}

func NewLoginHandler(repository app.Repository) *LoginHandler {
	return &LoginHandler{
		repository: repository,
	}
}

// @Summary Login a user
// @Description Login a user in the system
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login Request"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /login [post]
func (h *LoginHandler) Handle(ctx context.Context, req *LoginRequest) (*LoginResponse, int, error) {

	var user *domain.User
	var err error

	// Check login method
	if req.Email != "" {
		user, err = h.repository.GetUserByEmail(ctx, req.Email)
	} else {
		user, err = h.repository.GetUserByTelNumber(ctx, req.TelNumber)
	}

	if err != nil {
		return nil, http.StatusUnauthorized, fmt.Errorf("invalid credentials")
	}

	if !checkPasswordHash(req.Password, user.Password) {
		return nil, http.StatusUnauthorized, fmt.Errorf("invalid credentials")
	}

	session.SetSessionUserData(ctx, user)

	zap.L().Info("User logged in",
		zap.Int("user_id", user.Id),
		zap.String("login_method", getLoginMethod(req)))

	return &LoginResponse{
			UserData: domain.UserData{
				UserId: user.Id,
				RoleId: user.RoleId,
			}},
		http.StatusOK, nil
}

func getLoginMethod(req *LoginRequest) string {
	if req.Email != "" {
		return "email"
	}
	return "phone"
}
