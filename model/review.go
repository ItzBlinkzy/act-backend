package model

import "time"

type Review struct {
	ID          int        `json:"id" db:"id"`
	Description string     `json:"description" db:"description"`
	Stars       int        `json:"stars" db:"stars"`
	UserID      int        `json:"user_id" db:"user_id"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}
