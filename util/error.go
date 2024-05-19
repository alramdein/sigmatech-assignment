package util

import (
	"errors"

	"github.com/go-sql-driver/mysql"
)

func CheckMySQLError(err error) error {
	me, ok := err.(*mysql.MySQLError)
	if !ok {
		return err
	}
	if me.Number == 1062 {
		return errors.New("It already exists in a database.")
	}
	return err
}
