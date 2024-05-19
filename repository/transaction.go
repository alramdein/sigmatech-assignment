package repository

import (
	"alif-sigmatech/model"
	"database/sql"
)

type TransactionRepository interface {
	CreateTransaction(transaction *model.Transaction) error
}

type MySQLTransactionRepository struct {
	DB *sql.DB
}

// NewMySQLCustomerRepository creates a new instance of MySQLTransactionRepository
func NewMySQLTransactionRepository(db *sql.DB) *MySQLTransactionRepository {
	return &MySQLTransactionRepository{
		DB: db,
	}
}

func (repo *MySQLTransactionRepository) CreateTransaction(transaction *model.Transaction) error {
	query := "INSERT INTO transaction (customer_id, contract_number, otr, admin_fee, installment_amount, interest_amount, asset_name, tenor) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := repo.DB.Exec(query, transaction.CustomerID, transaction.ContractNumber, transaction.OTR, transaction.AdminFee, transaction.InstallmentAmount, transaction.InterestAmount, transaction.AssetName, transaction.Tenor)
	return err
}
