package handler

import (
	"alif-sigmatech/mocks"
	"alif-sigmatech/model"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type mockCustomerRepo struct {
	customers           map[int]*model.Customer
	GetCustomerByIDFunc func(id int) (*model.Customer, error)
}

func newMockCustomerRepo() *mockCustomerRepo {
	return &mockCustomerRepo{
		customers: map[int]*model.Customer{
			1: {
				ID: 1,
			},
		},
	}
}

func (m *mockCustomerRepo) GetCustomerByNIK(nik string) (*model.Customer, error) {
	for _, customer := range m.customers {
		if customer.NIK == nik {
			return customer, nil
		}
	}
	return nil, nil
}

func (m *mockCustomerRepo) RegisterCustomer(customer *model.Customer) error {
	m.customers[customer.ID] = customer
	return nil
}

func (m *mockCustomerRepo) GetCustomerByID(id int) (*model.Customer, error) {
	customer, exists := m.customers[id]
	if !exists {
		return nil, nil
	}
	return customer, nil
}

func hashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword)
}

func TestRegisterCustomer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCustomerRepo := mocks.NewMockCustomerRepository(ctrl)
	handler := &AuthHandler{
		CustomerRepo:  mockCustomerRepo,
		EncryptionKey: []byte("test-key"),
		JWTSecret:     []byte("test-secret"),
	}

	// Create a request body
	customer := model.Customer{
		NIK:       "182381283182",
		FullName:  "Alif Coba",
		LegalName: "John Doe",
		Password:  "password",
		KTPPhoto:  nil,
	}
	body, _ := json.Marshal(customer)

	mockCustomerRepo.EXPECT().GetCustomerByNIK(gomock.Any()).Times(1)
	mockCustomerRepo.EXPECT().RegisterCustomer(gomock.Any()).Times(1)

	// Create a request
	req, err := http.NewRequest("POST", "/auth/register", bytes.NewReader(body))
	assert.NoError(t, err)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler's RegisterCustomer method
	http.HandlerFunc(handler.RegisterCustomer).ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestLoginHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCustomerRepo := mocks.NewMockCustomerRepository(ctrl)
	handler := &AuthHandler{
		CustomerRepo:  mockCustomerRepo,
		EncryptionKey: []byte("test-key"),
		JWTSecret:     []byte("test-secret"),
	}

	// Create a request body
	credentials := model.AuthLogin{
		NIK:      "1231223",
		Password: "password",
	}
	body, _ := json.Marshal(credentials)

	mockCustomerRepo.EXPECT().GetCustomerByNIK(gomock.Any()).Return(&model.Customer{
		ID:       1,
		Password: hashPassword("password"),
	}, nil).Times(1)

	// Create a request
	req, err := http.NewRequest("POST", "/auth/login", bytes.NewReader(body))
	assert.NoError(t, err)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler's LoginHandler method
	http.HandlerFunc(handler.LoginHandler).ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check if token is present in the response body
	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "token")
}
