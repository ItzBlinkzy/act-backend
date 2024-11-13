package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/itzblinkzy/act-backend/database"
	"github.com/itzblinkzy/act-backend/model"
	"golang.org/x/crypto/bcrypt"
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
		"INSERT INTO users (first_name, last_name, email, password, type_user_id, created_at, updated_at, login_method) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		user.FirstName, user.LastName, user.Email, user.Password, user.TypeUserId, user.CreatedAt, user.UpdatedAt, user.LoginMethod)
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

// UpdateUser updates an existing user's information, optionally updating the password if provided.
func (r *UserRepository) UpdateUser(user *model.User, updatePassword bool) error {
	db := database.GetDB()

	// Prepare query based on provided data
	query := "UPDATE users SET first_name = $1, last_name = $2, email = $3, updated_at = $4"
	args := []interface{}{user.FirstName, user.LastName, user.Email, time.Now(), user.ID}

	if updatePassword && user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
		query += ", password = $5"
		args = append(args, user.Password)
	}

	query += " WHERE id = $6"

	result, err := db.Exec(query, args...)
	if err != nil {
		fmt.Printf("Database error during user update: %v\n", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Error checking rows affected: %v\n", err)
		return err
	}

	if rowsAffected == 0 {
		fmt.Println("No user found with the provided ID")
		return sql.ErrNoRows // or you can return a custom error if preferred
	}

	return nil
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

func (repo *UserRepository) UpdateUserCredit(userID int, amount float64) error {
	_, err := database.GetDB().Exec("UPDATE users SET credit = credit + $1 WHERE id = $2", amount, userID)
	if err != nil {
		log.Printf("Failed to update credit for user %d: %v", userID, err)
		return err
	}
	return nil
}
