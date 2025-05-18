package user

import (
	"context"
	"net/http"

	"github.com/donus-turkiye/backend/app"
	"github.com/donus-turkiye/backend/app/session"
	"github.com/donus-turkiye/backend/domain"
)

type UserDataRequest struct{}
type UserDataResponse struct {
	User domain.User `json:"user"`
}
type UserDataHandler struct {
	repository app.Repository
}

func NewUserDataHandler(repository app.Repository) *UserDataHandler {
	return &UserDataHandler{
		repository: repository,
	}
}

// @Summary Get user data
// @Description Get user data from the system
// @Tags user
// @Accept json
// @Produce json
// @Param request header session id "X-Session-ID"
// @Success 200 {object} UserDataResponse
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /user [get]
func (h *UserDataHandler) Handle(ctx context.Context, req *UserDataRequest) (*UserDataResponse, int, error) {
	userId, err := session.GetUserIdFromSession(ctx)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	// Get user data from repository
	user, err := h.repository.GetUserById(ctx, userId)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &UserDataResponse{
		User: *user,
	}, http.StatusOK, nil
}
