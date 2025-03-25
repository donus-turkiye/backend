package healthcheck

import (
	"context"
	"net/http"
)

type HealthCheckRequest struct {
}

type HealthCheckResponse struct {
	Status string `json:"status"`
}

type HealthCheckHandler struct {
}

func NewHealthCheckHandler() *HealthCheckHandler {
	return &HealthCheckHandler{}
}

func (h *HealthCheckHandler) Handle(ctx context.Context, req *HealthCheckRequest) (*HealthCheckResponse, int, error) {
	return &HealthCheckResponse{Status: "OK"}, http.StatusOK, nil
}
