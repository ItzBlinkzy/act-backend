package repository

import (
	"github.com/itzblinkzy/act-backend/database"
)

type PaymentRepository struct{}

var PaymentRepo = &PaymentRepository{}

func (repo *PaymentRepository) AddCreditToUser(userID string, credit float64) error {
	sql := "UPDATE users SET credit = credit + $1, updated_at = NOW() WHERE id = $2"

	_, err := database.GetDB().Exec(sql, credit, userID)
	if err != nil {
		return err
	}
	return nil
}
