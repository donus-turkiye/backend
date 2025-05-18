package waste

import (
	"context"
	"net/http"

	"github.com/donus-turkiye/backend/app"
	"github.com/donus-turkiye/backend/domain"
)

type GetCategoriesRequest struct{}

type GetCategoriesResponse struct {
	Categories []domain.Category `json:"categories"`
}

type GetCategoriesHandler struct {
	repository app.Repository
}

func NewGetCategoriesHandler(repository app.Repository) *GetCategoriesHandler {
	return &GetCategoriesHandler{
		repository: repository,
	}
}

// @Summary Get all categories
// @Description Get all categories
// @Tags waste
// @Accept json
// @Produce json
// @Success 200 {object} GetCategoriesResponse
// @Failure 500 {object} error
// @Router /waste/categories [get]
func (h *GetCategoriesHandler) Handle(ctx context.Context, req *GetCategoriesRequest) (*GetCategoriesResponse, int, error) {
	categories, err := h.repository.GetCategories(ctx)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &GetCategoriesResponse{
		Categories: categories,
	}, http.StatusOK, nil
}
