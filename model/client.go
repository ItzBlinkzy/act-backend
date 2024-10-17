package model

import "time"

type Client struct {
	ID          uint       `db:"id" json:"id"`
	CompanyName string     `db:"company_name" json:"company_name"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at"`
}

type ClientManagerAssociation struct {
	ID        uint       `db:"id" json:"id"`
	ManagerId uint       `db:"manager_id" json:"manager_id"`
	ClientId  uint       `db:"client_id" json:"client_id"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
}
