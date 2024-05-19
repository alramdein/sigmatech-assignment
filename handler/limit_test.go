package handler

import (
	"alif-sigmatech/model"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockLimitRepo is a mock implementation of LimitRepo
type MockLimitRepo struct {
	CreateLimitFunc          func(limit *model.Limit) error
	GetLimitByCustomerIDFunc func(customerID int) (*model.Limit, error)
}

func (m *MockLimitRepo) CreateLimit(limit *model.Limit) error {
	if m.CreateLimitFunc != nil {
		return m.CreateLimitFunc(limit)
	}
	return nil
}

func (m *MockLimitRepo) GetLimitByCustomerID(customerID int) (*model.Limit, error) {
	if m.GetLimitByCustomerIDFunc != nil {
		return m.GetLimitByCustomerIDFunc(customerID)
	}
	return nil, nil
}

func TestCreateLimit(t *testing.T) {
	tests := []struct {
		name                string
		input               model.Limit
		mockGetCustomerByID func(id int) (*model.Customer, error)
		mockCreateLimit     func(limit *model.Limit) error
		expectedStatusCode  int
		expectedResponse    interface{}
	}{
		{
			name: "Successful creation",
			input: model.Limit{
				CustomerID: 1,
				Tenor1:     1000,
				Tenor2:     2000,
				Tenor3:     3000,
				Tenor4:     4000,
			},
			mockGetCustomerByID: func(id int) (*model.Customer, error) {
				return &model.Customer{ID: 1}, nil
			},
			mockCreateLimit: func(limit *model.Limit) error {
				return nil
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponse: model.Limit{
				CustomerID: 1,
				Tenor1:     1000,
				Tenor2:     2000,
				Tenor3:     3000,
				Tenor4:     4000},
		},
		{
			name: "Customer not found",
			input: model.Limit{
				CustomerID: 2,
				Tenor1:     1000,
				Tenor2:     2000,
				Tenor3:     3000,
				Tenor4:     4000,
			},
			mockGetCustomerByID: func(id int) (*model.Customer, error) {
				return nil, nil
			},
			mockCreateLimit: func(limit *model.Limit) error {
				return nil
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   "Customer not found",
		},
		{
			name: "Failed to create limit",
			input: model.Limit{
				CustomerID: 1,
				Tenor1:     1000,
				Tenor2:     2000,
				Tenor3:     3000,
				Tenor4:     4000,
			},
			mockGetCustomerByID: func(id int) (*model.Customer, error) {
				return &model.Customer{ID: 1}, nil
			},
			mockCreateLimit: func(limit *model.Limit) error {
				return errors.New("some error")
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   "Failed to create limit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCustomerRepo := newMockCustomerRepo()
			mockLimitRepo := &MockLimitRepo{
				CreateLimitFunc: tt.mockCreateLimit,
			}

			handler := &LimitHandler{
				CustomerRepo: mockCustomerRepo,
				LimitRepo:    mockLimitRepo,
			}

			body, _ := json.Marshal(tt.input)
			req, err := http.NewRequest("POST", "/fund/limit", bytes.NewReader(body))
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			http.HandlerFunc(handler.CreateLimit).ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			if rr.Code == http.StatusCreated {
				var limit model.Limit
				err = json.Unmarshal(rr.Body.Bytes(), &limit)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponse, limit)
			} else {
				assert.Contains(t, rr.Body.String(), tt.expectedResponse.(string))
			}
		})
	}
}
