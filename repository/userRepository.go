package repository

import (
	"database/sql"
	"log"
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
			clauses = append(clauses, "firstName = ?")
			args = append(args, filters.FirstName)
		}
		if filters.LastName != "" {
			clauses = append(clauses, "lastName = ?")
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
		"INSERT INTO users (firstName, lastName, email, password, companyId, createdAt, updatedAt) VALUES (?, ?, ?, ?, ?, ?, ?)",
		user.FirstName, user.LastName, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	return err
}

// FindByEmail finds a user by their email.
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := database.GetDB().Get(&user, "SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No user found
		}
		return nil, err // Database error
	}
	return &user, nil
}

// UpdateUser updates an existing user's information.
func (r *UserRepository) UpdateUser(user *model.User) error {
	db := database.GetDB()
	_, err := db.Exec(
		"UPDATE users SET firstName = ?, lastName = ?, email = ?, password = ?, companyId = ?, updatedAt = NOW() WHERE id = ?",
		user.FirstName, user.LastName, user.Email, user.Password, user.ID)
	return err
}

// DeleteUser removes a user from the database.
func (r *UserRepository) DeleteUser(id uint) error {
	db := database.GetDB()
	_, err := db.Exec("DELETE FROM users WHERE id = ?", id)
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
func (r *UserRepository) GetAllUsersByCompanyId(companyId uint) ([]*model.User, error) {
	var users []*model.User
	query := "SELECT id, firstName, lastName, email, password, companyId, createdAt, updatedAt FROM users WHERE companyId = ?"
	err := database.GetDB().Select(&users, query, companyId)
	if err != nil {
		log.Printf("Error fetching users for company: %v", err)
		return nil, err
	}
	return users, nil
}

// FindById finds a user by their ID.
func (r *UserRepository) FindById(id uint) (*model.User, error) {
	var user model.User
	err := database.GetDB().Get(&user, "SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No user found
		}
		return nil, err // Database error
	}
	return &user, nil
}
