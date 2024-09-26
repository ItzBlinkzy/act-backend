package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/itzblinkzy/act-backend/database"
	"github.com/itzblinkzy/act-backend/model"
)

var UserRepo = &UserRepository{}

// UserRepository struct with no fields, using singleton database connection.
type UserRepository struct{}

// FetchAll retrieves all users with optional filters.
func (r *UserRepository) FetchAll(filters *model.User) ([]*model.User, int64, error) {
	var users []*model.User
	var query strings.Builder
	var args []interface{}

	// Constructing the query based on filters provided
	if filters != nil {
		query.WriteString(" WHERE ")
		clauses := []string{}
		if filters.FirstName != "" {
			clauses = append(clauses, "first_name = ?")
			args = append(args, filters.FirstName)
		}
		if filters.LastName != "" {
			clauses = append(clauses, "last_name = ?")
			args = append(args, filters.LastName)
		}
		if filters.Email != "" {
			clauses = append(clauses, "email = ?")
			args = append(args, filters.Email)
		}
		query.WriteString(strings.Join(clauses, " AND "))
	}

	var totalCount int64
	err := database.GetDB().Get(&totalCount, "SELECT COUNT(*) FROM users "+query.String(), args...)
	if err != nil {
		return nil, 0, err
	}

	err = database.GetDB().Select(&users, "SELECT * FROM users "+query.String(), args...)
	if err != nil {
		return nil, 0, err
	}

	return users, totalCount, nil
}

// CreateUser inserts a new user into the database.
func (r *UserRepository) CreateUser(user *model.User) error {
	db := database.GetDB()
	_, err := db.Exec(
		"INSERT INTO users (first_name, last_name, email, password, type_user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		user.FirstName, user.LastName, user.Email, user.Password, user.TypeUserId, user.CreatedAt, user.UpdatedAt)
	return err
}

// FindByEmail finds a user by their email.
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := database.GetDB().Get(&user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("No user found for email: %s\n", email)
			return nil, nil
		}
		fmt.Printf("Database error during user fetch: %v\n", err)
		return nil, err // Database error
	}
	return &user, nil
}

// UpdateUser updates an existing user's information.
func (r *UserRepository) UpdateUser(user *model.User) error {
	db := database.GetDB()
	_, err := db.Exec(
		"UPDATE users SET first_name = $1, last_name = $2, email = $3, password = $4, updated_at = NOW() WHERE id = $5",
		user.FirstName, user.LastName, user.Email, user.Password, user.ID)
	return err
}

// DeleteUser performs a soft delete on a user by setting the deleted_at timestamp.
func (r *UserRepository) DeleteUser(id uint) error {
	db := database.GetDB()
	// Use the current time to mark as deleted
	_, err := db.Exec("UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL", id)
	return err
}

// GetAllUsers retrieves all users from the database.
func (r *UserRepository) GetAllUsers() ([]*model.User, error) {
	var users []*model.User
	err := database.GetDB().Select(&users, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	return users, nil
}

// FindById finds a user by their ID.
func (r *UserRepository) FindById(id uint) (*model.User, error) {
	var user model.User
	err := database.GetDB().Get(&user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No user found
		}
		return nil, err // Database error
	}
	return &user, nil
}
