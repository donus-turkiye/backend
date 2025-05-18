package app

import (
	"context"

	"github.com/donus-turkiye/backend/domain"
)

type Repository interface {
	CreateUser(ctx context.Context, user *domain.User) (int, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByTelNumber(ctx context.Context, telNumber string) (*domain.User, error)
	GetUserById(ctx context.Context, userId int) (*domain.User, error)
	GetCategories(ctx context.Context) ([]domain.Category, error)
}
