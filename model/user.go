package model

import "time"

type User struct {
	ID         uint       `db:"id" json:"id"`
	FirstName  string     `db:"first_name" json:"first_name"`
	LastName   string     `db:"last_name" json:"last_name"`
	Email      string     `db:"email" json:"email"`
	Password   string     `db:"password" json:"password"`
	TypeUserId uint       `db:"type_user_id" json:"type_user_id"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt  *time.Time `db:"deleted_at" json:"deleted_at"`
}

type RegistrationPayload struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	TypeUserId uint   `json:"type_user_id"`
	Password   string `json:"password"`
}
