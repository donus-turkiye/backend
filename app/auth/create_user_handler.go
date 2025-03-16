package auth

import (
	"context"
	"github.com/donus-turkiye/backend/domain"
)

type CreateUserRequest struct {
	FullName   string `json:"full_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	RoleId     int    `json:"role_id"`
	TelNumber  int    `json:"tel_number"`
	Address    string `json:"address"`
	Coordinate string `json:"coordinate"`
}

type CreateUserResponse struct {
	ID string `json:"id"`
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

	userId, err := h.repository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &CreateUserResponse{ID: string(userId)}, nil
}
