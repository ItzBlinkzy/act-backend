package repository

import (
	"database/sql"
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

func (r *StockRepository) BuyStock(stock model.StockBuy) error {
	var stockId int

	err := database.GetDB().QueryRow(
		"INSERT INTO bought_stocks (user_id, client_id, ticker, quantity_owned, quantity_sold, created_at, updated_at) VALUES ($1, $2, $3, $4, 0, NOW(), NOW()) RETURNING id",
		stock.UserId, stock.Clientid, stock.Ticker, stock.BuyingQuantity).Scan(&stockId)
	if err != nil {
		log.Printf("Failed to buy stock: %v", err)
		return err
	}

	_, err = database.GetDB().Exec(
		"INSERT INTO logs_bought_stocks (user_id, bought_stock_id, client_id, quantity_bought, quantity_sold, created_at) VALUES ($1, $2, $3, $4, NULL, NOW())",
		stock.UserId, stockId, stock.Clientid, stock.BuyingQuantity)
	if err != nil {
		log.Printf("Failed to log stock purchase: %v", err)
		return err
	}

	return nil
}

func (r *StockRepository) UpdateStock(id uint, update model.StockUpdate) error {
	var userId uint
	var currentOwned int

	err := database.GetDB().QueryRow("SELECT user_id, quantity_owned FROM bought_stocks WHERE id = $1", id).Scan(&userId, &currentOwned)
	if err != nil {
		log.Printf("Error retrieving stock data: %v", err)
		return err
	}

	if currentOwned < update.SellingQuantity {
		return errors.New("not enough stock owned to sell the requested amount")
	}

	_, err = database.GetDB().Exec(
		"UPDATE bought_stocks SET quantity_owned = quantity_owned - $1, quantity_sold = quantity_sold + $1 WHERE id = $2",
		update.SellingQuantity, id)
	if err != nil {
		log.Printf("Failed to update stock: %v", err)
		return err
	}

	var clientId sql.NullInt64
	if update.ClientId != nil {
		clientId = sql.NullInt64{Int64: int64(*update.ClientId), Valid: true}
	} else {
		clientId = sql.NullInt64{Valid: false}
	}

	_, err = database.GetDB().Exec(
		"INSERT INTO logs_bought_stocks (user_id, bought_stock_id, client_id, quantity_bought, quantity_sold, created_at) VALUES ($1, $2, $3, NULL, $4, NOW())",
		userId, id, clientId, update.SellingQuantity)
	if err != nil {
		log.Printf("Failed to log stock sale: %v", err)
		return err
	}

	return nil
}

func (r *StockRepository) GetLogsByUserId(userId int) ([]model.Log, error) {
	var logs []model.Log
	err := database.GetDB().Select(&logs, "SELECT * FROM logs_bought_stocks WHERE user_id = $1", userId)
	if err != nil {
		log.Printf("Error fetching logs for user %d: %v", userId, err)
		return nil, err
	}
	return logs, nil
}

func (r *StockRepository) GetLogsByStockId(stockId int) ([]model.Log, error) {
	var logs []model.Log
	err := database.GetDB().Select(&logs, "SELECT * FROM logs_bought_stocks WHERE bought_stock_id = $1", stockId)
	if err != nil {
		return nil, err
	}
	return logs, nil
}
