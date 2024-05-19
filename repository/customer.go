package repository

import (
	"database/sql"
	"log"

	"alif-sigmatech/model"
)

// CustomerRepository defines the interface for customer data access
type CustomerRepository interface {
	RegisterCustomer(customer *model.Customer) error
	GetCustomerByNIK(nik string) (*model.Customer, error)
	GetCustomerByID(id int) (*model.Customer, error)
}

// MySQLCustomerRepository is a repository implementation using MySQL
type MySQLCustomerRepository struct {
	DB *sql.DB
}

// NewMySQLCustomerRepository creates a new instance of MySQLCustomerRepository
func NewMySQLCustomerRepository(db *sql.DB) *MySQLCustomerRepository {
	return &MySQLCustomerRepository{
		DB: db,
	}
}

// RegisterCustomer registers a new consumer
func (repo *MySQLCustomerRepository) RegisterCustomer(customer *model.Customer) error {
	query := "INSERT INTO customer (nik, full_name, password, legal_name, birth_place, birth_date, salary, ktp_photo, selfie_photo) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"

	_, err := repo.DB.Exec(query, customer.NIK, customer.FullName, customer.Password, customer.LegalName, customer.BirthPlace, customer.BirthDate, customer.Salary, customer.KTPPhoto, customer.SelfiePhoto)
	if err != nil {
		return err
	}

	return nil
}

// GetCustomerByNIK mengambil data pelanggan berdasarkan NIK dari database
func (repo *MySQLCustomerRepository) GetCustomerByNIK(nik string) (*model.Customer, error) {
	customer := &model.Customer{}
	query := "SELECT id, nik, full_name, password, legal_name, birth_place, birth_date, salary FROM customer WHERE nik = ?"

	err := repo.DB.QueryRow(query, nik).Scan(
		&customer.ID,
		&customer.NIK,
		&customer.FullName,
		&customer.Password,
		&customer.LegalName,
		&customer.BirthPlace,
		&customer.BirthDate,
		&customer.Salary,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No customer found with the given NIK
		}
		log.Printf("Error fetching customer by NIK: %v", err)
		return nil, err
	}

	return customer, nil
}

// GetCustomerByID mengambil data pelanggan berdasarkan ID dari database
func (repo *MySQLCustomerRepository) GetCustomerByID(id int) (*model.Customer, error) {
	customer := &model.Customer{}
	query := "SELECT id, nik, full_name, password, legal_name, birth_place, birth_date, salary FROM customer WHERE id = ?"

	err := repo.DB.QueryRow(query, id).Scan(
		&customer.ID,
		&customer.NIK,
		&customer.FullName,
		&customer.Password,
		&customer.LegalName,
		&customer.BirthPlace,
		&customer.BirthDate,
		&customer.Salary,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No customer found with the given ID
		}
		log.Printf("Error fetching customer by ID: %v", err)
		return nil, err
	}

	return customer, nil
}
