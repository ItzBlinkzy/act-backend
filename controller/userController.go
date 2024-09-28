// controller/user_controller.go

package controller

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/itzblinkzy/act-backend/model"
	"github.com/itzblinkzy/act-backend/repository"
	"github.com/labstack/echo/v4"
)

var userRepo = &repository.UserRepository{} // Using the global UserRepository instance

func RegisterUser(c echo.Context) error {
	var user model.User
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	// Check if user exists
	existingUser, err := userRepo.FindByEmail(user.Email)
	if err == nil && existingUser != nil {
		return echo.NewHTTPError(http.StatusConflict, "User already exists")
	}

	// Create user
	if err := userRepo.CreateUser(&user); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
	}

	return c.JSON(http.StatusCreated, user)
}

func GetUser(c echo.Context) error {
	email := c.Param("email")
	user, err := userRepo.FindByEmail(email)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	user.Password = "" // Do not expose the password

	return c.JSON(http.StatusOK, user)
}

func GetAllUsers(c echo.Context) error {
	users, err := userRepo.GetAllUsers()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch users")
	}

	for i := range users {
		users[i].Password = "" // Set each user's password to empty
	}

	return c.JSON(http.StatusOK, users)
}

func getUserFromContext(c echo.Context) (*model.User, error) {
	fmt.Println(os.Getenv("ENV"), "envv")
	if os.Getenv("ENV") == "development" {
		// Access request headers in development mode
		userEmailHeader := c.Request().Header.Get("User-email")
		fmt.Println("Received User-email header: ", userEmailHeader)
		if userEmailHeader == "" {
			return nil, echo.NewHTTPError(400, "User-email header is missing")
		}

		userEmail := userEmailHeader
		// Remove the err check here since `userEmailHeader` is already a string
		user, err := repository.UserRepo.FindByEmail(userEmail)
		if err != nil {
			return nil, err
		}
		// Optionally check if user is deleted (you commented this out)
		// if user.IsDeleted {
		// 	return nil, errors.New("user is deleted")
		// }
		return user, nil
	}
	jwtError := errors.New("JWT token missing or invalid")
	currentUser, ok := (c).Get("user").(*jwt.Token)
	if !ok {
		fmt.Println("-2")
		return nil, jwtError
	}

	fmt.Println("-1")

	claims, ok := currentUser.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwtError
	}

	fmt.Println("1")

	email, ok := claims["email"].(string)
	if !ok {
		return nil, jwtError
	}

	fmt.Println("2")

	user, err := userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	fmt.Println("3")

	// if user.IsDeleted {
	// 	return nil, errors.New("user is deleted")
	// }

	return user, nil
}

func GetCurrentUser(c echo.Context) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.Error{Message: "Unauthorized"})
	}

	user.Password = "" // Do not expose the password

	return c.JSON(http.StatusOK, user)
}

func UpdateUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	var user model.User
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	user.ID = uint(id)

	if err := repository.UserRepo.UpdateUser(&user); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update user")
	}

	return c.JSON(http.StatusOK, "User updated successfully")
}
