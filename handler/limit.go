package handler

import (
	"alif-sigmatech/model"
	"alif-sigmatech/repository"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

// LimitHandler handles HTTP requests related to limits
type LimitHandler struct {
	LimitRepo    repository.LimitRepository
	CustomerRepo repository.CustomerRepository
}

// NewLimitHandler creates a new instance of LimitHandler
func NewLimitHandler(limitRepo repository.LimitRepository,
	customerRepo repository.CustomerRepository) *LimitHandler {
	return &LimitHandler{
		LimitRepo:    limitRepo,
		CustomerRepo: customerRepo,
	}
}

// CreateLimit handles the creation of a new limit
func (h *LimitHandler) CreateLimit(w http.ResponseWriter, r *http.Request) {
	var limit model.Limit
	err := json.NewDecoder(r.Body).Decode(&limit)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	customer, err := h.CustomerRepo.GetCustomerByID(limit.CustomerID)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "Failed to create limit", http.StatusInternalServerError)
		return
	}
	if customer == nil {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	err = h.LimitRepo.CreateLimit(&limit)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "Failed to create limit", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(limit)
}
