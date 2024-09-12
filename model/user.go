package model

import "time"

type User struct {
	ID        uint       `db:"id" json:"id"`
	FirstName string     `db:"firstName" json:"firstName"`
	LastName  string     `db:"lastName" json:"lastName"`
	Email     string     `db:"email" json:"email"`
	Password  string     `db:"password" json:"password"`
	TypeId    uint       `db:"typeId" json:"typeId"`
	CreatedAt time.Time  `db:"createdAt" json:"createdAt"`
	UpdatedAt time.Time  `db:"updatedAt" json:"updatedAt"`
	DeletedAt *time.Time `db:"deletedAt" json:"deletedAt"`
}

type RegistrationPayload struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	TypeId    uint   `json:"typeId"`
}
