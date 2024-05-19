package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"alif-sigmatech/handler"
	"alif-sigmatech/middleware"
	"alif-sigmatech/repository"
)

// AppConfig contains the application configurations
type AppConfig struct {
	DB            *sql.DB
	jwtSecret     []byte
	encryptionKey []byte
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Open a connection to the database
	db, err := sql.Open("mysql", composeMySQLConnectionString())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize AppConfig with the database connection
	appConfig := &AppConfig{
		DB:            db,
		jwtSecret:     []byte(os.Getenv("JWT_SECRET")),
		encryptionKey: []byte(os.Getenv("ENCRYPTION_KEY")),
	}

	// Initialize router
	r := mux.NewRouter()

	// handlers
	registerHandlers(r, appConfig)

	// Start server
	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", r)
}

// registerHandlers registers all HTTP handlers
func registerHandlers(r *mux.Router, appConfig *AppConfig) {
	customerRepo := repository.NewMySQLCustomerRepository(appConfig.DB)
	transactionRepo := repository.NewMySQLTransactionRepository(appConfig.DB)
	limitRepo := repository.NewMySQLLimitRepository(appConfig.DB)

	authHandler := handler.NewAuthHandler(customerRepo, appConfig.jwtSecret, appConfig.encryptionKey)
	transactionhHandler := handler.NewTransactionHandler(transactionRepo, limitRepo)
	limitHandler := handler.NewLimitHandler(limitRepo, customerRepo)

	r.HandleFunc("/auth/register", authHandler.RegisterCustomer).Methods("POST")
	r.HandleFunc("/auth/login", authHandler.LoginHandler).Methods("POST")

	fundRouter := r.PathPrefix("/fund").Subrouter()
	fundRouter.Use(middleware.JWTMiddleware(appConfig.jwtSecret))

	fundRouter.HandleFunc("/transaction", transactionhHandler.CreateTransaction).Methods("POST")
	fundRouter.HandleFunc("/limit", limitHandler.CreateLimit).Methods("POST")
}

func composeMySQLConnectionString() string {
	// "root:password@tcp(localhost:3306)/mydb",
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
}
