package repository

import (
	"errors"
	"log"

	"github.com/itzblinkzy/act-backend/database"
	"github.com/itzblinkzy/act-backend/model"
)

var StockRepo = &StockRepository{}

type StockRepository struct{}

func (r *StockRepository) ListBoughtStocks(userId uint) ([]model.Stock, error) {
	var stocks []model.Stock
	query := "SELECT * FROM bought_stocks WHERE user_id = $1 AND deleted_at IS NULL"
	log.Printf("Executing query: %s with userID: %d", query, userId)
	err := database.GetDB().Select(&stocks, query, userId)
	if err != nil {
		log.Printf("Error fetching stocks for user %d: %v", userId, err)
		return nil, err
	}
	log.Printf("Successfully fetched stocks for user %d", userId)
	return stocks, nil
}

func (r *StockRepository) BuyStock(stock model.StockBuy, quantity int) error {
	_, err := database.GetDB().Exec("INSERT INTO bought_stocks (user_id, ticker, quantity_owned, quantity_sold, created_at, updated_at) VALUES ($1, $2, $3, 0, NOW(), NOW())",
		stock.UserId, stock.Ticker, quantity)
	if err != nil {
		log.Printf("Failed to buy stock: %v", err)
		return err
	}
	return nil
}

func (r *StockRepository) UpdateStock(id uint, update model.StockUpdate) error {
	if update.SellingQuantity > 0 {
		// Check if sufficient quantity owned before selling
		var currentOwned int
		err := database.GetDB().QueryRow("SELECT quantity_owned FROM bought_stocks WHERE id = $1", id).Scan(&currentOwned)
		if err != nil {
			log.Printf("Error checking current stock quantity: %v", err)
			return err
		}
		if currentOwned < update.SellingQuantity {
			return errors.New("not enough stock owned to sell the requested amount")
		}

		// Proceed with update if sufficient
		_, err = database.GetDB().Exec("UPDATE bought_stocks SET quantity_owned = quantity_owned - $1, quantity_sold = quantity_sold + $1 WHERE id = $2",
			update.SellingQuantity, id)
		if err != nil {
			log.Printf("Failed to update stock: %v", err)
			return err
		}
	} else if update.BuyingQuantity > 0 {
		_, err := database.GetDB().Exec("UPDATE bought_stocks SET quantity_owned = quantity_owned + $1 WHERE id = $2",
			update.BuyingQuantity, id)
		if err != nil {
			log.Printf("Failed to update stock: %v", err)
			return err
		}
	}
	return nil
}
