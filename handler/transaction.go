package handler

import (
	"alif-sigmatech/model"
	"alif-sigmatech/repository"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

// TransactionHandler handles HTTP requests related to Transaction
type TransactionHandler struct {
	TransactionRepo repository.TransactionRepository
	LimitRepo       repository.LimitRepository
}

// NewTransactionHandler creates a new instance of TransactionHandler
func NewTransactionHandler(repo repository.TransactionRepository,
	limitRepo repository.LimitRepository) *TransactionHandler {
	return &TransactionHandler{
		TransactionRepo: repo,
		LimitRepo:       limitRepo,
	}
}

// CreateTransaction handles the creation of a new transaction
func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction model.Transaction
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Check customer limit
	limit, err := h.LimitRepo.GetLimitByCustomerID(transaction.CustomerID)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "Failed to get customer limit", http.StatusInternalServerError)
		return
	}
	if limit == nil {
		http.Error(w, "Customer limit not found", http.StatusNotFound)
		return
	}

	if !isWithinLimit(transaction, limit) {
		http.Error(w, "Transaction exceeds limit", http.StatusBadRequest)
		return
	}

	err = h.TransactionRepo.CreateTransaction(&transaction)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "Failed to create transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}

// isWithinLimit checks if the transaction is within the customer's limit based on the tenor
func isWithinLimit(transaction model.Transaction, limit *model.Limit) bool {
	switch transaction.Tenor {
	case 1:
		return transaction.InstallmentAmount <= limit.Tenor1
	case 2:
		return transaction.InstallmentAmount <= limit.Tenor2
	case 3:
		return transaction.InstallmentAmount <= limit.Tenor3
	case 4:
		return transaction.InstallmentAmount <= limit.Tenor4
	default:
		return false
	}
}
