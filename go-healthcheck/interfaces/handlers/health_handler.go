package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/nattatorn-dev/go-healthcheck/domain/services"
)

type ReadinessHandler struct {
	service *services.HealthService
}

func NewReadinessHandler(service *services.HealthService) *ReadinessHandler {
	return &ReadinessHandler{service: service}
}

func (h *ReadinessHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	statuses := h.service.GetReadinessStatuses()
	response := map[string]interface{}{
		"status": "healthy",
		"checks": statuses,
	}
	for _, status := range statuses {
		if status.Status == "DOWN" {
			response["status"] = "unhealthy"
			break
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type LivenessHandler struct {
	service *services.HealthService
}

func NewLivenessHandler(service *services.HealthService) *LivenessHandler {
	return &LivenessHandler{service: service}
}

func (h *LivenessHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	statuses := h.service.GetLivenessStatuses()
	response := map[string]interface{}{
		"status": "healthy",
		"checks": statuses,
	}
	for _, status := range statuses {
		if status.Status == "DOWN" {
			response["status"] = "unhealthy"
			break
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
