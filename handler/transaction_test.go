package handler

import (
	"alif-sigmatech/model"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"alif-sigmatech/mocks"
)

func TestCreateTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLimitRepo := mocks.NewMockLimitRepository(ctrl)
	mockTransactionRepo := mocks.NewMockTransactionRepository(ctrl)

	h := NewTransactionHandler(mockTransactionRepo, mockLimitRepo)

	t.Run("Success", func(t *testing.T) {
		mockLimit := &model.Limit{
			CustomerID: 1,
			Tenor1:     500000,
			Tenor2:     700000,
			Tenor3:     900000,
			Tenor4:     1100000,
		}

		mockLimitRepo.EXPECT().GetLimitByCustomerID(1).Return(mockLimit, nil)
		mockTransactionRepo.EXPECT().CreateTransaction(gomock.Any()).Return(nil)

		transaction := &model.Transaction{
			CustomerID:        1,
			InstallmentAmount: 300000,
			Tenor:             1,
		}

		body, _ := json.Marshal(transaction)
		req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
		recorder := httptest.NewRecorder()

		h.CreateTransaction(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code)
	})

	t.Run("Invalid payload", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer([]byte("invalid json")))
		recorder := httptest.NewRecorder()

		h.CreateTransaction(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("Limit not found", func(t *testing.T) {
		mockLimitRepo.EXPECT().GetLimitByCustomerID(1).Return(nil, nil)

		transaction := &model.Transaction{
			CustomerID:        1,
			InstallmentAmount: 300000,
			Tenor:             1,
		}

		body, _ := json.Marshal(transaction)
		req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
		recorder := httptest.NewRecorder()

		h.CreateTransaction(recorder, req)

		assert.Equal(t, http.StatusNotFound, recorder.Code)
	})

	t.Run("Transaction exceeds limit", func(t *testing.T) {
		mockLimit := &model.Limit{
			CustomerID: 1,
			Tenor1:     100000,
			Tenor2:     100000,
			Tenor3:     100000,
			Tenor4:     100000,
		}

		mockLimitRepo.EXPECT().GetLimitByCustomerID(1).Return(mockLimit, nil)

		transaction := &model.Transaction{
			CustomerID:        1,
			InstallmentAmount: 300000,
			Tenor:             1,
		}

		body, _ := json.Marshal(transaction)
		req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
		recorder := httptest.NewRecorder()

		h.CreateTransaction(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("Error from GetLimitByCustomerID", func(t *testing.T) {
		mockLimitRepo.EXPECT().GetLimitByCustomerID(1).Return(nil, errors.New("database error"))

		transaction := &model.Transaction{
			CustomerID:        1,
			InstallmentAmount: 300000,
			Tenor:             1,
		}

		body, _ := json.Marshal(transaction)
		req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
		recorder := httptest.NewRecorder()

		h.CreateTransaction(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})

	t.Run("Error from CreateTransaction", func(t *testing.T) {
		mockLimit := &model.Limit{
			CustomerID: 1,
			Tenor1:     500000,
			Tenor2:     700000,
			Tenor3:     900000,
			Tenor4:     1100000,
		}

		mockLimitRepo.EXPECT().GetLimitByCustomerID(1).Return(mockLimit, nil)
		mockTransactionRepo.EXPECT().CreateTransaction(gomock.Any()).Return(errors.New("database error"))

		transaction := &model.Transaction{
			CustomerID:        1,
			InstallmentAmount: 300000,
			Tenor:             1,
		}

		body, _ := json.Marshal(transaction)
		req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
		recorder := httptest.NewRecorder()

		h.CreateTransaction(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})
}
