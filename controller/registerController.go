package controller

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/itzblinkzy/act-backend/model"
	"github.com/itzblinkzy/act-backend/repository"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Register(c echo.Context) error {
	var payload model.RegistrationPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request payload")
	}


	if payload.Email == "" || payload.Password == "" || payload.FirstName == "" || payload.LastName == "" || payload.TypeUserId == 0 || payload.LoginMethod == "" {
		return c.JSON(http.StatusBadRequest, "Tutti i campi sono obbligatori")
	}

	if len(payload.Password) < 6 || !containsDigit(payload.Password) || !containsLetter(payload.Password) || len(payload.Password) > 20 {
		return c.JSON(http.StatusBadRequest, "La password deve contenere da 6 a 20 caratteri e includere almeno un numero e una lettera")
	}

	// Check for existing email
	existingUser, err := repository.UserRepo.FindByEmail(payload.Email)
	if err == nil && existingUser != nil {
		return c.JSON(http.StatusBadRequest, "Utente gia esistente")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Impossibile eseguire l'hashing della password")
	}

	// Create the user
	newUser := &model.User{
		FirstName:  payload.FirstName,
		LastName:   payload.LastName,
		Email:      payload.Email,
		Password:   string(hashedPassword),
		TypeUserId: payload.TypeUserId,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		LoginMethod: payload.LoginMethod,
	}

	if err := repository.UserRepo.CreateUser(newUser); err != nil {
		log.Printf("Failed to create user: %v", err)
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to create user: %v", err))
	}

	return c.JSON(http.StatusOK, "Utente creato con successo!")
}

func containsDigit(s string) bool {
	for _, r := range s {
		if '0' <= r && r <= '9' {
			return true
		}
	}
	return false
}

func containsLetter(s string) bool {
	for _, r := range s {
		if ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') {
			return true
		}
	}
	return false
}
