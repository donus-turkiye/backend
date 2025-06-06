package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/donus-turkiye/backend/app"
	"github.com/donus-turkiye/backend/app/session"
	"github.com/donus-turkiye/backend/domain"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	FullName   string `json:"full_name" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=6"`
	RoleId     int    `json:"role_id" validate:"required"`
	TelNumber  string `json:"tel_number" validate:"required"`
	Address    string `json:"address" validate:"required"`
	Coordinate string `json:"coordinate" validate:"required"`
}

type RegisterResponse struct{}

type RegisterHandler struct {
	repository app.Repository
}

func NewRegisterHandler(repository app.Repository) *RegisterHandler {
	return &RegisterHandler{
		repository: repository,
	}
}

// @Summary Register a new user
// @Description Register a new user in the system
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Register Request"
// @Success 200 {object} RegisterResponse
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /user [post]
func (h *RegisterHandler) Handle(ctx context.Context, req *RegisterRequest) (*RegisterResponse, int, error) {

	user := &domain.User{
		FullName:   req.FullName,
		Email:      req.Email,
		Password:   req.Password,
		RoleId:     req.RoleId,
		TelNumber:  req.TelNumber,
		Address:    req.Address,
		Coordinate: req.Coordinate,
	}

	userId, err := h.register(ctx, user)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to register user: %w", err) // TODO: Check error type
	}

	session.SetSessionUserData(ctx, user)

	zap.L().Info("User registered", zap.Int("user_id", userId))

	return &RegisterResponse{}, http.StatusOK, nil
}

func (h *RegisterHandler) register(ctx context.Context, user *domain.User) (int, error) {

	var err error
	user.Password, err = hashPassword(user.Password)
	if err != nil {
		return 0, err
	}

	// Check if user already exists
	_, err = h.repository.GetUserByEmail(ctx, user.Email)
	if err == nil {
		return 0, fmt.Errorf("user already exists")
	}

	userId, err := h.repository.CreateUser(ctx, user)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
