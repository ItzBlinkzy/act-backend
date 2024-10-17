package model

import "time"

type Stock struct {
	ID            uint       `db:"id" json:"id"`
	UserId        uint       `db:"user_id" json:"user_id"`
	Ticker        string     `db:"ticker" json:"ticker"`
	QuantityOwned *int       `db:"quantity_owned" json:"quantity_owned"`
	QuantitySold  *int       `db:"quantity_sold" json:"quantity_sold"`
	ClientId      *uint      `db:"client_id" json:"client_id"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt     *time.Time `db:"deleted_at" json:"deleted_at"`
}

type StockBuy struct {
	UserId         uint   `db:"user_id" json:"user_id"`
	Clientid       *uint  `db:"client_id" json:"client_id"`
	Ticker         string `json:"ticker"`
	BuyingQuantity int    `json:"buying_quantity"`
}

type StockUpdate struct {
	Ticker          string `json:"ticker"`
	BuyingQuantity  int    `json:"buying_quantity"`
	SellingQuantity int    `json:"selling_quantity"`
	ClientId        *uint  `db:"client_id" json:"client_id"`
}
