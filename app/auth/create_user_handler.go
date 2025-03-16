package auth

import (
	"context"
	"github.com/donus-turkiye/backend/domain"
	"go.uber.org/zap"
)

type CreateUserRequest struct {
	FullName   string `json:"full_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	RoleId     int    `json:"role_id"`
	TelNumber  string `json:"tel_number"`
	Address    string `json:"address"`
	Coordinate string `json:"coordinate"`
}

type CreateUserResponse struct {
	ID int `json:"id"`
}

type CreateUserHandler struct {
	repository Repository
}

func NewCreateUserHandler(repository Repository) *CreateUserHandler {
	return &CreateUserHandler{
		repository: repository,
	}
}

func (h *CreateUserHandler) Handle(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {

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

	return &CreateUserResponse{ID: userId}, nil
}
