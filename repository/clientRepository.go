package repository

import (
	"fmt"
	"time"

	"github.com/itzblinkzy/act-backend/database"
	"github.com/itzblinkzy/act-backend/model"
	"github.com/lib/pq"
)

var ClientRepo = &ClientRepository{}

// ClientRepository struct with no fields, using singleton database connection.
type ClientRepository struct{}

// CreateClient inserts a new client into the database and returns the new client's ID.
func CreateClient(client model.Client) (uint, error) {
	var clientId uint
	err := database.GetDB().QueryRow(
		"INSERT INTO clients (company_name, created_at, updated_at) VALUES ($1, $2, $3) RETURNING id",
		client.CompanyName, time.Now(), time.Now(),
	).Scan(&clientId)
	if err != nil {
		return 0, err
	}
	return clientId, nil
}

// CreateAssociation tries to create a client-manager association and handles potential foreign key errors.
func CreateAssociation(association model.ClientManagerAssociation) error {
	tx, err := database.GetDB().Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		"INSERT INTO client_manager_association (manager_id, client_id, created_at, updated_at) VALUES ($1, $2, NOW(), NOW())",
		association.ManagerId, association.ClientId,
	)
	if err != nil {
		tx.Rollback()                       // Roll back the transaction on error
		if err, ok := err.(*pq.Error); ok { // Check if it is a pq error
			if err.Code == "23503" { // PostgreSQL foreign key violation error code
				return fmt.Errorf("foreign key violation: %v", err)
			}
		}
		return err
	}

	return tx.Commit()
}

// UpdateClient updates the company name for an existing client.
func UpdateClient(clientId uint, companyName string) error {
	_, err := database.GetDB().Exec(
		"UPDATE clients SET company_name = $1, updated_at = $2 WHERE id = $3",
		companyName, time.Now(), clientId,
	)
	return err
}

func SoftDeleteClient(clientId uint) error {
	now := time.Now()
	tx, err := database.GetDB().Begin()
	if err != nil {
		return err
	}

	// Soft delete client manager associations
	if _, err := tx.Exec("UPDATE client_manager_association SET deleted_at = $1 WHERE client_id = $2", now, clientId); err != nil {
		tx.Rollback()
		return err
	}

	// Soft delete the client
	if _, err := tx.Exec("UPDATE clients SET deleted_at = $1 WHERE id = $2", now, clientId); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// GetClient retrieves a single client by its ID.
func GetClient(clientId uint) (*model.Client, error) {
	var client model.Client
	err := database.GetDB().QueryRow(
		"SELECT id, company_name, created_at, updated_at FROM clients WHERE id = $1 AND deleted_at IS NULL",
		clientId,
	).Scan(&client.ID, &client.CompanyName, &client.CreatedAt, &client.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &client, nil
}

// GetAllClients retrieves all clients associated with a given manager.
func GetAllClients(managerId uint) ([]model.Client, error) {
	var clients []model.Client
	rows, err := database.GetDB().Query(
		"SELECT c.id, c.company_name, c.created_at, c.updated_at FROM clients c JOIN client_manager_association cma ON c.id = cma.client_id WHERE cma.manager_id = $1 AND c.deleted_at IS NULL",
		managerId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var client model.Client
		if err := rows.Scan(&client.ID, &client.CompanyName, &client.CreatedAt, &client.UpdatedAt); err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return clients, nil
}

func GetStocksOfClient(clientId uint) ([]model.Stock, error) {
	var stocks []model.Stock
	rows, err := database.GetDB().Query(
		"SELECT id, user_id, ticker, quantity_owned, quantity_sold, client_id, created_at, updated_at FROM bought_stocks WHERE client_id = $1 AND deleted_at IS NULL",
		clientId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var stock model.Stock
		if err := rows.Scan(&stock.ID, &stock.UserId, &stock.Ticker, &stock.QuantityOwned, &stock.QuantitySold, &stock.ClientId, &stock.CreatedAt, &stock.UpdatedAt); err != nil {
			return nil, err
		}
		stocks = append(stocks, stock)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return stocks, nil
}
