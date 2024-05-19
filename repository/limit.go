package repository

import (
	"alif-sigmatech/model"
	"database/sql"
)

type LimitRepository interface {
	GetLimitByCustomerID(customerID int) (*model.Limit, error)
	CreateLimit(limit *model.Limit) error
}

type MySQLLimitRepository struct {
	DB *sql.DB
}

// NewMySQLLimitRepository creates a new instance of MySQLLimitRepository
func NewMySQLLimitRepository(db *sql.DB) *MySQLLimitRepository {
	return &MySQLLimitRepository{
		DB: db,
	}
}

func (repo *MySQLLimitRepository) GetLimitByCustomerID(customerID int) (*model.Limit, error) {
	query := "SELECT customer_id, tenor_1, tenor_2, tenor_3, tenor_4 FROM `limit` WHERE customer_id = ?"
	row := repo.DB.QueryRow(query, customerID)

	var limit model.Limit
	err := row.Scan(&limit.CustomerID, &limit.Tenor1, &limit.Tenor2, &limit.Tenor3, &limit.Tenor4)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No limit found with the given ID
		}
		return nil, err
	}

	return &limit, nil
}

func (repo *MySQLLimitRepository) CreateLimit(limit *model.Limit) error {
	query := "INSERT INTO `limit` (customer_id, tenor_1, tenor_2, tenor_3, tenor_4) VALUES (?, ?, ?, ?, ?)"
	_, err := repo.DB.Exec(query, limit.CustomerID, limit.Tenor1, limit.Tenor2, limit.Tenor3, limit.Tenor4)
	return err
}
