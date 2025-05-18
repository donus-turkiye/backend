package session

import (
	"context"
	"fmt"

	"github.com/donus-turkiye/backend/domain"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func SetSessionUserData(ctx context.Context, user *domain.User) error {
	sess, err := getSessionFromContext(ctx)
	if err != nil {
		return err
	}
	// Set user ID in session
	sess.Set(string(domain.UserDataKey), &domain.UserData{
		UserId: user.Id,
		RoleId: user.RoleId,
	})
	return nil
}

func GetUserIdFromSession(ctx context.Context) (int, error) {
	sess, err := getSessionFromContext(ctx)
	if err != nil {
		return 0, err
	}

	value := sess.Get(string(domain.UserDataKey))
	if value == nil {
		return 0, fmt.Errorf("userData not found in session")
	}
	userData, ok := value.(domain.UserData)
	if !ok {
		return 0, fmt.Errorf("invalid user data format in session")
	}

	return userData.UserId, nil
}

func getSessionFromContext(ctx context.Context) (*session.Session, error) {
	sess, ok := ctx.Value("session").(*session.Session)
	if !ok {
		return nil, fmt.Errorf("session not found from context")
	}
	return sess, nil
}
