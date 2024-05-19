package handler

import (
	"alif-sigmatech/model"
	"alif-sigmatech/repository"
	"alif-sigmatech/util"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	CustomerRepo  repository.CustomerRepository
	JWTSecret     []byte
	EncryptionKey []byte
}

// NewAuthHandler creates a new instance of AuthHandler
func NewAuthHandler(repo repository.CustomerRepository, jwtSecret []byte,
	EncryptionKey []byte) *AuthHandler {
	return &AuthHandler{
		CustomerRepo:  repo,
		JWTSecret:     jwtSecret,
		EncryptionKey: EncryptionKey,
	}
}

// RegisterCustomer handles registration of a new consumer
func (h *AuthHandler) RegisterCustomer(w http.ResponseWriter, r *http.Request) {
	var customer model.Customer
	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = h.validateRegisterCustomerInput(customer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if customer.KTPPhoto != nil && len(customer.KTPPhoto) > 0 {
		// Encrypt sensitive data before storing it on the database
		encryptedKTPPhoto, err := util.EncryptData(customer.KTPPhoto, h.EncryptionKey)
		if err != nil {
			logrus.Error(err)
			http.Error(w, "Failed to encrypt KTP photo", http.StatusInternalServerError)
			return
		}
		customer.KTPPhoto = encryptedKTPPhoto
	}

	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(customer.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	customer.Password = string(hashedPassword)

	cust, err := h.CustomerRepo.GetCustomerByNIK(customer.NIK)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "Failed to register consumer", http.StatusInternalServerError)
		return
	}
	if cust != nil {
		err := errors.New("NIK already exist")
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	err = h.CustomerRepo.RegisterCustomer(&customer)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "Failed to register consumer", http.StatusInternalServerError)
		return
	}
	customer.Password = "" // obfuscate

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials model.AuthLogin
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Fetch the customer by username
	customer, err := h.CustomerRepo.GetCustomerByNIK(credentials.NIK)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if customer == nil {
		fmt.Println("sinikah?")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Unauthorized")
		return
	}

	// Compare the provided password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(credentials.Password))
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Unauthorized")
		return
	}

	// Create JWT token
	expirationTime := time.Now().Add(30 * time.Minute)
	claims := &model.Claims{
		NIK:      customer.NIK,
		FullName: customer.FullName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.JWTSecret)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}

	// Send the token to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func (h *AuthHandler) validateRegisterCustomerInput(customer model.Customer) error {
	if customer.NIK == "" {
		return errors.New("NIK is required")
	}
	if customer.Password == "" {
		return errors.New("Password is required")
	}
	if customer.FullName == "" {
		return errors.New("FullName is required")
	}
	if customer.LegalName == "" {
		return errors.New("LegalName is required")
	}
	return nil
}
