package auth

import (
	"context"
	"github.com/donus-turkiye/backend/app"
	"github.com/donus-turkiye/backend/domain"
	"go.uber.org/zap"
)

type RegisterRequest struct {
	FullName   string `json:"full_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	RoleId     int    `json:"role_id"`
	TelNumber  string `json:"tel_number"`
	Address    string `json:"address"`
	Coordinate string `json:"coordinate"`
}

type RegisterResponse struct {
	ID int `json:"id"`
}

type RegisterHandler struct {
	repository app.Repository
}

func NewRegisterHandler(repository app.Repository) *RegisterHandler {
	return &RegisterHandler{
		repository: repository,
	}
}

func (h *RegisterHandler) Handle(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {

	user := &domain.User{
		FullName:   req.FullName,
		Email:      req.Email,
		Password:   req.Password,
		RoleId:     req.RoleId,
		TelNumber:  req.TelNumber,
		Address:    req.Address,
		Coordinate: req.Coordinate,
	}
	zap.L().Info("User", zap.Any("user", user))
	userId, err := h.repository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	zap.L().Info("User created", zap.Any("userId", userId))

	return &RegisterResponse{ID: userId}, nil
}
